package command

import (
	"context"
	"github.com/antchfx/jsonquery"
	"github.com/antchfx/xpath"
	"github.com/urfave/cli/v3"
	"github.com/vearne/autotest/internal/resource"
	"github.com/vearne/autotest/internal/rule"
	slog "github.com/vearne/simplelog"
	"github.com/vearne/zaplog"
	"strings"
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
	// 2. 初始化logger & RestyClient
	slog.Info("2. Initialize logger&RestyClient")
	loggerConfig := resource.GlobalConfig.Global.Logger
	zaplog.InitLogger(loggerConfig.FilePath, loggerConfig.Level)
	resource.InitRestyClient(resource.GlobalConfig.Global.Debug)

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
	slog.Info("1. check xpath")
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
					rule := r.(*rule.HttpBodyAtLeastOneRule)
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

	slog.Info("2. check if ID is duplicate")
	for filePath, testcases := range resource.HttpTestCases {
		slog.Info("filePath:%v, len(testcases):%v", filePath, len(testcases))
		exist := make(map[uint64]struct{})
		for _, tc := range testcases {
			_, ok := exist[tc.ID]
			if ok {
				slog.Error("filePath:%v, ID [%v] is duplicate", filePath, tc.ID)
				break
			}
			exist[tc.ID] = struct{}{}
		}
	}
	return nil
}

func ExtractXpath(ctx context.Context, cmd *cli.Command) error {
	// 检查testcase的xpath语法是否正确

	xpathStr := cmd.String("xpath")
	slog.Info("xpathStr:%v", xpathStr)

	jsonStr := cmd.String("json")
	slog.Info("jsonStr:%v", jsonStr)

	_, err := xpath.Compile(xpathStr)
	if err != nil {
		slog.Error("xpath syntax error")
		return nil
	}

	doc, err := jsonquery.Parse(strings.NewReader(jsonStr))
	if err != nil {
		slog.Error("jsonStr format error")
		return nil
	}
	nodes := jsonquery.Find(doc, xpathStr)
	for idx, node := range nodes {
		if node != nil {
			slog.Info("[%v] = %v", idx, node.Value())
		}
	}
	return nil
}
