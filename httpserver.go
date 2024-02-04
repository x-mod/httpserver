package httpserver

import (
	"context"
	"crypto/tls"
	"fmt"
	"io"
	"log"
	"net"
	"net/http"
	"runtime"
	"strings"
	"sync"
	"time"

	"github.com/gorilla/mux"
	"github.com/x-mod/event"
	"github.com/x-mod/glog"
	"golang.org/x/net/trace"
)

type Server struct {
	name    string
	addr    string
	http    *http.Server
	tls     *tls.Config
	routes  *mux.Router
	handler http.Handler
	stopped *event.Event
	serving *event.Event
	traced  bool
	events  trace.EventLog
	mu      sync.Mutex
}

type ServerOpt func(*Server)

func Name(name string) ServerOpt {
	return func(srv *Server) {
		srv.name = name
	}
}

func Address(addr string) ServerOpt {
	return func(srv *Server) {
		srv.addr = addr
	}
}

func TLSConfig(cf *tls.Config) ServerOpt {
	return func(srv *Server) {
		srv.tls = cf
	}
}

func NetTrace(flag bool) ServerOpt {
	return func(srv *Server) {
		srv.traced = flag
	}
}

func HTTPHandler(handler http.Handler) ServerOpt {
	return func(srv *Server) {
		srv.handler = handler
	}
}

type MiddlewareFunc func(http.Handler) http.Handler

func Middleware(m MiddlewareFunc) ServerOpt {
	return func(srv *Server) {
		srv.routes.Use(mux.MiddlewareFunc(m))
	}
}

func New(opts ...ServerOpt) *Server {
	srv := &Server{
		name: "httpserver",
		http: &http.Server{
			ReadTimeout:  15 * time.Second,
			WriteTimeout: 15 * time.Second,
			IdleTimeout:  60 * time.Second,
		},
		routes:  mux.NewRouter(),
		stopped: event.New(),
		serving: event.New(),
	}
	for _, opt := range opts {
		opt(srv)
	}
	if srv.traced {
		_, file, line, _ := runtime.Caller(1)
		srv.events = trace.NewEventLog(srv.name, fmt.Sprintf("%s:%d", file, line))
	}
	// Disable net/http internal logging
	srv.http.ErrorLog = log.New(io.Discard, "", 0)
	return srv
}

func NewServer(opts ...ServerOpt) *Server {
	return New(opts...)
}

type RouteCfg struct {
	schemes []string
	host    string
	methods []string
	prefix  string
	pattern string
	headers []string
	queries []string
	handler http.Handler
}
type RouteOpt func(cf *RouteCfg)

func Schema(schemas ...string) RouteOpt {
	return func(cf *RouteCfg) {
		cf.schemes = append(cf.schemes, schemas...)
	}
}
func Host(host string) RouteOpt {
	return func(cf *RouteCfg) {
		cf.host = host
	}
}
func Method(methods ...string) RouteOpt {
	return func(cf *RouteCfg) {
		for _, m := range methods {
			cf.methods = append(cf.methods, strings.ToUpper(m))
		}
	}
}
func Prefix(prefix string) RouteOpt {
	return func(cf *RouteCfg) {
		cf.prefix = prefix
	}
}
func Pattern(pattern string) RouteOpt {
	return func(cf *RouteCfg) {
		cf.pattern = pattern
	}
}
func Header(headers ...string) RouteOpt {
	return func(cf *RouteCfg) {
		cf.headers = append(cf.headers, headers...)
	}
}
func Query(queries ...string) RouteOpt {
	return func(cf *RouteCfg) {
		cf.queries = append(cf.queries, queries...)
	}
}
func Handler(h http.Handler) RouteOpt {
	return func(cf *RouteCfg) {
		cf.handler = h
	}
}

func (srv *Server) Route(opts ...RouteOpt) {
	cf := &RouteCfg{
		schemes: []string{},
		methods: []string{},
		headers: []string{},
	}
	for _, opt := range opts {
		opt(cf)
	}
	if cf.handler != nil {
		r := srv.routes.NewRoute().Handler(cf.handler)
		if cf.pattern != "" {
			r.Path(cf.pattern)
			srv.printf("route pattern %s", cf.pattern)
		}
		if cf.prefix != "" {
			r.PathPrefix(cf.prefix)
			srv.printf("route prefix %s", cf.prefix)
		}
		if len(cf.schemes) > 0 {
			r.Schemes(cf.schemes...)
			srv.printf("route schemas %s", strings.Join(cf.schemes, "|"))
		}
		if len(cf.methods) > 0 {
			r.Methods(cf.methods...)
			srv.printf("route methods %s", strings.Join(cf.methods, "|"))
		}
		if len(cf.host) > 0 {
			r.Host(cf.host)
			srv.printf("route host %s", cf.host)
		}
		if len(cf.headers) > 0 && len(cf.headers)%2 == 0 {
			r.Headers(cf.headers...)
			srv.printf("route headers %s", strings.Join(cf.headers, "|"))
		}
		if len(cf.queries) > 0 && len(cf.queries)%2 == 0 {
			r.Queries(cf.queries...)
			srv.printf("route queries %s", strings.Join(cf.queries, "|"))
		}
	}
}

// ServeHTTP Implement http.Handler interface
func (srv *Server) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	if srv.handler != nil {
		srv.handler.ServeHTTP(w, req)
		return
	}
	srv.routes.ServeHTTP(w, req)
}

func (srv *Server) Serve(ctx context.Context) error {
	defer func() {
		if srv.events != nil {
			srv.events.Finish()
		}
	}()
	srv.http.Handler = srv
	srv.http.BaseContext = func(net.Listener) context.Context {
		return ctx
	}

	ln, err := net.Listen("tcp", srv.addr)
	if err != nil {
		return err
	}
	if srv.tls != nil {
		ln = tls.NewListener(ln, srv.tls)
	}

	defer srv.stopped.Fire()
	srv.serving.Fire()

	glog.Info(srv.name, " serving at ", srv.addr)
	return srv.http.Serve(ln)
}

func (srv *Server) Shutdown(ctx context.Context) error {
	return srv.http.Shutdown(ctx)
}

func (srv *Server) Serving() <-chan struct{} {
	return srv.serving.Done()
}

func (srv *Server) Close() <-chan struct{} {
	if srv.serving.HasFired() {
		srv.http.Close()
		return srv.stopped.Done()
	}
	return event.Done()
}

func (srv *Server) printf(format string, a ...interface{}) {
	srv.mu.Lock()
	defer srv.mu.Unlock()
	if srv.events != nil {
		srv.events.Printf(format, a...)
	}
	glog.V(2).Infof(format, a...)
}

func (srv *Server) errorf(format string, a ...interface{}) {
	srv.mu.Lock()
	defer srv.mu.Unlock()
	if srv.events != nil {
		srv.events.Errorf(format, a...)
	}
	glog.Errorf(format, a...)
}
