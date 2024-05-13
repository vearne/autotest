package rule

import (
	"github.com/go-resty/resty/v2"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestHttpBodyAtLeastOneRule(t *testing.T) {
	cases := []struct {
		xpath    string
		expected any
	}{
		// Notice: the first element is 1, not zero
		{"//title", "Effective Go"},
		{"//author", "Alan A. A. Donovan and Brian W. Kernighan"},
	}
	var resp resty.Response
	resp.SetBody([]byte(jsonStr1))
	for _, item := range cases {
		rule := HttpBodyAtLeastOneRule{item.xpath, item.expected}
		assert.True(t, rule.Verify(&resp))
	}
}

func TestHttpBodyEqualRule(t *testing.T) {
	cases := []struct {
		xpath    string
		expected any
	}{
		// Notice: the first element is 1, not zero
		{"(//title)[2]", "Effective Go"},
		{"(//author)[1]", "Alan A. A. Donovan and Brian W. Kernighan"},
	}
	var resp resty.Response
	resp.SetBody([]byte(jsonStr1))
	for _, item := range cases {
		rule := HttpBodyEqualRule{item.xpath, item.expected}
		assert.True(t, rule.Verify(&resp))
	}

}

func TestHttpBodyEqualRule2(t *testing.T) {
	cases := []struct {
		xpath    string
		expected any
	}{
		// Notice: the first element is 1, not zero
		{"/person/name", "John"},
		{"/person/female", false},
		{"/person/hobbies", []string{"coding", "eating", "football"}},
	}
	var resp resty.Response
	resp.SetBody([]byte(jsonStr2))
	for _, item := range cases {
		rule := HttpBodyEqualRule{item.xpath, item.expected}
		assert.True(t, rule.Verify(&resp))
	}
}
