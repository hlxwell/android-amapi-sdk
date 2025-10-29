package amapi

// 这个文件用于导出所有子包的功能
// 实际的实现分散在各个子包中

import (
	"time"

	"amapi-pkg/pkgs/amapi/client"
	"amapi-pkg/pkgs/amapi/config"
	"amapi-pkg/pkgs/amapi/types"
)

// Client 是主要的API客户端
type Client = client.Client

// Config 是配置结构
type Config = config.Config

// 导出所有类型
type (
	// 企业相关类型
	Enterprise = types.Enterprise

	// 策略相关类型
	Policy = types.Policy

	// 设备相关类型
	Device = types.Device

	// 注册令牌相关类型
	EnrollmentToken = types.EnrollmentToken

	// 迁移令牌相关类型
	MigrationToken = types.MigrationToken

	// Web应用相关类型
	WebApp = types.WebApp

	// Web令牌相关类型
	WebToken = types.WebToken

	// 配置信息相关类型
	ProvisioningInfo = types.ProvisioningInfo

	// 错误类型
	APIError = types.Error
)

// 请求和响应类型（简化版本）
type CreateEnterpriseRequest struct {
	DisplayName  string `json:"display_name"`
	ProjectID    string `json:"project_id"`
	CallbackURL  string `json:"callback_url,omitempty"`
	PrimaryColor string `json:"primary_color,omitempty"`
}

type UpdateEnterpriseRequest struct {
	DisplayName  string `json:"display_name,omitempty"`
	PrimaryColor string `json:"primary_color,omitempty"`
}

type CreatePolicyRequest struct {
	Name string `json:"name"`
}

type UpdatePolicyRequest struct {
	Name                     string `json:"name,omitempty"`
	CameraDisabled           *bool  `json:"camera_disabled,omitempty"`
	KioskMode                *bool  `json:"kiosk_mode,omitempty"`
	BluetoothDisabled        *bool  `json:"bluetooth_disabled,omitempty"`
	WifiDisabled             *bool  `json:"wifi_disabled,omitempty"`
	UsbStorageDisabled       *bool  `json:"usb_storage_disabled,omitempty"`
	InstallUnknownSourcesDisabled *bool `json:"install_unknown_sources_disabled,omitempty"`
	DebuggingDisabled        *bool  `json:"debugging_disabled,omitempty"`
	ScreenCaptureDisabled    *bool  `json:"screen_capture_disabled,omitempty"`
	LocationDisabled         *bool  `json:"location_disabled,omitempty"`
	MicrophoneDisabled       *bool  `json:"microphone_disabled,omitempty"`
}

type ListDevicesRequest struct {
	PageSize int    `json:"page_size,omitempty"`
	Filter   string `json:"filter,omitempty"`
}

type CreateEnrollmentTokenRequest struct {
	Enterprise  string        `json:"enterprise"`
	Policy      string        `json:"policy"`
	Duration    time.Duration `json:"duration,omitempty"`
	WorkProfile bool          `json:"work_profile,omitempty"`
	OneTime     bool          `json:"one_time,omitempty"`
}

type ListEnrollmentTokensRequest struct {
	ActiveOnly bool `json:"active_only,omitempty"`
}

type QRCodeRequest struct {
	WifiSSID     string `json:"wifi_ssid,omitempty"`
	WifiPassword string `json:"wifi_password,omitempty"`
	WifiSecurity string `json:"wifi_security,omitempty"`
	SkipSetup    bool   `json:"skip_setup,omitempty"`
	Locale       string `json:"locale,omitempty"`
}

type StartLostModeRequest struct {
	Message     string `json:"message,omitempty"`
	PhoneNumber string `json:"phone_number,omitempty"`
}

// NewClient 创建新的API客户端
func NewClient(cfg *Config) (*Client, error) {
	return client.New(cfg)
}
