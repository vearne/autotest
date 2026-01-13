package resource

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"sync"
	"sync/atomic"
	"time"

	"github.com/go-resty/resty/v2"
	"github.com/vearne/autotest/internal/config"
	"github.com/vearne/autotest/internal/luavm"
	"github.com/vearne/autotest/internal/model"
	"github.com/vearne/autotest/internal/rule"
	"github.com/vearne/autotest/internal/util"
	slog "github.com/vearne/simplelog"
	"golang.org/x/sync/singleflight"
	"gopkg.in/yaml.v3"
)

var GlobalConfig config.AutoTestConfig
var HttpTestCases map[string][]*config.TestCaseHttp
var GrpcTestCases map[string][]*config.TestCaseGrpc

var EnvVars map[string]string
var CustomerVars sync.Map

var RestyClient *resty.Client
var RetryClient *util.RetryableHTTPClient
var TerminationFlag atomic.Bool

var DescSourceCache *model.DescSourceCache
var CacheManager *util.CacheManager
var RateLimiter *util.RateLimiter
var EnvironmentManager *util.EnvironmentManager
var ReportGenerator *util.ReportGenerator
var NotificationService *util.NotificationService

var SingleFlightGroup singleflight.Group

func init() {
	EnvVars = make(map[string]string, 10)
	HttpTestCases = make(map[string][]*config.TestCaseHttp, 10)
	GrpcTestCases = make(map[string][]*config.TestCaseGrpc, 10)
	TerminationFlag.Store(false)
	DescSourceCache = model.NewDescSourceCache()
}

// InitCacheManager 初始化缓存管理器
func InitCacheManager() {
	CacheManager = util.NewCacheManager(GlobalConfig)
	slog.Info("Cache manager initialized")
}

// InitRateLimiter 初始化并发控制器
func InitRateLimiter() {
	// 设置合理的默认值
	maxConcurrent := 20 // 默认允许20个并发请求
	rateLimit := 50     // 默认每秒50个请求

	// 如果用户配置了，使用用户的配置
	if GlobalConfig.Global.Concurrency.MaxConcurrentRequests > 0 {
		maxConcurrent = GlobalConfig.Global.Concurrency.MaxConcurrentRequests
	}
	if GlobalConfig.Global.Concurrency.RateLimitPerSecond > 0 {
		rateLimit = GlobalConfig.Global.Concurrency.RateLimitPerSecond
	}

	RateLimiter = util.NewRateLimiter(maxConcurrent, rateLimit)
	slog.Info("Rate limiter initialized: MaxConcurrent=%d, RateLimit=%d/s", maxConcurrent, rateLimit)
}

// InitEnvironmentManager 初始化环境管理器
func InitEnvironmentManager() {
	EnvironmentManager = util.NewEnvironmentManager(GlobalConfig)
	slog.Info("Environment manager initialized")
}

// LoadEnvironment 加载指定环境的配置
func LoadEnvironment(envName string) error {
	if EnvironmentManager == nil {
		return fmt.Errorf("environment manager not initialized")
	}

	// 显示可用环境列表（用于调试和用户友好提示）
	availableEnvs := EnvironmentManager.ListAvailableEnvironments()
	if len(availableEnvs) > 0 {
		slog.Debug("Available environments: %v", availableEnvs)
	}

	if envName == "" {
		if len(availableEnvs) > 0 {
			slog.Info("No environment specified, using default configuration. Available environments: %v", availableEnvs)
		} else {
			slog.Info("No environment specified, using default configuration")
		}
		return nil
	}

	// 验证环境是否存在
	if err := EnvironmentManager.ValidateEnvironment(envName); err != nil {
		return fmt.Errorf("%w. Available environments: %v", err, availableEnvs)
	}

	// 加载环境配置
	if err := EnvironmentManager.LoadEnvironment(envName); err != nil {
		return fmt.Errorf("failed to load environment '%s': %w", envName, err)
	}

	// 将环境变量合并到全局 EnvVars 中
	envVars := EnvironmentManager.GetAllVars()
	for key, value := range envVars {
		EnvVars[key] = value
	}

	slog.Info("Environment '%s' loaded successfully with %d variables", envName, len(envVars))
	return nil
}

