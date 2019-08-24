package demo

import (
	"context"
	"net/http"

	server "github.com/x-mod/httpserver/grpc"
	"github.com/x-mod/options"
)

var _Demo_HTTP_serviceDesc = server.ServiceDescription{
	PackageName: "demo",
	ServiceName: "Demo",
	Option: &options.ServiceOption{
		Version: "v1",
	},
	Implemention: (*DemoServer)(nil),
	Methods: []server.MethodDescription{
		{
			MethodName: "Hello",
			Option: &options.HttpOption{
				Method: "post",
				Uri:    "/v1/hello",
			},
			Handler: _Demo_Hello_HTTP_Handler,
		},
	},
}

func RegisterDemoHTTPServer(s *server.HTTPServer, srv DemoServer) error {
	return s.RegisterService(&_Demo_HTTP_serviceDesc, srv)
}

func _Demo_Hello_HTTP_Handler(srv interface{}, ctx context.Context, wr http.ResponseWriter, req *http.Request) {
	in := new(HelloReq)
	err := server.PBRequest(req, in)
	if err != nil {
		server.PBResponse(wr, nil, err)
		return
	}
	out, err := srv.(DemoServer).Hello(server.PBContext(req, ctx), in)
	if err != nil {
		server.PBResponse(wr, nil, err)
		return
	}
	server.PBResponse(wr, out, err)
}
