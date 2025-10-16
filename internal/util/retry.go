package util

import (
	"fmt"
	"time"

	"github.com/go-resty/resty/v2"
	"github.com/vearne/autotest/internal/config"
	slog "github.com/vearne/simplelog"
)

// RetryableHTTPClient 支持重试的HTTP客户端
type RetryableHTTPClient struct {
	client      *resty.Client
	retryConfig config.AutoTestConfig
}

// NewRetryableHTTPClient 创建支持重试的HTTP客户端
func NewRetryableHTTPClient(client *resty.Client, cfg config.AutoTestConfig) *RetryableHTTPClient {
	return &RetryableHTTPClient{
		client:      client,
		retryConfig: cfg,
	}
}

// ExecuteWithRetry 执行HTTP请求并支持重试
func (r *RetryableHTTPClient) ExecuteWithRetry(req *resty.Request) (*resty.Response, error) {
	maxAttempts := r.retryConfig.Global.Retry.MaxAttempts
	if maxAttempts <= 0 {
		maxAttempts = 1 // 默认至少执行一次
	}

	retryDelay := r.retryConfig.Global.Retry.RetryDelay
	if retryDelay <= 0 {
		retryDelay = time.Second // 默认1秒重试间隔
	}

	retryStatusCodes := make(map[int]bool)
	for _, code := range r.retryConfig.Global.Retry.RetryOnStatusCodes {
		retryStatusCodes[code] = true
	}

	var lastErr error
	var resp *resty.Response

	for attempt := 1; attempt <= maxAttempts; attempt++ {
		slog.Debug("HTTP request attempt %d/%d", attempt, maxAttempts)

		resp, lastErr = req.Send()

		if lastErr == nil {
			// 检查状态码是否需要重试
			if len(retryStatusCodes) == 0 || !retryStatusCodes[resp.StatusCode()] {
				slog.Debug("Request succeeded on attempt %d", attempt)
				return resp, nil
			}
			slog.Info("Status code %d requires retry, attempt %d/%d", resp.StatusCode(), attempt, maxAttempts)
		} else {
			slog.Error("Request failed on attempt %d/%d: %v", attempt, maxAttempts, lastErr)
		}

		// 如果不是最后一次尝试，等待重试间隔
		if attempt < maxAttempts {
			slog.Debug("Waiting %v before retry", retryDelay)
			time.Sleep(retryDelay)
		}
	}

	if lastErr != nil {
		return nil, fmt.Errorf("request failed after %d attempts: %w", maxAttempts, lastErr)
	}

	return resp, fmt.Errorf("request failed after %d attempts with status code %d", maxAttempts, resp.StatusCode())
}

// IsRetryableStatusCode 检查状态码是否需要重试
func IsRetryableStatusCode(statusCode int, retryCodes []int) bool {
	if len(retryCodes) == 0 {
		// 默认重试的状态码
		defaultRetryCodes := []int{500, 502, 503, 504, 408, 429}
		for _, code := range defaultRetryCodes {
			if statusCode == code {
				return true
			}
		}
		return false
	}

	for _, code := range retryCodes {
		if statusCode == code {
			return true
		}
	}
	return false
}
