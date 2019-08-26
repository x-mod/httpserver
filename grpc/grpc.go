package grpc

import (
	"context"
	"crypto/tls"
	"fmt"
	"net/http"
	"reflect"

	"github.com/x-mod/httpserver"
	"github.com/x-mod/options"
	"github.com/x-mod/tlsconfig"
)

//MethodHandler handler format
type MethodHandler func(interface{}, context.Context, http.ResponseWriter, *http.Request)

type MethodDescription struct {
	MethodName string
	Handler    MethodHandler
	Option     *options.HttpOption
}

type ServiceDescription struct {
	PackageName  string
	ServiceName  string
	Implemention interface{}
	Methods      []MethodDescription
	Option       *options.ServiceOption
}

type HTTPServerCfg struct {
	host string
	tls  *tls.Config
}
type HTTPServer struct {
	*httpserver.Server
	cfg *HTTPServerCfg
}
type HTTPServerOpt func(*HTTPServerCfg)

func Host(host string) HTTPServerOpt {
	return func(cfg *HTTPServerCfg) {
		cfg.host = host
	}
}

func TLSConfig(opts ...tlsconfig.Option) HTTPServerOpt {
	return func(cfg *HTTPServerCfg) {
		if len(opts) > 0 {
			cfg.tls = tlsconfig.New(opts...)
		}
	}
}

func NewHTTPServer(opts ...HTTPServerOpt) *HTTPServer {
	cfg := &HTTPServerCfg{
		host: "127.0.0.1",
	}
	for _, opt := range opts {
		opt(cfg)
	}
	srvopts := []httpserver.ServerOpt{
		httpserver.ListenAddress(cfg.host),
	}
	if cfg.tls != nil {
		srvopts = append(srvopts, httpserver.TLSConfig(cfg.tls))
	}
	return &HTTPServer{
		Server: httpserver.NewServer(srvopts...),
		cfg:    cfg,
	}
}

//RegistService
func (srv *HTTPServer) RegisterService(sd *ServiceDescription, impl interface{}) error {
	destT := reflect.TypeOf(sd.Implemention).Elem()
	implT := reflect.TypeOf(impl)
	if !implT.Implements(destT) {
		return fmt.Errorf("implemention type %v does not satisfy %v", implT, destT)
	}
	sd.Implemention = impl
	return srv.register(sd)
}

func (srv *HTTPServer) register(sd *ServiceDescription) error {
	version := "v1"
	if sd.Option != nil {
		if sd.Option.Version != "" {
			version = sd.Option.Version
		}
	}
	for _, m := range sd.Methods {
		opts := []httpserver.RouteOpt{
			httpserver.Handler(func(ctx context.Context, wr http.ResponseWriter, req *http.Request) {
				m.Handler(sd.Implemention, ctx, wr, req)
			}),
			httpserver.Pattern(URIFormat(version, sd.PackageName, sd.ServiceName, m.MethodName)),
		}
		if m.Option != nil {
			if m.Option.Method != "" {
				opts = append(opts, httpserver.Method(m.Option.Method))
			}
			if m.Option.Uri != "" {
				opts = append(opts, httpserver.Pattern(m.Option.Uri))
			}
		}
		srv.Route(opts...)
	}
	return nil
}

//default URIFormat: /v1/pkg.Service/Method
func defaultURIFormat(version string, pkg string, service string, method string) string {
	return fmt.Sprintf("/%s/%s.%s/%s", version, pkg, service, method)
}

type URIFormatFunc func(version string, pkg string, service string, method string) string

var URIFormat URIFormatFunc

func init() {
	URIFormat = defaultURIFormat
}
