package grpc

import (
	"bytes"
	"context"
	"io/ioutil"
	"net/http"

	"github.com/x-mod/errors"

	"github.com/golang/protobuf/jsonpb"
	"github.com/golang/protobuf/proto"
	spb "google.golang.org/genproto/googleapis/rpc/status"
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
		st := &spb.Status{}
		st.Code = errors.ValueFrom(err)
		st.Message = err.Error()
		marshaler := &jsonpb.Marshaler{EmitDefaults: true}
		_ = marshaler.Marshal(wr, st)
	} else {
		wr.Header().Set("Content-Type", "application/json")
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
