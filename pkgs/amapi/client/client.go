// Package client provides a high-level client for the Android Management API.
//
// 这个包实现了 Android Management API 的核心客户端功能，包括：
//
//   - 认证和连接管理
//   - 自动重试机制（支持本地和分布式 Redis 实现）
//   - 速率限制（支持本地和分布式 Redis 实现）
//   - 错误处理和包装
//   - 资源清理
//
// # 快速开始
//
//	cfg := &config.Config{
//	    ProjectID:       "your-project-id",
//	    CredentialsFile: "./sa-key.json",
//	}
//
//	client, err := New(cfg)
//	if err != nil {
//	    log.Fatal(err)
//	}
//	defer client.Close()
//
// # 分布式 Rate Limiting 和 Retry
//
// 当您的应用程序运行多个进程时，可以使用 Redis 实现分布式的 rate limiting 和 retry 管理：
//
//	cfg := &config.Config{
//	    ProjectID:        "your-project-id",
//	    CredentialsFile:  "./sa-key.json",
//	    RedisAddress:     "localhost:6379",
//	    UseRedisRateLimit: true,  // 所有进程共享同一个 rate limit
//	    UseRedisRetry:     true,  // 防止多个进程同时重试同一操作
//	}
//
// 这样所有进程会共享同一个 rate limit，确保不会超过 API 的限制。
//
// # 服务访问
//
// 客户端提供了多个服务访问方法：
//
//	enterprises := client.Enterprises()
//	policies := client.Policies()
//	devices := client.Devices()
//	enrollment := client.EnrollmentTokens()
//
// 每个服务都有完整的 CRUD 操作方法。
//
// 更多详细信息请参考各服务类型的文档。
package client

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"strings"
	"time"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/androidmanagement/v1"
	"google.golang.org/api/googleapi"
	"google.golang.org/api/option"

	"amapi-pkg/pkgs/amapi/config"
	"amapi-pkg/pkgs/amapi/types"
	"amapi-pkg/pkgs/amapi/utils"

	"github.com/redis/go-redis/v9"
)

// 常量定义
const (
	// ClientVersion 是客户端版本号
	ClientVersion = "1.0.0"

	// DefaultRetryMaxDelay 是默认的最大重试延迟时间
	DefaultRetryMaxDelay = 30 * time.Second

	// DefaultRedisTimeout 是默认的Redis连接超时时间
	DefaultRedisTimeout = 5 * time.Second

	// DefaultHealthCheckTimeout 是默认的健康检查超时时间
	DefaultHealthCheckTimeout = 10 * time.Second
)

// Client represents the Android Management API client.
//
// Client 提供了访问 Android Management API 的所有方法。
// 它是线程安全的，可以在多个 goroutine 中并发使用。
// 线程安全性来自底层的 http.Client 和 androidmanagement.Service，它们本身都是线程安全的。
//
// 使用 New 或 NewWithContext 创建客户端实例。
// 在使用完毕后，务必调用 Close() 方法释放资源（包括 Redis 连接）。
//
// 示例：
//
//	client, err := New(cfg)
//	if err != nil {
//	    return err
//	}
//	defer client.Close()  // 确保释放资源
//
//	// 使用客户端
//	enterprises, err := client.Enterprises().List(nil)
//
// # 分布式支持
//
// 如果配置了 Redis，Client 会自动使用 Redis 实现分布式的 rate limiting 和 retry 管理。
// 这对于多进程部署非常重要，可以确保所有进程共享同一个 rate limit。
type Client struct {
	// service is the underlying Android Management API service
	service *androidmanagement.Service

	// config contains the client configuration
	config *config.Config

	// ctx is the context for API calls
	ctx context.Context

	// httpClient is the HTTP client used for requests
	httpClient *http.Client

	// retryHandler handles retry logic (local or Redis-based)
	retryHandler utils.RetryHandlerInterface

	// rateLimiter handles rate limiting (local or Redis-based)
	rateLimiter utils.RateLimiterInterface

	// redisClient is the Redis client (if using Redis for distributed rate limiting/retry)
	redisClient *redis.Client

	// info contains client information
	info *types.ClientInfo
}

// New creates a new Android Management API client.
//
// 根据配置创建客户端实例，支持：
//   - 本地 rate limiting 和 retry（默认）
//   - 分布式 Redis rate limiting 和 retry（如果配置了 Redis）
//
// 如果配置了 RedisAddress，会自动尝试连接 Redis 并验证连接。
// 如果 Redis 连接失败，会返回错误。
//
// 使用示例：
//
//	cfg := &config.Config{
//	    ProjectID:       "your-project-id",
//	    CredentialsFile: "./sa-key.json",
//	}
//
//	client, err := New(cfg)
//	if err != nil {
//	    log.Fatal(err)
//	}
//	defer client.Close()
//
// 使用分布式 Redis：
//
//	cfg := &config.Config{
//	    ProjectID:        "your-project-id",
//	    CredentialsFile:  "./sa-key.json",
//	    RedisAddress:     "localhost:6379",
//	    UseRedisRateLimit: true,
//	    UseRedisRetry:     true,
//	}
//
//	client, err := New(cfg)
func New(cfg *config.Config) (*Client, error) {
	return newClientWithContext(context.Background(), cfg)
}

