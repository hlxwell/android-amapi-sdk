package config

import "os"

// Environment variable names used by the amapi package.
const (
	// Google Cloud configuration
	EnvProjectID               = "GOOGLE_CLOUD_PROJECT"
	EnvCredentialsFile         = "GOOGLE_APPLICATION_CREDENTIALS"
	EnvCredentialsJSON         = "GOOGLE_APPLICATION_CREDENTIALS_JSON"
	EnvServiceAccountEmail     = "AMAPI_SERVICE_ACCOUNT_EMAIL"

	// API configuration
	EnvScopes                  = "AMAPI_SCOPES"

	// Client configuration
	EnvTimeout                 = "AMAPI_TIMEOUT"
	EnvRetryAttempts          = "AMAPI_RETRY_ATTEMPTS"
	EnvRetryDelay             = "AMAPI_RETRY_DELAY"
	EnvEnableRetry            = "AMAPI_ENABLE_RETRY"

	// Callback configuration
	EnvCallbackURL            = "AMAPI_CALLBACK_URL"

	// Cache configuration
	EnvEnableCache            = "AMAPI_ENABLE_CACHE"
	EnvCacheTTL               = "AMAPI_CACHE_TTL"

	// Logging configuration
	EnvLogLevel               = "AMAPI_LOG_LEVEL"
	EnvEnableDebugLogging     = "AMAPI_ENABLE_DEBUG_LOGGING"

	// Rate limiting
	EnvRateLimit              = "AMAPI_RATE_LIMIT"
	EnvRateBurst              = "AMAPI_RATE_BURST"

	// Alternative environment variable names for compatibility
	AltEnvProjectID           = "AMAPI_PROJECT_ID"
	AltEnvCredentialsFile     = "AMAPI_CREDENTIALS_FILE"
	AltEnvCredentialsJSON     = "AMAPI_CREDENTIALS_JSON"
)

// GetEnvVar returns the value of an environment variable, trying multiple possible names.
func GetEnvVar(primary string, alternatives ...string) string {
	if value := getEnv(primary); value != "" {
		return value
	}

	for _, alt := range alternatives {
		if value := getEnv(alt); value != "" {
			return value
		}
	}

	return ""
}

// getEnv is a helper function to get environment variable value.
func getEnv(key string) string {
	return os.Getenv(key)
}