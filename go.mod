module github.com/x-mod/httpserver

go 1.12

require (
	github.com/golang/protobuf v1.3.2
	github.com/gorilla/mux v1.7.3
	github.com/x-mod/errors v0.1.6
	github.com/x-mod/httpclient v0.2.1
	github.com/x-mod/options v0.1.0
	github.com/x-mod/routine v1.1.2
	github.com/x-mod/tlsconfig v0.0.1
	golang.org/x/net v0.0.0-20190509222800-a4d6f7feada5
	google.golang.org/genproto v0.0.0-20180817151627-c66870c02cf8
	google.golang.org/grpc v1.19.1
)

replace github.com/x-mod/httpclient v0.2.1 => ../httpclient
