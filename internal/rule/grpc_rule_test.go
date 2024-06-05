package rule

import (
	"github.com/stretchr/testify/assert"
	"github.com/vearne/autotest/internal/model"
	"testing"
)

func TestGrpcBodyAtLeastOneRule(t *testing.T) {
	cases := []struct {
		xpath    string
		expected any
	}{
		// Notice: the first element is 1, not zero
		{"//title", "Effective Go"},
		{"//author", "Alan A. A. Donovan and Brian W. Kernighan"},
	}
	var resp model.GrpcResp
	resp.Body = jsonStr1
	for _, item := range cases {
		rule := GrpcBodyAtLeastOneRule{item.xpath, item.expected}
		assert.True(t, rule.Verify(&resp))
	}
}

func TestGrpcBodyEqualRule(t *testing.T) {
	cases := []struct {
		xpath    string
		expected any
	}{
		// Notice: the first element is 1, not zero
		{"(//title)[2]", "Effective Go"},
		{"(//author)[1]", "Alan A. A. Donovan and Brian W. Kernighan"},
	}
	var resp model.GrpcResp
	resp.Body = jsonStr1
	for _, item := range cases {
		rule := GrpcBodyEqualRule{item.xpath, item.expected}
		assert.True(t, rule.Verify(&resp))
	}

}

func TestGrpcBodyEqualRule2(t *testing.T) {
	cases := []struct {
		xpath    string
		expected any
	}{
		// Notice: the first element is 1, not zero
		{"/person/name", "John"},
		{"/person/female", false},
		{"/person/hobbies", []string{"coding", "eating", "football"}},
	}
	var resp model.GrpcResp
	resp.Body = jsonStr2
	for _, item := range cases {
		rule := GrpcBodyEqualRule{item.xpath, item.expected}
		assert.True(t, rule.Verify(&resp))
	}
}
