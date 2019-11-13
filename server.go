package httpserver

import (
	"context"
	"crypto/tls"
	"log"
	"net"
	"net/http"
	"strings"
	"time"

	"github.com/gorilla/mux"
)

type ContextHandler func(context.Context, http.ResponseWriter, *http.Request)

type Server struct {
	addr   string
	rctx   context.Context
	http   *http.Server
	tls    *tls.Config
	routes *mux.Router
}

type ServerOpt func(*Server)

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

func NewServer(opts ...ServerOpt) *Server {
	hsrv := &Server{
		rctx: context.TODO(),
		http: &http.Server{
			ReadTimeout:  15 * time.Second,
			WriteTimeout: 15 * time.Second,
			IdleTimeout:  60 * time.Second,
		},
		routes: mux.NewRouter(),
	}
	for _, opt := range opts {
		opt(hsrv)
	}
	return hsrv
}

type RouteCfg struct {
	schemes []string
	host    string
	methods []string
	prefix  string
	pattern string
	headers []string
	queries []string
	handler ContextHandler
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
func Handler(h ContextHandler) RouteOpt {
	return func(cf *RouteCfg) {
		cf.handler = h
	}
}

func (srv *Server) wrapHandler(h ContextHandler) http.HandlerFunc {
	return func(wr http.ResponseWriter, req *http.Request) {
		h(srv.rctx, wr, req)
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
		r := srv.routes.NewRoute().Handler(srv.wrapHandler(cf.handler))
		if cf.pattern != "" {
			r.Path(cf.pattern)
		}
		if cf.prefix != "" {
			r.PathPrefix(cf.prefix)
		}
		if len(cf.schemes) > 0 {
			r.Schemes(cf.schemes...)
		}
		if len(cf.methods) > 0 {
			r.Methods(cf.methods...)
		}
		if len(cf.host) > 0 {
			r.Host(cf.host)
		}
		if len(cf.headers) > 0 && len(cf.headers)%2 == 0 {
			r.Headers(cf.headers...)
		}
		if len(cf.queries) > 0 && len(cf.queries)%2 == 0 {
			r.Queries(cf.queries...)
		}
	}
}

func (srv *Server) Serve(ctx context.Context) error {
	srv.rctx = ctx
	srv.http.Handler = srv.routes

	ln, err := net.Listen("tcp", srv.addr)
	if err != nil {
		return err
	}
	log.Println("serving at ", srv.addr)
	if srv.tls != nil {
		ln = tls.NewListener(ln, srv.tls)
	}
	return srv.http.Serve(ln)
}

func (srv *Server) Shutdown(ctx context.Context) error {
	return srv.http.Shutdown(ctx)
}
