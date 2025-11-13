package command

import (
	"context"
	"fmt"

	"github.com/vearne/autotest/internal/luavm"
	lua "github.com/yuin/gopher-lua"

	"github.com/fullstorydev/grpcurl"

	// ignore SA1019 we have to import this because it appears in exported API
	"os"
	"strings"
	"time"

	"github.com/golang/protobuf/proto"   //nolint:staticcheck
	"github.com/jhump/protoreflect/desc" //nolint:staticcheck
	"github.com/jhump/protoreflect/grpcreflect"
	"github.com/vearne/autotest/internal/config"
	"github.com/vearne/autotest/internal/model"
	"github.com/vearne/autotest/internal/resource"
	"github.com/vearne/autotest/internal/util"
	"github.com/vearne/executor"
	slog "github.com/vearne/simplelog"
	"github.com/vearne/zaplog"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

type GrpcTestCaseResult struct {
	State  model.State
	ID     uint64
	Desc   string
	Reason model.Reason
	// actual request
	Request   config.RequestGrpc
	TestCase  *config.TestCaseGrpc
	KeyValues map[string]any
	Error     error
	Response  *model.GrpcResp
}

func (t *GrpcTestCaseResult) ReqDetail() string {
	var builder strings.Builder
	builder.WriteString(fmt.Sprintf("ADDRESS: %v\n", t.Request.Address))
	builder.WriteString(fmt.Sprintf("SYMBOL: %v\n", t.Request.Symbol))
	builder.WriteString("HEADERS:\n")
	for _, item := range t.Request.Headers {
		builder.WriteString(fmt.Sprintf("%v\n", item))
	}
	builder.WriteString("BODY:\n")
	builder.WriteString(fmt.Sprintf("%v\n", t.Request.Body))
	return builder.String()
}

func (t *GrpcTestCaseResult) RespDetail() string {
	if t.Response == nil {
		return ""
	}

	var builder strings.Builder
	builder.WriteString(fmt.Sprintf("GRPC.CODE: %v\n", t.Response.Code))
	builder.WriteString(fmt.Sprintf("GRPC.MESSAGE %v\n", t.Response.Message))
	builder.WriteString("HEADERS:\n")
	for _, item := range t.Response.Headers {
		builder.WriteString(fmt.Sprintf("%v\n", item))
	}
	builder.WriteString("BODY:\n")
	builder.WriteString(fmt.Sprintf("%v\n", t.Response.Body))
	return builder.String()
}

type GrpcTestCallable struct {
	testcase   *config.TestCaseGrpc
	stateGroup *model.StateGroup
}

func (m *GrpcTestCallable) Call(ctx context.Context) *executor.GPResult {
	r := executor.GPResult{}
	var cc *grpc.ClientConn
	var rf grpcurl.RequestParser
	var formatter grpcurl.Formatter
	var in *strings.Reader
	var req config.RequestGrpc
	var timeout time.Duration
	var rCtx context.Context
	var cancel context.CancelFunc
	var handler *EventHandler
	var err error

	tcResult := GrpcTestCaseResult{
		ID:        m.testcase.ID,
		Desc:      m.testcase.Desc,
		State:     model.StateSuccessFul,
		Reason:    model.ReasonSuccess,
		TestCase:  m.testcase,
		KeyValues: map[string]any{},
	}

	// 1. Check other test cases of dependencies
	for _, id := range m.testcase.DependOnIDs {
		if m.stateGroup.GetState(id) == model.StateNotExecuted { // a certain dependency is not yet completed
			tcResult.State = model.StateNotExecuted
			tcResult.Reason = model.ReasonDependentItemNotCompleted
			r.Value = tcResult
			r.Err = nil
			return &r
		} else if m.stateGroup.GetState(id) == model.StateFailed {
			tcResult.State = model.StateFailed
			tcResult.Reason = model.ReasonDependentItemFailed
			r.Value = tcResult
			r.Err = nil
			return &r
		}
	}

	// 2. deal delay
	if m.testcase.Delay > 0 {
		zaplog.Debug("sleep", zap.Any("delay", m.testcase.Delay))
		time.Sleep(m.testcase.Delay)
	}

	// 3. render
	zaplog.Info("before render()", zap.Uint64("testCaseId", m.testcase.ID),
		zap.Any("request", m.testcase.Request))
	req, err = renderRequestGrpc(m.testcase.Request)
	tcResult.Request = req
	zaplog.Info("after render()", zap.Uint64("testCaseId", m.testcase.ID),
		zap.Any("request", tcResult.Request))
	if err != nil {
		tcResult.State = model.StateFailed
		tcResult.Reason = model.ReasonTemplateRenderError
		tcResult.Error = err
		r.Value = tcResult
		r.Err = err
		return &r
	}

	// 4. get description source
	options := grpcurl.FormatOptions{
		EmitJSONDefaultFields: true,
		IncludeTextSeparator:  true,
		AllowUnknownFields:    true,
	}

	reqInfo := tcResult.Request
	descSource, err := getDescSourceWitchCache(ctx, reqInfo.Address)
	if err != nil {
		zaplog.Error("GrpcTestCallable-get desc source",
			zap.Uint64("testCaseId", m.testcase.ID),
			zap.String("address", reqInfo.Address),
			zap.Error(err),
		)
		goto ERROR
	}

	cc, err = dial(reqInfo.Address)
	if err != nil {
		zaplog.Error("GrpcTestCallable-dial",
			zap.Uint64("testCaseId", m.testcase.ID),
			zap.String("address", reqInfo.Address),
			zap.Error(err),
		)
		goto ERROR
	}

	// 5. timeout
	timeout = resource.GlobalConfig.Global.RequestTimeout
	rCtx, cancel = context.WithTimeout(ctx, timeout)
	defer cancel()

	// 6. trigger remote request with rate limiting
	in = strings.NewReader(reqInfo.Body)
	rf, formatter, err = grpcurl.RequestParserAndFormatter(grpcurl.FormatJSON, descSource, in, options)
	if err != nil {
		zaplog.Error("GrpcTestCallable-RequestParserAndFormatter",
			zap.Uint64("testCaseId", m.testcase.ID),
			zap.String("address", reqInfo.Address),
			zap.Error(err),
		)
		goto ERROR
	}
	handler = NewEventHandler(formatter)

	// 使用限流器和重试机制控制gRPC请求的并发、速率和稳定性
	err = resource.RateLimiter.ExecuteWithLimit(rCtx, func() error {
		return util.ExecuteGrpcWithRetry(rCtx, resource.GlobalConfig, func() error {
			return grpcurl.InvokeRPC(rCtx, descSource, cc, reqInfo.Symbol, reqInfo.Headers, handler, rf.Next)
		})
	})
	if err != nil {
		zaplog.Error("GrpcTestCallable-invokeRPC",
			zap.Uint64("testCaseId", m.testcase.ID),
			zap.String("address", reqInfo.Address),
			zap.Error(err),
		)
		goto ERROR
	}

	tcResult.Response = &handler.resp

	if resource.GlobalConfig.Global.Debug {
		debugPrint(reqInfo, handler.resp)
	}

	// 8. export
	if m.testcase.Export != nil {
		exportConfig := m.testcase.Export
		// TODO handle error
		value, _ := exportTo(handler.resp.Body, exportConfig)
		tcResult.KeyValues[exportConfig.ExportTo] = value
	}

	// 9. verify
	for idx, rule := range m.testcase.VerifyRules {
		VerifyResult := rule.Verify(&handler.resp)
		if !VerifyResult {
			zaplog.Error("GrpcTestCallable rules validate failed",
				zap.Uint64("testCaseId", m.testcase.ID),
				zap.Int("ruleIdx", idx+1),
				zap.Any("rule", m.testcase.VerifyRules[idx]))

			tcResult.State = model.StateFailed
			tcResult.Reason = model.ReasonRuleVerifyFailed
			break
		}
	}
	r.Value = tcResult
	r.Err = nil
	return &r

ERROR:
	tcResult.State = model.StateFailed
	tcResult.Reason = model.ReasonRequestFailed
	tcResult.Error = err
	r.Value = tcResult
	r.Err = err
	return &r
}

func renderRequestGrpc(req config.RequestGrpc) (config.RequestGrpc, error) {
	var err error
	// address
	req.Address, err = templateRender(req.Address)
	if err != nil {
		return req, err
	}

	// headers
	for i := 0; i < len(req.Headers); i++ {
		req.Headers[i], err = templateRender(req.Headers[i])
		if err != nil {
			return req, err
		}
	}

	// body
	if len(req.LuaBody) > 0 {
		source := req.LuaBody +
			`
		
		return body();
	`
		var value lua.LValue
		zaplog.Info("renderRequestGrpc", zap.String("source", source))
		value, err = luavm.ExecuteLuaWithGlobalsPool(nil, nil, source)
		if err != nil {
			zaplog.Error("renderRequestGrpc-luaBody",
				zap.String("LuaStr", req.LuaBody),
				zap.Error(err))
			return req, err
		}
		req.Body = value.String()
	} else {
		req.Body, err = templateRender(req.Body)
		if err != nil {
			return req, err
		}
	}

	return req, nil
}

func getDescSourceWitchCache(ctx context.Context, address string) (grpcurl.DescriptorSource, error) {
	// 尝试从缓存获取
	if cached, found := resource.CacheManager.GrpcDescriptorCache.Get(address); found {
		if desc, ok := cached.(grpcurl.DescriptorSource); ok {
			slog.Debug("Using cached gRPC descriptor for %s", address)
			return desc, nil
		}
	}

	// 缓存未命中，获取新的描述符
	v, err, _ := resource.SingleFlightGroup.Do(address, func() (interface{}, error) {
		md := grpcurl.MetadataFromHeaders([]string{})
		refCtx := metadata.NewOutgoingContext(ctx, md)
		cc, dialErr := dial(address)
		if dialErr != nil {
			return nil, dialErr
		}
		refClient := grpcreflect.NewClientAuto(refCtx, cc)
		refClient.AllowMissingFileDescriptors()
		descSource := grpcurl.DescriptorSourceFromServer(ctx, refClient)
		return descSource, nil
	})
	if err != nil {
		zaplog.Error("getDescSourceWitchCache", zap.Error(err))
		return nil, err
	}

	s := v.(grpcurl.DescriptorSource)
	// 存储到缓存
	resource.CacheManager.GrpcDescriptorCache.Set(address, s)
	return s, err
}

func dial(target string) (*grpc.ClientConn, error) {
	var opts []grpc.DialOption
	network := "tcp"
	cc, err := grpcurl.BlockingDial(context.Background(), network, target, nil, opts...)
	if err != nil {
		return nil, err
	}
	return cc, nil
}

func debugPrint(reqInfo config.RequestGrpc, resp model.GrpcResp) {
	var b strings.Builder
	fmt.Fprintln(&b, "==============================================================================")
	fmt.Fprintln(&b, "~~~ REQUEST ~~~")
	fmt.Fprintln(&b, "ADDRESS\t:", reqInfo.Address)
	fmt.Fprintln(&b, "SYMBOL\t:", reqInfo.Symbol)
	fmt.Fprintln(&b, "HEADERS\t:")
	for _, item := range reqInfo.Headers {
		fmt.Fprintln(&b, "\t", item)
	}
	fmt.Fprintln(&b, "BODY\t:")
	fmt.Fprintln(&b, reqInfo.Body)
	fmt.Fprintln(&b, "------------------------------------------------------------------------------")
	fmt.Fprintln(&b, "~~~ RESPONSE ~~~")
	fmt.Fprintln(&b, "GRPC.CODE:", resp.Code)
	fmt.Fprintln(&b, "GRPC.MESSAGE:", resp.Message)
	fmt.Fprintln(&b, "HEADERS\t:")
	for _, item := range resp.Headers {
		fmt.Fprintln(&b, "\t", item)
	}
	fmt.Fprintln(&b, "BODY\t:")
	fmt.Fprintln(&b, resp.Body)

	os.Stderr.WriteString(b.String())
}

type EventHandler struct {
	resp      model.GrpcResp
	formatter grpcurl.Formatter
}

func NewEventHandler(formatter grpcurl.Formatter) *EventHandler {
	var handler EventHandler
	handler.formatter = formatter
	return &handler
}

func (m *EventHandler) OnResolveMethod(descriptor *desc.MethodDescriptor) {

}

func (m *EventHandler) OnSendHeaders(md metadata.MD) {

}

func (m *EventHandler) OnReceiveHeaders(md metadata.MD) {
	for key, values := range md {
		for _, value := range values {
			m.resp.Headers = append(m.resp.Headers, fmt.Sprintf("%v:%v", key, value))
		}
	}
}

func (m *EventHandler) OnReceiveResponse(msg proto.Message) {
	// Convert the message to the format expected by the formatter
	// The formatter expects github.com/golang/protobuf/proto.Message
	m.resp.Body, _ = m.formatter(msg)
}

func (m *EventHandler) OnReceiveTrailers(status *status.Status, md metadata.MD) {
	for key, values := range md {
		for _, value := range values {
			m.resp.Headers = append(m.resp.Headers, fmt.Sprintf("%v:%v", key, value))
		}
	}
	m.resp.Code = status.Code().String()
	m.resp.Message = status.Message()
}
