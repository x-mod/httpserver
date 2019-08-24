package main

import (
	"context"
	"log"

	"github.com/x-mod/errors"

	"github.com/x-mod/httpclient/grpc"
	demo "github.com/x-mod/httpserver/example/grpc/proto"
)

func main() {
	c := demo.NewHTTPDemoClient(grpc.Host("127.0.0.1:8080"))
	rsp, err := c.Hello(context.TODO(), &demo.HelloReq{Name: "JayL"})
	if err != nil {
		log.Println("err: ", err)
		log.Println("err code: ", errors.ValueFrom(err))
		return
	}
	log.Println("rsp: ", rsp.Greet)
}
