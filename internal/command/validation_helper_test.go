package command

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// MockTestCase 用于测试的模拟测试用例
type MockTestCase struct {
	ID          uint64
	DependOnIDs []uint64
}

func (m MockTestCase) GetID() uint64 {
	return m.ID
}

func (m MockTestCase) GetDependOnIDs() []uint64 {
	return m.DependOnIDs
}

func TestValidateTestCaseIDs(t *testing.T) {
	tests := []struct {
		name      string
		testcases []MockTestCase
		wantError bool
	}{
		{
			name: "无重复ID",
			testcases: []MockTestCase{
				{ID: 1},
				{ID: 2},
				{ID: 3},
			},
			wantError: false,
		},
		{
			name: "有重复ID",
			testcases: []MockTestCase{
				{ID: 1},
				{ID: 2},
				{ID: 1}, // 重复
			},
			wantError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateTestCaseIDs("test.yml", tt.testcases)
			if tt.wantError {
				assert.Error(t, err)
				assert.Equal(t, ErrorIDduplicate, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestValidateDependencies(t *testing.T) {
	tests := []struct {
		name      string
		testcases []MockTestCase
		wantError error
	}{
		{
			name: "正常依赖",
			testcases: []MockTestCase{
				{ID: 1, DependOnIDs: []uint64{}},
				{ID: 2, DependOnIDs: []uint64{1}},
				{ID: 3, DependOnIDs: []uint64{1, 2}},
			},
			wantError: nil,
		},
		{
			name: "依赖不存在",
			testcases: []MockTestCase{
				{ID: 1, DependOnIDs: []uint64{}},
				{ID: 2, DependOnIDs: []uint64{99}}, // 依赖不存在的ID
			},
			wantError: ErrorDependencyNotExist,
		},
		{
			name: "依赖ID大于等于当前ID",
			testcases: []MockTestCase{
				{ID: 1, DependOnIDs: []uint64{2}}, // 依赖更大的ID
				{ID: 2, DependOnIDs: []uint64{}},
			},
			wantError: ErrorIDRestrict,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateDependencies("test.yml", tt.testcases)
			if tt.wantError != nil {
				assert.Error(t, err)
				assert.Equal(t, tt.wantError, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestValidateXPath(t *testing.T) {
	tests := []struct {
		name      string
		xpath     string
		wantError bool
	}{
		{
			name:      "有效的XPath",
			xpath:     "//title",
			wantError: false,
		},
		{
			name:      "复杂的XPath",
			xpath:     "/data/books[1]/title",
			wantError: false,
		},
		{
			name:      "无效的XPath",
			xpath:     "//[invalid",
			wantError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateXPath(1, tt.xpath)
			if tt.wantError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestValidateBodyFields(t *testing.T) {
	tests := []struct {
		name      string
		body      string
		luaBody   string
		wantError bool
	}{
		{
			name:      "只有body",
			body:      `{"key": "value"}`,
			luaBody:   "",
			wantError: false,
		},
		{
			name:      "只有luaBody",
			body:      "",
			luaBody:   "function body() return '{}' end",
			wantError: false,
		},
		{
			name:      "都为空",
			body:      "",
			luaBody:   "",
			wantError: false,
		},
		{
			name:      "同时有body和luaBody",
			body:      `{"key": "value"}`,
			luaBody:   "function body() return '{}' end",
			wantError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateBodyFields(1, tt.body, tt.luaBody)
			if tt.wantError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
