// Package utils provides utility functions for the amapi package.
package utils

import (
	"context"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
)

// RedisRateLimiter provides distributed rate limiting using Redis.
//
// 使用 Redis 的滑动窗口计数器算法实现分布式的 rate limiting。
// 所有使用同一个 Redis 实例的进程会共享同一个 rate limit。
//
// # 工作原理
//
// 1. 每个请求在 Redis sorted set 中记录一个带时间戳的条目
// 2. 定期清理超出时间窗口（1分钟）的旧条目
// 3. 统计当前时间窗口内的请求数
// 4. 如果超过限制，等待直到有足够的配额
//
// # 使用示例
//
//	client := redis.NewClient(&redis.Options{
//	    Addr: "localhost:6379",
//	})
//
//	limiter := NewRedisRateLimiter(client, "amapi:", 100, 20)
//	defer limiter.Close()
//
//	// 等待直到允许请求
//	err := limiter.Wait(ctx)
//	if err != nil {
//	    return err
//	}
//
//	// 或者检查是否允许（不等待）
//	if limiter.Allow(ctx) {
//	    // 执行请求
//	}
type RedisRateLimiter struct {
	client    *redis.Client
	keyPrefix string
	rateLimit int // requests per minute
	burst     int
	window    time.Duration // time window for rate limiting
}

// NewRedisRateLimiter creates a new Redis-based rate limiter.
// rateLimit is requests per minute, burst is the burst capacity.
func NewRedisRateLimiter(client *redis.Client, keyPrefix string, rateLimit, burst int) *RedisRateLimiter {
	return NewRedisRateLimiterWithWindow(client, keyPrefix, rateLimit, burst, 60*time.Second)
}

// NewRedisRateLimiterWithWindow creates a new Redis-based rate limiter with custom window.
// rateLimit is requests per window, burst is the burst capacity, window is the time window.
func NewRedisRateLimiterWithWindow(client *redis.Client, keyPrefix string, rateLimit, burst int, window time.Duration) *RedisRateLimiter {
	if rateLimit <= 0 {
		rateLimit = 100 // Default to 100 requests per window
	}
	if burst <= 0 {
		burst = 10 // Default burst of 10
	}
	if window <= 0 {
		window = 60 * time.Second // Default 1 minute window
	}

	rl := &RedisRateLimiter{
		client:    client,
		keyPrefix: keyPrefix,
		rateLimit: rateLimit,
		burst:     burst,
		window:    window,
	}

	return rl
}

// Wait waits until the rate limiter allows the request.
// Uses Redis sliding window counter algorithm to ensure distributed rate limiting.
func (rl *RedisRateLimiter) Wait(ctx context.Context) error {
	key := rl.keyPrefix + "ratelimit:requests"

	// Get current time in seconds
	now := time.Now().Unix()
	windowStart := now - int64(rl.window.Seconds())

	// Use Redis pipeline for atomic operations
	pipe := rl.client.Pipeline()

	// Remove old entries outside the window
	pipe.ZRemRangeByScore(ctx, key, "0", fmt.Sprintf("%d", windowStart))

	// Count current requests in the window
	countCmd := pipe.ZCard(ctx, key)

	// Add current request with score = current timestamp
	pipe.ZAdd(ctx, key, redis.Z{
		Score:  float64(now),
		Member: fmt.Sprintf("%d", now),
	})

	// Set expiry on the sorted set
	pipe.Expire(ctx, key, rl.window+10*time.Second)

	// Execute pipeline
	_, err := pipe.Exec(ctx)
	if err != nil {
		return fmt.Errorf("redis rate limit error: %w", err)
	}

	// Check if we've exceeded the limit
	currentCount := countCmd.Val()
	if currentCount >= int64(rl.rateLimit) {
		// Calculate wait time until the oldest request expires
		oldestCmd := rl.client.ZRangeWithScores(ctx, key, 0, 0)
		if oldestCmd.Err() == nil && len(oldestCmd.Val()) > 0 {
			oldestScore := int64(oldestCmd.Val()[0].Score)
			waitTime := time.Duration(oldestScore-int64(windowStart)) * time.Second
			if waitTime > 0 {
				select {
				case <-ctx.Done():
					return ctx.Err()
				case <-time.After(waitTime):
					// Retry after waiting
					return rl.Wait(ctx)
				}
			}
		}
		return fmt.Errorf("rate limit exceeded: %d requests in window", currentCount)
	}

	return nil
}

// Allow checks if a request is allowed without waiting.
// Implements RateLimiterInterface.
func (rl *RedisRateLimiter) Allow(ctx context.Context) bool {
	key := rl.keyPrefix + "ratelimit:requests"
	now := time.Now().Unix()
	windowStart := now - int64(rl.window.Seconds())

	// Remove old entries and count current requests
	pipe := rl.client.Pipeline()
	pipe.ZRemRangeByScore(ctx, key, "0", fmt.Sprintf("%d", windowStart))
	countCmd := pipe.ZCard(ctx, key)
	pipe.Exec(ctx)

	currentCount := countCmd.Val()
	if currentCount >= int64(rl.rateLimit) {
		return false
	}

	// Add current request
	rl.client.ZAdd(ctx, key, redis.Z{
		Score:  float64(now),
		Member: fmt.Sprintf("%d", now),
	})
	rl.client.Expire(ctx, key, rl.window+10*time.Second)

	return true
}

// SetLimit changes the rate limit.
func (rl *RedisRateLimiter) SetLimit(rateLimit int) {
	rl.rateLimit = rateLimit
}

// SetBurst changes the burst capacity (not used in Redis implementation, but kept for compatibility).
func (rl *RedisRateLimiter) SetBurst(burst int) {
	rl.burst = burst
}

// Close closes the Redis client connection.
func (rl *RedisRateLimiter) Close() error {
	if rl.client != nil {
		return rl.client.Close()
	}
	return nil
}
