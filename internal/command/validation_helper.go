package command

import (
	"fmt"

	"github.com/antchfx/xpath"
	"github.com/vearne/autotest/internal/model"
	slog "github.com/vearne/simplelog"
)

// ValidateTestCaseIDs 验证测试用例ID是否重复
func ValidateTestCaseIDs[T model.IdItem](filePath string, testcases []T) error {
	slog.Info("filePath:%v, len(testcases):%v", filePath, len(testcases))
	exist := make(map[uint64]struct{})
	for _, tc := range testcases {
		id := tc.GetID()
		if _, ok := exist[id]; ok {
			slog.Error("filePath:%v, ID [%v] is duplicate", filePath, id)
			return ErrorIDduplicate
		}
		exist[id] = struct{}{}
	}
	return nil
}

// ValidateDependencies 验证依赖关系
func ValidateDependencies[T interface {
	GetID() uint64
	GetDependOnIDs() []uint64
}](filePath string, testcases []T) error {
	slog.Info("filePath:%v, len(testcases):%v", filePath, len(testcases))
	exist := make(map[uint64]struct{})
	for _, tc := range testcases {
		exist[tc.GetID()] = struct{}{}
	}

	for _, tc := range testcases {
		for _, dID := range tc.GetDependOnIDs() {
			if _, ok := exist[dID]; !ok {
				slog.Error("For testcase %v, the testcase %v that it depend on does not exist",
					tc.GetID(), dID)
				return ErrorDependencyNotExist
			}
			if dID >= tc.GetID() {
				slog.Error("The ID of the dependent testcase must be smaller "+
					"than the ID of the current testcase, testcase:%v, dependent testcase:%v",
					tc.GetID(), dID)
				return ErrorIDRestrict
			}
		}
	}
	return nil
}

// ValidateXPath 验证XPath语法
func ValidateXPath(testCaseId uint64, xpathStr string) error {
	_, err := xpath.Compile(xpathStr)
	if err != nil {
		slog.Error("rule error, testCaseId:%v, xpath:%v", testCaseId, xpathStr)
		return err
	}
	return nil
}

// ValidateBodyFields 验证body和luaBody字段
func ValidateBodyFields(testCaseId uint64, body, luaBody string) error {
	if len(body) > 0 && len(luaBody) > 0 {
		slog.Error("request error, testCaseId:%v", testCaseId)
		return fmt.Errorf("body and luaBody cannot have values at the same time, testCaseId:%v", testCaseId)
	}
	return nil
}