func InitRestyClient(debug bool) {
	httpClient := http.Client{
		Transport: &http.Transport{
			MaxIdleConnsPerHost: 500,
		},
	}
	RestyClient = resty.NewWithClient(&httpClient)
	RestyClient.SetDebug(debug)
}

// InitRetryClient 初始化带重试功能的HTTP客户端
func InitRetryClient() {
	// 设置重试机制的默认值
	if GlobalConfig.Global.Retry.MaxAttempts <= 0 {
		GlobalConfig.Global.Retry.MaxAttempts = 3 // 默认最多重试3次
	}
	if GlobalConfig.Global.Retry.RetryDelay <= 0 {
		GlobalConfig.Global.Retry.RetryDelay = time.Second // 默认重试间隔1秒
	}
	if len(GlobalConfig.Global.Retry.RetryOnStatusCodes) == 0 {
		// 默认重试的状态码：服务器错误、网关错误、服务不可用、网关超时、请求超时、请求过多
		GlobalConfig.Global.Retry.RetryOnStatusCodes = []int{500, 502, 503, 504, 408, 429}
	}

	RetryClient = util.NewRetryableHTTPClient(RestyClient, GlobalConfig)

	slog.Info("RetryClient initialized: MaxAttempts=%d, RetryDelay=%v, RetryCodes=%v",
		GlobalConfig.Global.Retry.MaxAttempts,
		GlobalConfig.Global.Retry.RetryDelay,
		GlobalConfig.Global.Retry.RetryOnStatusCodes)
}

// InitReportGenerator 初始化报告生成器
func InitReportGenerator() {
	ReportGenerator = util.NewReportGenerator(GlobalConfig)
	slog.Info("ReportGenerator initialized")
}

// InitNotificationService 初始化通知服务
func InitNotificationService() {
	NotificationService = util.NewNotificationService(GlobalConfig)

	if GlobalConfig.Global.Notifications.Enabled {
		slog.Info("NotificationService initialized: Enabled=%v, WebhookURL=%s, OnFailure=%v, OnSuccess=%v",
			GlobalConfig.Global.Notifications.Enabled,
			GlobalConfig.Global.Notifications.WebhookURL,
			GlobalConfig.Global.Notifications.OnFailure,
			GlobalConfig.Global.Notifications.OnSuccess)
	} else {
		slog.Info("NotificationService initialized: Disabled")
	}
}

// InitLuaVM 初始化Lua虚拟机并加载预加载文件
func InitLuaVM() error {
	if len(GlobalConfig.Global.Lua.PreloadFiles) == 0 {
		slog.Info("No Lua preload files configured")
		return nil
	}

	slog.Info("Initializing Lua VM with %d preload files", len(GlobalConfig.Global.Lua.PreloadFiles))

	configDir := ""

	if len(GlobalConfig.HttpRuleFiles) > 0 {
		configDir = filepath.Dir(GlobalConfig.HttpRuleFiles[0])
	} else if len(GlobalConfig.GrpcRuleFiles) > 0 {
		configDir = filepath.Dir(GlobalConfig.GrpcRuleFiles[0])
	}

	absolutePaths := make([]string, 0, len(GlobalConfig.Global.Lua.PreloadFiles))
	for _, file := range GlobalConfig.Global.Lua.PreloadFiles {
		absPath := file
		if !filepath.IsAbs(file) && configDir != "" {
			absPath = filepath.Join(configDir, file)
		}
		absolutePaths = append(absolutePaths, absPath)

		if _, err := os.Stat(absPath); os.IsNotExist(err) {
			return fmt.Errorf("lua preload file not found: %s (resolved to %s)", file, absPath)
		}
	}

	if err := luavm.LoadPreloadLuaFiles(absolutePaths); err != nil {
		return fmt.Errorf("failed to load Lua preload files: %w", err)
	}

	slog.Info("Lua VM initialized successfully with %d preload files", len(absolutePaths))
	return nil
}

