// Package utils provides utility functions and interfaces for rate limiting and retry handling.
package utils

import (
	"context"
)

// RateLimiterInterface defines the interface for rate limiters.
//
// 此接口允许使用不同的 rate limiting 实现：
//   - 本地实现：使用 golang.org/x/time/rate 包
//   - 分布式实现：使用 Redis 实现跨进程的 rate limiting
//
// 实现此接口的类型包括：
//   - RateLimiter: 本地 rate limiter（用于单进程应用）
//   - RedisRateLimiter: 分布式 rate limiter（用于多进程应用）
type RateLimiterInterface interface {
	// Wait waits until the rate limiter allows the request.
	// 如果超过 rate limit，会阻塞直到有足够的配额。
	Wait(ctx context.Context) error

	// Allow checks if a request is allowed without waiting.
	// 如果允许，返回 true；如果不允许（超过限制），返回 false。
	Allow(ctx context.Context) bool

	// SetLimit changes the rate limit (requests per minute).
	SetLimit(rateLimit int)

	// SetBurst changes the burst capacity.
	SetBurst(burst int)

	// Close closes the rate limiter and releases resources.
	// 对于本地实现，这是一个 no-op。
	// 对于 Redis 实现，会关闭 Redis 连接。
	Close() error
}

// RetryHandlerInterface defines the interface for retry handlers.
//
// 此接口允许使用不同的 retry 实现：
//   - 本地实现：在单个进程内管理重试
//   - 分布式实现：使用 Redis 防止多个进程同时重试同一操作
//
// 实现此接口的类型包括：
//   - RetryHandler: 本地 retry handler（用于单进程应用）
//   - RedisRetryHandler: 分布式 retry handler（用于多进程应用）
type RetryHandlerInterface interface {
	// Execute executes an operation with retry logic.
	//
	// operationID 是操作的唯一标识符。对于分布式 retry handler，
	// 此 ID 用于在 Redis 中创建锁，防止多个进程同时重试同一操作。
	// 对于本地 retry handler，此参数会被忽略。
	//
	// operation 是要执行的操作函数。如果操作失败且错误是可重试的，
	// handler 会根据配置进行重试（使用指数退避）。
	Execute(ctx context.Context, operationID string, operation func() error) error

	// Close closes the retry handler and releases resources.
	// 对于本地实现，这是一个 no-op。
	// 对于 Redis 实现，会关闭 Redis 连接。
	Close() error
}
