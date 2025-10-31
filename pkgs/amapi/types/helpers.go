// Package types provides helper functions for androidmanagement types.
//
// 这个包提供了操作 google.golang.org/api/androidmanagement/v1 包中类型的辅助函数。
// 由于我们直接使用 androidmanagement 包的原生类型（而不是自定义类型），
// 所有的工具函数都是独立的函数，而不是类型的方法。
//
// # 使用方式
//
//	import "amapi-pkg/pkgs/amapi/types"
//
//	// 操作策略
//	hasApp := types.HasApplication(policy, "com.example.app")
//	types.AddApplication(policy, app)
//	types.RemoveApplication(policy, "com.example.app")
//	err := types.ValidatePolicy(policy)
//
//	// 提取资源 ID
//	enterpriseID := types.GetEnterpriseID(enterprise)
//	policyID := types.GetPolicyID(policy)
//	deviceID := types.GetDeviceID(device)
//
//	// 检查令牌状态
//	isExpired := types.IsEnrollmentTokenExpired(token)
//
//	// 生成 QR 码数据
//	qrData := types.GenerateQRCodeData(token, options)
//
// # 主要功能
//
//   - 资源 ID 提取：从资源名称中提取 ID
//   - 策略操作：添加/移除应用、验证策略等
//   - 令牌管理：检查过期、生成 QR 码等
//   - 设备状态：检查设备状态、合规性等
//
// 更多函数请参考具体的函数文档。
package types

import (
	"encoding/json"
	"time"

	"google.golang.org/api/androidmanagement/v1"
)

// Helper functions for androidmanagement.Enterprise

// GetEnterpriseID extracts the enterprise ID from the resource name.
func GetEnterpriseID(enterprise *androidmanagement.Enterprise) string {
	if enterprise == nil || enterprise.Name == "" {
		return ""
	}

	const prefix = "enterprises/"
	if len(enterprise.Name) > len(prefix) && enterprise.Name[:len(prefix)] == prefix {
		return enterprise.Name[len(prefix):]
	}

	return enterprise.Name
}

// Helper functions for androidmanagement.EnrollmentToken

// GetEnrollmentTokenID extracts the token ID from the resource name.
func GetEnrollmentTokenID(token *androidmanagement.EnrollmentToken) string {
	if token == nil || token.Name == "" {
		return ""
	}

	// Extract ID from name format: enterprises/{enterpriseId}/enrollmentTokens/{tokenId}
	for i := len(token.Name) - 1; i >= 0; i-- {
		if token.Name[i] == '/' {
			return token.Name[i+1:]
		}
	}

	return token.Name
}

// GetEnrollmentTokenEnterpriseID extracts the enterprise ID from the token resource name.
func GetEnrollmentTokenEnterpriseID(token *androidmanagement.EnrollmentToken) string {
	if token == nil || token.Name == "" {
		return ""
	}

	const prefix = "enterprises/"
	if len(token.Name) <= len(prefix) || token.Name[:len(prefix)] != prefix {
		return ""
	}

	remaining := token.Name[len(prefix):]
	for i, char := range remaining {
		if char == '/' {
			return remaining[:i]
		}
	}

	return ""
}

// IsEnrollmentTokenExpired checks if the enrollment token has expired.
func IsEnrollmentTokenExpired(token *androidmanagement.EnrollmentToken) bool {
	if token == nil || token.ExpirationTimestamp == "" {
		return false
	}

	expiration, err := time.Parse(time.RFC3339, token.ExpirationTimestamp)
	if err != nil {
		return false
	}

	return time.Now().After(expiration)
}

// GetEnrollmentTokenAllowPersonalUsageBool converts the string AllowPersonalUsage to bool.
func GetEnrollmentTokenAllowPersonalUsageBool(token *androidmanagement.EnrollmentToken) bool {
	if token == nil {
		return false
	}
	return token.AllowPersonalUsage == "PERSONAL_USAGE_ALLOWED"
}

// QRCodeOptions provides options for QR code generation.
type QRCodeOptions struct {
	WiFiSSID                  string                 `json:"wifi_ssid,omitempty"`
	WiFiPassword              string                 `json:"wifi_password,omitempty"`
	WiFiSecurityType          string                 `json:"wifi_security_type,omitempty"`
	WiFiHidden                bool                   `json:"wifi_hidden,omitempty"`
	TimeZone                  string                 `json:"time_zone,omitempty"`
	Locale                    string                 `json:"locale,omitempty"`
	SkipSetupWizard           bool                   `json:"skip_setup_wizard,omitempty"`
	LeaveAllSystemAppsEnabled bool                   `json:"leave_all_system_apps_enabled,omitempty"`
	AdminExtrasBundle         map[string]interface{} `json:"admin_extras_bundle,omitempty"`
}

// QRCodeData represents the data encoded in enrollment QR codes.
type QRCodeData struct {
	EnrollmentToken             string                 `json:"android.app.extra.PROVISIONING_DEVICE_ADMIN_COMPONENT_NAME,omitempty"`
	WiFiSSID                    string                 `json:"android.app.extra.PROVISIONING_WIFI_SSID,omitempty"`
	WiFiPassword                string                 `json:"android.app.extra.PROVISIONING_WIFI_PASSWORD,omitempty"`
	WiFiSecurityType            string                 `json:"android.app.extra.PROVISIONING_WIFI_SECURITY_TYPE,omitempty"`
	WiFiHidden                  bool                   `json:"android.app.extra.PROVISIONING_WIFI_HIDDEN,omitempty"`
	TimeZone                    string                 `json:"android.app.extra.PROVISIONING_TIME_ZONE,omitempty"`
	Locale                      string                 `json:"android.app.extra.PROVISIONING_LOCALE,omitempty"`
	SkipSetupWizard             bool                   `json:"android.app.extra.PROVISIONING_SKIP_SETUP_WIZARD,omitempty"`
	LeaveAllSystemAppsEnabled   bool                   `json:"android.app.extra.PROVISIONING_LEAVE_ALL_SYSTEM_APPS_ENABLED,omitempty"`
	AdminExtrasBundle           map[string]interface{} `json:"android.app.extra.PROVISIONING_ADMIN_EXTRAS_BUNDLE,omitempty"`
}

// GenerateQRCodeData generates QR code data for an enrollment token.
func GenerateQRCodeData(token *androidmanagement.EnrollmentToken, options *QRCodeOptions) *QRCodeData {
	data := &QRCodeData{
		EnrollmentToken: token.Value,
	}

	if options != nil {
		data.WiFiSSID = options.WiFiSSID
		data.WiFiPassword = options.WiFiPassword
		data.WiFiSecurityType = options.WiFiSecurityType
		data.WiFiHidden = options.WiFiHidden
		data.TimeZone = options.TimeZone
		data.Locale = options.Locale
		data.SkipSetupWizard = options.SkipSetupWizard
		data.LeaveAllSystemAppsEnabled = options.LeaveAllSystemAppsEnabled
		data.AdminExtrasBundle = options.AdminExtrasBundle
	}

	return data
}

// ToJSON converts QR code data to JSON string for encoding.
func (qr *QRCodeData) ToJSON() (string, error) {
	data, err := json.Marshal(qr)
	if err != nil {
		return "", err
	}
	return string(data), nil
}

