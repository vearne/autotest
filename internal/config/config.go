package config

import (
	"github.com/vearne/autotest/internal/rule"
	"time"
)

type AutoTestConfig struct {
	Global struct {
		WorkerNum          int           `yaml:"worker_num"`
		IgnoreTestCaseFail bool          `yaml:"ignore_testcase_fail"`
		Debug              bool          `yaml:"debug"`
		RequestTimeout     time.Duration `yaml:"request_timeout"`
		Logger             struct {
			Level    string `yaml:"level"`
			FilePath string `yaml:"file_path"`
		} `yaml:"logger"`
		Report struct {
			DirPath string `yaml:"dir_path"`
		} `yaml:"report"`
	} `yaml:"global"`

	HttpRuleFiles []string `yaml:"http_rule_files"`
	GrpcRuleFiles []string `yaml:"grpc_rule_files"`
}

type TestCaseHttp struct {
	ID   uint64 `yaml:"id"`
	Desc string `yaml:"desc"`
	// Delay for a while before executing
	Delay       time.Duration    `yaml:"delay,omitempty"`
	Request     RequestHttp      `yaml:"request"`
	OriginRules []map[string]any `yaml:"rules" json:"-"`
	DependOnIDs []uint64         `yaml:"dependOnIDs,omitempty"`
	Export      *Export          `yaml:"export"`
	VerifyRules []rule.VerifyRule
}

func (t *TestCaseHttp) GetID() uint64 {
	return t.ID
}

type Export struct {
	Xpath    string `yaml:"xpath"`
	ExportTo string `yaml:"exportTo"`
	Type     string `yaml:"type"`
}

type RequestHttp struct {
	Method  string   `yaml:"method"`
	URL     string   `yaml:"url"`
	Headers []string `yaml:"headers"`
	Body    string   `yaml:"body"`
	LuaBody string   `yaml:"luaBody"`
}

type TestCaseGrpc struct {
	ID   uint64 `yaml:"id"`
	Desc string `yaml:"desc"`
	// Delay for a while before executing
	Delay       time.Duration    `yaml:"delay,omitempty"`
	Request     RequestGrpc      `yaml:"request"`
	OriginRules []map[string]any `yaml:"rules" json:"-"`
	DependOnIDs []uint64         `yaml:"dependOnIDs,omitempty"`
	Export      *Export          `yaml:"export"`
	VerifyRules []rule.VerifyRuleGrpc
}

func (t *TestCaseGrpc) GetID() uint64 {
	return t.ID
}

type RequestGrpc struct {
	Address string   `yaml:"address"`
	Symbol  string   `yaml:"symbol"`
	Headers []string `yaml:"headers"`
	Body    string   `yaml:"body"`
	LuaBody string   `yaml:"luaBody"`
}
