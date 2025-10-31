// Package config 提供 Android Management API 客户端的配置管理功能。
//
// 此包支持多种配置加载方式：
//   - 环境变量（最高优先级）
//   - YAML 配置文件
//   - JSON 配置文件
//   - 程序化配置（代码中直接构造）
//
// # 快速开始
//
// 使用自动配置加载（推荐）：
//
//	cfg, err := config.AutoLoadConfig()
//	if err != nil {
//	    log.Fatal(err)
//	}
//
// 从环境变量加载：
//
//	cfg, err := config.LoadFromEnv()
//
// 从文件加载：
//
//	cfg, err := config.LoadFromFile("./config.yaml")
//
// 手动配置：
//
//	cfg := &config.Config{
//	    ProjectID:       "your-project-id",
//	    CredentialsFile: "./sa-key.json",
//	    Timeout:         30 * time.Second,
//	    RetryAttempts:   3,
//	}
//
// # 配置文件搜索路径
//
// AutoLoadConfig 会按以下顺序搜索配置文件：
//   1. ./config.yaml
//   2. ./config.yml
//   3. ./amapi.yaml
//   4. ./amapi.yml
//   5. ~/.config/amapi/config.yaml
//   6. ~/.config/amapi/config.yml
//   7. /etc/amapi/config.yaml
//   8. /etc/amapi/config.yml
//
// # 环境变量
//
// 支持的环境变量：
//   - GOOGLE_CLOUD_PROJECT: Google Cloud 项目 ID
//   - GOOGLE_APPLICATION_CREDENTIALS: 服务账号密钥文件路径
//   - AMAPI_CALLBACK_URL: 企业注册回调 URL
//   - AMAPI_TIMEOUT: API 请求超时时间
//   - AMAPI_RETRY_ATTEMPTS: 重试次数
//   - AMAPI_ENABLE_RETRY: 是否启用重试
//   - AMAPI_LOG_LEVEL: 日志级别 (debug/info/warn/error)
//
// 详见各配置字段的文档。
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

// Config 包含 Android Management API 客户端的所有配置选项。
//
// 配置可以通过多种方式提供：环境变量、配置文件或程序化创建。
// 使用 Validate() 方法可以验证配置的完整性和有效性。
type Config struct {
	// Google Cloud 配置

	// ProjectID 是 Google Cloud 项目 ID（必需）。
	// 可通过环境变量 GOOGLE_CLOUD_PROJECT 设置。
	ProjectID string `yaml:"project_id" json:"project_id"`

	// CredentialsFile 是服务账号密钥 JSON 文件的路径。
	// 与 CredentialsJSON 二选一，优先使用 CredentialsFile。
	// 可通过环境变量 GOOGLE_APPLICATION_CREDENTIALS 设置。
	CredentialsFile string `yaml:"credentials_file" json:"credentials_file"`

	// CredentialsJSON 是服务账号密钥的 JSON 内容。
	// 与 CredentialsFile 二选一。
	CredentialsJSON string `yaml:"credentials_json" json:"credentials_json"`

	// API 配置

	// ServiceAccountEmail 是服务账号的邮箱地址（可选）。
	ServiceAccountEmail string `yaml:"service_account_email" json:"service_account_email"`

	// Scopes 是 OAuth2 权限范围列表。
	// 默认为 ["https://www.googleapis.com/auth/androidmanagement"]
	Scopes []string `yaml:"scopes" json:"scopes"`

	// 客户端配置

	// Timeout 是 API 请求的超时时间。
	// 默认为 30 秒。
	// 可通过环境变量 AMAPI_TIMEOUT 设置（如 "30s"）。
	Timeout time.Duration `yaml:"timeout" json:"timeout"`

	// RetryAttempts 是失败请求的最大重试次数。
	// 默认为 3 次。
	// 可通过环境变量 AMAPI_RETRY_ATTEMPTS 设置。
	RetryAttempts int `yaml:"retry_attempts" json:"retry_attempts"`

	// RetryDelay 是重试之间的基础延迟时间。
	// 实际延迟使用指数退避算法计算。
	// 默认为 1 秒。
	// 可通过环境变量 AMAPI_RETRY_DELAY 设置（如 "1s"）。
	RetryDelay time.Duration `yaml:"retry_delay" json:"retry_delay"`

	// EnableRetry 控制是否启用自动重试。
	// 默认为 true。
	// 可通过环境变量 AMAPI_ENABLE_RETRY 设置。
	EnableRetry bool `yaml:"enable_retry" json:"enable_retry"`

	// 回调配置

	// CallbackURL 是企业注册完成后的回调 URL。
	// 可通过环境变量 AMAPI_CALLBACK_URL 设置。
	CallbackURL string `yaml:"callback_url" json:"callback_url"`

	// 缓存配置

	// EnableCache 控制是否启用响应缓存（实验性功能）。
	// 默认为 false。
	EnableCache bool `yaml:"enable_cache" json:"enable_cache"`

	// CacheTTL 是缓存的有效期。
	// 默认为 5 分钟。
	CacheTTL time.Duration `yaml:"cache_ttl" json:"cache_ttl"`

	// 日志配置

	// LogLevel 是日志级别，可选值：debug, info, warn, error。
	// 默认为 "info"。
	// 可通过环境变量 AMAPI_LOG_LEVEL 设置。
	LogLevel string `yaml:"log_level" json:"log_level"`

	// EnableDebugLogging 控制是否启用详细的调试日志。
	// 默认为 false。
	// 可通过环境变量 AMAPI_ENABLE_DEBUG_LOGGING 设置。
	EnableDebugLogging bool `yaml:"enable_debug_logging" json:"enable_debug_logging"`

	// 速率限制

	// RateLimit 是每分钟允许的最大请求数。
	// 默认为 100。
	// 可通过环境变量 AMAPI_RATE_LIMIT 设置。
	RateLimit int `yaml:"rate_limit" json:"rate_limit"`

	// RateBurst 是允许的突发请求数量。
	// 默认为 10。
	// 可通过环境变量 AMAPI_RATE_BURST 设置。
	RateBurst int `yaml:"rate_burst" json:"rate_burst"`

	// Redis 配置（用于分布式 rate limiting 和 retry 管理）

	// RedisAddress 是 Redis 服务器地址（格式：host:port）。
	// 如果设置，将使用 Redis 实现分布式的 rate limiting 和 retry 管理。
	// 可通过环境变量 AMAPI_REDIS_ADDRESS 设置。
	RedisAddress string `yaml:"redis_address" json:"redis_address"`

	// RedisPassword 是 Redis 服务器密码（可选）。
	// 可通过环境变量 AMAPI_REDIS_PASSWORD 设置。
	RedisPassword string `yaml:"redis_password" json:"redis_password"`

	// RedisDB 是 Redis 数据库编号。
	// 默认为 0。
	// 可通过环境变量 AMAPI_REDIS_DB 设置。
	RedisDB int `yaml:"redis_db" json:"redis_db"`

	// RedisKeyPrefix 是 Redis key 的前缀。
	// 用于区分不同项目或环境的 key。
	// 默认为 "amapi:"。
	// 可通过环境变量 AMAPI_REDIS_KEY_PREFIX 设置。
	RedisKeyPrefix string `yaml:"redis_key_prefix" json:"redis_key_prefix"`

	// UseRedisRateLimit 控制是否使用 Redis 进行分布式 rate limiting。
	// 如果 RedisAddress 未设置，此选项无效。
	// 默认为 false。
	// 可通过环境变量 AMAPI_USE_REDIS_RATE_LIMIT 设置。
	UseRedisRateLimit bool `yaml:"use_redis_rate_limit" json:"use_redis_rate_limit"`

	// UseRedisRetry 控制是否使用 Redis 进行分布式 retry 管理。
	// 如果 RedisAddress 未设置，此选项无效。
	// 默认为 false。
	// 可通过环境变量 AMAPI_USE_REDIS_RETRY 设置。
	UseRedisRetry bool `yaml:"use_redis_retry" json:"use_redis_retry"`
}