// ValidateLuaFiles 验证Lua预加载文件
func ValidateLuaFiles() error {
	if len(GlobalConfig.Global.Lua.PreloadFiles) == 0 {
		return nil
	}

	slog.Info("Validating %d Lua preload files", len(GlobalConfig.Global.Lua.PreloadFiles))

	configDir := ""

	if len(GlobalConfig.HttpRuleFiles) > 0 {
		configDir = filepath.Dir(GlobalConfig.HttpRuleFiles[0])
	} else if len(GlobalConfig.GrpcRuleFiles) > 0 {
		configDir = filepath.Dir(GlobalConfig.GrpcRuleFiles[0])
	}

	for _, file := range GlobalConfig.Global.Lua.PreloadFiles {
		absPath := file
		if !filepath.IsAbs(file) && configDir != "" {
			absPath = filepath.Join(configDir, file)
		}

		if _, err := os.Stat(absPath); os.IsNotExist(err) {
			return fmt.Errorf("lua preload file not found: %s (resolved to %s)", file, absPath)
		}
	}

	return nil
}

func ParseConfigFile(filePath string) error {
	slog.Info("[start]ParseConfigFile:%v", filePath)

	b, err := readFile(filePath)
	if err != nil {
		return err
	}

	slog.Info("1) parse global file")
	err = yaml.Unmarshal(b, &GlobalConfig)
	if err != nil {
		return err
	}

	if GlobalConfig.Global.RequestTimeout < time.Second {
		GlobalConfig.Global.RequestTimeout = time.Second
	}

	slog.Info("2) parse http rule files")
	// read testcase
	// 1) http testcase
	for idx, f := range GlobalConfig.HttpRuleFiles {
		slog.Info("parse http rule file:%v", f)

		b, err = readFile(f)
		if err != nil {
			slog.Error("readFile:%v, error:%v", f, err)
			return err
		}

		var testcases []*config.TestCaseHttp
		err = yaml.Unmarshal(b, &testcases)
		if err != nil {
			slog.Error("file:%v parse error, %v", f, err)
			return err
		}

		for i := 0; i < len(testcases); i++ {
			c := testcases[i]
			c.Request.Body = strings.ReplaceAll(c.Request.Body, "\n", "")
			c.VerifyRules = make([]rule.VerifyRule, 0)
			for _, r := range c.OriginRules {
				b, _ = json.Marshal(r)
				switch r["name"] {
				case "HttpStatusEqualRule":
					var item rule.HttpStatusEqualRule
					err = json.Unmarshal(b, &item)
					if err != nil {
						slog.Error("parse rule[HttpStatusEqualRule], %v", err)
						return err
					}
					c.VerifyRules = append(c.VerifyRules, &item)
				case "HttpBodyEqualRule":
					var item rule.HttpBodyEqualRule
					err = json.Unmarshal(b, &item)
					if err != nil {
						slog.Error("parse rule[HttpBodyEqualRule], %v", err)
						return err
					}
					c.VerifyRules = append(c.VerifyRules, &item)
				case "HttpBodyAtLeastOneRule":
					var item rule.HttpBodyAtLeastOneRule
					err = json.Unmarshal(b, &item)
					if err != nil {
						slog.Error("parse rule[HttpBodyAtLeastOneRule], %v", err)
						return err
					}
					c.VerifyRules = append(c.VerifyRules, &item)
				case "HttpLuaRule":
					var item rule.HttpLuaRule
					err = json.Unmarshal(b, &item)
					if err != nil {
						slog.Error("parse rule[HttpLuaRule], %v", err)
						return err
					}
					c.VerifyRules = append(c.VerifyRules, &item)
				default:
					return fmt.Errorf("unknow http-VerifyRule:%v", r["name"])
				}
			}

			if c.Export != nil {
				if len(c.Export.Type) <= 0 {
					c.Export.Type = "string"
				}
			}
		}

		slog.Info("parse file:%v, len(testcases):%v", f, len(testcases))
		absolutePath, _ := filepath.Abs(f)
		HttpTestCases[absolutePath] = testcases
		GlobalConfig.HttpRuleFiles[idx] = absolutePath
	}

	slog.Info("3) parse grpc rule files")
	// 2) grpc testcase
	for idx, f := range GlobalConfig.GrpcRuleFiles {
		slog.Info("parse grpc rule:%v", f)

		b, err = readFile(f)
		if err != nil {
			slog.Error("readFile:%v, error:%v", f, err)
			return err
		}

		var testcases []*config.TestCaseGrpc
		err = yaml.Unmarshal(b, &testcases)
		if err != nil {
			slog.Error("file:%v parse error, %v", f, err)
			return err
		}

		for i := 0; i < len(testcases); i++ {
			c := testcases[i]
			c.Request.Body = strings.ReplaceAll(c.Request.Body, "\n", "")
			c.VerifyRules = make([]rule.VerifyRuleGrpc, 0)
			for _, r := range c.OriginRules {
				b, _ := json.Marshal(r)
				switch r["name"] {
				case "GrpcCodeEqualRule":
					var item rule.GrpcCodeEqualRule
					err = json.Unmarshal(b, &item)
					if err != nil {
						slog.Error("parse rule[GrpcCodeEqualRule], %v", err)
						return err
					}
					c.VerifyRules = append(c.VerifyRules, &item)
				case "GrpcBodyEqualRule":
					var item rule.GrpcBodyEqualRule
					err = json.Unmarshal(b, &item)
					if err != nil {
						slog.Error("parse rule[GrpcBodyEqualRule], %v", err)
						return err
					}
					c.VerifyRules = append(c.VerifyRules, &item)

				case "GrpcBodyAtLeastOneRule":
					var item rule.GrpcBodyAtLeastOneRule
					err = json.Unmarshal(b, &item)
					if err != nil {
						slog.Error("parse rule[GrpcBodyAtLeastOneRule], %v", err)
						return err
					}
					c.VerifyRules = append(c.VerifyRules, &item)

				case "GrpcLuaRule":
					var item rule.GrpcLuaRule
					err = json.Unmarshal(b, &item)
					if err != nil {
						slog.Error("parse rule[GrpcLuaRule], %v", err)
						return err
					}
					c.VerifyRules = append(c.VerifyRules, &item)
				default:
					return fmt.Errorf("unknow Grpc-VerifyRule:%v", r["name"])
				}
			}

			if c.Export != nil {
				if len(c.Export.Type) <= 0 {
					c.Export.Type = "string"
				}
			}
		}

		slog.Info("parse file:%v, len(testcases):%v", f, len(testcases))
		absolutePath, _ := filepath.Abs(f)
		GrpcTestCases[absolutePath] = testcases
		GlobalConfig.GrpcRuleFiles[idx] = absolutePath
	}

	slog.Info("4) parse lua preload files")
	for idx, f := range GlobalConfig.Global.Lua.PreloadFiles {
		slog.Info("parse lua preload file:%v", f)
		absolutePath, err := filepath.Abs(f)
		if err != nil {
			slog.Error("convert to absolute path failed, file:%v, error:%v", f, err)
			return err
		}
		GlobalConfig.Global.Lua.PreloadFiles[idx] = absolutePath
		slog.Info("lua preload file converted to absolute path:%v", absolutePath)
	}

	// 5) modify report path
	newPath := filepath.Join(GlobalConfig.Global.Report.DirPath, "autotest_"+strconv.Itoa(int(time.Now().Unix())))
	// 创建单个文件夹
	err = os.Mkdir(newPath, 0755)
	if err != nil {
		return err
	}

	GlobalConfig.Global.Report.DirPath = newPath
	slog.Info("[end]ParseConfigFile")
	return nil
}

func readFile(filePath string) ([]byte, error) {
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		return nil, os.ErrNotExist
	}
	return os.ReadFile(filePath)
}
