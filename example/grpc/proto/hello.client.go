package demo

import (
	"context"

	"github.com/x-mod/httpclient"
	client "github.com/x-mod/httpclient/grpc"
	"google.golang.org/grpc"
)

type HTTPDemoClient struct {
	*client.HTTPClient
}

func NewHTTPDemoClient(opts ...client.HTTPClientOpt) DemoClient {
	gopts := []client.HTTPClientOpt{}
	gopts = append(gopts, client.Version("v1"))
	gopts = append(gopts, client.PackageName("demo"))
	gopts = append(gopts, client.ServiceName("Demo"))
	gopts = append(gopts, opts...)
	return &HTTPDemoClient{
		HTTPClient: client.NewHTTPClient(gopts...),
	}
}

func (c *HTTPDemoClient) Hello(ctx context.Context, in *HelloReq, opts ...grpc.CallOption) (*HelloResp, error) {
	req, err := c.MakeRequest(
		"Hello",
		httpclient.Method("post"),
		httpclient.URL(httpclient.URI("/v1/hello")),
		httpclient.Content(httpclient.PBJSON(in)),
	)
	if err != nil {
		return nil, err
	}
	out := new(HelloResp)
	if err := c.Execute(ctx, req, client.PBResponse(out)); err != nil {
		return nil, err
	}
	return out, nil
}
