// Package utils provides utility functions for the amapi package.
package utils

import (
	"math/rand"
	"time"

	"amapi-pkg/pkgs/amapi/types"
)

// RetryConfig contains configuration for retry behavior.
type RetryConfig struct {
	// MaxAttempts is the maximum number of retry attempts
	MaxAttempts int

	// BaseDelay is the base delay between retry attempts
	BaseDelay time.Duration

	// MaxDelay is the maximum delay between retry attempts
	MaxDelay time.Duration

	// EnableRetry indicates if retry is enabled
	EnableRetry bool

	// Jitter adds randomness to delays to prevent thundering herd
	Jitter bool
}

// RetryHandler handles retry logic for API operations.
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
func (r *RetryHandler) Execute(operation func() error) error {
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