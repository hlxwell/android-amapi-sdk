// Package types exposes shared constants, helper structs, and utility types for the Android Management API client.
package types

import (
	"time"
)

// ListResult represents a paginated list result.
type ListResult[T any] struct {
	Items         []T    `json:"items"`
	NextPageToken string `json:"next_page_token,omitempty"`
	TotalCount    int    `json:"total_count,omitempty"`
}

// ClientInfo provides information about the client and its capabilities.
type ClientInfo struct {
	Version      string    `json:"version"`
	ProjectID    string    `json:"project_id"`
	UserAgent    string    `json:"user_agent"`
	Capabilities []string  `json:"capabilities"`
	CreatedAt    time.Time `json:"created_at"`
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
	DeviceStateActive       DeviceState = "ACTIVE"
	DeviceStateDisabled     DeviceState = "DISABLED"
	DeviceStateDeleted      DeviceState = "DELETED"
	DeviceStateProvisioning DeviceState = "PROVISIONING"
)

// ApplicationInstallType represents how an application should be installed.
type ApplicationInstallType string

const (
	InstallTypeRequired         ApplicationInstallType = "REQUIRED"
	InstallTypePreinstalled     ApplicationInstallType = "PREINSTALLED"
	InstallTypeBlocked          ApplicationInstallType = "BLOCKED"
	InstallTypeAvailable        ApplicationInstallType = "AVAILABLE"
	InstallTypeRequiredForSetup ApplicationInstallType = "REQUIRED_FOR_SETUP"
	InstallTypeKiosk            ApplicationInstallType = "KIOSK"
)

// CommandType represents the type of command that can be issued to a device.
type CommandType string

const (
	CommandTypeLock           CommandType = "LOCK"
	CommandTypeReset          CommandType = "RESET"
	CommandTypeReboot         CommandType = "REBOOT"
	CommandTypeRemovePassword CommandType = "REMOVE_PASSWORD"
	CommandTypeClearAppData   CommandType = "CLEAR_APP_DATA"
	CommandTypeStartLostMode  CommandType = "START_LOST_MODE"
	CommandTypeStopLostMode   CommandType = "STOP_LOST_MODE"
)

// EnrollmentTokenType represents the type of enrollment token.
type EnrollmentTokenType string

const (
	EnrollmentTypeDefault      EnrollmentTokenType = "userlessDeviceProvisioning"
	EnrollmentTypeUserless     EnrollmentTokenType = "userlessDeviceProvisioning"
	EnrollmentTypePersonalWork EnrollmentTokenType = "personalWorkDeviceProvisioning"
)
