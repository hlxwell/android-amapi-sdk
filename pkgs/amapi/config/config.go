// Package config provides configuration management for the amapi package.
// It supports loading configuration from environment variables, YAML files, and JSON files.
package config

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"time"

	"gopkg.in/yaml.v3"
)

// Config represents the configuration for the Android Management API client.
type Config struct {
	// Google Cloud configuration
	ProjectID               string `yaml:"project_id" json:"project_id"`
	CredentialsFile         string `yaml:"credentials_file" json:"credentials_file"`
	CredentialsJSON         string `yaml:"credentials_json" json:"credentials_json"`

	// API configuration
	ServiceAccountEmail     string `yaml:"service_account_email" json:"service_account_email"`
	Scopes                  []string `yaml:"scopes" json:"scopes"`

	// Client configuration
	Timeout                 time.Duration `yaml:"timeout" json:"timeout"`
	RetryAttempts          int    `yaml:"retry_attempts" json:"retry_attempts"`
	RetryDelay             time.Duration `yaml:"retry_delay" json:"retry_delay"`
	EnableRetry            bool   `yaml:"enable_retry" json:"enable_retry"`

	// Callback configuration
	CallbackURL            string `yaml:"callback_url" json:"callback_url"`

	// Cache configuration
	EnableCache            bool   `yaml:"enable_cache" json:"enable_cache"`
	CacheTTL               time.Duration `yaml:"cache_ttl" json:"cache_ttl"`

	// Logging configuration
	LogLevel               string `yaml:"log_level" json:"log_level"`
	EnableDebugLogging     bool   `yaml:"enable_debug_logging" json:"enable_debug_logging"`

	// Rate limiting
	RateLimit              int    `yaml:"rate_limit" json:"rate_limit"`
	RateBurst              int    `yaml:"rate_burst" json:"rate_burst"`
}

// DefaultConfig returns a configuration with sensible defaults.
func DefaultConfig() *Config {
	return &Config{
		Scopes: []string{
			"https://www.googleapis.com/auth/androidmanagement",
		},
		Timeout:                30 * time.Second,
		RetryAttempts:         3,
		RetryDelay:            1 * time.Second,
		EnableRetry:           true,
		EnableCache:           false,
		CacheTTL:              5 * time.Minute,
		LogLevel:              "info",
		EnableDebugLogging:    false,
		RateLimit:             100,
		RateBurst:             10,
	}
}

// Validate validates the configuration and returns an error if invalid.
func (c *Config) Validate() error {
	if c.ProjectID == "" {
		return fmt.Errorf("project_id is required")
	}

	if c.CredentialsFile == "" && c.CredentialsJSON == "" {
		return fmt.Errorf("either credentials_file or credentials_json must be specified")
	}

	if c.CredentialsFile != "" {
		if _, err := os.Stat(c.CredentialsFile); os.IsNotExist(err) {
			return fmt.Errorf("credentials file not found: %s", c.CredentialsFile)
		}
	}

	if c.Timeout <= 0 {
		return fmt.Errorf("timeout must be positive")
	}

	if c.RetryAttempts < 0 {
		return fmt.Errorf("retry_attempts must be non-negative")
	}

	if c.RetryDelay < 0 {
		return fmt.Errorf("retry_delay must be non-negative")
	}

	validLogLevels := map[string]bool{
		"debug": true,
		"info":  true,
		"warn":  true,
		"error": true,
	}
	if !validLogLevels[c.LogLevel] {
		return fmt.Errorf("invalid log_level: %s (must be debug, info, warn, or error)", c.LogLevel)
	}

	return nil
}

// LoadFromFile loads configuration from a YAML or JSON file.
func LoadFromFile(path string) (*Config, error) {
	config := DefaultConfig()

	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("failed to read config file: %w", err)
	}

	ext := filepath.Ext(path)
	switch ext {
	case ".yaml", ".yml":
		if err := yaml.Unmarshal(data, config); err != nil {
			return nil, fmt.Errorf("failed to parse YAML config: %w", err)
		}
	case ".json":
		if err := json.Unmarshal(data, config); err != nil {
			return nil, fmt.Errorf("failed to parse JSON config: %w", err)
		}
	default:
		return nil, fmt.Errorf("unsupported config file format: %s (supported: .yaml, .yml, .json)", ext)
	}

	if err := config.Validate(); err != nil {
		return nil, fmt.Errorf("invalid configuration: %w", err)
	}

	return config, nil
}

// SaveToFile saves the configuration to a YAML or JSON file.
func (c *Config) SaveToFile(path string) error {
	if err := c.Validate(); err != nil {
		return fmt.Errorf("invalid configuration: %w", err)
	}

	var data []byte
	var err error

	ext := filepath.Ext(path)
	switch ext {
	case ".yaml", ".yml":
		data, err = yaml.Marshal(c)
		if err != nil {
			return fmt.Errorf("failed to marshal YAML: %w", err)
		}
	case ".json":
		data, err = json.MarshalIndent(c, "", "  ")
		if err != nil {
			return fmt.Errorf("failed to marshal JSON: %w", err)
		}
	default:
		return fmt.Errorf("unsupported config file format: %s (supported: .yaml, .yml, .json)", ext)
	}

	if err := os.WriteFile(path, data, 0644); err != nil {
		return fmt.Errorf("failed to write config file: %w", err)
	}

	return nil
}

// Clone creates a deep copy of the configuration.
func (c *Config) Clone() *Config {
	clone := *c

	// Deep copy slices
	if c.Scopes != nil {
		clone.Scopes = make([]string, len(c.Scopes))
		copy(clone.Scopes, c.Scopes)
	}

	return &clone
}

// parseDuration safely parses a duration from environment variable.
func parseDuration(value string, defaultValue time.Duration) time.Duration {
	if value == "" {
		return defaultValue
	}

	if d, err := time.ParseDuration(value); err == nil {
		return d
	}

	// Try parsing as seconds if no unit specified
	if seconds, err := strconv.Atoi(value); err == nil {
		return time.Duration(seconds) * time.Second
	}

	return defaultValue
}

// parseInt safely parses an integer from environment variable.
func parseInt(value string, defaultValue int) int {
	if value == "" {
		return defaultValue
	}

	if i, err := strconv.Atoi(value); err == nil {
		return i
	}

	return defaultValue
}

// parseBool safely parses a boolean from environment variable.
func parseBool(value string, defaultValue bool) bool {
	if value == "" {
		return defaultValue
	}

	if b, err := strconv.ParseBool(value); err == nil {
		return b
	}

	return defaultValue
}