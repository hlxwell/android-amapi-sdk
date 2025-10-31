// Package amapi 提供了 Android Management API 的 Go 客户端库。
//
// 这个包是 Android Management API 的完整实现，提供了企业移动设备管理的所有功能，
// 包括企业管理、策略配置、设备控制、注册令牌管理等。
//
// # 设计理念
//
// 本 SDK 直接使用 google.golang.org/api/androidmanagement/v1 包中的原生类型，
// 避免不必要的类型转换，提高代码效率和可维护性。所有核心类型（Enterprise、Policy、
// Device 等）都是 androidmanagement 包类型的别名，可以直接使用。
//
// 辅助功能通过 types 包中的工具函数提供，而不是通过自定义类型的方法。
//
// # 快速开始
//
// 基本用法：
//
//	cfg, err := config.AutoLoadConfig()
//	if err != nil {
//	    log.Fatal(err)
//	}
//
//	c, err := amapi.NewClient(cfg)
//	if err != nil {
//	    log.Fatal(err)
//	}
//	defer c.Close()
//
//	// 列出企业
//	result, err := c.Enterprises().List(nil)
//	if err != nil {
//	    log.Fatal(err)
//	}
//
//	for _, enterprise := range result.Items {
//	    fmt.Printf("Enterprise: %s\n", enterprise.Name)
//	}
//
// # 配置
//
// 支持多种配置方式：
//   - 环境变量
//   - YAML/JSON 配置文件
//   - 程序化配置
//
// 支持分布式 rate limiting 和 retry 管理（使用 Redis）：
//
//	cfg := &config.Config{
//	    ProjectID:       "your-project-id",
//	    CredentialsFile: "./sa-key.json",
//	    RedisAddress:    "localhost:6379",
//	    UseRedisRateLimit: true,  // 启用分布式 rate limiting
//	    UseRedisRetry:     true,  // 启用分布式 retry 管理
//	}
//
// 详见 config 包的文档。
//
// # 核心功能
//
// 企业管理：
//   - 创建和管理企业
//   - 生成注册 URL
//   - 配置通知和 Pub/Sub
//
// 策略管理：
//   - 创建和更新设备策略
//   - 使用默认策略模板
//   - 应用管理和限制
//
// 设备管理：
//   - 查询设备信息
//   - 远程控制设备（锁定、重置、重启等）
//   - 监控设备合规性
//
// 注册令牌：
//   - 创建注册令牌
//   - 生成 QR 码
//   - 管理令牌生命周期
//
// # 类型使用
//
// 所有核心类型都直接使用 androidmanagement 包的类型：
//
//	import (
//	    "amapi-pkg/pkgs/amapi"
//	    "google.golang.org/api/androidmanagement/v1"
//	)
//
//	var policy *androidmanagement.Policy
//	// 或者使用别名
//	var policy amapi.Policy  // 等同于 androidmanagement.Policy
//
// 使用辅助函数操作类型：
//
//	import "amapi-pkg/pkgs/amapi/types"
//
//	// 添加应用到策略
//	types.AddApplication(policy, &androidmanagement.ApplicationPolicy{
//	    PackageName: "com.example.app",
//	    InstallType: "REQUIRED",
//	})
//
// 更多信息请参考各子包的文档。
package amapi

import (
	"google.golang.org/api/androidmanagement/v1"

	"amapi-pkg/pkgs/amapi/client"
	"amapi-pkg/pkgs/amapi/config"
	"amapi-pkg/pkgs/amapi/types"
)

// Client 是 Android Management API 的主要客户端接口。
// 它提供了访问所有 API 功能的方法，包括企业、策略、设备、令牌等管理功能。
//
// 使用 NewClient 或 client.New 创建客户端实例。
// 客户端是线程安全的，可以在多个 goroutine 中共享使用。
//
// 示例：
//
//	cfg := &Config{
//	    ProjectID:       "your-project-id",
//	    CredentialsFile: "./sa-key.json",
//	}
//	c, err := NewClient(cfg)
//	if err != nil {
//	    log.Fatal(err)
//	}
//	defer c.Close()
type Client = client.Client

// Config 包含 Android Management API 客户端的配置选项。
//
// 配置可以通过以下方式提供：
//   - 使用 config.LoadFromEnv() 从环境变量加载
//   - 使用 config.LoadFromFile(path) 从 YAML/JSON 文件加载
//   - 使用 config.AutoLoadConfig() 自动检测并加载
//   - 手动构造 Config 结构
//
// 示例：
//
//	cfg := &Config{
//	    ProjectID:       "my-project",
//	    CredentialsFile: "./key.json",
//	    Timeout:         30 * time.Second,
//	    RetryAttempts:   3,
//	    EnableRetry:     true,
//	}
type Config = config.Config

// 核心类型定义
//
// 以下类型是 google.golang.org/api/androidmanagement/v1 包中原生类型的别名。
// 直接使用原生类型可以：
//   - 避免不必要的类型转换
//   - 提高代码效率和性能
//   - 保持与 Google SDK 的兼容性
//   - 简化维护工作
//
// 操作这些类型时，请使用 types 包中的辅助函数：
//   - types.AddApplication(policy, app)
//   - types.RemoveApplication(policy, packageName)
//   - types.ValidatePolicy(policy)
//   - types.IsEnrollmentTokenExpired(token)
//   - types.GenerateQRCodeData(token, options)
//
// 更多辅助函数请参考 types 包的文档。

// 类型别名 - 直接使用 androidmanagement 包的类型
type (
	// Enterprise 表示一个企业实体，是设备管理的顶层组织单位。
	// 每个企业可以包含多个策略、设备和注册令牌。
	Enterprise = androidmanagement.Enterprise

	// Policy 定义了设备的管理策略和限制。
	// 策略控制设备的功能、应用、网络设置等各个方面。
	Policy = androidmanagement.Policy

	// Device 表示一个已注册的 Android 设备。
	// 包含设备信息、状态、合规性等详细数据。
	Device = androidmanagement.Device

	// EnrollmentToken 是用于注册新设备的令牌。
	// 设备使用此令牌完成初始注册并应用指定的策略。
	EnrollmentToken = androidmanagement.EnrollmentToken

	// MigrationToken 用于从其他 EMM 系统迁移设备。
	// 允许将现有的托管设备迁移到 Android Management API。
	MigrationToken = androidmanagement.MigrationToken

	// WebApp 表示企业托管的 Web 应用。
	// 可以部署到托管设备上作为 Web 快捷方式。
	WebApp = androidmanagement.WebApp

	// WebToken 用于生成管理员访问企业管理界面的临时令牌。
	// 提供安全的 Web UI 访问权限。
	WebToken = androidmanagement.WebToken

	// ProvisioningInfo 包含设备配置和预配置信息。
	// 用于查询设备的配置状态和要求。
	ProvisioningInfo = androidmanagement.ProvisioningInfo

	// Command 表示可以发送到设备的命令。
	Command = androidmanagement.Command

	// APIError 表示 API 操作中发生的错误。
	// 提供详细的错误代码、消息和重试信息。
	APIError = types.Error
)

// NewClient 创建一个新的 Android Management API 客户端。
//
// 参数 cfg 包含客户端配置，包括项目 ID、认证凭证等。
// 返回的客户端是线程安全的，可以在多个 goroutine 中共享使用。
//
// 在使用完毕后应该调用 Close() 方法释放资源：
//
//	c, err := NewClient(cfg)
//	if err != nil {
//	    return err
//	}
//	defer c.Close()
//
// 如果配置无效或认证失败，将返回错误。
func NewClient(cfg *Config) (*Client, error) {
	return client.New(cfg)
}
