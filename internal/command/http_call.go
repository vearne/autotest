package command

import (
	"context"
	"fmt"
	"github.com/go-resty/resty/v2"
	"github.com/vearne/autotest/internal/config"
	"github.com/vearne/autotest/internal/model"
	"github.com/vearne/autotest/internal/resource"
	"github.com/vearne/executor"
	"github.com/vearne/zaplog"
	"go.uber.org/zap"
	"net/url"
	"strings"
	"time"
)

type HttpTestCaseResult struct {
	State  model.State
	ID     uint64
	Desc   string
	Reason model.Reason
	// actual request
	Request   config.RequestHttp
	TestCase  *config.TestCaseHttp
	KeyValues map[string]any
	Error     error
	Response  *resty.Response
}

/*
	POST /api/order/subscribe HTTP/1.1
	HOST: localhost:8080
	HEADERS:
	Content-Type: application/json
	User-Agent: go-resty/2.12.0 (https://github.com/go-resty/resty)
	BODY:
	{
	  "title": "book3_title",
	  "author": "book3_author"
	}
*/

func (t *HttpTestCaseResult) ReqDetail() string {
	var builder strings.Builder
	builder.WriteString(fmt.Sprintf("%v %v\n", t.Request.Method, t.Request.URL))
	u, _ := url.Parse(t.Request.URL)
	builder.WriteString(fmt.Sprintf("HOST: %v\n", u.Host))
	builder.WriteString("HEADERS:\n")
	for _, item := range t.Request.Headers {
		builder.WriteString(fmt.Sprintf("%v\n", item))
	}
	builder.WriteString("BODY:\n")
	builder.WriteString(fmt.Sprintf("%v\n", t.Request.Body))
	return builder.String()
}

/*
	STATUS       : 200 OK
	PROTO        : HTTP/1.1
	RECEIVED AT  : 2024-10-17T17:05:05.156315+08:00
	TIME DURATION: 38.354958ms
	HEADERS      :
		Content-Length: 17
		Content-Type: application/json; charset=utf-8
		Date: Thu, 17 Oct 2024 09:05:05 GMT
	BODY         :
	{
		"code": "E000"
	}
*/

func (t *HttpTestCaseResult) RespDetail() string {
	if t.Response == nil {
		return ""
	}

	var builder strings.Builder
	builder.WriteString(fmt.Sprintf("STATUS: %v %v\n", t.Response.StatusCode, t.Response.Status))
	builder.WriteString("HEADERS:\n")
	for key, values := range t.Response.Header() {
		builder.WriteString(fmt.Sprintf("%v: %v\n", key, strings.Join(values, ",")))
	}
	builder.WriteString("BODY:\n")
	builder.WriteString(fmt.Sprintf("%v\n", t.Response.String()))
	return builder.String()
}

type HttpTestCallable struct {
	testcase   *config.TestCaseHttp
	stateGroup *model.StateGroup
}

func (m *HttpTestCallable) Call(ctx context.Context) *executor.GPResult {
	zaplog.Debug("Call()", zap.Any("VerifyRules", m.testcase.VerifyRules))

	r := executor.GPResult{}
	tcResult := HttpTestCaseResult{
		ID:        m.testcase.ID,
		Desc:      m.testcase.Desc,
		State:     model.StateSuccessFul,
		Reason:    model.ReasonSuccess,
		TestCase:  m.testcase,
		KeyValues: map[string]any{},
		Error:     nil,
		Response:  nil,
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
	zaplog.Debug("before render()", zap.Uint64("testCaseId", m.testcase.ID),
		zap.Any("request", m.testcase.Request))
	req, err := renderRequestHttp(m.testcase.Request)
	tcResult.Request = req
	zaplog.Debug("after render()", zap.Uint64("testCaseId", m.testcase.ID),
		zap.Any("request", tcResult.Request))
	if err != nil {
		tcResult.State = model.StateFailed
		tcResult.Reason = model.ReasonTemplateRenderError
		tcResult.Error = err
		r.Value = tcResult
		r.Err = err
		return &r
	}

	// 4. timeout
	timeout := resource.GlobalConfig.Global.RequestTimeout
	rCtx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	// 5. trigger remote request
	in := resource.RestyClient.R().SetContext(rCtx)
	for _, item := range req.Headers {
		strList := strings.Split(item, ":")
		in.SetHeader(strings.TrimSpace(strList[0]), strings.TrimSpace(strList[1]))
	}

	in.SetHeader("Accept", "*/*")

	if len(req.Body) > 0 {
		in.SetBody(req.Body)
	}

	var out *resty.Response
	method := strings.ToUpper(req.Method)

	switch method {
	case "POST":
		out, err = in.Post(req.URL)
	case "PUT":
		out, err = in.Put(req.URL)
	case "DELETE":
		out, err = in.Delete(req.URL)
	default: // get and others
		out, err = in.Get(req.URL)
	}

	if err != nil {
		zaplog.Error("HttpTestCallable rules verify failed",
			zap.Uint64("testCaseId", m.testcase.ID),
			zap.Error(err),
		)
		tcResult.State = model.StateFailed
		tcResult.Reason = model.ReasonRequestFailed
		tcResult.Error = err
		r.Value = tcResult
		r.Err = err
		return &r
	}

	tcResult.Response = out

	// 6. export
	if m.testcase.Export != nil {
		exportConfig := m.testcase.Export
		// TODO handle error
		value, _ := exportTo(out.String(), exportConfig)
		tcResult.KeyValues[exportConfig.ExportTo] = value
	}

	// 7. verify
	for idx, rule := range m.testcase.VerifyRules {
		VerifyResult := rule.Verify(out)
		if !VerifyResult {
			zaplog.Error("HttpTestCallable rules validate failed",
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
}

func renderRequestHttp(req config.RequestHttp) (config.RequestHttp, error) {
	var err error
	// url
	req.URL, err = templateRender(req.URL)
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
	req.Body, err = templateRender(req.Body)
	if err != nil {
		return req, err
	}

	return req, nil
}
