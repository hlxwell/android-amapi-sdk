// Package client provides a high-level client for the Android Management API.
package client

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"time"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/androidmanagement/v1"
	"google.golang.org/api/googleapi"
	"google.golang.org/api/option"

	"amapi-pkg/pkgs/amapi/config"
	"amapi-pkg/pkgs/amapi/types"
	"amapi-pkg/pkgs/amapi/utils"
)

// Client represents the Android Management API client.
type Client struct {
	// service is the underlying Android Management API service
	service *androidmanagement.Service

	// config contains the client configuration
	config *config.Config

	// ctx is the context for API calls
	ctx context.Context

	// httpClient is the HTTP client used for requests
	httpClient *http.Client

	// retryHandler handles retry logic
	retryHandler *utils.RetryHandler

	// rateLimiter handles rate limiting
	rateLimiter *utils.RateLimiter

	// info contains client information
	info *types.ClientInfo
}

// New creates a new Android Management API client.
func New(cfg *config.Config) (*Client, error) {
	if cfg == nil {
		return nil, types.NewError(types.ErrCodeConfiguration, "configuration is required")
	}

	if err := cfg.Validate(); err != nil {
		return nil, types.WrapError(err, types.ErrCodeConfiguration, "invalid configuration")
	}

	ctx := context.Background()

	// Create HTTP client with authentication
	httpClient, err := createHTTPClient(ctx, cfg)
	if err != nil {
		return nil, types.WrapError(err, types.ErrCodeAuthentication, "failed to create HTTP client")
	}

	// Create Android Management API service
	service, err := androidmanagement.NewService(ctx, option.WithHTTPClient(httpClient))
	if err != nil {
		return nil, types.WrapError(err, types.ErrCodeConfiguration, "failed to create Android Management service")
	}

	// Create retry handler
	retryHandler := utils.NewRetryHandler(utils.RetryConfig{
		MaxAttempts:  cfg.RetryAttempts,
		BaseDelay:    cfg.RetryDelay,
		MaxDelay:     30 * time.Second,
		EnableRetry:  cfg.EnableRetry,
	})

	// Create rate limiter
	rateLimiter := utils.NewRateLimiter(cfg.RateLimit, cfg.RateBurst)

	// Create client info
	clientInfo := &types.ClientInfo{
		Version:     "1.0.0",
		ProjectID:   cfg.ProjectID,
		UserAgent:   fmt.Sprintf("amapi-client/1.0.0 (project=%s)", cfg.ProjectID),
		Capabilities: []string{
			"enterprises",
			"policies",
			"devices",
			"enrollment_tokens",
			"applications",
		},
		CreatedAt: time.Now(),
	}

	client := &Client{
		service:      service,
		config:       cfg,
		ctx:          ctx,
		httpClient:   httpClient,
		retryHandler: retryHandler,
		rateLimiter:  rateLimiter,
		info:         clientInfo,
	}

	return client, nil
}

// NewWithContext creates a new Android Management API client with the specified context.
func NewWithContext(ctx context.Context, cfg *config.Config) (*Client, error) {
	client, err := New(cfg)
	if err != nil {
		return nil, err
	}

	client.ctx = ctx
	return client, nil
}

// createHTTPClient creates an authenticated HTTP client.
func createHTTPClient(ctx context.Context, cfg *config.Config) (*http.Client, error) {
	var creds *google.Credentials
	var err error

	// Load credentials
	if cfg.CredentialsJSON != "" {
		creds, err = google.CredentialsFromJSON(ctx, []byte(cfg.CredentialsJSON), cfg.Scopes...)
	} else if cfg.CredentialsFile != "" {
		// Read file and use CredentialsFromJSON
		jsonData, readErr := os.ReadFile(cfg.CredentialsFile)
		if readErr != nil {
			return nil, fmt.Errorf("failed to read credentials file: %w", readErr)
		}
		creds, err = google.CredentialsFromJSON(ctx, jsonData, cfg.Scopes...)
	} else {
		// Try default credentials
		creds, err = google.FindDefaultCredentials(ctx, cfg.Scopes...)
	}

	if err != nil {
		return nil, fmt.Errorf("failed to load credentials: %w", err)
	}

	// Create OAuth2 token source
	tokenSource := creds.TokenSource

	// Create HTTP client with authentication
	httpClient := oauth2.NewClient(ctx, tokenSource)

	// Set timeout
	httpClient.Timeout = cfg.Timeout

	return httpClient, nil
}

// GetInfo returns information about the client.
func (c *Client) GetInfo() *types.ClientInfo {
	return c.info
}

// GetConfig returns the client configuration.
func (c *Client) GetConfig() *config.Config {
	return c.config.Clone()
}

// Close closes the client and releases resources.
func (c *Client) Close() error {
	// Close any open connections
	if c.httpClient != nil {
		c.httpClient.CloseIdleConnections()
	}

	return nil
}

// Health checks the health of the client and API connectivity.
func (c *Client) Health() error {
	ctx, cancel := context.WithTimeout(c.ctx, 10*time.Second)
	defer cancel()

	// Try to list enterprises to test connectivity
	_, err := c.service.Enterprises.List().Context(ctx).Do()
	if err != nil {
		return types.WrapError(err, types.ErrCodeServiceUnavailable, "health check failed")
	}

	return nil
}

// executeWithRetry executes a function with retry logic.
func (c *Client) executeWithRetry(operation func() error) error {
	if !c.config.EnableRetry {
		return operation()
	}

	return c.retryHandler.Execute(operation)
}

