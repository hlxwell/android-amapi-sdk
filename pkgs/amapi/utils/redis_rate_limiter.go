// Package utils provides utility functions for the amapi package.
package utils

import (
	"context"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
)

// RedisRateLimiter provides distributed rate limiting using Redis.
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
	if rateLimit <= 0 {
		rateLimit = 100 // Default to 100 requests per minute
	}
	if burst <= 0 {
		burst = 10 // Default burst of 10
	}

	rl := &RedisRateLimiter{
		client:    client,
		keyPrefix: keyPrefix,
		rateLimit: rateLimit,
		burst:     burst,
		window:    60 * time.Second, // 1 minute window
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
