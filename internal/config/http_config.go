package autotest

import "github.com/vearne/autotest/internal/rule"

type Request struct {
	URL     string   `yaml:"url"`
	Headers []string `yaml:"headers"`
	Body    string   `yaml:"body"`
}

type TestCase struct {
	ID          uint64            `yaml:"id"`
	Request     Request           `yaml:"request"`
	Rules       []rule.VerifyRule `yaml:"rules"`
	DependOnIDs []uint64          `yaml:"dependOnIDs,omitempty"`
}
