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

	"github.com/x-mod/httpserver"
)

func main() {
	srv := httpserver.NewServer(
		httpserver.Address(":8080"),
	)
	srv.Route(
		httpserver.Pattern("/hello"),
		httpserver.Handler(http.HandlerFunc(Hello)),
	)
	log.Println("httpserver:", srv.Serve(context.TODO()))
}

func Hello(wr http.ResponseWriter, req *http.Request) {
	wr.WriteHeader(http.StatusOK)
	_, _ = wr.Write([]byte("I'm OK"))
}
````
