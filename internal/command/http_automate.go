package command

import (
	"github.com/vearne/autotest/internal/config"
	slog "github.com/vearne/simplelog"
)

func HttpAutomateTest(httpTestCases map[string][]config.TestCase) {
	total := 0
	for _, testcases := range httpTestCases {
		total += len(testcases)
	}
	slog.Info("HttpTestCases, total:%v", total)
}
