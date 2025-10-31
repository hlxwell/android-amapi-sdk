package utils

import (
	"context"

	"golang.org/x/time/rate"
)

// RateLimiter provides rate limiting functionality.
type RateLimiter struct {
	limiter *rate.Limiter
}

// NewRateLimiter creates a new rate limiter.
// rateLimit is requests per minute, burst is the burst capacity.
func NewRateLimiter(rateLimit, burst int) *RateLimiter {
	if rateLimit <= 0 {
		rateLimit = 100 // Default to 100 requests per minute
	}
	if burst <= 0 {
		burst = 10 // Default burst of 10
	}

	// Convert rate limit from per-minute to per-second
	r := rate.Limit(float64(rateLimit) / 60.0)

	return &RateLimiter{
		limiter: rate.NewLimiter(r, burst),
	}
}

// Wait waits until the rate limiter allows the request.
func (rl *RateLimiter) Wait(ctx context.Context) error {
	return rl.limiter.Wait(ctx)
}

// Close closes the rate limiter (no-op for local limiter).
func (rl *RateLimiter) Close() error {
	return nil
}

// Allow checks if a request is allowed without waiting.
// For compatibility with interface, accepts context but ignores it for local limiter.
func (rl *RateLimiter) Allow(ctx context.Context) bool {
	return rl.limiter.Allow()
}

// Reserve reserves a request and returns a reservation.
func (rl *RateLimiter) Reserve() *rate.Reservation {
	return rl.limiter.Reserve()
}

// SetLimit changes the rate limit.
func (rl *RateLimiter) SetLimit(rateLimit int) {
	r := rate.Limit(float64(rateLimit) / 60.0)
	rl.limiter.SetLimit(r)
}

// SetBurst changes the burst capacity.
func (rl *RateLimiter) SetBurst(burst int) {
	rl.limiter.SetBurst(burst)
}