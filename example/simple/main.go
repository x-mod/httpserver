package main

import (
	"context"
	"log"
	"net/http"
	"syscall"
	"time"

	_ "net/http/pprof"

	"golang.org/x/net/trace"

	"github.com/x-mod/glog"
	"github.com/x-mod/httpserver"
	"github.com/x-mod/routine"
)

func main() {
	glog.Open(
		glog.LogToStderr(true),
		glog.Verbosity(2),
	)
	defer glog.Close()

	srv := httpserver.NewServer(
		httpserver.Address(":8080"),
		httpserver.NetTrace(true),
	)
	//优先匹配放在前面定义
	srv.Route(
		httpserver.Pattern("/hello"),
		httpserver.Handler(http.HandlerFunc(Query)),
		httpserver.Query("foo", "bar"),
	)
	srv.Route(
		httpserver.Pattern("/hello"),
		httpserver.Handler(http.HandlerFunc(Hello)),
	)
	srv.Route(
		httpserver.Prefix("/foo"),
		httpserver.Handler(http.HandlerFunc(Helo)),
	)

	ctx := context.WithValue(context.TODO(), "x", "y")
	err := routine.Main(
		ctx,
		routine.ExecutorFunc(srv.Serve),
		// routine.Go(routine.Profiling(":6060")),
		routine.Signal(syscall.SIGINT, routine.SigHandler(func() {
			log.Println("SIGINT ...")
			if err := srv.Shutdown(ctx); err != nil {
				log.Println("httpserver shutdown:", err)
			}
		})),
		// routine.Cleanup(routine.ExecutorFunc(srv.Shutdown)),
	)
	if err != nil {
		log.Println("failed: ", err)
	}
}

func Hello(wr http.ResponseWriter, req *http.Request) {
	if tr, ok := trace.FromContext(req.Context()); ok {
		tr.LazyPrintf("hello : %s", req.Context().Value("x"))
	}
	log.Println("hello handler ... ok", req.Context().Value("x"))
	wr.WriteHeader(http.StatusOK)
	_, _ = wr.Write([]byte("I'm OK"))
}

func Helo(wr http.ResponseWriter, req *http.Request) {
	if tr, ok := trace.FromContext(req.Context()); ok {
		tr.LazyPrintf("helo :", req.Context().Value("x"))
	}
	log.Println("Helo handler ... ok", req.Context().Value("x"))
	wr.WriteHeader(http.StatusOK)
	_, _ = wr.Write([]byte("I'm OK"))
}

func Query(wr http.ResponseWriter, req *http.Request) {
	log.Println("Query handler ... ok", req.Context().Value("x"))
	time.Sleep(5 * time.Second)
	wr.WriteHeader(http.StatusOK)
	_, _ = wr.Write([]byte("Query OK"))
}
