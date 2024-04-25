package rule

import (
	"github.com/antchfx/jsonquery"
	"github.com/go-resty/resty/v2"
	"strings"
)

// 实现 VerifyRule
type HttpStatusEqualRule struct {
	ExpectedStatus int `json:"expectedStatus"`
}

func (r *HttpStatusEqualRule) Name() string {
	return "HttpStatusEqualRule"
}

func (r *HttpStatusEqualRule) Valid(resp *resty.Response) bool {
	return resp.StatusCode() == r.ExpectedStatus
}

// 实现 VerifyRule
type HttpBodyEqualRule struct {
	Xpath    string `json:"xpath"`
	Expected any    `json:"expected"`
}

func (r *HttpBodyEqualRule) Name() string {
	return "HttpBodyEqualRule"
}

func (r *HttpBodyEqualRule) Valid(resp *resty.Response) bool {
	doc, err := jsonquery.Parse(strings.NewReader(resp.String()))
	if err != nil {
		return false
	}
	node := jsonquery.FindOne(doc, r.Xpath)
	return node != nil && convStr(r.Expected) == convStr(node.Value())
}

// 实现 VerifyRule
// 至少找到一个满足条件的元素
type HttpBodyAtLeastOneRule struct {
	Xpath    string `json:"xpath"`
	Expected any    `json:"expected"`
}

func (r *HttpBodyAtLeastOneRule) Name() string {
	return "HttpBodyAtLeastOneRule"
}

func (r *HttpBodyAtLeastOneRule) Valid(resp *resty.Response) bool {
	doc, err := jsonquery.Parse(strings.NewReader(resp.String()))
	if err != nil {
		return false
	}
	nodes := jsonquery.Find(doc, r.Xpath)
	for _, node := range nodes {
		if node != nil && convStr(r.Expected) == convStr(node.Value()) {
			return true
		}
	}
	return false
}
