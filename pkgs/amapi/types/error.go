package types

import (
	"fmt"
	"net/http"
	"time"
)

// Error represents an API error with additional context.
type Error struct {
	// Code is the error code (HTTP status code or custom error code)
	Code int `json:"code"`

	// Message is the human-readable error message
	Message string `json:"message"`

	// Details provides additional error information
	Details string `json:"details,omitempty"`

	// Retryable indicates if the operation can be retried
	Retryable bool `json:"retryable"`

	// Timestamp when the error occurred
	Timestamp time.Time `json:"timestamp"`

	// RequestID for tracking purposes
	RequestID string `json:"request_id,omitempty"`

	// Underlying error (not serialized)
	Cause error `json:"-"`
}

// Error implements the error interface.
func (e *Error) Error() string {
	if e.Details != "" {
		return fmt.Sprintf("[%d] %s: %s", e.Code, e.Message, e.Details)
	}
	return fmt.Sprintf("[%d] %s", e.Code, e.Message)
}

// Unwrap returns the underlying error.
func (e *Error) Unwrap() error {
	return e.Cause
}

// IsRetryable returns true if the error is retryable.
func (e *Error) IsRetryable() bool {
	return e.Retryable
}

// Predefined error codes
const (
	// Client errors (4xx)
	ErrCodeBadRequest          = http.StatusBadRequest          // 400
	ErrCodeUnauthorized        = http.StatusUnauthorized        // 401
	ErrCodeForbidden           = http.StatusForbidden           // 403
	ErrCodeNotFound            = http.StatusNotFound            // 404
	ErrCodeConflict            = http.StatusConflict            // 409
	ErrCodePreconditionFailed  = http.StatusPreconditionFailed  // 412
	ErrCodeTooManyRequests     = http.StatusTooManyRequests     // 429

	// Server errors (5xx)
	ErrCodeInternalServerError = http.StatusInternalServerError // 500
	ErrCodeBadGateway          = http.StatusBadGateway          // 502
	ErrCodeServiceUnavailable  = http.StatusServiceUnavailable  // 503
	ErrCodeGatewayTimeout      = http.StatusGatewayTimeout      // 504

	// Custom error codes (6xx)
	ErrCodeConfiguration       = 600 // Configuration error
	ErrCodeAuthentication      = 601 // Authentication setup error
	ErrCodeInvalidResponse     = 602 // Invalid API response
	ErrCodeTimeout             = 603 // Operation timeout
	ErrCodeRetryExhausted      = 604 // Retry attempts exhausted
	ErrCodeInvalidInput        = 605 // Invalid input parameters
	ErrCodeResourceNotReady    = 606 // Resource not ready for operation
)

// Common error creators

// NewError creates a new Error with the given code and message.
func NewError(code int, message string) *Error {
	return &Error{
		Code:      code,
		Message:   message,
		Retryable: isRetryableCode(code),
		Timestamp: time.Now(),
	}
}

// NewErrorWithDetails creates a new Error with additional details.
func NewErrorWithDetails(code int, message, details string) *Error {
	return &Error{
		Code:      code,
		Message:   message,
		Details:   details,
		Retryable: isRetryableCode(code),
		Timestamp: time.Now(),
	}
}

// NewErrorWithCause creates a new Error wrapping an underlying error.
func NewErrorWithCause(code int, message string, cause error) *Error {
	return &Error{
		Code:      code,
		Message:   message,
		Cause:     cause,
		Retryable: isRetryableCode(code),
		Timestamp: time.Now(),
	}
}

// WrapError wraps an existing error with additional context.
func WrapError(err error, code int, message string) *Error {
	return &Error{
		Code:      code,
		Message:   message,
		Cause:     err,
		Retryable: isRetryableCode(code),
		Timestamp: time.Now(),
	}
}

// isRetryableCode determines if an error code represents a retryable error.
func isRetryableCode(code int) bool {
	switch code {
	case ErrCodeTooManyRequests,
		ErrCodeInternalServerError,
		ErrCodeBadGateway,
		ErrCodeServiceUnavailable,
		ErrCodeGatewayTimeout,
		ErrCodeTimeout:
		return true
	default:
		return false
	}
}

