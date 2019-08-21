package httpserver

import (
	"context"
	"fmt"
	"net/http"
	"reflect"
)

//MethodHandler handler format
type MethodHandler func(interface{}, context.Context, http.ResponseWriter, *http.Request)

type MethodDescription struct {
	MethodName string
	Handler    MethodHandler
}

type ServiceDescription struct {
	PackageName  string
	ServiceName  string
	Implemention interface{}
	Methods      []MethodDescription
}

//RegistService
func (srv *Server) RegisterService(sd *ServiceDescription, impl interface{}) error {
	destT := reflect.TypeOf(sd.Implemention).Elem()
	implT := reflect.TypeOf(impl)
	if !implT.Implements(destT) {
		return fmt.Errorf("implemention type %v does not satisfy %v", implT, destT)
	}
	sd.Implemention = impl
	return srv.register(sd)
}

func (srv *Server) register(sd *ServiceDescription) error {
	for _, m := range sd.Methods {
		srv.Route(
			Pattern(fmt.Sprintf("/%s/%s/%s", sd.PackageName, sd.ServiceName, m.MethodName)),
			Handler(func(ctx context.Context, wr http.ResponseWriter, req *http.Request) {
				m.Handler(sd.Implemention, ctx, wr, req)
			}),
		)
	}
	return nil
}
