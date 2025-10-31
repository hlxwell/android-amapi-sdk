// Package amapi provides helper functions that work directly with androidmanagement types.
package amapi

import (
	"strings"
	"time"

	"google.golang.org/api/androidmanagement/v1"

	"amapi-pkg/pkgs/amapi/types"
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

// GetEnrollmentTokenPolicyID extracts the policy ID from the policy name.
func GetEnrollmentTokenPolicyID(token *androidmanagement.EnrollmentToken) string {
	if token == nil || token.PolicyName == "" {
		return ""
	}

	// Extract ID from policy name format: enterprises/{enterpriseId}/policies/{policyId}
	for i := len(token.PolicyName) - 1; i >= 0; i-- {
		if token.PolicyName[i] == '/' {
			return token.PolicyName[i+1:]
		}
	}

	return token.PolicyName
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

// TimeUntilEnrollmentTokenExpiration returns the duration until the token expires.
func TimeUntilEnrollmentTokenExpiration(token *androidmanagement.EnrollmentToken) time.Duration {
	if token == nil || token.ExpirationTimestamp == "" {
		return 0
	}

	expiration, err := time.Parse(time.RFC3339, token.ExpirationTimestamp)
	if err != nil {
		return 0
	}

	if time.Now().After(expiration) {
		return 0
	}

	return time.Until(expiration)
}

// GetEnrollmentTokenAllowPersonalUsageBool converts the string AllowPersonalUsage to bool.
func GetEnrollmentTokenAllowPersonalUsageBool(token *androidmanagement.EnrollmentToken) bool {
	if token == nil {
		return false
	}
	return token.AllowPersonalUsage == "PERSONAL_USAGE_ALLOWED"
}

// SetEnrollmentTokenAllowPersonalUsage sets the AllowPersonalUsage field from bool.
func SetEnrollmentTokenAllowPersonalUsage(token *androidmanagement.EnrollmentToken, allow bool) {
	if token == nil {
		return
	}
	if allow {
		token.AllowPersonalUsage = "PERSONAL_USAGE_ALLOWED"
	} else {
		token.AllowPersonalUsage = "PERSONAL_USAGE_DISALLOWED"
	}
}

// Helper functions for androidmanagement.Policy

// GetPolicyID extracts the policy ID from the resource name.
//
// This is a convenience wrapper around types.ExtractResourceField.
func GetPolicyID(policy *androidmanagement.Policy) string {
	if policy == nil || policy.Name == "" {
		return ""
	}
	return types.ExtractResourceField(policy.Name, "PolicyID")
}

// GetPolicyEnterpriseID extracts the enterprise ID from the policy resource name.
//
// This is a convenience wrapper around types.ExtractResourceField.
func GetPolicyEnterpriseID(policy *androidmanagement.Policy) string {
	if policy == nil || policy.Name == "" {
		return ""
	}
	return types.ExtractResourceField(policy.Name, "EnterpriseID")
}

// Helper functions for androidmanagement.Device

// GetDeviceID extracts the device ID from the resource name.
func GetDeviceID(device *androidmanagement.Device) string {
	if device == nil || device.Name == "" {
		return ""
	}

	// Extract ID from device name format: enterprises/{enterpriseId}/devices/{deviceId}
	for i := len(device.Name) - 1; i >= 0; i-- {
		if device.Name[i] == '/' {
			return device.Name[i+1:]
		}
	}

	return device.Name
}

// GetDeviceEnterpriseID extracts the enterprise ID from the device resource name.
func GetDeviceEnterpriseID(device *androidmanagement.Device) string {
	if device == nil || device.Name == "" {
		return ""
	}

	const prefix = "enterprises/"
	if len(device.Name) <= len(prefix) || device.Name[:len(prefix)] != prefix {
		return ""
	}

	remaining := device.Name[len(prefix):]
	for i, char := range remaining {
		if char == '/' {
			return remaining[:i]
		}
	}

	return ""
}

// GetDeviceState returns the device state as a string constant.
func GetDeviceState(device *androidmanagement.Device) string {
	if device == nil {
		return ""
	}
	return device.State
}

// IsDeviceActive checks if the device is active.
func IsDeviceActive(device *androidmanagement.Device) bool {
	if device == nil {
		return false
	}
	return device.State == "ACTIVE"
}

// Helper functions for parsing resource names

// ParseResourceName parses a resource name and returns the type, enterprise ID, and resource ID.
// Example: "enterprises/E123/policies/P456" -> ("policies", "E123", "P456")
func ParseResourceName(name string) (resourceType, enterpriseID, resourceID string) {
	if name == "" {
		return "", "", ""
	}

	const prefix = "enterprises/"
	if len(name) <= len(prefix) || name[:len(prefix)] != prefix {
		return "", "", ""
	}

	parts := strings.Split(name[len(prefix):], "/")
	if len(parts) < 2 {
		return "", "", ""
	}

	enterpriseID = parts[0]
	resourceType = parts[1]
	if len(parts) > 2 {
		resourceID = parts[2]
	} else if len(parts) == 2 {
		resourceID = parts[1]
	}

	return resourceType, enterpriseID, resourceID
}

// QRCodeData represents the data encoded in enrollment QR codes.
// This is kept in helpers.go as it's a helper type for QR code generation.
type QRCodeData struct {
	EnrollmentToken           string                 `json:"android.app.extra.PROVISIONING_DEVICE_ADMIN_COMPONENT_NAME,omitempty"`
	WiFiSSID                  string                 `json:"android.app.extra.PROVISIONING_WIFI_SSID,omitempty"`
	WiFiPassword              string                 `json:"android.app.extra.PROVISIONING_WIFI_PASSWORD,omitempty"`
	WiFiSecurityType          string                 `json:"android.app.extra.PROVISIONING_WIFI_SECURITY_TYPE,omitempty"`
	WiFiHidden                bool                   `json:"android.app.extra.PROVISIONING_WIFI_HIDDEN,omitempty"`
	TimeZone                  string                 `json:"android.app.extra.PROVISIONING_TIME_ZONE,omitempty"`
	Locale                    string                 `json:"android.app.extra.PROVISIONING_LOCALE,omitempty"`
	SkipSetupWizard           bool                   `json:"android.app.extra.PROVISIONING_SKIP_SETUP_WIZARD,omitempty"`
	LeaveAllSystemAppsEnabled bool                   `json:"android.app.extra.PROVISIONING_LEAVE_ALL_SYSTEM_APPS_ENABLED,omitempty"`
	AdminExtrasBundle         map[string]interface{} `json:"android.app.extra.PROVISIONING_ADMIN_EXTRAS_BUNDLE,omitempty"`
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