// withRateLimit applies rate limiting to an operation.
func (c *Client) withRateLimit(operation func() error) error {
	if err := c.rateLimiter.Wait(c.ctx); err != nil {
		return types.WrapError(err, types.ErrCodeTooManyRequests, "rate limit exceeded")
	}

	return operation()
}

// executeAPICall executes an API call with rate limiting and retry logic.
func (c *Client) executeAPICall(operation func() error) error {
	return c.withRateLimit(func() error {
		return c.executeWithRetry(operation)
	})
}

// wrapAPIError wraps API errors with additional context.
func (c *Client) wrapAPIError(err error, operation string) error {
	if err == nil {
		return nil
	}

	// Check if it's already our error type
	if apiErr, ok := err.(*types.Error); ok {
		return apiErr
	}

	// Determine error code based on error type
	code := types.ErrCodeInternalServerError
	message := fmt.Sprintf("%s failed", operation)

	// Try to extract HTTP status code
	if httpErr, ok := err.(*googleapi.Error); ok {
		code = httpErr.Code
		message = httpErr.Message
	}

	return types.NewErrorWithCause(code, message, err)
}

// Utility methods

// buildResourceName builds a resource name from components.
func buildResourceName(components ...string) string {
	if len(components) == 0 {
		return ""
	}

	result := components[0]
	for i := 1; i < len(components); i++ {
		result += "/" + components[i]
	}

	return result
}

// parseResourceName parses a resource name into components.
func parseResourceName(resourceName string) []string {
	if resourceName == "" {
		return nil
	}

	var components []string
	start := 0

	for i, char := range resourceName {
		if char == '/' {
			if i > start {
				components = append(components, resourceName[start:i])
			}
			start = i + 1
		}
	}

	// Add the last component
	if start < len(resourceName) {
		components = append(components, resourceName[start:])
	}

	return components
}

// validateEnterpriseID validates an enterprise ID.
func validateEnterpriseID(enterpriseID string) error {
	if enterpriseID == "" {
		return types.ErrInvalidEnterpriseID
	}
	return nil
}

// validateDeviceID validates a device ID.
func validateDeviceID(deviceID string) error {
	if deviceID == "" {
		return types.ErrInvalidDeviceID
	}
	return nil
}

// validatePolicyID validates a policy ID.
func validatePolicyID(policyID string) error {
	if policyID == "" {
		return types.ErrInvalidPolicyID
	}
	return nil
}

// validateTokenID validates an enrollment token ID.
func validateTokenID(tokenID string) error {
	if tokenID == "" {
		return types.ErrInvalidTokenID
	}
	return nil
}

// buildEnterpriseName builds an enterprise resource name.
func buildEnterpriseName(enterpriseID string) string {
	return buildResourceName("enterprises", enterpriseID)
}

// buildDeviceName builds a device resource name.
func buildDeviceName(enterpriseID, deviceID string) string {
	return buildResourceName("enterprises", enterpriseID, "devices", deviceID)
}

// buildPolicyName builds a policy resource name.
func buildPolicyName(enterpriseID, policyID string) string {
	return buildResourceName("enterprises", enterpriseID, "policies", policyID)
}

// buildEnrollmentTokenName builds an enrollment token resource name.
func buildEnrollmentTokenName(enterpriseID, tokenID string) string {
	return buildResourceName("enterprises", enterpriseID, "enrollmentTokens", tokenID)
}

// parseEnterpriseName extracts the enterprise ID from an enterprise resource name.
func parseEnterpriseName(enterpriseName string) (string, error) {
	components := parseResourceName(enterpriseName)
	if len(components) != 2 || components[0] != "enterprises" {
		return "", types.NewErrorWithDetails(types.ErrCodeInvalidInput,
			"invalid enterprise name format", "expected format: enterprises/{enterpriseId}")
	}
	return components[1], nil
}

// parseDeviceName extracts enterprise and device IDs from a device resource name.
func parseDeviceName(deviceName string) (string, string, error) {
	components := parseResourceName(deviceName)
	if len(components) != 4 || components[0] != "enterprises" || components[2] != "devices" {
		return "", "", types.NewErrorWithDetails(types.ErrCodeInvalidInput,
			"invalid device name format", "expected format: enterprises/{enterpriseId}/devices/{deviceId}")
	}
	return components[1], components[3], nil
}

// parsePolicyName extracts enterprise and policy IDs from a policy resource name.
func parsePolicyName(policyName string) (string, string, error) {
	components := parseResourceName(policyName)
	if len(components) != 4 || components[0] != "enterprises" || components[2] != "policies" {
		return "", "", types.NewErrorWithDetails(types.ErrCodeInvalidInput,
			"invalid policy name format", "expected format: enterprises/{enterpriseId}/policies/{policyId}")
	}
	return components[1], components[3], nil
}

// parseEnrollmentTokenName extracts enterprise and token IDs from a token resource name.
func parseEnrollmentTokenName(tokenName string) (string, string, error) {
	components := parseResourceName(tokenName)
	if len(components) != 4 || components[0] != "enterprises" || components[2] != "enrollmentTokens" {
		return "", "", types.NewErrorWithDetails(types.ErrCodeInvalidInput,
			"invalid enrollment token name format", "expected format: enterprises/{enterpriseId}/enrollmentTokens/{tokenId}")
	}
	return components[1], components[3], nil
}