// newClientWithContext 是内部的客户端创建函数，支持自定义 context
func newClientWithContext(ctx context.Context, cfg *config.Config) (*Client, error) {
	if cfg == nil {
		return nil, types.NewError(types.ErrCodeConfiguration, "configuration is required")
	}

	if err := cfg.Validate(); err != nil {
		return nil, types.WrapError(err, types.ErrCodeConfiguration, "invalid configuration")
	}

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

	// Initialize Redis client if configured
	var redisClient *redis.Client
	if cfg.RedisAddress != "" {
		redisClient = redis.NewClient(&redis.Options{
			Addr:     cfg.RedisAddress,
			Password: cfg.RedisPassword,
			DB:       cfg.RedisDB,
		})

		// Test Redis connection
		pingCtx, cancel := context.WithTimeout(context.Background(), DefaultRedisTimeout)
		defer cancel()
		if err := redisClient.Ping(pingCtx).Err(); err != nil {
			// 确保在连接失败时清理 Redis 客户端资源
			redisClient.Close()
			return nil, types.WrapError(err, types.ErrCodeConfiguration, "failed to connect to Redis")
		}
	}

	// Create retry handler (Redis or local)
	var retryHandler utils.RetryHandlerInterface
	retryConfig := utils.RetryConfig{
		MaxAttempts: cfg.RetryAttempts,
		BaseDelay:   cfg.RetryDelay,
		MaxDelay:    DefaultRetryMaxDelay,
		EnableRetry: cfg.EnableRetry,
	}

	if redisClient != nil && cfg.UseRedisRetry {
		retryHandler = utils.NewRedisRetryHandler(redisClient, cfg.RedisKeyPrefix, retryConfig)
	} else {
		retryHandler = utils.NewRetryHandler(retryConfig)
	}

	// Create rate limiter (Redis or local)
	var rateLimiter utils.RateLimiterInterface
	if redisClient != nil && cfg.UseRedisRateLimit {
		rateLimiter = utils.NewRedisRateLimiter(redisClient, cfg.RedisKeyPrefix, cfg.RateLimit, cfg.RateBurst)
	} else {
		rateLimiter = utils.NewRateLimiter(cfg.RateLimit, cfg.RateBurst)
	}

	// Create client info
	clientInfo := &types.ClientInfo{
		Version:   ClientVersion,
		ProjectID: cfg.ProjectID,
		UserAgent: fmt.Sprintf("amapi-client/%s (project=%s)", ClientVersion, cfg.ProjectID),
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
		redisClient:  redisClient,
		info:         clientInfo,
	}

	return client, nil
}

// NewWithContext creates a new Android Management API client with the specified context.
//
// 与 New 功能相同，但使用指定的 context 而不是默认的 Background context。
// 这对于需要超时控制或取消操作的场景很有用。
//
// 示例：
//
//	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
//	defer cancel()
//
//	client, err := NewWithContext(ctx, cfg)
func NewWithContext(ctx context.Context, cfg *config.Config) (*Client, error) {
	return newClientWithContext(ctx, cfg)
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
			return nil, types.WrapError(readErr, types.ErrCodeConfiguration, "failed to read credentials file")
		}
		creds, err = google.CredentialsFromJSON(ctx, jsonData, cfg.Scopes...)
	} else {
		// Try default credentials
		creds, err = google.FindDefaultCredentials(ctx, cfg.Scopes...)
	}

	if err != nil {
		return nil, types.WrapError(err, types.ErrCodeAuthentication, "failed to load credentials")
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
	// Close rate limiter
	if c.rateLimiter != nil {
		if err := c.rateLimiter.Close(); err != nil {
			return err
		}
	}

	// Close retry handler
	if c.retryHandler != nil {
		if err := c.retryHandler.Close(); err != nil {
			return err
		}
	}

	// Close Redis client
	if c.redisClient != nil {
		if err := c.redisClient.Close(); err != nil {
			return err
		}
	}

	// Close any open HTTP connections
	if c.httpClient != nil {
		c.httpClient.CloseIdleConnections()
	}

	return nil
}

