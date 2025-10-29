package config

import (
	"fmt"
	"os"
	"strings"
)

// LoadConfig loads configuration from multiple sources in the following priority order:
// 1. Environment variables
// 2. Configuration file (if specified)
// 3. Default values
func LoadConfig(configPath string) (*Config, error) {
	config := DefaultConfig()

	// Load from file if specified
	if configPath != "" {
		fileConfig, err := LoadFromFile(configPath)
		if err != nil {
			return nil, fmt.Errorf("failed to load config from file: %w", err)
		}
		config = fileConfig
	}

	// Override with environment variables
	loadFromEnv(config)

	// Validate final configuration
	if err := config.Validate(); err != nil {
		return nil, fmt.Errorf("invalid configuration: %w", err)
	}

	return config, nil
}

// LoadFromEnv loads configuration entirely from environment variables.
func LoadFromEnv() (*Config, error) {
	config := DefaultConfig()
	loadFromEnv(config)

	if err := config.Validate(); err != nil {
		return nil, fmt.Errorf("invalid configuration: %w", err)
	}

	return config, nil
}

// loadFromEnv loads configuration from environment variables, overriding existing values.
func loadFromEnv(config *Config) {
	// Google Cloud configuration
	if projectID := GetEnvVar(EnvProjectID, AltEnvProjectID); projectID != "" {
		config.ProjectID = projectID
	}

	if credFile := GetEnvVar(EnvCredentialsFile, AltEnvCredentialsFile); credFile != "" {
		config.CredentialsFile = credFile
	}

	if credJSON := GetEnvVar(EnvCredentialsJSON, AltEnvCredentialsJSON); credJSON != "" {
		config.CredentialsJSON = credJSON
	}

	if serviceAccount := GetEnvVar(EnvServiceAccountEmail); serviceAccount != "" {
		config.ServiceAccountEmail = serviceAccount
	}

	// API configuration
	if scopes := GetEnvVar(EnvScopes); scopes != "" {
		config.Scopes = strings.Split(scopes, ",")
		// Trim whitespace from each scope
		for i, scope := range config.Scopes {
			config.Scopes[i] = strings.TrimSpace(scope)
		}
	}

	// Client configuration
	if timeout := GetEnvVar(EnvTimeout); timeout != "" {
		config.Timeout = parseDuration(timeout, config.Timeout)
	}

	if retryAttempts := GetEnvVar(EnvRetryAttempts); retryAttempts != "" {
		config.RetryAttempts = parseInt(retryAttempts, config.RetryAttempts)
	}

	if retryDelay := GetEnvVar(EnvRetryDelay); retryDelay != "" {
		config.RetryDelay = parseDuration(retryDelay, config.RetryDelay)
	}

	if enableRetry := GetEnvVar(EnvEnableRetry); enableRetry != "" {
		config.EnableRetry = parseBool(enableRetry, config.EnableRetry)
	}

	// Callback configuration
	if callbackURL := GetEnvVar(EnvCallbackURL); callbackURL != "" {
		config.CallbackURL = callbackURL
	}

	// Cache configuration
	if enableCache := GetEnvVar(EnvEnableCache); enableCache != "" {
		config.EnableCache = parseBool(enableCache, config.EnableCache)
	}

	if cacheTTL := GetEnvVar(EnvCacheTTL); cacheTTL != "" {
		config.CacheTTL = parseDuration(cacheTTL, config.CacheTTL)
	}

	// Logging configuration
	if logLevel := GetEnvVar(EnvLogLevel); logLevel != "" {
		config.LogLevel = strings.ToLower(logLevel)
	}

	if enableDebugLogging := GetEnvVar(EnvEnableDebugLogging); enableDebugLogging != "" {
		config.EnableDebugLogging = parseBool(enableDebugLogging, config.EnableDebugLogging)
	}

	// Rate limiting
	if rateLimit := GetEnvVar(EnvRateLimit); rateLimit != "" {
		config.RateLimit = parseInt(rateLimit, config.RateLimit)
	}

	if rateBurst := GetEnvVar(EnvRateBurst); rateBurst != "" {
		config.RateBurst = parseInt(rateBurst, config.RateBurst)
	}
}

// AutoLoadConfig attempts to automatically load configuration from common locations.
// It searches for configuration files in the following order:
// 1. ./config.yaml
// 2. ./config.yml
// 3. ./amapi.yaml
// 4. ./amapi.yml
// 5. ~/.config/amapi/config.yaml
// 6. ~/.config/amapi/config.yml
// 7. /etc/amapi/config.yaml
// 8. /etc/amapi/config.yml
func AutoLoadConfig() (*Config, error) {
	searchPaths := []string{
		"./config.yaml",
		"./config.yml",
		"./amapi.yaml",
		"./amapi.yml",
	}

	// Add user config directory paths
	if homeDir, err := os.UserHomeDir(); err == nil {
		userConfigDir := homeDir + "/.config/amapi"
		searchPaths = append(searchPaths,
			userConfigDir+"/config.yaml",
			userConfigDir+"/config.yml",
			userConfigDir+"/config.json",
		)
	}

	// Add system config directory paths
	searchPaths = append(searchPaths,
		"/etc/amapi/config.yaml",
		"/etc/amapi/config.yml",
		"/etc/amapi/config.json",
	)

	// Try to find and load a config file
	for _, path := range searchPaths {
		if _, err := os.Stat(path); err == nil {
			return LoadConfig(path)
		}
	}

	// If no config file found, load from environment variables only
	return LoadFromEnv()
}

// CreateExampleConfigFile creates an example configuration file at the specified path.
func CreateExampleConfigFile(path string) error {
	config := DefaultConfig()

	// Set example values
	config.ProjectID = "your-gcp-project-id"
	config.CredentialsFile = "./service-account-key.json"
	config.CallbackURL = "https://your-app.example.com/api/v1/amapi/callback"
	config.LogLevel = "info"

	return config.SaveToFile(path)
}

// GetConfigSummary returns a human-readable summary of the current configuration.
func (c *Config) GetConfigSummary() string {
	var summary strings.Builder

	summary.WriteString("Android Management API Configuration Summary:\n")
	summary.WriteString("==========================================\n")
	summary.WriteString(fmt.Sprintf("Project ID: %s\n", c.ProjectID))

	if c.CredentialsFile != "" {
		summary.WriteString(fmt.Sprintf("Credentials File: %s\n", c.CredentialsFile))
	} else {
		summary.WriteString("Credentials: JSON string provided\n")
	}

	summary.WriteString(fmt.Sprintf("Scopes: %s\n", strings.Join(c.Scopes, ", ")))
	summary.WriteString(fmt.Sprintf("Timeout: %v\n", c.Timeout))
	summary.WriteString(fmt.Sprintf("Retry Attempts: %d\n", c.RetryAttempts))
	summary.WriteString(fmt.Sprintf("Retry Enabled: %t\n", c.EnableRetry))

	if c.CallbackURL != "" {
		summary.WriteString(fmt.Sprintf("Callback URL: %s\n", c.CallbackURL))
	}

	summary.WriteString(fmt.Sprintf("Cache Enabled: %t\n", c.EnableCache))
	if c.EnableCache {
		summary.WriteString(fmt.Sprintf("Cache TTL: %v\n", c.CacheTTL))
	}

	summary.WriteString(fmt.Sprintf("Log Level: %s\n", c.LogLevel))
	summary.WriteString(fmt.Sprintf("Debug Logging: %t\n", c.EnableDebugLogging))
	summary.WriteString(fmt.Sprintf("Rate Limit: %d/min, Burst: %d\n", c.RateLimit, c.RateBurst))

	return summary.String()
}