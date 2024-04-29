package config

import (
	"github.com/vearne/autotest/internal/rule"
)

type AutoTestConfig struct {
	Global struct {
		WorkerNum          int  `yaml:"worker_num"`
		IgnoreTestCaseFail bool `yaml:"ignore_testcase_fail"`
		Debug              bool `yaml:"debug"`
		Logger             struct {
			Level    string `yaml:"level"`
			FilePath string `yaml:"filepath"`
		}
	} `yaml:"global"`

	HttpRuleFiles []string `yaml:"http_rule_files"`
	GrpcRuleFiles []string `yaml:"grpc_rule_files"`
}

type TestCase struct {
	ID          uint64           `yaml:"id"`
	Request     Request          `yaml:"request"`
	OriginRules []map[string]any `yaml:"rules" json:"-"`
	DependOnIDs []uint64         `yaml:"dependOnIDs,omitempty"`
	Export      *Export          `yaml:"export"`
	VerifyRules []rule.VerifyRule
}

type Export struct {
	Xpath    string `yaml:"xpath"`
	ExportTo string `yaml:"exportTo"`
	Type     string `yaml:"type"`
}

type Request struct {
	Method  string   `yaml:"method"`
	URL     string   `yaml:"url"`
	Headers []string `yaml:"headers"`
	Body    string   `yaml:"body"`
}