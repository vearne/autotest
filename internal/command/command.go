package command

import (
	"context"
	"errors"
	"github.com/antchfx/jsonquery"
	"github.com/antchfx/xpath"
	"github.com/urfave/cli/v3"
	"github.com/vearne/autotest/internal/resource"
	"github.com/vearne/autotest/internal/rule"
	slog "github.com/vearne/simplelog"
	"github.com/vearne/zaplog"
	"strings"
)

var (
	ErrorIDduplicate        = errors.New("testcase's ID duplicate")
	ErrorDependencyNotExist = errors.New("dependency does not exist")
	ErrorIDRestrict         = errors.New("the ID of the dependent testcase must be smaller than the ID of the current testcase")
)

func RunTestCases(ctx context.Context, cmd *cli.Command) error {
	confFilePath := cmd.String("config-file")
	slog.Info("config-file:%v", confFilePath)

	envFilePath := cmd.String("env-file")
	slog.Info("env-file:%v", envFilePath)

	// 1. Parsing configuration files
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
	// 2. validate config
	slog.Info("2. validate config file")
	err = AllCheck()
	if err != nil {
		slog.Error("validate config file, error:%v", err)
		return err
	}
	// 3. initialize logger & RestyClient
	slog.Info("3. Initialize logger&RestyClient")
	loggerConfig := resource.GlobalConfig.Global.Logger
	slog.Info("loggerConfig:FilePath:%v, level:%v", loggerConfig.FilePath, loggerConfig.Level)
	zaplog.InitLogger(loggerConfig.FilePath, loggerConfig.Level)
	resource.InitRestyClient(resource.GlobalConfig.Global.Debug)

	// 4. Initialize the executor and execute the testcase concurrently
	// (if the execution fails, you may need to explain the reason for the failure)
	slog.Info("4. Execute test cases")
	HttpAutomateTest(resource.HttpTestCases)
	GrpcAutomateTest(resource.GrpcTestCases)
	// 6. generate report
	slog.Info("6. output report to file")
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
	return AllCheck()
}

func AllCheck() error {
	var err error
	err = CheckTestCaseHttp()
	if err != nil {
		return err
	}

	err = CheckTestCaseGrpc()
	if err != nil {
		return err
	}

	return nil
}

func CheckTestCaseHttp() error {
	slog.Info("CheckTestCaseHttp")

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
				return ErrorIDduplicate
			}
			exist[tc.ID] = struct{}{}
		}
	}
	slog.Info("3. check dependencies")
	for filePath, testcases := range resource.HttpTestCases {
		slog.Info("filePath:%v, len(testcases):%v", filePath, len(testcases))
		exist := make(map[uint64]struct{})
		for _, tc := range testcases {
			exist[tc.ID] = struct{}{}
		}
		for _, tc := range testcases {
			for _, dID := range tc.DependOnIDs {
				if _, ok := exist[dID]; !ok {
					slog.Error("For testcase %v, the testcase %v that it depend on does not exist",
						tc.ID, dID)
					return ErrorDependencyNotExist
				}
				if dID >= tc.ID {
					slog.Error("The ID of the dependent testcase must be smaller "+
						"than the ID of the current testcase, testcase:%v, dependent testcase:%v",
						tc.ID, dID)
					return ErrorIDRestrict
				}
			}
		}
	}
	return nil
}

func CheckTestCaseGrpc() error {
	slog.Info("CheckTestCaseGrpc")

	slog.Info("1. check xpath")
	for filePath, testcases := range resource.GrpcTestCases {
		slog.Info("filePath:%v, len(testcases):%v", filePath, len(testcases))
		for _, tc := range testcases {
			for _, r := range tc.VerifyRules {
				switch r.Name() {
				case "GrpcBodyEqualRule":
					rule := r.(*rule.GrpcBodyEqualRule)
					_, err := xpath.Compile(rule.Xpath)
					if err != nil {
						slog.Error("rule error, testCaseId:%v, xpath:%v", tc.ID, rule.Xpath)
						return err
					}
				case "GrpcBodyAtLeastOneRule":
					rule := r.(*rule.GrpcBodyAtLeastOneRule)
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
	for filePath, testcases := range resource.GrpcTestCases {
		slog.Info("filePath:%v, len(testcases):%v", filePath, len(testcases))
		exist := make(map[uint64]struct{})
		for _, tc := range testcases {
			_, ok := exist[tc.ID]
			if ok {
				slog.Error("filePath:%v, ID [%v] is duplicate", filePath, tc.ID)
				return ErrorIDduplicate
			}
			exist[tc.ID] = struct{}{}
		}
	}
	slog.Info("3. check dependencies")
	for filePath, testcases := range resource.GrpcTestCases {
		slog.Info("filePath:%v, len(testcases):%v", filePath, len(testcases))
		exist := make(map[uint64]struct{})
		for _, tc := range testcases {
			exist[tc.ID] = struct{}{}
		}
		for _, tc := range testcases {
			for _, dID := range tc.DependOnIDs {
				if _, ok := exist[dID]; !ok {
					slog.Error("For testcase %v, the testcase %v that it depend on does not exist",
						tc.ID, dID)
					return ErrorDependencyNotExist
				}
				if dID >= tc.ID {
					slog.Error("The ID of the dependent testcase must be smaller "+
						"than the ID of the current testcase, testcase:%v, dependent testcase:%v",
						tc.ID, dID)
					return ErrorIDRestrict
				}
			}
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
