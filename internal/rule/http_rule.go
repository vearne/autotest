package rule

import "github.com/go-resty/resty/v2"

// 实现 VerifyRule
type HttpStatusEqualRule struct {
	ExpectedStatus int `json:"expectedStatus"`
}

func (r *HttpStatusEqualRule) Name() string {
	return "HttpStatusEqualRule"
}

func (r *HttpStatusEqualRule) Valid(response *resty.Response) bool {
	return response.StatusCode() == r.ExpectedStatus
}

// 实现 VerifyRule
type HttpBodyEqualRule struct {
	Xpath string `json:"xpath"`
	Value any    `json:"value"`
}

func (r *HttpBodyEqualRule) Name() string {
	return "HttpBodyEqualRule"
}

func (r *HttpBodyEqualRule) Valid(response *resty.Response) bool {
	return true
}
