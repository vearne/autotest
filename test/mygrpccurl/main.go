package main

import (
	"context"
	"encoding/json"
	"log/slog"
	"strings"

	"github.com/fullstorydev/grpcurl"
	"github.com/golang/protobuf/proto"
	"github.com/jhump/protoreflect/desc"
	"github.com/jhump/protoreflect/grpcreflect"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

func dial(target string) *grpc.ClientConn {
	var opts []grpc.DialOption
	network := "tcp"
	cc, err := grpcurl.BlockingDial(context.Background(), network, target, nil, opts...)
	if err != nil {
		panic(err)
	}
	return cc
}

func main() {
	ctx := context.Background()
	md := grpcurl.MetadataFromHeaders([]string{})
	refCtx := metadata.NewOutgoingContext(ctx, md)
	cc := dial("localhost:50031")
	refClient := grpcreflect.NewClientAuto(refCtx, cc)
	refClient.AllowMissingFileDescriptors()
	descSource := grpcurl.DescriptorSourceFromServer(ctx, refClient)

	var err error
	handler := new(myEventHandler)
	symbol := "Bookstore/ListBook"

	options := grpcurl.FormatOptions{
		EmitJSONDefaultFields: true,
		IncludeTextSeparator:  true,
		AllowUnknownFields:    true,
	}

	in := strings.NewReader("{}")
	rf, _, err := grpcurl.RequestParserAndFormatter(grpcurl.FormatJSON, descSource, in, options)
	if err != nil {
		slog.Error("grpcurl.RequestParserAndFormatter", "error", err)
	}
	err = grpcurl.InvokeRPC(ctx, descSource, cc, symbol, []string{}, handler, rf.Next)
	if err != nil {
		slog.Error("grpcurl.InvokeRPC", "error", err)
	}

}

type myEventHandler struct {
}

func (m myEventHandler) OnResolveMethod(descriptor *desc.MethodDescriptor) {

}

func (m myEventHandler) OnSendHeaders(md metadata.MD) {

}

func (m myEventHandler) OnReceiveHeaders(md metadata.MD) {
	b, _ := json.Marshal(md)
	slog.Info("OnReceiveHeaders", "data", string(b))
}

func (m myEventHandler) OnReceiveResponse(msg proto.Message) {
	b, _ := json.Marshal(msg)
	slog.Info("OnReceiveResponse", "message", string(b))
}

func (m myEventHandler) OnReceiveTrailers(status *status.Status, md metadata.MD) {
	slog.Info("OnReceiveTrailers", "code", status.Code().String(), "message", status.Message())
}
