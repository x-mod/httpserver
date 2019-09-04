package main

import (
	"context"
	"log"
	"net/http"
	"time"

	"github.com/x-mod/httpserver"
	"github.com/x-mod/routine"
)

func main() {
	srv := httpserver.NewServer(
		httpserver.ListenAddress(":8080"),
	)
	//优先匹配放在前面定义
	srv.Route(
		httpserver.Pattern("/hello"),
		httpserver.Handler(Query),
		httpserver.Query("foo", "bar"),
	)
	srv.Route(
		httpserver.Pattern("/hello"),
		httpserver.Handler(Hello),
	)
	srv.Route(
		httpserver.Prefix("/foo"),
		httpserver.Handler(Helo),
	)

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

func Hello(ctx context.Context, wr http.ResponseWriter, req *http.Request) {
	log.Println("hello handler ... ok", ctx.Value("x"))
	wr.WriteHeader(http.StatusOK)
	_, _ = wr.Write([]byte("I'm OK"))
}

func Helo(ctx context.Context, wr http.ResponseWriter, req *http.Request) {
	log.Println("Helo handler ... ok", ctx.Value("x"))
	wr.WriteHeader(http.StatusOK)
	_, _ = wr.Write([]byte("I'm OK"))
}

func Query(ctx context.Context, wr http.ResponseWriter, req *http.Request) {
	log.Println("Query handler ... ok", ctx.Value("x"))
	time.Sleep(5 * time.Second)
	wr.WriteHeader(http.StatusOK)
	_, _ = wr.Write([]byte("Query OK"))
}
