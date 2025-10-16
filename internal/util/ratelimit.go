package util

import (
	"context"
	"sync"
	"time"

	"golang.org/x/sync/semaphore"
)

// RateLimiter 限流器
type RateLimiter struct {
	// 并发控制
	semaphore *semaphore.Weighted
	// 速率控制
	ticker *time.Ticker
	tokens chan struct{}
	stopCh chan struct{}
	once   sync.Once
}

// NewRateLimiter 创建限流器
func NewRateLimiter(maxConcurrent int, ratePerSecond int) *RateLimiter {
	rl := &RateLimiter{
		semaphore: semaphore.NewWeighted(int64(maxConcurrent)),
		stopCh:    make(chan struct{}),
	}

	if ratePerSecond > 0 {
		rl.tokens = make(chan struct{}, ratePerSecond)
		rl.ticker = time.NewTicker(time.Second / time.Duration(ratePerSecond))

		// 启动令牌生成器
		go rl.tokenGenerator()

		// 初始填充令牌桶
		for i := 0; i < ratePerSecond && i < cap(rl.tokens); i++ {
			select {
			case rl.tokens <- struct{}{}:
			default:
			}
		}
	}

	return rl
}

// tokenGenerator 令牌生成器
func (rl *RateLimiter) tokenGenerator() {
	for {
		select {
		case <-rl.ticker.C:
			select {
			case rl.tokens <- struct{}{}:
			default:
				// 令牌桶已满，丢弃令牌
			}
		case <-rl.stopCh:
			return
		}
	}
}

// Acquire 获取执行权限
func (rl *RateLimiter) Acquire(ctx context.Context) error {
	// 1. 获取并发控制权限
	if err := rl.semaphore.Acquire(ctx, 1); err != nil {
		return err
	}

	// 2. 获取速率控制令牌
	if rl.tokens != nil {
		select {
		case <-rl.tokens:
			// 获得令牌
		case <-ctx.Done():
			rl.semaphore.Release(1)
			return ctx.Err()
		}
	}

	return nil
}

// Release 释放执行权限
func (rl *RateLimiter) Release() {
	rl.semaphore.Release(1)
}

// Stop 停止限流器
func (rl *RateLimiter) Stop() {
	rl.once.Do(func() {
		close(rl.stopCh)
		if rl.ticker != nil {
			rl.ticker.Stop()
		}
	})
}

// ExecuteWithLimit 在限流控制下执行函数
func (rl *RateLimiter) ExecuteWithLimit(ctx context.Context, fn func() error) error {
	if err := rl.Acquire(ctx); err != nil {
		return err
	}
	defer rl.Release()

	return fn()
}
