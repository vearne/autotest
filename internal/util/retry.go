package util

import (
	"context"
	"fmt"
	"time"

	"github.com/go-resty/resty/v2"
	"github.com/vearne/autotest/internal/config"
	slog "github.com/vearne/simplelog"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
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

// ExecuteHttpWithRetry 执行HTTP请求并支持重试
func ExecuteHttpWithRetry(ctx context.Context, cfg config.AutoTestConfig, operation func() error) error {
	maxAttempts := cfg.Global.Retry.MaxAttempts
	if maxAttempts <= 0 {
		maxAttempts = 3 // 默认最多重试3次
	}

	retryDelay := cfg.Global.Retry.RetryDelay
	if retryDelay <= 0 {
		retryDelay = time.Second // 默认重试间隔1秒
	}

	retryStatusCodes := make(map[int]bool)
	if len(cfg.Global.Retry.RetryOnStatusCodes) == 0 {
		// 默认重试的状态码
		defaultRetryCodes := []int{500, 502, 503, 504, 408, 429}
		for _, code := range defaultRetryCodes {
			retryStatusCodes[code] = true
		}
	} else {
		for _, code := range cfg.Global.Retry.RetryOnStatusCodes {
			retryStatusCodes[code] = true
		}
	}

	var lastErr error

	for attempt := 1; attempt <= maxAttempts; attempt++ {
		slog.Debug("HTTP request attempt %d/%d", attempt, maxAttempts)

		lastErr = operation()

		if lastErr == nil {
			slog.Debug("HTTP request succeeded on attempt %d", attempt)
			return nil
		}

		// TODO: 这里需要检查具体的 HTTP 状态码，但当前架构下比较困难
		// 先简单地重试所有错误，后续可以优化

		slog.Info("HTTP request failed (retryable), attempt %d/%d: %v", attempt, maxAttempts, lastErr)

		// 如果不是最后一次尝试，等待重试间隔
		if attempt < maxAttempts {
			slog.Debug("Waiting %v before HTTP retry", retryDelay)
			select {
			case <-time.After(retryDelay):
				// 继续重试
			case <-ctx.Done():
				return ctx.Err()
			}
		}
	}

	return fmt.Errorf("HTTP request failed after %d attempts: %w", maxAttempts, lastErr)
}

// ExecuteGrpcWithRetry 执行gRPC请求并支持重试
func ExecuteGrpcWithRetry(ctx context.Context, cfg config.AutoTestConfig, operation func() error) error {
	maxAttempts := cfg.Global.Retry.MaxAttempts
	if maxAttempts <= 0 {
		maxAttempts = 3 // 默认最多重试3次
	}

	retryDelay := cfg.Global.Retry.RetryDelay
	if retryDelay <= 0 {
		retryDelay = time.Second // 默认重试间隔1秒
	}

	var lastErr error

	for attempt := 1; attempt <= maxAttempts; attempt++ {
		slog.Debug("gRPC request attempt %d/%d", attempt, maxAttempts)

		lastErr = operation()

		if lastErr == nil {
			slog.Debug("gRPC request succeeded on attempt %d", attempt)
			return nil
		}

		// 检查是否为可重试的 gRPC 错误
		if !isRetryableGrpcError(lastErr) {
			slog.Info("gRPC error not retryable: %v", lastErr)
			return lastErr
		}

		slog.Info("gRPC request failed (retryable), attempt %d/%d: %v", attempt, maxAttempts, lastErr)

		// 如果不是最后一次尝试，等待重试间隔
		if attempt < maxAttempts {
			slog.Debug("Waiting %v before gRPC retry", retryDelay)
			select {
			case <-time.After(retryDelay):
				// 继续重试
			case <-ctx.Done():
				return ctx.Err()
			}
		}
	}

	return fmt.Errorf("gRPC request failed after %d attempts: %w", maxAttempts, lastErr)
}

// isRetryableGrpcError 检查gRPC错误是否可重试
func isRetryableGrpcError(err error) bool {
	if err == nil {
		return false
	}

	// 检查是否为gRPC状态错误
	if st, ok := status.FromError(err); ok {
		switch st.Code() {
		case codes.Unavailable, // 服务不可用
			codes.DeadlineExceeded,  // 超时
			codes.ResourceExhausted, // 资源耗尽（如速率限制）
			codes.Aborted,           // 操作被中止
			codes.Internal:          // 内部错误
			return true
		default:
			return false
		}
	}

	// 对于网络等其他错误，也尝试重试
	return true
}
