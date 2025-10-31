package utils

import (
	"context"
)

// RateLimiterInterface defines the interface for rate limiters.
type RateLimiterInterface interface {
	// Wait waits until the rate limiter allows the request.
	Wait(ctx context.Context) error

	// Allow checks if a request is allowed without waiting.
	Allow(ctx context.Context) bool

	// SetLimit changes the rate limit.
	SetLimit(rateLimit int)

	// SetBurst changes the burst capacity.
	SetBurst(burst int)

	// Close closes the rate limiter and releases resources.
	Close() error
}

// RetryHandlerInterface defines the interface for retry handlers.
type RetryHandlerInterface interface {
	// Execute executes an operation with retry logic.
	// For Redis retry handler, operationID should be a unique identifier for the operation.
	Execute(ctx context.Context, operationID string, operation func() error) error

	// Close closes the retry handler and releases resources.
	Close() error
}
