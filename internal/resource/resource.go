package resource

import (
	"bufio"
	"encoding/json"
	"github.com/go-resty/resty/v2"
	"github.com/vearne/autotest/internal/config"
	"github.com/vearne/autotest/internal/rule"
	slog "github.com/vearne/simplelog"
	"gopkg.in/yaml.v3"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"sync"
)

var GlobalConfig config.AutoTestConfig
var HttpTestCases map[string][]*config.TestCase

var EnvVars map[string]string
var CustomerVars sync.Map

var RestyClient *resty.Client

func init() {
	EnvVars = make(map[string]string, 10)
	HttpTestCases = make(map[string][]*config.TestCase, 10)
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

	slog.Info("2) parse http rule files")
	// read testcase
	// 1) http testcase
	for idx, f := range GlobalConfig.HttpRuleFiles {
		slog.Info("parse file:%v", f)

		b, err := readFile(f)
		if err != nil {
			slog.Error("readFile:%v, error:%v", f, err)
			return err
		}

		var testcases []*config.TestCase
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
