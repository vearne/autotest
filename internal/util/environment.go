package util

import (
	"fmt"
	"os"
	"strings"

	"github.com/vearne/autotest/internal/config"
	slog "github.com/vearne/simplelog"
)

// EnvironmentManager 环境管理器
type EnvironmentManager struct {
	config config.AutoTestConfig
	vars   map[string]string
}

// NewEnvironmentManager 创建环境管理器
func NewEnvironmentManager(cfg config.AutoTestConfig) *EnvironmentManager {
	return &EnvironmentManager{
		config: cfg,
		vars:   make(map[string]string),
	}
}

// LoadEnvironment 加载指定环境的配置
func (em *EnvironmentManager) LoadEnvironment(envName string) error {
	if envName == "" {
		slog.Info("No environment specified, using default configuration")
		return nil
	}

	envConfig, exists := em.config.Environments[envName]
	if !exists {
		return fmt.Errorf("environment '%s' not found in configuration", envName)
	}

	slog.Info("Loading environment: %s", envName)

	// 加载环境变量
	for key, value := range envConfig {
		em.vars[key] = value
		// 同时设置到系统环境变量中，以便模板引擎使用
		os.Setenv(key, value)
		slog.Debug("Set environment variable: %s=%s", key, value)
	}

	return nil
}

// LoadFromFile 从文件加载环境变量
func (em *EnvironmentManager) LoadFromFile(filePath string) error {
	if filePath == "" {
		return nil
	}

	content, err := os.ReadFile(filePath)
	if err != nil {
		return fmt.Errorf("failed to read environment file %s: %w", filePath, err)
	}

	lines := strings.Split(string(content), "\n")
	for lineNum, line := range lines {
		line = strings.TrimSpace(line)

		// 跳过空行和注释
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}

		// 解析 KEY=VALUE 格式
		parts := strings.SplitN(line, "=", 2)
		if len(parts) != 2 {
			slog.Warn("Invalid line %d in environment file %s: %s", lineNum+1, filePath, line)
			continue
		}

		key := strings.TrimSpace(parts[0])
		value := strings.TrimSpace(parts[1])

		// 移除引号
		if (strings.HasPrefix(value, "\"") && strings.HasSuffix(value, "\"")) ||
			(strings.HasPrefix(value, "'") && strings.HasSuffix(value, "'")) {
			value = value[1 : len(value)-1]
		}

		em.vars[key] = value
		os.Setenv(key, value)
		slog.Debug("Loaded from file: %s=%s", key, value)
	}

	slog.Info("Loaded %d variables from environment file: %s", len(em.vars), filePath)
	return nil
}

// SetVariable 设置环境变量
func (em *EnvironmentManager) SetVariable(key, value string) {
	em.vars[key] = value
	os.Setenv(key, value)
	slog.Debug("Set variable: %s=%s", key, value)
}

// GetVariable 获取环境变量
func (em *EnvironmentManager) GetVariable(key string) (string, bool) {
	// 优先从内部变量获取
	if value, exists := em.vars[key]; exists {
		return value, true
	}

	// 然后从系统环境变量获取
	value := os.Getenv(key)
	return value, value != ""
}

// GetAllVariables 获取所有环境变量
func (em *EnvironmentManager) GetAllVariables() map[string]string {
	result := make(map[string]string)

	// 复制内部变量
	for k, v := range em.vars {
		result[k] = v
	}

	return result
}

// ListAvailableEnvironments 列出可用的环境
func (em *EnvironmentManager) ListAvailableEnvironments() []string {
	var envs []string
	for envName := range em.config.Environments {
		envs = append(envs, envName)
	}
	return envs
}

// ValidateEnvironment 验证环境配置
func (em *EnvironmentManager) ValidateEnvironment(envName string) error {
	if envName == "" {
		return nil
	}

	envConfig, exists := em.config.Environments[envName]
	if !exists {
		availableEnvs := em.ListAvailableEnvironments()
		return fmt.Errorf("environment '%s' not found. Available environments: %v", envName, availableEnvs)
	}

	// 验证必需的环境变量
	requiredVars := []string{"HOST"} // 可以根据需要扩展
	for _, requiredVar := range requiredVars {
		if _, exists := envConfig[requiredVar]; !exists {
			slog.Warn("Required variable '%s' not found in environment '%s'", requiredVar, envName)
		}
	}

	return nil
}

// ExportToFile 导出环境变量到文件
func (em *EnvironmentManager) ExportToFile(filePath string) error {
	file, err := os.Create(filePath)
	if err != nil {
		return fmt.Errorf("failed to create environment file: %w", err)
	}
	defer file.Close()

	file.WriteString("# AutoTest Environment Variables\n")             //nolint:errcheck
	fmt.Fprintf(file, "# Generated at: %s\n\n", "2006-01-02 15:04:05") //nolint:errcheck

	for key, value := range em.vars {
		// 如果值包含空格或特殊字符，添加引号
		if strings.ContainsAny(value, " \t\n\"'") {
			value = fmt.Sprintf("\"%s\"", strings.ReplaceAll(value, "\"", "\\\""))
		}

		if _, err := fmt.Fprintf(file, "%s=%s\n", key, value); err != nil {
			return fmt.Errorf("failed to write environment variable: %w", err)
		}
	}

	slog.Info("Environment variables exported to: %s", filePath)
	return nil
}
