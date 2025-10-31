// Package utils provides utility functions for the amapi package.
package utils

import (
	"context"
	"math/rand"
	"time"

	"amapi-pkg/pkgs/amapi/types"
)

// RetryConfig contains configuration for retry behavior.
//
// 配置项说明：
//   - MaxAttempts: 最大重试次数（包括首次尝试）
//   - BaseDelay: 基础延迟时间（使用指数退避）
//   - MaxDelay: 最大延迟时间（延迟不会超过此值）
//   - EnableRetry: 是否启用重试
//   - Jitter: 是否添加随机抖动（防止惊群效应）
type RetryConfig struct {
	// MaxAttempts is the maximum number of retry attempts (including the first attempt).
	// Default: 3
	MaxAttempts int

	// BaseDelay is the base delay between retry attempts.
	// Actual delay = baseDelay * 2^attempt (with jitter).
	// Default: 1 second
	BaseDelay time.Duration

	// MaxDelay is the maximum delay between retry attempts.
	// The calculated delay will be capped at this value.
	// Default: 30 seconds
	MaxDelay time.Duration

	// EnableRetry indicates if retry is enabled.
	// If false, operations will only be attempted once.
	// Default: true
	EnableRetry bool

	// Jitter adds randomness to delays to prevent thundering herd.
	// Adds up to 10% random jitter to the delay.
	// Default: true
	Jitter bool
}

// RetryHandler handles retry logic for API operations.
//
// 本地 retry handler，适用于单进程应用。
// 使用指数退避算法，带有随机抖动以防止惊群效应。
//
// # 使用示例
//
//	handler := NewRetryHandler(utils.RetryConfig{
//	    MaxAttempts: 3,
//	    BaseDelay:   1 * time.Second,
//	    MaxDelay:    30 * time.Second,
//	    EnableRetry: true,
//	})
//	defer handler.Close()
//
//	operationID := "unique-operation-id"
//	err := handler.Execute(ctx, operationID, func() error {
//	    // 执行可能失败的操作
//	    return someAPI.Call()
//	})
//
// # 分布式场景
//
// 如果您的应用程序运行多个进程，请使用 RedisRetryHandler 代替，
// 以防止多个进程同时重试同一操作。
type RetryHandler struct {
	config RetryConfig
}

// NewRetryHandler creates a new retry handler with the given configuration.
func NewRetryHandler(config RetryConfig) *RetryHandler {
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

	return &RetryHandler{
		config: config,
	}
}

// Execute executes an operation with retry logic.
// For compatibility with interface, operationID is ignored for local retry handler.
func (r *RetryHandler) Execute(ctx context.Context, operationID string, operation func() error) error {
	if !r.config.EnableRetry {
		return operation()
	}

	var lastErr error

	for attempt := 0; attempt < r.config.MaxAttempts; attempt++ {
		err := operation()
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
		time.Sleep(delay)
	}

	// All attempts failed
	if apiErr, ok := lastErr.(*types.Error); ok {
		return types.NewErrorWithCause(types.ErrCodeRetryExhausted,
			"retry attempts exhausted", apiErr)
	}

	return types.NewErrorWithCause(types.ErrCodeRetryExhausted,
		"retry attempts exhausted", lastErr)
}

// Close closes the retry handler (no-op for local handler).
func (r *RetryHandler) Close() error {
	return nil
}

// calculateDelay calculates the delay for the given attempt using exponential backoff.
func (r *RetryHandler) calculateDelay(attempt int) time.Duration {
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

// IsRetryableError checks if an error is retryable.
func IsRetryableError(err error) bool {
	if err == nil {
		return false
	}

	if apiErr, ok := err.(*types.Error); ok {
		return apiErr.IsRetryable()
	}

	return false
}

// GetRetryDelay calculates the retry delay for an error.
func GetRetryDelay(err error, attempt int, baseDelay time.Duration) time.Duration {
	if apiErr, ok := err.(*types.Error); ok {
		return apiErr.RetryDelay(attempt, baseDelay)
	}

	// Default exponential backoff
	delay := baseDelay * time.Duration(1<<uint(attempt))
	if delay > 30*time.Second {
		delay = 30 * time.Second
	}

	return delay
}