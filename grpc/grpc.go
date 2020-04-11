package grpc

import (
	"bytes"
	"context"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/golang/protobuf/jsonpb"
	"github.com/golang/protobuf/proto"
	"github.com/x-mod/httpserver/render"
)

func defaultPBContext(req *http.Request, ctx context.Context) context.Context {
	return req.Context()
}

func defaultPBRequest(req *http.Request, in proto.Message) error {
	data, err := ioutil.ReadAll(req.Body)
	if err != nil {
		return err
	}
	defer req.Body.Close()
	return jsonpb.Unmarshal(bytes.NewBuffer(data), in)
}

func defaultPBResponse(wr http.ResponseWriter, out proto.Message, err error) {
	if err != nil {
		render.Error(err).Response(wr, render.StatusCode(http.StatusExpectationFailed))
	} else {
		render.PBMessage(out).Response(wr)
	}
}

//default URIFormat: /v1/pkg.Service/Method
func defaultURIFormat(version string, pkg string, service string, method string) string {
	return fmt.Sprintf("/%s/%s.%s/%s", version, pkg, service, method)
}

type URIFormatFunc func(version string, pkg string, service string, method string) string
type PBContextFunc func(req *http.Request, ctx context.Context) context.Context
type PBRequestFunc func(req *http.Request, in proto.Message) error
type PBResponseFunc func(wr http.ResponseWriter, out proto.Message, err error)

var PBContext PBContextFunc
var PBRequest PBRequestFunc
var PBResponse PBResponseFunc
var URIFormat URIFormatFunc

func init() {
	PBContext = defaultPBContext
	PBRequest = defaultPBRequest
	PBResponse = defaultPBResponse
	URIFormat = defaultURIFormat
}