// Predefined errors for common scenarios

var (
	// Configuration errors
	ErrMissingProjectID = NewError(ErrCodeConfiguration, "project ID is required")
	ErrMissingCredentials = NewError(ErrCodeConfiguration, "credentials are required")
	ErrInvalidCredentials = NewError(ErrCodeAuthentication, "invalid credentials provided")

	// Input validation errors
	ErrInvalidEnterpriseID = NewError(ErrCodeInvalidInput, "invalid enterprise ID")
	ErrInvalidDeviceID     = NewError(ErrCodeInvalidInput, "invalid device ID")
	ErrInvalidPolicyID     = NewError(ErrCodeInvalidInput, "invalid policy ID")
	ErrInvalidTokenID      = NewError(ErrCodeInvalidInput, "invalid enrollment token ID")

	// Resource errors
	ErrEnterpriseNotFound = NewError(ErrCodeNotFound, "enterprise not found")
	ErrDeviceNotFound     = NewError(ErrCodeNotFound, "device not found")
	ErrPolicyNotFound     = NewError(ErrCodeNotFound, "policy not found")
	ErrTokenNotFound      = NewError(ErrCodeNotFound, "enrollment token not found")

	// Operation errors
	ErrOperationTimeout     = NewError(ErrCodeTimeout, "operation timed out")
	ErrRetryExhausted      = NewError(ErrCodeRetryExhausted, "retry attempts exhausted")
	ErrResourceNotReady    = NewError(ErrCodeResourceNotReady, "resource is not ready for this operation")

	// API errors
	ErrRateLimitExceeded   = NewError(ErrCodeTooManyRequests, "rate limit exceeded")
	ErrInvalidResponse     = NewError(ErrCodeInvalidResponse, "received invalid response from API")
	ErrServiceUnavailable  = NewError(ErrCodeServiceUnavailable, "service temporarily unavailable")
)

// ErrorType represents different categories of errors.
type ErrorType string

const (
	ErrorTypeClient        ErrorType = "client_error"
	ErrorTypeServer        ErrorType = "server_error"
	ErrorTypeConfiguration ErrorType = "configuration_error"
	ErrorTypeValidation    ErrorType = "validation_error"
	ErrorTypeNetwork       ErrorType = "network_error"
	ErrorTypeTimeout       ErrorType = "timeout_error"
	ErrorTypeAuth          ErrorType = "authentication_error"
)

// GetErrorType returns the category of the error based on its code.
func (e *Error) GetErrorType() ErrorType {
	switch {
	case e.Code >= 400 && e.Code < 500:
		if e.Code == ErrCodeUnauthorized || e.Code == ErrCodeForbidden {
			return ErrorTypeAuth
		}
		return ErrorTypeClient
	case e.Code >= 500 && e.Code < 600:
		return ErrorTypeServer
	case e.Code == ErrCodeConfiguration || e.Code == ErrCodeAuthentication:
		return ErrorTypeConfiguration
	case e.Code == ErrCodeInvalidInput:
		return ErrorTypeValidation
	case e.Code == ErrCodeTimeout:
		return ErrorTypeTimeout
	default:
		return ErrorTypeClient
	}
}

// ShouldRetry determines if an operation should be retried based on the error.
func (e *Error) ShouldRetry(attempt int, maxAttempts int) bool {
	if attempt >= maxAttempts {
		return false
	}

	return e.Retryable
}

// RetryDelay calculates the delay before the next retry attempt using exponential backoff.
func (e *Error) RetryDelay(attempt int, baseDelay time.Duration) time.Duration {
	if !e.Retryable {
		return 0
	}

	// Exponential backoff with jitter
	delay := baseDelay * time.Duration(1<<uint(attempt))
	if delay > 30*time.Second {
		delay = 30 * time.Second
	}

	return delay
}