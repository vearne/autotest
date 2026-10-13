package util

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPrettyJSON(t *testing.T) {
	cases := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:  "compact object is indented",
			input: `{"resultCode":1000,"data":{"id":2,"account":"account00000006"}}`,
			expected: `{
  "resultCode": 1000,
  "data": {
    "id": 2,
    "account": "account00000006"
  }
}`,
		},
		{
			name:  "array is indented",
			input: `[{"id":1},{"id":2}]`,
			expected: `[
  {
    "id": 1
  },
  {
    "id": 2
  }
]`,
		},
		{
			name:     "invalid JSON is returned unchanged",
			input:    `not a json {body`,
			expected: `not a json {body`,
		},
		{
			name:     "empty string is returned unchanged",
			input:    ``,
			expected: ``,
		},
		{
			name:     "plain text is returned unchanged",
			input:    `hello world`,
			expected: `hello world`,
		},
		{
			name:     "scalar JSON keeps its value",
			input:    `123`,
			expected: `123`,
		},
		{
			name:     "key order is preserved",
			input:    `{"z":1,"a":2,"m":3}`,
			expected: "{\n  \"z\": 1,\n  \"a\": 2,\n  \"m\": 3\n}",
		},
	}

	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			got := PrettyJSON(tt.input)
			assert.Equal(t, tt.expected, got)
		})
	}
}
