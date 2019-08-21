package main

import (
	"context"
	"log"
	"time"

	"github.com/x-mod/httpserver"
	demo "github.com/x-mod/httpserver/example/grpc/proto"
	"github.com/x-mod/routine"
)

type HelloImpl struct{}

func (*HelloImpl) Hello(ctx context.Context, req *demo.HelloReq) (*demo.HelloResp, error) {
	return &demo.HelloResp{
		Greet: "Hello " + req.Name,
	}, nil
}

func main() {
	srv := httpserver.NewServer(
		httpserver.ListenAddress(":8080"),
	)

	if err := demo.RegisterDemoHTTPServer(srv, &HelloImpl{}); err != nil {
		log.Println("register failed: ", err)
		return
	}

	err := routine.Main(
		routine.ExecutorFunc(srv.Serve),
		routine.Context(context.WithValue(context.TODO(), "x", "y")),
		routine.Interrupts(routine.DefaultCancelInterruptors...),
		routine.Cleanup(
			routine.ExecutorFunc(func(ctx context.Context) error {
				//graceful shutdown MaxTime 15s
				tmctx, cancel := context.WithTimeout(ctx, 3*time.Second)
				defer cancel()
				return srv.Shutdown(tmctx)
			})),
	)
	if err != nil {
		log.Println("failed: ", err)
	}
}
