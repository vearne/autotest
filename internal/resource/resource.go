package resource

import (
	"bufio"
	"encoding/json"
	"fmt"
	"github.com/go-resty/resty/v2"
	"github.com/vearne/autotest/internal/config"
	"github.com/vearne/autotest/internal/model"
	"github.com/vearne/autotest/internal/rule"
	slog "github.com/vearne/simplelog"
	"golang.org/x/sync/singleflight"
	"gopkg.in/yaml.v3"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"sync"
	"sync/atomic"
	"time"
)

var GlobalConfig config.AutoTestConfig
var HttpTestCases map[string][]*config.TestCaseHttp
var GrpcTestCases map[string][]*config.TestCaseGrpc

var EnvVars map[string]string
var CustomerVars sync.Map

var RestyClient *resty.Client
var TerminationFlag atomic.Bool

var DescSourceCache *model.DescSourceCache

var SingleFlightGroup singleflight.Group

func init() {
	EnvVars = make(map[string]string, 10)
	HttpTestCases = make(map[string][]*config.TestCaseHttp, 10)
	GrpcTestCases = make(map[string][]*config.TestCaseGrpc, 10)
	TerminationFlag.Store(false)
	DescSourceCache = model.NewDescSourceCache()
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

func ParseConfigFile(filePath string) error {
	slog.Info("[start]ParseConfigFile")

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

		b, err := readFile(f)
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
				b, _ := json.Marshal(r)
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

		b, err := readFile(f)
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

	// 3) modify report path
	newPath := filepath.Join(GlobalConfig.Global.Report.DirPath, "autotest_"+strconv.Itoa(int(time.Now().Unix())))
	// 创建单个文件夹
	err = os.Mkdir("newPath", 0755)
	if err != nil {
		return err
	}

	GlobalConfig.Global.Report.DirPath = newPath
	slog.Info("[end]ParseConfigFile")
	return nil
}

func ParseEnvFile(filePath string) error {
	slog.Info("[start]ParseEnvFile")

	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		return os.ErrNotExist
	}

	lines, err := ReadLines(filePath)
	if err != nil {
		return err
	}
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if strings.HasPrefix(line, "#") {
			continue
		}
		strlist := strings.Split(line, "=")
		if len(strlist) == 2 {
			EnvVars[strlist[0]] = strlist[1]
		}
	}

	slog.Info("[end]ParseEnvFile")
	return nil
}

func readFile(filePath string) ([]byte, error) {
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		return nil, os.ErrNotExist
	}
	return os.ReadFile(filePath)
}

// ReadLines reads all lines of the file.
func ReadLines(path string) ([]string, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var lines []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}
	return lines, scanner.Err()
}
