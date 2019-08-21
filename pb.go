package httpserver

import (
	"bytes"
	"context"
	"io/ioutil"
	"net/http"

	"github.com/golang/protobuf/jsonpb"
	"github.com/golang/protobuf/proto"
)

func defaultPBContext(req *http.Request, ctx context.Context) context.Context {
	return ctx
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
		wr.WriteHeader(http.StatusExpectationFailed)
		_, _ = wr.Write([]byte(err.Error()))
	} else {
		wr.WriteHeader(http.StatusOK)
		marshaler := &jsonpb.Marshaler{EmitDefaults: true}
		_ = marshaler.Marshal(wr, out)
	}
}

type PBContextFunc func(req *http.Request, ctx context.Context) context.Context
type PBRequestFunc func(req *http.Request, in proto.Message) error
type PBResponseFunc func(wr http.ResponseWriter, out proto.Message, err error)

var PBContext PBContextFunc
var PBRequest PBRequestFunc
var PBResponse PBResponseFunc

func init() {
	PBContext = defaultPBContext
	PBRequest = defaultPBRequest
	PBResponse = defaultPBResponse
}
