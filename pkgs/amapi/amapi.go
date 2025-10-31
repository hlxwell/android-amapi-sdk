// Package amapi 提供了 Android Management API 的 Go 客户端库。
//
// 这个包是 Android Management API 的完整实现，提供了企业移动设备管理的所有功能，
// 包括企业管理、策略配置、设备控制、注册令牌管理等。
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
//	enterprises, err := c.Enterprises().List(nil)
//	if err != nil {
//	    log.Fatal(err)
//	}
//
// # 配置
//
// 支持多种配置方式：
//   - 环境变量
//   - YAML/JSON 配置文件
//   - 程序化配置
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
//   - 使用预设策略模板
//   - 应用管理和限制
//
// 设备管理：
//   - 查询设备信息
//   - 远程控制设备
//   - 监控设备合规性
//
// 注册令牌：
//   - 创建注册令牌
//   - 生成 QR 码
//   - 管理令牌生命周期
//
// 更多信息请参考各子包的文档。
package amapi

import (
	"time"

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
// 直接使用 google.golang.org/api/androidmanagement/v1 包中的类型，
// 避免不必要的类型转换，提高代码效率和可维护性。
//
// 使用 helpers.go 中的辅助函数来操作这些类型。

// 类型别名，直接使用 androidmanagement 包的类型
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

// 请求和响应类型
//
// 以下类型用于构造 API 请求。

// CreateEnterpriseRequest 包含创建企业所需的参数。
type CreateEnterpriseRequest struct {
	DisplayName  string `json:"display_name"`            // 企业显示名称
	ProjectID    string `json:"project_id"`              // Google Cloud 项目 ID
	CallbackURL  string `json:"callback_url,omitempty"`  // 注册完成后的回调 URL
	PrimaryColor string `json:"primary_color,omitempty"` // 主题颜色（十六进制格式）
}

// UpdateEnterpriseRequest 包含更新企业信息的参数。
type UpdateEnterpriseRequest struct {
	DisplayName  string `json:"display_name,omitempty"`  // 企业显示名称
	PrimaryColor string `json:"primary_color,omitempty"` // 主题颜色（十六进制格式）
}

// CreatePolicyRequest 包含创建策略所需的参数。
type CreatePolicyRequest struct {
	Name string `json:"name"` // 策略名称/ID
}

// UpdatePolicyRequest 包含更新策略的参数。
// 所有字段都是可选的，只更新提供的字段。
type UpdatePolicyRequest struct {
	Name                          string `json:"name,omitempty"`                             // 策略名称
	CameraDisabled                *bool  `json:"camera_disabled,omitempty"`                  // 是否禁用摄像头
	KioskMode                     *bool  `json:"kiosk_mode,omitempty"`                       // 是否启用 Kiosk 模式
	BluetoothDisabled             *bool  `json:"bluetooth_disabled,omitempty"`               // 是否禁用蓝牙
	WifiDisabled                  *bool  `json:"wifi_disabled,omitempty"`                    // 是否禁用 WiFi
	UsbStorageDisabled            *bool  `json:"usb_storage_disabled,omitempty"`             // 是否禁用 USB 存储
	InstallUnknownSourcesDisabled *bool  `json:"install_unknown_sources_disabled,omitempty"` // 是否禁止安装未知来源应用
	DebuggingDisabled             *bool  `json:"debugging_disabled,omitempty"`               // 是否禁用调试功能
	ScreenCaptureDisabled         *bool  `json:"screen_capture_disabled,omitempty"`          // 是否禁用屏幕截图
	LocationDisabled              *bool  `json:"location_disabled,omitempty"`                // 是否禁用位置服务
	MicrophoneDisabled            *bool  `json:"microphone_disabled,omitempty"`              // 是否禁用麦克风
}

// ListDevicesRequest 包含列出设备的查询参数。
type ListDevicesRequest struct {
	PageSize int    `json:"page_size,omitempty"` // 每页返回的设备数量
	Filter   string `json:"filter,omitempty"`    // 过滤条件（遵循 Google API 过滤语法）
}

// CreateEnrollmentTokenRequest 包含创建注册令牌的参数。
type CreateEnrollmentTokenRequest struct {
	Enterprise  string        `json:"enterprise"`            // 企业名称（格式：enterprises/{enterpriseId}）
	Policy      string        `json:"policy"`                // 策略名称（格式：enterprises/{enterpriseId}/policies/{policyId}）
	Duration    time.Duration `json:"duration,omitempty"`    // 令牌有效期（默认 1 小时）
	WorkProfile bool          `json:"work_profile,omitempty"` // 是否使用工作配置文件模式
	OneTime     bool          `json:"one_time,omitempty"`     // 是否为一次性令牌
}

// ListEnrollmentTokensRequest 包含列出注册令牌的查询参数。
type ListEnrollmentTokensRequest struct {
	ActiveOnly bool `json:"active_only,omitempty"` // 是否只返回活动的令牌
}

// QRCodeRequest 包含生成 QR 码的参数。
type QRCodeRequest struct {
	WifiSSID     string `json:"wifi_ssid,omitempty"`     // WiFi 网络 SSID
	WifiPassword string `json:"wifi_password,omitempty"` // WiFi 密码
	WifiSecurity string `json:"wifi_security,omitempty"` // WiFi 安全类型（WPA2/WPA/开放）
	SkipSetup    bool   `json:"skip_setup,omitempty"`    // 是否跳过设置向导
	Locale       string `json:"locale,omitempty"`        // 语言区域设置（如 zh_CN）
}

// StartLostModeRequest 包含启动丢失模式的参数。
type StartLostModeRequest struct {
	Message     string `json:"message,omitempty"`      // 显示在锁屏上的消息
	PhoneNumber string `json:"phone_number,omitempty"` // 联系电话号码
}

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
