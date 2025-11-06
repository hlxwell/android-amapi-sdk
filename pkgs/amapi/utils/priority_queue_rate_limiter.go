// Package utils provides utility functions for the amapi package.
package utils

import (
	"context"
	"fmt"
	"time"

	"amapi-pkg/pkgs/amapi/types"
)

// PriorityQueueRateLimiter implements rate limiting using a priority queue.
//
// 使用优先级队列实现 rate limiting。
// 任务会被放入队列，由 TaskWorker 按优先级执行。
// 执行时会在 worker 中应用实际的 rate limiting。
type PriorityQueueRateLimiter struct {
	queue  *RedisPriorityQueue
	worker *TaskWorker
	config PriorityQueueRateLimiterConfig
}

// PriorityQueueRateLimiterConfig contains configuration for the priority queue rate limiter.
type PriorityQueueRateLimiterConfig struct {
	// DefaultPriority is the default priority for tasks (0-1000).
	DefaultPriority int

	// MaxQueueSize is the maximum queue size.
	MaxQueueSize int64

	// Timeout is the timeout for waiting for task completion.
	Timeout time.Duration
}

// DefaultPriorityQueueRateLimiterConfig returns default configuration.
func DefaultPriorityQueueRateLimiterConfig() PriorityQueueRateLimiterConfig {
	return PriorityQueueRateLimiterConfig{
		DefaultPriority: 500,
		MaxQueueSize:    10000,
		Timeout:          5 * time.Minute,
	}
}

// NewPriorityQueueRateLimiter creates a new priority queue rate limiter.
func NewPriorityQueueRateLimiter(queue *RedisPriorityQueue, worker *TaskWorker, config PriorityQueueRateLimiterConfig) *PriorityQueueRateLimiter {
	if config.DefaultPriority < 0 {
		config.DefaultPriority = 0
	} else if config.DefaultPriority > 1000 {
		config.DefaultPriority = 1000
	}
	if config.Timeout <= 0 {
		config.Timeout = 5 * time.Minute
	}

	return &PriorityQueueRateLimiter{
		queue:  queue,
		worker: worker,
		config: config,
	}
}

// Wait waits until the rate limiter allows the request.
//
// 在优先级队列模式下，rate limiting 由 TaskWorker 在执行任务时应用。
// 这个方法只是检查队列是否已满，实际的 rate limiting 在 worker 中处理。
func (pqrl *PriorityQueueRateLimiter) Wait(ctx context.Context) error {
	// Check queue size
	size, err := pqrl.queue.Size(ctx)
	if err != nil {
		return fmt.Errorf("failed to check queue size: %w", err)
	}

	if size >= pqrl.config.MaxQueueSize {
		return types.NewError(types.ErrCodeTooManyRequests, "priority queue is full")
	}

	// In priority queue mode, rate limiting is handled by the worker
	// We just check if the queue has space
	return nil
}

// Allow checks if a request is allowed without waiting.
//
// 检查队列大小和 rate limit 状态。
// 注意：这不会实际执行任务，只是检查是否可以入队。
func (pqrl *PriorityQueueRateLimiter) Allow(ctx context.Context) bool {
	// Check queue size
	size, err := pqrl.queue.Size(ctx)
	if err != nil {
		return false
	}

	if size >= pqrl.config.MaxQueueSize {
		return false
	}

	// TODO: Could check rate limiter status here
	// For now, just check queue size
	return true
}

// SetLimit changes the rate limit (not directly applicable to priority queue).
func (pqrl *PriorityQueueRateLimiter) SetLimit(rateLimit int) {
	// Rate limit is controlled by the worker's rate limiter
	// This is a no-op for compatibility
}

// SetBurst changes the burst capacity (not directly applicable to priority queue).
func (pqrl *PriorityQueueRateLimiter) SetBurst(burst int) {
	// Burst is controlled by the worker's rate limiter
	// This is a no-op for compatibility
}

// Close closes the priority queue rate limiter.
func (pqrl *PriorityQueueRateLimiter) Close() error {
	// Don't close queue or worker as they may be shared
	return nil
}

