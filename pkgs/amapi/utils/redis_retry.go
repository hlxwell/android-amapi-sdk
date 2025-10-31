// Package utils provides utility functions for the amapi package.
package utils

import (
	"context"
	"fmt"
	"math/rand"
	"time"

	"github.com/redis/go-redis/v9"

	"amapi-pkg/pkgs/amapi/types"
)

// RedisRetryHandler handles distributed retry logic using Redis to prevent concurrent retries.
//
// 使用 Redis 分布式锁防止多个进程同时重试同一操作。
// 这可以减少重复的 API 调用，特别是在高并发场景下。
//
// # 工作原理
//
// 1. 每个重试操作生成唯一的 operation ID
// 2. 尝试获取 Redis 分布式锁（使用 SETNX）
// 3. 如果获取成功，执行重试操作
// 4. 如果获取失败，等待一小段时间后检查操作是否已成功
// 5. 操作完成后释放锁
//
// # 使用示例
//
//	client := redis.NewClient(&redis.Options{
//	    Addr: "localhost:6379",
//	})
//
//	handler := NewRedisRetryHandler(client, "amapi:", utils.RetryConfig{
//	    MaxAttempts: 3,
//	    BaseDelay:   1 * time.Second,
//	    MaxDelay:    30 * time.Second,
//	    EnableRetry: true,
//	})
//	defer handler.Close()
//
//	operationID := fmt.Sprintf("operation-%d", time.Now().UnixNano())
//	err := handler.Execute(ctx, operationID, func() error {
//	    // 执行可能失败的操作
//	    return someOperation()
//	})
//
// # 重试统计
//
// handler 会在 Redis 中记录每个操作的重试次数，可以通过 GetRetryCount 查询。
type RedisRetryHandler struct {
	client    *redis.Client
	keyPrefix string
	config    RetryConfig
}

// NewRedisRetryHandler creates a new Redis-based retry handler.
func NewRedisRetryHandler(client *redis.Client, keyPrefix string, config RetryConfig) *RedisRetryHandler {
	// Set defaults
	if config.MaxAttempts <= 0 {
		config.MaxAttempts = 3
	}
	if config.BaseDelay <= 0 {
		config.BaseDelay = 1 * time.Second
	}
	if config.MaxDelay <= 0 {
		config.MaxDelay = 30 * time.Second
	}
	config.Jitter = true // Enable jitter by default

	return &RedisRetryHandler{
		client:    client,
		keyPrefix: keyPrefix,
		config:    config,
	}
}

// Execute executes an operation with retry logic using Redis to coordinate retries across processes.
func (r *RedisRetryHandler) Execute(ctx context.Context, operationID string, operation func() error) error {
	if !r.config.EnableRetry {
		return operation()
	}

	var lastErr error

	for attempt := 0; attempt < r.config.MaxAttempts; attempt++ {
		// Check if another process is already retrying this operation
		retryKey := fmt.Sprintf("%sretry:lock:%s", r.keyPrefix, operationID)

		// Try to acquire lock to prevent concurrent retries
		lockAcquired, err := r.client.SetNX(ctx, retryKey, "1", time.Minute).Result()
		if err != nil {
			// If we can't acquire lock, proceed anyway (failover to local retry)
			lockAcquired = true
		}

		if !lockAcquired {
			// Another process is handling this, wait a bit and check if it succeeded
			time.Sleep(500 * time.Millisecond)
			// Try operation once more
			err := operation()
			if err == nil {
				return nil
			}
			lastErr = err
			continue
		}

		// Execute operation
		err = operation()

		// Release lock immediately after operation
		r.client.Del(ctx, retryKey)

		if err == nil {
			return nil
		}

		lastErr = err

		// Check if error is retryable
		if apiErr, ok := err.(*types.Error); ok {
			if !apiErr.IsRetryable() {
				return err
			}
		} else {
			// For non-API errors, only retry on the first attempt
			if attempt > 0 {
				return err
			}
		}

		// Don't sleep after the last attempt
		if attempt == r.config.MaxAttempts-1 {
			break
		}

		// Calculate delay
		delay := r.calculateDelay(attempt)

		// Track retry in Redis for monitoring
		retryCountKey := fmt.Sprintf("%sretry:count:%s", r.keyPrefix, operationID)
		r.client.Incr(ctx, retryCountKey)
		r.client.Expire(ctx, retryCountKey, time.Hour)

		// Wait before next attempt
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-time.After(delay):
			// Continue to next attempt
		}
	}

	// All attempts failed
	if apiErr, ok := lastErr.(*types.Error); ok {
		return types.NewErrorWithCause(types.ErrCodeRetryExhausted,
			"retry attempts exhausted", apiErr)
	}

	return types.NewErrorWithCause(types.ErrCodeRetryExhausted,
		"retry attempts exhausted", lastErr)
}

// calculateDelay calculates the delay for the given attempt using exponential backoff.
func (r *RedisRetryHandler) calculateDelay(attempt int) time.Duration {
	// Exponential backoff: baseDelay * 2^attempt
	delay := r.config.BaseDelay * time.Duration(1<<uint(attempt))

	// Cap at maximum delay
	if delay > r.config.MaxDelay {
		delay = r.config.MaxDelay
	}

	// Add jitter to prevent thundering herd
	if r.config.Jitter {
		jitter := time.Duration(rand.Float64() * float64(delay) * 0.1)
		delay += jitter
	}

	return delay
}

// GetRetryCount returns the number of retries for a given operation ID.
func (r *RedisRetryHandler) GetRetryCount(ctx context.Context, operationID string) (int64, error) {
	retryCountKey := fmt.Sprintf("%sretry:count:%s", r.keyPrefix, operationID)
	count, err := r.client.Get(ctx, retryCountKey).Int64()
	if err == redis.Nil {
		return 0, nil
	}
	return count, err
}

// Close closes the Redis client connection.
func (r *RedisRetryHandler) Close() error {
	if r.client != nil {
		return r.client.Close()
	}
	return nil
}
