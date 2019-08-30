httpserver
===

Another HTTP Server use handler with context & Response Render.

# Quick Start

````go
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
		httpserver.Address(":8080"),
	)
	srv.Route(
		httpserver.Pattern("/hello"),
		httpserver.Handler(Hello),
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
````