// DefaultConfig 返回一个包含合理默认值的配置对象。
//
// 返回的配置包含以下默认值：
//   - Timeout: 30 秒
//   - RetryAttempts: 3 次
//   - RetryDelay: 1 秒
//   - EnableRetry: true
//   - LogLevel: "info"
//   - RateLimit: 100 次/分钟
//   - RateBurst: 10 次
//
// 你仍然需要设置必需的字段（ProjectID 和认证信息）。
//
// 示例：
//
//	cfg := config.DefaultConfig()
//	cfg.ProjectID = "your-project-id"
//	cfg.CredentialsFile = "./sa-key.json"
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
		RedisDB:               0,
		RedisKeyPrefix:        "amapi:",
		UseRedisRateLimit:     false,
		UseRedisRetry:         false,
	}
}

// Validate 验证配置的完整性和有效性。
//
// 执行以下验证：
//   - ProjectID 不能为空
//   - 必须提供 CredentialsFile 或 CredentialsJSON 之一
//   - 如果提供了 CredentialsFile，检查文件是否存在
//   - Timeout 必须大于 0
//   - RetryAttempts 必须非负
//   - RetryDelay 必须非负
//   - LogLevel 必须是 debug/info/warn/error 之一
//
// 返回第一个发现的验证错误，如果配置有效则返回 nil。
//
// 示例：
//
//	cfg := &Config{...}
//	if err := cfg.Validate(); err != nil {
//	    log.Fatalf("配置无效: %v", err)
//	}
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

// LoadFromFile 从 YAML 或 JSON 文件加载配置。
//
// 支持的文件格式：
//   - .yaml, .yml (YAML 格式)
//   - .json (JSON 格式)
//
// 文件格式由扩展名自动识别。
// 加载的配置会与默认配置合并，文件中的值覆盖默认值。
//
// 参数：
//   - path: 配置文件的路径
//
// 返回：
//   - 加载并验证后的配置对象
//   - 如果文件不存在、格式错误或验证失败，返回错误
//
// 示例：
//
//	cfg, err := config.LoadFromFile("./config.yaml")
//	if err != nil {
//	    log.Fatal(err)
//	}
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

// SaveToFile 将配置保存到 YAML 或 JSON 文件。
//
// 支持的文件格式：
//   - .yaml, .yml (YAML 格式)
//   - .json (JSON 格式，带缩进美化)
//
// 文件格式由扩展名自动识别。
// 在保存前会先验证配置的有效性。
//
// 参数：
//   - path: 目标文件的路径
//
// 返回：
//   - 如果配置无效、序列化失败或写入失败，返回错误
//
// 示例：
//
//	cfg := config.DefaultConfig()
//	cfg.ProjectID = "my-project"
//	cfg.CredentialsFile = "./key.json"
//
//	if err := cfg.SaveToFile("./config.yaml"); err != nil {
//	    log.Fatal(err)
//	}
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

// Clone 创建配置的深拷贝。
//
// 返回的配置对象是完全独立的副本，修改副本不会影响原配置。
// 包括所有嵌套的切片（如 Scopes）都会被深拷贝。
//
// 这在需要基于现有配置创建变体时很有用。
//
// 示例：
//
//	originalCfg := config.DefaultConfig()
//
//	// 创建一个副本用于测试
//	testCfg := originalCfg.Clone()
//	testCfg.RetryAttempts = 10
//	testCfg.EnableDebugLogging = true
//
//	// originalCfg 保持不变
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