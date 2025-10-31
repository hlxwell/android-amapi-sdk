// Package types provides type definitions for the Android Management API client.
package types

import (
	"time"

	"google.golang.org/api/androidmanagement/v1"
)

// Result represents a generic operation result with optional error handling.
type Result[T any] struct {
	Data  *T     `json:"data,omitempty"`
	Error *Error `json:"error,omitempty"`
}

// ListResult represents a paginated list result.
type ListResult[T any] struct {
	Items         []T    `json:"items"`
	NextPageToken string `json:"next_page_token,omitempty"`
	TotalCount    int    `json:"total_count,omitempty"`
}

// OperationOptions provides common options for API operations.
type OperationOptions struct {
	// Timeout for the operation (overrides client default)
	Timeout *time.Duration `json:"timeout,omitempty"`

	// RetryCount for this specific operation (overrides client default)
	RetryCount *int `json:"retry_count,omitempty"`

	// Context values to pass through
	Context map[string]interface{} `json:"context,omitempty"`
}

// ListOptions provides common options for list operations.
type ListOptions struct {
	// PageSize specifies the maximum number of items to return
	PageSize int `json:"page_size,omitempty"`

	// PageToken for pagination
	PageToken string `json:"page_token,omitempty"`

	// Filter expression (if supported by the endpoint)
	Filter string `json:"filter,omitempty"`

	// OrderBy field name (if supported by the endpoint)
	OrderBy string `json:"order_by,omitempty"`
}

// ClientInfo provides information about the client and its capabilities.
type ClientInfo struct {
	Version     string    `json:"version"`
	ProjectID   string    `json:"project_id"`
	UserAgent   string    `json:"user_agent"`
	Capabilities []string `json:"capabilities"`
	CreatedAt   time.Time `json:"created_at"`
}

// CallbackData represents data structure for webhook callbacks.
type CallbackData struct {
	EnterpriseID   string                 `json:"enterprise_id"`
	CallbackURL    string                 `json:"callback_url"`
	CompletionTime time.Time              `json:"completion_time"`
	Data           map[string]interface{} `json:"data"`
}

// PolicyMode represents the different policy modes available.
type PolicyMode string

const (
	PolicyModeFullyManaged PolicyMode = "fully_managed"
	PolicyModeDedicated    PolicyMode = "dedicated"
	PolicyModeWorkProfile  PolicyMode = "work_profile"
)

// DeviceState represents the state of a device.
type DeviceState string

const (
	DeviceStateActive      DeviceState = "ACTIVE"
	DeviceStateDisabled    DeviceState = "DISABLED"
	DeviceStateDeleted     DeviceState = "DELETED"
	DeviceStateProvisioning DeviceState = "PROVISIONING"
)

// ApplicationInstallType represents how an application should be installed.
type ApplicationInstallType string

const (
	InstallTypeRequired           ApplicationInstallType = "REQUIRED"
	InstallTypePreinstalled      ApplicationInstallType = "PREINSTALLED"
	InstallTypeBlocked           ApplicationInstallType = "BLOCKED"
	InstallTypeAvailable         ApplicationInstallType = "AVAILABLE"
	InstallTypeRequiredForSetup  ApplicationInstallType = "REQUIRED_FOR_SETUP"
	InstallTypeKiosk             ApplicationInstallType = "KIOSK"
)

// CommandType represents the type of command that can be issued to a device.
type CommandType string

const (
	CommandTypeLock               CommandType = "LOCK"
	CommandTypeReset              CommandType = "RESET"
	CommandTypeReboot             CommandType = "REBOOT"
	CommandTypeRemovePassword     CommandType = "REMOVE_PASSWORD"
	CommandTypeClearAppData       CommandType = "CLEAR_APP_DATA"
	CommandTypeStartLostMode      CommandType = "START_LOST_MODE"
	CommandTypeStopLostMode       CommandType = "STOP_LOST_MODE"
)

// EnrollmentTokenType represents the type of enrollment token.
type EnrollmentTokenType string

const (
	EnrollmentTypeDefault     EnrollmentTokenType = "userlessDeviceProvisioning"
	EnrollmentTypeUserless    EnrollmentTokenType = "userlessDeviceProvisioning"
	EnrollmentTypePersonalWork EnrollmentTokenType = "personalWorkDeviceProvisioning"
)

// Convert between AMAPI types and our types

// Note: ToAMAPIPolicy and FromAMAPIPolicy removed - use androidmanagement.Policy directly

// ToAMAPIDevice converts our device representation to the AMAPI format.
func ToAMAPIDevice(device *Device) *androidmanagement.Device {
	if device == nil {
		return nil
	}

	return &androidmanagement.Device{
		Name:                    device.Name,
		ApplicationReports:      device.ApplicationReports,
		AppliedPolicyName:       device.AppliedPolicyName,
		AppliedPolicyVersion:    device.AppliedPolicyVersion,
		ApiLevel:                device.APILevel,
		EnrollmentTime:          device.EnrollmentTime,
		LastPolicyComplianceReportTime: device.LastPolicyComplianceReportTime,
		LastPolicySyncTime:      device.LastPolicySyncTime,
		LastStatusReportTime:    device.LastStatusReportTime,
		MemoryInfo:              device.MemoryInfo,
		NetworkInfo:             device.NetworkInfo,
		PolicyCompliant:         device.PolicyCompliant,
		SoftwareInfo:            device.SoftwareInfo,
		State:                   string(device.State),
		UserName:                device.UserName,
		// Add other fields as needed
	}
}

// FromAMAPIDevice converts AMAPI device to our representation.
func FromAMAPIDevice(device *androidmanagement.Device) *Device {
	if device == nil {
		return nil
	}

	return &Device{
		Name:                    device.Name,
		ApplicationReports:      device.ApplicationReports,
		AppliedPolicyName:       device.AppliedPolicyName,
		AppliedPolicyVersion:    device.AppliedPolicyVersion,
		APILevel:                device.ApiLevel,
		EnrollmentTime:          device.EnrollmentTime,
		LastPolicyComplianceReportTime: device.LastPolicyComplianceReportTime,
		LastPolicySyncTime:      device.LastPolicySyncTime,
		LastStatusReportTime:    device.LastStatusReportTime,
		MemoryInfo:              device.MemoryInfo,
		NetworkInfo:             device.NetworkInfo,
		PolicyCompliant:         device.PolicyCompliant,
		SoftwareInfo:            device.SoftwareInfo,
		State:                   DeviceState(device.State),
		UserName:                device.UserName,
		// Add other fields as needed
	}
}