// Health checks the health of the client and API connectivity.
func (c *Client) Health() error {
	ctx, cancel := context.WithTimeout(c.ctx, DefaultHealthCheckTimeout)
	defer cancel()

	// Try to list enterprises to test connectivity
	// ProjectId is required for the API call
	if c.config.ProjectID == "" {
		return types.NewError(types.ErrCodeConfiguration, "project ID is required for health check")
	}

	call := c.service.Enterprises.List().ProjectId(c.config.ProjectID)
	_, err := call.Context(ctx).Do()
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

	// Generate operation ID for distributed retry coordination
	operationID := fmt.Sprintf("%d", time.Now().UnixNano())

	return c.retryHandler.Execute(c.ctx, operationID, operation)
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
	// Apply rate limiting first
	if err := c.rateLimiter.Wait(c.ctx); err != nil {
		return types.WrapError(err, types.ErrCodeTooManyRequests, "rate limit exceeded")
	}

	// Then apply retry logic
	if !c.config.EnableRetry {
		return operation()
	}

	// Generate operation ID for distributed retry coordination
	operationID := fmt.Sprintf("%d", time.Now().UnixNano())
	return c.retryHandler.Execute(c.ctx, operationID, operation)
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

// validateResourceName 验证资源名称格式并返回解析后的组件
// expectedParts 定义了期望的资源名称格式，例如 []string{"enterprises", "{enterpriseId}", "devices", "{deviceId}"}
func validateResourceName(resourceName string, expectedParts []string, resourceType string) ([]string, error) {
	if resourceName == "" {
		return nil, types.NewErrorWithDetails(types.ErrCodeInvalidInput,
			fmt.Sprintf("empty %s name", resourceType),
			fmt.Sprintf("expected format: %s", strings.Join(expectedParts, "/")))
	}

	components := parseResourceName(resourceName)
	if len(components) != len(expectedParts) {
		return nil, types.NewErrorWithDetails(types.ErrCodeInvalidInput,
			fmt.Sprintf("invalid %s name format", resourceType),
			fmt.Sprintf("expected format: %s", strings.Join(expectedParts, "/")))
	}

	// 验证固定部分（非参数部分）
	for i, expected := range expectedParts {
		if !strings.HasPrefix(expected, "{") && !strings.HasSuffix(expected, "}") {
			if components[i] != expected {
				return nil, types.NewErrorWithDetails(types.ErrCodeInvalidInput,
					fmt.Sprintf("invalid %s name format", resourceType),
					fmt.Sprintf("expected format: %s", strings.Join(expectedParts, "/")))
			}
		}
	}

	return components, nil
}

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
	expectedParts := []string{"enterprises", "{enterpriseId}"}
	components, err := validateResourceName(enterpriseName, expectedParts, "enterprise")
	if err != nil {
		return "", err
	}
	return components[1], nil
}

// parseDeviceName extracts enterprise and device IDs from a device resource name.
func parseDeviceName(deviceName string) (string, string, error) {
	expectedParts := []string{"enterprises", "{enterpriseId}", "devices", "{deviceId}"}
	components, err := validateResourceName(deviceName, expectedParts, "device")
	if err != nil {
		return "", "", err
	}
	return components[1], components[3], nil
}

// parsePolicyName extracts enterprise and policy IDs from a policy resource name.
func parsePolicyName(policyName string) (string, string, error) {
	expectedParts := []string{"enterprises", "{enterpriseId}", "policies", "{policyId}"}
	components, err := validateResourceName(policyName, expectedParts, "policy")
	if err != nil {
		return "", "", err
	}
	return components[1], components[3], nil
}

// parseEnrollmentTokenName extracts enterprise and token IDs from a token resource name.
func parseEnrollmentTokenName(tokenName string) (string, string, error) {
	expectedParts := []string{"enterprises", "{enterpriseId}", "enrollmentTokens", "{tokenId}"}
	components, err := validateResourceName(tokenName, expectedParts, "enrollment token")
	if err != nil {
		return "", "", err
	}
	return components[1], components[3], nil
}

// parseMigrationTokenName extracts enterprise and token IDs from a migration token resource name.
func parseMigrationTokenName(tokenName string) (string, string, error) {
	expectedParts := []string{"enterprises", "{enterpriseId}", "migrationTokens", "{tokenId}"}
	components, err := validateResourceName(tokenName, expectedParts, "migration token")
	if err != nil {
		return "", "", err
	}
	return components[1], components[3], nil
}

// parseWebAppName extracts enterprise and web app IDs from a web app resource name.
func parseWebAppName(webAppName string) (string, string, error) {
	expectedParts := []string{"enterprises", "{enterpriseId}", "webApps", "{webAppId}"}
	components, err := validateResourceName(webAppName, expectedParts, "web app")
	if err != nil {
		return "", "", err
	}
	return components[1], components[3], nil
}

// parseWebTokenName extracts enterprise and token IDs from a web token resource name.
func parseWebTokenName(tokenName string) (string, string, error) {
	expectedParts := []string{"enterprises", "{enterpriseId}", "webTokens", "{tokenId}"}
	components, err := validateResourceName(tokenName, expectedParts, "web token")
	if err != nil {
		return "", "", err
	}
	return components[1], components[3], nil
}
