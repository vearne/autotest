package command

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/antchfx/jsonquery"
	"github.com/antchfx/xpath"
	"github.com/urfave/cli/v3"
	"github.com/vearne/autotest/internal/resource"
	"github.com/vearne/autotest/internal/rule"
	"github.com/vearne/autotest/internal/util"
	slog "github.com/vearne/simplelog"
	"github.com/vearne/zaplog"
)

var (
	ErrorIDduplicate        = errors.New("testcase's ID duplicate")
	ErrorDependencyNotExist = errors.New("dependency does not exist")
	ErrorIDRestrict         = errors.New("the ID of the dependent testcase must be smaller than the ID of the current testcase")
)

// UnifiedTestResults 统一的测试结果
type UnifiedTestResults struct {
	TotalTests  int
	PassedTests int
	FailedTests int
	FailedCases []string
}

// CombineResults 合并HTTP和gRPC测试结果
func CombineResults(httpResults, grpcResults *UnifiedTestResults) *UnifiedTestResults {
	if httpResults == nil {
		httpResults = &UnifiedTestResults{}
	}
	if grpcResults == nil {
		grpcResults = &UnifiedTestResults{}
	}

	combined := &UnifiedTestResults{
		TotalTests:  httpResults.TotalTests + grpcResults.TotalTests,
		PassedTests: httpResults.PassedTests + grpcResults.PassedTests,
		FailedTests: httpResults.FailedTests + grpcResults.FailedTests,
	}

	// 合并失败用例
	combined.FailedCases = append(combined.FailedCases, httpResults.FailedCases...)
	combined.FailedCases = append(combined.FailedCases, grpcResults.FailedCases...)

	return combined
}

