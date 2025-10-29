package config

import (
	"os"
	"testing"
	"time"
)

func TestDefaultConfig(t *testing.T) {
	cfg := DefaultConfig()

	if cfg.Timeout != 30*time.Second {
		t.Errorf("Expected timeout 30s, got %v", cfg.Timeout)
	}

	if cfg.RetryAttempts != 3 {
		t.Errorf("Expected retry attempts 3, got %d", cfg.RetryAttempts)
	}

	if !cfg.EnableRetry {
		t.Error("Expected retry to be enabled by default")
	}

	if len(cfg.Scopes) == 0 {
		t.Error("Expected default scopes to be set")
	}
}

func TestConfigValidation(t *testing.T) {
	tests := []struct {
		name        string
		config      *Config
		expectError bool
	}{
		{
			name: "valid config",
			config: &Config{
				ProjectID:       "test-project",
				CredentialsFile: "test-credentials.json",
				Timeout:         10 * time.Second,
				RetryAttempts:   3,
				LogLevel:        "info",
			},
			expectError: false,
		},
		{
			name: "missing project ID",
			config: &Config{
				CredentialsFile: "test-credentials.json",
				Timeout:         10 * time.Second,
			},
			expectError: true,
		},
		{
			name: "missing credentials",
			config: &Config{
				ProjectID: "test-project",
				Timeout:   10 * time.Second,
			},
			expectError: true,
		},
		{
			name: "invalid timeout",
			config: &Config{
				ProjectID:       "test-project",
				CredentialsFile: "test-credentials.json",
				Timeout:         -1 * time.Second,
			},
			expectError: true,
		},
		{
			name: "invalid log level",
			config: &Config{
				ProjectID:       "test-project",
				CredentialsFile: "test-credentials.json",
				Timeout:         10 * time.Second,
				LogLevel:        "invalid",
			},
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.config.Validate()
			if tt.expectError && err == nil {
				t.Error("Expected validation error, got nil")
			}
			if !tt.expectError && err != nil {
				t.Errorf("Expected no validation error, got %v", err)
			}
		})
	}
}

func TestLoadFromEnv(t *testing.T) {
	// Save original env vars
	originalVars := map[string]string{
		EnvProjectID:       os.Getenv(EnvProjectID),
		EnvCredentialsFile: os.Getenv(EnvCredentialsFile),
		EnvTimeout:         os.Getenv(EnvTimeout),
		EnvLogLevel:        os.Getenv(EnvLogLevel),
	}

	// Clean up after test
	defer func() {
		for key, value := range originalVars {
			if value == "" {
				os.Unsetenv(key)
			} else {
				os.Setenv(key, value)
			}
		}
	}()

	// Set test environment variables
	os.Setenv(EnvProjectID, "test-project")
	os.Setenv(EnvCredentialsFile, "test-credentials.json")
	os.Setenv(EnvTimeout, "45s")
	os.Setenv(EnvLogLevel, "debug")

	cfg, err := LoadFromEnv()
	if err != nil {
		t.Fatalf("Failed to load from env: %v", err)
	}

	if cfg.ProjectID != "test-project" {
		t.Errorf("Expected project ID 'test-project', got %s", cfg.ProjectID)
	}

	if cfg.CredentialsFile != "test-credentials.json" {
		t.Errorf("Expected credentials file 'test-credentials.json', got %s", cfg.CredentialsFile)
	}

	if cfg.Timeout != 45*time.Second {
		t.Errorf("Expected timeout 45s, got %v", cfg.Timeout)
	}

	if cfg.LogLevel != "debug" {
		t.Errorf("Expected log level 'debug', got %s", cfg.LogLevel)
	}
}

func TestConfigClone(t *testing.T) {
	original := &Config{
		ProjectID:       "test-project",
		CredentialsFile: "test-credentials.json",
		Scopes:          []string{"scope1", "scope2"},
		Timeout:         30 * time.Second,
	}

	cloned := original.Clone()

	// Verify values are copied
	if cloned.ProjectID != original.ProjectID {
		t.Error("ProjectID not cloned correctly")
	}

	// Verify slices are deep copied
	if &cloned.Scopes[0] == &original.Scopes[0] {
		t.Error("Scopes slice not deep copied")
	}

	// Modify clone and verify original is unchanged
	cloned.ProjectID = "modified-project"
	cloned.Scopes[0] = "modified-scope"

	if original.ProjectID == "modified-project" {
		t.Error("Original config was modified when clone was changed")
	}

	if original.Scopes[0] == "modified-scope" {
		t.Error("Original scopes were modified when clone was changed")
	}
}

func TestGetEnvVar(t *testing.T) {
	// Set up test environment variables
	os.Setenv("TEST_PRIMARY", "primary_value")
	os.Setenv("TEST_ALT1", "alt1_value")
	defer func() {
		os.Unsetenv("TEST_PRIMARY")
		os.Unsetenv("TEST_ALT1")
	}()

	// Test primary variable
	value := GetEnvVar("TEST_PRIMARY", "TEST_ALT1", "TEST_ALT2")
	if value != "primary_value" {
		t.Errorf("Expected 'primary_value', got %s", value)
	}

	// Test fallback to alternative
	os.Unsetenv("TEST_PRIMARY")
	value = GetEnvVar("TEST_PRIMARY", "TEST_ALT1", "TEST_ALT2")
	if value != "alt1_value" {
		t.Errorf("Expected 'alt1_value', got %s", value)
	}

	// Test no variables set
	os.Unsetenv("TEST_ALT1")
	value = GetEnvVar("TEST_PRIMARY", "TEST_ALT1", "TEST_ALT2")
	if value != "" {
		t.Errorf("Expected empty string, got %s", value)
	}
}

func TestParseDuration(t *testing.T) {
	tests := []struct {
		input    string
		expected time.Duration
		fallback time.Duration
	}{
		{"30s", 30 * time.Second, time.Minute},
		{"5m", 5 * time.Minute, time.Minute},
		{"30", 30 * time.Second, time.Minute}, // Parse as seconds
		{"", time.Minute, time.Minute},        // Use fallback
		{"invalid", time.Minute, time.Minute}, // Use fallback
	}

	for _, tt := range tests {
		result := parseDuration(tt.input, tt.fallback)
		if result != tt.expected {
			t.Errorf("parseDuration(%q, %v) = %v, expected %v",
				tt.input, tt.fallback, result, tt.expected)
		}
	}
}

func TestParseInt(t *testing.T) {
	tests := []struct {
		input    string
		expected int
		fallback int
	}{
		{"42", 42, 100},
		{"0", 0, 100},
		{"", 100, 100},     // Use fallback
		{"abc", 100, 100},  // Use fallback
		{"-5", -5, 100},    // Negative numbers
	}

	for _, tt := range tests {
		result := parseInt(tt.input, tt.fallback)
		if result != tt.expected {
			t.Errorf("parseInt(%q, %d) = %d, expected %d",
				tt.input, tt.fallback, result, tt.expected)
		}
	}
}

func TestParseBool(t *testing.T) {
	tests := []struct {
		input    string
		expected bool
		fallback bool
	}{
		{"true", true, false},
		{"false", false, true},
		{"1", true, false},
		{"0", false, true},
		{"", false, false},        // Use fallback
		{"invalid", false, false}, // Use fallback
	}

	for _, tt := range tests {
		result := parseBool(tt.input, tt.fallback)
		if result != tt.expected {
			t.Errorf("parseBool(%q, %t) = %t, expected %t",
				tt.input, tt.fallback, result, tt.expected)
		}
	}
}