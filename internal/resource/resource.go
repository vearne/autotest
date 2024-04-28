package resource

import (
	"bufio"
	"encoding/json"
	"fmt"
	"github.com/vearne/autotest/internal/config"
	"github.com/vearne/autotest/internal/rule"
	slog "github.com/vearne/simplelog"
	"gopkg.in/yaml.v3"
	"os"
	"path/filepath"
	"strings"
	"sync"
)

var GlobalConfig config.AutoTestConfig
var HttpTestCases map[string][]config.TestCase

var EnvVars map[string]string
var CustomerVars sync.Map

func init() {
	EnvVars = make(map[string]string, 10)
	HttpTestCases = make(map[string][]config.TestCase, 10)
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
	for _, f := range GlobalConfig.HttpRuleFiles {
		slog.Info("parse file:%v", f)

		b, err := readFile(f)
		if err != nil {
			slog.Error("readFile:%v, error:%v", f, err)
			return err
		}

		var testcases []config.TestCase
		err = yaml.Unmarshal(b, &testcases)
		if err != nil {
			fmt.Println(err)
			slog.Error("file:%v parse error, %v", f, err)
			return err
		}

		for i := 0; i < len(testcases); i++ {
			c := testcases[i]
			c.VerifyRules = make([]rule.VerifyRule, 0)
			for _, r := range c.OriginRules {
				b, _ := json.Marshal(r)
				switch r["name"] {
				case "HttpStatusEqualRule":
					var r rule.HttpStatusEqualRule
					err = json.Unmarshal(b, &r)
					if err != nil {
						slog.Error("parse rule[HttpStatusEqualRule], %v", err)
						return err
					}
					c.VerifyRules = append(c.VerifyRules, &r)
				case "HttpBodyEqualRule":
					var r rule.HttpBodyEqualRule
					err = json.Unmarshal(b, &r)
					if err != nil {
						slog.Error("parse rule[HttpBodyEqualRule], %v", err)
						return err
					}
					c.VerifyRules = append(c.VerifyRules, &r)

				case "HttpBodyAtLeastOneRule":
					var r rule.HttpBodyAtLeastOneRule
					err = json.Unmarshal(b, &r)
					if err != nil {
						slog.Error("parse rule[HttpBodyAtLeastOneRule], %v", err)
						return err
					}
					c.VerifyRules = append(c.VerifyRules, &r)
				}
			}
		}

		slog.Info("parse file:%v, len(testcases):%v", f, len(testcases))
		absolutePath, _ := filepath.Abs(f)
		HttpTestCases[absolutePath] = testcases
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
