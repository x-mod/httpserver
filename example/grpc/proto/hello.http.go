package demo

import (
	"context"
	"net/http"

	"github.com/x-mod/httpserver"
)

var _Demo_HTTP_serviceDesc = httpserver.ServiceDescription{
	PackageName:  "demo",
	ServiceName:  "Demo",
	Implemention: (*DemoServer)(nil),
	Methods: []httpserver.MethodDescription{
		{
			MethodName: "Hello",
			Handler:    _Demo_Hello_HTTP_Handler,
		},
	},
}

func RegisterDemoHTTPServer(s *httpserver.Server, srv DemoServer) error {
	return s.RegisterService(&_Demo_HTTP_serviceDesc, srv)
}

func _Demo_Hello_HTTP_Handler(srv interface{}, ctx context.Context, wr http.ResponseWriter, req *http.Request) {
	in := new(HelloReq)
	err := httpserver.PBRequest(req, in)
	if err != nil {
		httpserver.PBResponse(wr, nil, err)
		return
	}
	out, err := srv.(DemoServer).Hello(httpserver.PBContext(req, ctx), in)
	if err != nil {
		httpserver.PBResponse(wr, nil, err)
		return
	}
	httpserver.PBResponse(wr, out, err)
}
