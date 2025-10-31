package types

import (
	"time"

	"google.golang.org/api/androidmanagement/v1"
)

// EnrollmentToken 相关类型和函数
//
// 注意：EnrollmentToken 类型直接使用 androidmanagement.EnrollmentToken。
// 此文件包含注册令牌相关的请求类型和辅助函数。
//
// 使用方式：
//
//	import "amapi-pkg/pkgs/amapi/types"
//
//	// 创建注册令牌请求
//	req := &types.EnrollmentTokenCreateRequest{
//	    EnterpriseName: "enterprises/LC00abc123",
//	    PolicyName:     "enterprises/LC00abc123/policies/default",
//	    Duration:       24 * time.Hour,
//	    User:           &androidmanagement.User{
//	        AccountIdentifier: "user@example.com",
//	    },
//	}
//
//	// 生成 QR 码数据
//	options := types.NewBasicQRCodeOptions()
//	qrData := types.GenerateQRCodeData(token, options)

// EnrollmentTokenCreateRequest represents a request to create an enrollment token.
type EnrollmentTokenCreateRequest struct {
	// EnterpriseName is the enterprise to create the token for
	EnterpriseName string `json:"enterprise_name"`

	// PolicyName is the policy to associate with the token
	PolicyName string `json:"policy_name"`

	// Duration specifies how long the token should be valid
	Duration time.Duration `json:"duration,omitempty"`

	// AllowPersonalUsage indicates if personal usage is allowed
	AllowPersonalUsage bool `json:"allow_personal_usage,omitempty"`

	// OneTimeOnly indicates if the token can only be used once
	OneTimeOnly bool `json:"one_time_only,omitempty"`

	// User information for the enrollment
	User *androidmanagement.User `json:"user,omitempty"`

	// AdditionalData contains additional enrollment data
	AdditionalData map[string]interface{} `json:"additional_data,omitempty"`
}

// EnrollmentTokenListRequest represents a request to list enrollment tokens.
type EnrollmentTokenListRequest struct {
	ListOptions

	// EnterpriseName is the enterprise to list tokens for
	EnterpriseName string `json:"enterprise_name"`

	// PolicyName filters tokens by policy
	PolicyName string `json:"policy_name,omitempty"`

	// IncludeExpired indicates whether to include expired tokens
	IncludeExpired bool `json:"include_expired,omitempty"`
}

// EnrollmentTokenGetRequest represents a request to get a specific enrollment token.
type EnrollmentTokenGetRequest struct {
	// Name is the enrollment token resource name
	Name string `json:"name"`
}

// EnrollmentTokenDeleteRequest represents a request to delete an enrollment token.
type EnrollmentTokenDeleteRequest struct {
	// Name is the enrollment token resource name
	Name string `json:"name"`
}

// WiFi security type constants
const (
	WiFiSecurityTypeNone = "NONE"
	WiFiSecurityTypeWEP  = "WEP"
	WiFiSecurityTypeWPA  = "WPA"
	WiFiSecurityTypeWPA2 = "WPA2"
)

// NewBasicQRCodeOptions creates basic QR code options.
func NewBasicQRCodeOptions() *QRCodeOptions {
	return &QRCodeOptions{
		SkipSetupWizard: true,
		Locale:          "en_US",
	}
}

// NewQRCodeOptionsWithWiFi creates QR code options with WiFi configuration.
func NewQRCodeOptionsWithWiFi(ssid, password, securityType string) *QRCodeOptions {
	options := NewBasicQRCodeOptions()
	options.WiFiSSID = ssid
	options.WiFiPassword = password
	options.WiFiSecurityType = securityType
	return options
}

// Validate checks if the enrollment token create request is valid.
func (req *EnrollmentTokenCreateRequest) Validate() error {
	if req.EnterpriseName == "" {
		return NewError(ErrCodeInvalidInput, "enterprise name is required")
	}

	if req.PolicyName == "" {
		return NewError(ErrCodeInvalidInput, "policy name is required")
	}

	if req.Duration < 0 {
		return NewError(ErrCodeInvalidInput, "duration cannot be negative")
	}

	// Default duration if not specified
	if req.Duration == 0 {
		req.Duration = 24 * time.Hour // Default to 24 hours
	}

	// Maximum duration check (API limit is typically 30 days)
	if req.Duration > 30*24*time.Hour {
		return NewError(ErrCodeInvalidInput, "duration cannot exceed 30 days")
	}

	return nil
}