func RunTestCases(ctx context.Context, cmd *cli.Command) error {
	confFilePath := cmd.String("config-file")
	slog.Info("config-file:%v", confFilePath)

	envFilePath := cmd.String("env-file")
	slog.Info("env-file:%v", envFilePath)

	environment := cmd.String("environment")
	slog.Info("environment:%v", environment)

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
	// 3. initialize logger & RestyClient & RetryClient & Cache & RateLimiter & EnvironmentManager & ReportGenerator & NotificationService
	slog.Info("3. Initialize logger&RestyClient&RetryClient&Cache&RateLimiter&EnvironmentManager&ReportGenerator&NotificationService")
	loggerConfig := resource.GlobalConfig.Global.Logger
	slog.Info("loggerConfig:FilePath:%v, level:%v", loggerConfig.FilePath, loggerConfig.Level)
	zaplog.InitLogger(loggerConfig.FilePath, loggerConfig.Level)
	resource.InitRestyClient(resource.GlobalConfig.Global.Debug)
	resource.InitRetryClient()
	resource.InitCacheManager()
	resource.InitRateLimiter()
	resource.InitEnvironmentManager()
	resource.InitReportGenerator()
	resource.InitNotificationService()

	// 4. Load specified environment
	slog.Info("4. Load environment configuration")
	err = resource.LoadEnvironment(environment)
	if err != nil {
		slog.Error("failed to load environment '%s': %v", environment, err)
		return err
	}

	// 5. Initialize the executor and execute the testcase concurrently
	// (if the execution fails, you may need to explain the reason for the failure)
	slog.Info("5. Execute test cases")
	startTime := time.Now()

	httpResults := HttpAutomateTest(resource.HttpTestCases)
	grpcResults := GrpcAutomateTest(resource.GrpcTestCases)

	endTime := time.Now()
	totalDuration := endTime.Sub(startTime)

	// 6. output cache statistics
	slog.Info("6. Cache statistics")
	outputCacheStats()

	// 7. generate unified reports and send notifications
	slog.Info("7. Generate reports and send notifications")
	err = generateUnifiedReportsAndNotifications(httpResults, grpcResults, startTime, endTime, totalDuration)
	if err != nil {
		slog.Error("Failed to generate reports or send notifications: %v", err)
		// 不返回错误，因为测试已经完成，报告生成失败不应该影响整体结果
	}

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
			// 1.1 body and luaBody cannot have values at the same time
			if len(tc.Request.Body) > 0 && len(tc.Request.LuaBody) > 0 {
				slog.Error("request error, testCaseId:%v", tc.ID)
				return fmt.Errorf("body and luaBody cannot have values at the same time, testCaseId:%v", tc.ID)
			}

			// 1.2 verify rule
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
			// 1.1 body and luaBody cannot have values at the same time
			if len(tc.Request.Body) > 0 && len(tc.Request.LuaBody) > 0 {
				slog.Error("request error, testCaseId:%v", tc.ID)
				return fmt.Errorf("body and luaBody cannot have values at the same time, testCaseId:%v", tc.ID)
			}

			// 1.2 verify rule
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

// outputCacheStats 输出缓存统计信息
func outputCacheStats() {
	if resource.CacheManager == nil {
		return
	}

	stats := resource.CacheManager.GetStats()
	slog.Info("=== Cache Statistics ===")

	for name, stat := range stats {
		hits := stat["hits"].(int64)
		misses := stat["misses"].(int64)
		hitRate := stat["hit_rate"].(float64)
		size := stat["size"].(int)

		slog.Info("Cache [%s]: Hits=%d, Misses=%d, HitRate=%.2f%%, Size=%d",
			name, hits, misses, hitRate, size)
	}
}

// generateUnifiedReportsAndNotifications 生成统一报告并发送通知
func generateUnifiedReportsAndNotifications(httpResults, grpcResults *UnifiedTestResults, startTime, endTime time.Time, totalDuration time.Duration) error {
	// 合并测试结果
	combinedResults := CombineResults(httpResults, grpcResults)

	if combinedResults.TotalTests == 0 {
		slog.Info("No test cases to report")
		return nil
	}

	// 构建报告数据
	reportData := util.ReportData{}
	reportData.Summary.TotalTests = combinedResults.TotalTests
	reportData.Summary.PassedTests = combinedResults.PassedTests
	reportData.Summary.FailedTests = combinedResults.FailedTests
	reportData.Summary.SkippedTests = 0 // 当前架构没有跳过的测试
	reportData.Summary.Duration = totalDuration
	reportData.Summary.StartTime = startTime
	reportData.Summary.EndTime = endTime

	if combinedResults.TotalTests > 0 {
		reportData.Summary.PassRate = float64(combinedResults.PassedTests) / float64(combinedResults.TotalTests) * 100
	}

	// 转换测试用例详情（这里简化处理，实际项目中可能需要更详细的信息）
	for i, failedCase := range combinedResults.FailedCases {
		testCaseResult := util.TestCaseResult{
			ID:          uint64(i + 1),
			Description: failedCase,
			Status:      "failed",
			Duration:    time.Second, // 简化处理
			StartTime:   startTime,
			EndTime:     endTime,
			ErrorMsg:    "Test case failed", // 简化处理
		}
		reportData.TestCases = append(reportData.TestCases, testCaseResult)
	}

	// 生成报告
	if resource.ReportGenerator != nil {
		slog.Info("Generating reports with formats: %v", resource.GlobalConfig.Global.Report.Formats)
		err := resource.ReportGenerator.GenerateReports(reportData)
		if err != nil {
			slog.Error("Failed to generate reports: %v", err)
			return fmt.Errorf("failed to generate reports: %w", err)
		}
		slog.Info("Reports generated successfully")
	} else {
		slog.Warn("ReportGenerator not initialized, skipping report generation")
	}

	// 发送通知
	if resource.NotificationService != nil {
		notificationResult := util.TestResult{
			TotalTests:  combinedResults.TotalTests,
			PassedTests: combinedResults.PassedTests,
			FailedTests: combinedResults.FailedTests,
			Duration:    totalDuration,
			StartTime:   startTime,
			EndTime:     endTime,
			FailedCases: combinedResults.FailedCases,
		}

		err := resource.NotificationService.SendTestResult(notificationResult)
		if err != nil {
			slog.Error("Failed to send notification: %v", err)
			return fmt.Errorf("failed to send notification: %w", err)
		}
		slog.Info("Notification sent successfully")
	} else {
		slog.Warn("NotificationService not initialized, skipping notification")
	}

	return nil
}
