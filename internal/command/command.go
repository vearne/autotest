package command

import (
	"context"
	"github.com/antchfx/xpath"
	"github.com/urfave/cli/v3"
	"github.com/vearne/autotest/internal/resource"
	"github.com/vearne/autotest/internal/rule"
	slog "github.com/vearne/simplelog"
	"github.com/vearne/zaplog"
)

func RunTestCases(ctx context.Context, cmd *cli.Command) error {
	confFilePath := cmd.String("config-file")
	slog.Info("config-file:%v", confFilePath)

	envFilePath := cmd.String("env-file")
	slog.Info("env-file:%v", envFilePath)

	// 1. 解析配置文件
	slog.Info("1. Parse config file")
	// 1.1 config file
	err := resource.ParseConfigFile(confFilePath)
	if err != nil {
		slog.Error("config file parse error, %v", err)
		return err
	}
	// 1.2 env file
	err = resource.ParseEnvFile(envFilePath)
	if err != nil {
		slog.Error("env file parse error, %v", err)
		return err
	}
	// 2. 初始化logger
	slog.Info("2. Initialize logger")
	loggerConfig := resource.GlobalConfig.Global.Logger
	zaplog.InitLogger(loggerConfig.FilePath, loggerConfig.Level)

	// 3. 初始化执行器, 并发执行testcase (执行失败可能需要解释失败的原因)
	slog.Info("3. Execute test cases")
	HttpAutomateTest(resource.HttpTestCases)
	// 4. 生成报告
	return nil
}

func ValidateConfig(ctx context.Context, cmd *cli.Command) error {
	// 检查testcase的xpath语法是否正确

	filePath := cmd.String("config-file")
	slog.Info("config-file:%v", filePath)

	err := resource.ParseConfigFile(filePath)
	if err != nil {
		slog.Error("config file parse error, %v", err)
		return err
	}

	slog.Info("=== validate config file ===")
	for filePath, testcases := range resource.HttpTestCases {
		slog.Info("filePath:%v, len(testcases):%v", filePath, len(testcases))
		for _, tc := range testcases {
			for _, r := range tc.VerifyRules {
				switch r.Name() {
				case "HttpBodyEqualRule":
					rule := r.(*rule.HttpBodyEqualRule)
					_, err := xpath.Compile(rule.Xpath)
					if err != nil {
						slog.Error("rule error, testCaseId:%v, xpath:%v", tc.ID, rule.Xpath)
						return err
					}
				case "HttpBodyAtLeastOneRule":
					rule := r.(*rule.HttpBodyEqualRule)
					_, err := xpath.Compile(rule.Xpath)
					if err != nil {
						slog.Error("rule error, testCaseId:%v, xpath:%v", tc.ID, rule.Xpath)
						return err
					}
				default:
					slog.Debug("ignore rule:%v", r.Name())
				}
			}
		}
	}
	return nil
}
