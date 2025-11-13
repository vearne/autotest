package config

import (
	"time"

	"github.com/vearne/autotest/internal/rule"
)

type AutoTestConfig struct {
	Global struct {
		WorkerNum          int           `yaml:"worker_num"`
		IgnoreTestCaseFail bool          `yaml:"ignore_testcase_fail"`
		Debug              bool          `yaml:"debug"`
		RequestTimeout     time.Duration `yaml:"request_timeout"`

		// 重试配置
		Retry struct {
			MaxAttempts        int           `yaml:"max_attempts"`
			RetryDelay         time.Duration `yaml:"retry_delay"`
			RetryOnStatusCodes []int         `yaml:"retry_on_status_codes"`
		} `yaml:"retry"`

		// 并发控制
		Concurrency struct {
			MaxConcurrentRequests int `yaml:"max_concurrent_requests"`
			RateLimitPerSecond    int `yaml:"rate_limit_per_second"`
		} `yaml:"concurrency"`

		// 缓存配置
		Cache struct {
			Enabled bool          `yaml:"enabled"`
			TTL     time.Duration `yaml:"ttl"`
			MaxSize int           `yaml:"max_size"`
		} `yaml:"cache"`

		Logger struct {
			Level    string `yaml:"level"`
			FilePath string `yaml:"file_path"`
		} `yaml:"logger"`

		Report struct {
			DirPath      string   `yaml:"dir_path"`
			Formats      []string `yaml:"formats"`
			TemplatePath string   `yaml:"template_path"`
		} `yaml:"report"`

		// 通知配置
		Notifications struct {
			Enabled    bool   `yaml:"enabled"`
			WebhookURL string `yaml:"webhook_url"`
			OnFailure  bool   `yaml:"on_failure"`
			OnSuccess  bool   `yaml:"on_success"`
		} `yaml:"notifications"`
	} `yaml:"global"`

	HttpRuleFiles []string                     `yaml:"http_rule_files"`
	GrpcRuleFiles []string                     `yaml:"grpc_rule_files"`
	Environments  map[string]map[string]string `yaml:"environments"`
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
