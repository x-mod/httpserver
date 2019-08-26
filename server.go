package httpserver

import (
	"context"
	"crypto/tls"
	"net/http"
	"strings"
	"time"

	"github.com/gorilla/mux"
)

type ContextHandler func(context.Context, http.ResponseWriter, *http.Request)

type Server struct {
	http.Server
	tls   *tls.Config
	ctx   context.Context
	route *mux.Router
}

type ServerOpt func(*Server)

func ListenAddress(addr string) ServerOpt {
	return func(srv *Server) {
		srv.Server.Addr = addr
	}
}

func TLSConfig(cf *tls.Config) ServerOpt {
	return func(srv *Server) {
		srv.tls = cf
	}
}

func ReadTimeout(rd time.Duration) ServerOpt {
	return func(srv *Server) {
		srv.Server.ReadTimeout = rd
	}
}
func WriteTimeout(wr time.Duration) ServerOpt {
	return func(srv *Server) {
		srv.Server.WriteTimeout = wr
	}
}
func IdleTimeout(idle time.Duration) ServerOpt {
	return func(srv *Server) {
		srv.Server.IdleTimeout = idle
	}
}

func NewServer(opts ...ServerOpt) *Server {
	srv := &Server{
		ctx:   context.TODO(),
		route: mux.NewRouter(),
		Server: http.Server{
			ReadTimeout:  15 * time.Second,
			WriteTimeout: 15 * time.Second,
			IdleTimeout:  60 * time.Second,
		},
	}
	for _, opt := range opts {
		opt(srv)
	}
	return srv
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
		h(srv.ctx, wr, req)
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
	if cf.pattern != "" && cf.handler != nil {
		r := srv.route.Handle(cf.pattern, srv.wrapHandler(cf.handler))
		if len(cf.schemes) > 0 {
			r = r.Schemes(cf.schemes...)
		}
		if len(cf.methods) > 0 {
			r = r.Methods(cf.methods...)
		}
		if len(cf.host) > 0 {
			r = r.Host(cf.host)
		}
		if len(cf.prefix) > 0 {
			r = r.PathPrefix(cf.prefix)
		}
		if len(cf.headers) > 0 && len(cf.headers)%2 == 0 {
			r = r.Headers(cf.headers...)
		}
		if len(cf.queries) > 0 && len(cf.queries)%2 == 0 {
			r = r.Queries(cf.queries...)
		}
	}
}

func (srv *Server) Serve(ctx context.Context) error {
	srv.ctx = ctx
	srv.Handler = srv.route
	if srv.tls != nil {
		srv.Server.TLSConfig = srv.tls
		return srv.ListenAndServeTLS("", "")
	}
	return srv.Server.ListenAndServe()
}
