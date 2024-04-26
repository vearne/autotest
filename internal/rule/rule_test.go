package rule

import (
	"fmt"
	//"fmt"
	"github.com/antchfx/jsonquery"
	"github.com/stretchr/testify/assert"

	"strings"
	"testing"
)

const jsonStr1 = `[
  {
    "id": 1,
    "title": "The Go Programming Language",
    "author": "Alan A. A. Donovan and Brian W. Kernighan"
  },
  {
    "id": 2,
    "title": "Effective Go",
    "author": "The Go Authors"
  }
]`
const jsonStr2 = `
{
	"person": {
		"name": "John",
		"age": 31,
		"female": false,
		"city": null,
		"hobbies": [
			"coding",
			"eating",
			"football"
		]
	}
}`

func TestXpathFind(t *testing.T) {
	cases := []struct {
		xpath    string
		expected any
	}{
		{"//title", "Effective Go"},
		{"//id", 2},
	}
	for _, item := range cases {
		doc, _ := jsonquery.Parse(strings.NewReader(jsonStr1))
		nodes := jsonquery.Find(doc, item.xpath)
		exist := false
		for _, node := range nodes {
			if convStr(item.expected) == fmt.Sprintf("%v", node.Value()) {
				exist = true
			}
		}
		assert.True(t, exist)
	}
}

func TestXpathFindOne(t *testing.T) {
	cases := []struct {
		xpath    string
		expected any
	}{
		// Notice: the first element is 1, not zero
		{"(//title)[2]", "Effective Go"},
		{"(//id)[1]", 1},
	}
	for _, item := range cases {
		doc, _ := jsonquery.Parse(strings.NewReader(jsonStr1))
		node := jsonquery.FindOne(doc, item.xpath)
		assert.Equal(t, convStr(item.expected), convStr(node.Value()))
	}
}

func TestXpathFindOne2(t *testing.T) {
	cases := []struct {
		xpath    string
		expected any
	}{
		// Notice: the first element is 1, not zero
		{"/person/name", "John"},
		{"/person/female", false},
		{"/person/hobbies", []string{"coding", "eating", "football"}},
	}
	for _, item := range cases {
		doc, _ := jsonquery.Parse(strings.NewReader(jsonStr2))
		node := jsonquery.FindOne(doc, item.xpath)
		assert.Equal(t, convStr(item.expected), convStr(node.Value()))
	}
}

func TestConvStr(t *testing.T) {
	cases := []struct {
		v        any
		expected string
	}{
		// Notice: the first element is 1, not zero
		{15, "15"},
		{"hello", "hello"},
		{[]int{1, 2, 3}, "[1 2 3]"},
		{[]string{"1", "2", "3"}, "[1 2 3]"},
		{[]string{"hello", "world"}, "[hello world]"},
	}
	for _, item := range cases {
		assert.Equal(t, item.expected, convStr(item.v))
	}
}
