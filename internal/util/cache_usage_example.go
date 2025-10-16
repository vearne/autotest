package util

import (
	"crypto/md5"
	"fmt"
	"strings"

	"github.com/fullstorydev/grpcurl"
	"github.com/go-resty/resty/v2"
	slog "github.com/vearne/simplelog"
	lua "github.com/yuin/gopher-lua"
)

// 这个文件展示了如何在autotest项目中使用缓存

// CachedGrpcDescriptor 缓存gRPC描述符
func CachedGrpcDescriptor(cache *Cache, address string, getDescriptor func() (grpcurl.DescriptorSource, error)) (grpcurl.DescriptorSource, error) {
	// 尝试从缓存获取
	if cached, found := cache.Get(address); found {
		if desc, ok := cached.(grpcurl.DescriptorSource); ok {
			slog.Debug("Using cached gRPC descriptor for %s", address)
			return desc, nil
		}
	}

	// 缓存未命中，获取新的描述符
	slog.Debug("Fetching new gRPC descriptor for %s", address)
	descriptor, err := getDescriptor()
	if err != nil {
		return nil, err
	}

	// 存入缓存
	cache.Set(address, descriptor)
	return descriptor, nil
}

// CachedHttpResponse 缓存HTTP响应（仅适用于GET请求）
func CachedHttpResponse(cache *Cache, method, url string, headers []string, executeRequest func() (*resty.Response, error)) (*resty.Response, error) {
	// 只缓存GET请求
	if strings.ToUpper(method) != "GET" {
		return executeRequest()
	}

	// 生成缓存键
	cacheKey := generateHttpCacheKey(method, url, headers)

	// 尝试从缓存获取
	if cached, found := cache.Get(cacheKey); found {
		if resp, ok := cached.(*resty.Response); ok {
			slog.Debug("Using cached HTTP response for %s %s", method, url)
			return resp, nil
		}
	}

	// 缓存未命中，执行请求
	slog.Debug("Executing new HTTP request for %s %s", method, url)
	response, err := executeRequest()
	if err != nil {
		return nil, err
	}

	// 只缓存成功的响应
	if response.StatusCode() >= 200 && response.StatusCode() < 300 {
		cache.Set(cacheKey, response)
	}

	return response, nil
}

// CachedLuaScript 缓存编译后的Lua脚本
func CachedLuaScript(cache *Cache, scriptContent string, compileScript func() (*lua.LFunction, error)) (*lua.LFunction, error) {
	// 生成脚本的哈希作为缓存键
	cacheKey := generateScriptHash(scriptContent)

	// 尝试从缓存获取
	if cached, found := cache.Get(cacheKey); found {
		if fn, ok := cached.(*lua.LFunction); ok {
			slog.Debug("Using cached Lua script")
			return fn, nil
		}
	}

	// 缓存未命中，编译脚本
	slog.Debug("Compiling new Lua script")
	compiledScript, err := compileScript()
	if err != nil {
		return nil, err
	}

	// 存入缓存
	cache.Set(cacheKey, compiledScript)
	return compiledScript, nil
}

// CachedTemplate 缓存模板渲染结果
func CachedTemplate(cache *Cache, templateStr string, variables map[string]string, renderTemplate func() (string, error)) (string, error) {
	// 生成缓存键（模板内容 + 变量）
	cacheKey := generateTemplateCacheKey(templateStr, variables)

	// 尝试从缓存获取
	if cached, found := cache.Get(cacheKey); found {
		if result, ok := cached.(string); ok {
			slog.Debug("Using cached template result")
			return result, nil
		}
	}

	// 缓存未命中，渲染模板
	slog.Debug("Rendering new template")
	result, err := renderTemplate()
	if err != nil {
		return "", err
	}

	// 存入缓存
	cache.Set(cacheKey, result)
	return result, nil
}

// generateHttpCacheKey 生成HTTP请求的缓存键
func generateHttpCacheKey(method, url string, headers []string) string {
	key := fmt.Sprintf("%s:%s", method, url)

	// 包含重要的头信息
	for _, header := range headers {
		if strings.Contains(strings.ToLower(header), "authorization") ||
			strings.Contains(strings.ToLower(header), "content-type") {
			key += ":" + header
		}
	}

	return fmt.Sprintf("http:%x", md5.Sum([]byte(key)))
}

// generateScriptHash 生成脚本内容的哈希
func generateScriptHash(script string) string {
	return fmt.Sprintf("lua:%x", md5.Sum([]byte(script)))
}

// generateTemplateCacheKey 生成模板缓存键
func generateTemplateCacheKey(template string, variables map[string]string) string {
	key := template

	// 按键排序添加变量，确保缓存键的一致性
	for k, v := range variables {
		key += fmt.Sprintf(":%s=%s", k, v)
	}

	return fmt.Sprintf("template:%x", md5.Sum([]byte(key)))
}

// CacheStats 缓存统计信息
type CacheStats struct {
	Name    string  `json:"name"`
	Hits    int64   `json:"hits"`
	Misses  int64   `json:"misses"`
	HitRate float64 `json:"hit_rate"`
	Size    int     `json:"size"`
	Enabled bool    `json:"enabled"`
}

// GetCacheReport 获取缓存使用报告
func GetCacheReport(cacheManager *CacheManager) []CacheStats {
	stats := cacheManager.GetStats()

	var report []CacheStats
	for name, stat := range stats {
		report = append(report, CacheStats{
			Name:    name,
			Hits:    stat["hits"].(int64),
			Misses:  stat["misses"].(int64),
			HitRate: stat["hit_rate"].(float64),
			Size:    stat["size"].(int),
			Enabled: true,
		})
	}

	return report
}

// WarmupCache 预热缓存
func WarmupCache(cacheManager *CacheManager) {
	slog.Info("Starting cache warmup...")

	// 这里可以预加载一些常用的数据
	// 例如：常用的gRPC服务描述符、模板等

	// 示例：预热一些常用模板
	commonTemplates := []string{
		"http://{{ HOST }}/api/books",
		"http://{{ HOST }}/api/users/{{ USER_ID }}",
		"{{ GRPC_SERVER }}",
	}

	commonVars := map[string]string{
		"HOST":        "localhost:8080",
		"GRPC_SERVER": "localhost:50051",
		"USER_ID":     "123",
	}

	for _, template := range commonTemplates {
		cacheKey := generateTemplateCacheKey(template, commonVars)
		// 这里可以预渲染模板并存入缓存
		cacheManager.TemplateCache.Set(cacheKey, template) // 简化示例
	}

	slog.Info("Cache warmup completed")
}
