package rule

import (
	"github.com/antchfx/jsonquery"
	"strings"
)

type GrpcResp struct {
	Code    int
	JsonStr string
}

// implement VerifyRule
type GrpcCodeEqualRule struct {
	Expected int `json:"expected"`
}

func (r *GrpcCodeEqualRule) Name() string {
	return "GrpcCodeEqualRule"
}

func (r *GrpcCodeEqualRule) Verify(resp *GrpcResp) bool {
	return resp.Code == r.Expected
}

type GrpcBodyEqualRule struct {
	Xpath    string `json:"xpath"`
	Expected any    `json:"expected"`
}

func (r *GrpcBodyEqualRule) Name() string {
	return "GrpcBodyEqualRule"
}

func (r *GrpcBodyEqualRule) Verify(resp *GrpcResp) bool {
	doc, err := jsonquery.Parse(strings.NewReader(resp.JsonStr))
	if err != nil {
		return false
	}
	node := jsonquery.FindOne(doc, r.Xpath)
	return node != nil && convStr(r.Expected) == convStr(node.Value())
}

// implement VerifyRule
// Find at least one element that satisfies the condition
type GrpcBodyAtLeastOneRule struct {
	Xpath    string `json:"xpath"`
	Expected any    `json:"expected"`
}

func (r *GrpcBodyAtLeastOneRule) Name() string {
	return "GrpcBodyAtLeastOneRule"
}

func (r *GrpcBodyAtLeastOneRule) Verify(resp *GrpcResp) bool {
	doc, err := jsonquery.Parse(strings.NewReader(resp.JsonStr))
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
