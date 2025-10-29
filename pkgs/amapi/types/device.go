package types

import (
	"time"

	"google.golang.org/api/androidmanagement/v1"
)

// Device represents a device in the Android Management API.
type Device struct {
	// Name is the resource name of the device
	Name string `json:"name"`

	// ApplicationReports provides information about installed applications
	ApplicationReports []*androidmanagement.ApplicationReport `json:"application_reports,omitempty"`

	// AppliedPolicyName is the name of the policy currently applied
	AppliedPolicyName string `json:"applied_policy_name,omitempty"`

	// AppliedPolicyVersion is the version of the applied policy
	AppliedPolicyVersion int64 `json:"applied_policy_version,omitempty"`

	// APILevel is the API level of the Android version running on the device
	APILevel int64 `json:"api_level,omitempty"`

	// EnrollmentTime is when the device was enrolled
	EnrollmentTime string `json:"enrollment_time,omitempty"`

	// LastPolicyComplianceReportTime is the last compliance report time
	LastPolicyComplianceReportTime string `json:"last_policy_compliance_report_time,omitempty"`

	// LastPolicySyncTime is the last policy sync time
	LastPolicySyncTime string `json:"last_policy_sync_time,omitempty"`

	// LastStatusReportTime is the last status report time
	LastStatusReportTime string `json:"last_status_report_time,omitempty"`

	// MemoryInfo provides information about device memory
	MemoryInfo *androidmanagement.MemoryInfo `json:"memory_info,omitempty"`

	// NetworkInfo provides information about network connectivity
	NetworkInfo *androidmanagement.NetworkInfo `json:"network_info,omitempty"`

	// PolicyCompliant indicates if the device is policy compliant
	PolicyCompliant bool `json:"policy_compliant,omitempty"`

	// SoftwareInfo provides information about device software
	SoftwareInfo *androidmanagement.SoftwareInfo `json:"software_info,omitempty"`

	// State is the current state of the device
	State DeviceState `json:"state,omitempty"`

	// UserName is the name of the user associated with the device
	UserName string `json:"user_name,omitempty"`

	// HardwareInfo provides information about device hardware
	HardwareInfo *androidmanagement.HardwareInfo `json:"hardware_info,omitempty"`

	// PowerManagementEvents contains power management related events
	PowerManagementEvents []*androidmanagement.PowerManagementEvent `json:"power_management_events,omitempty"`

	// PreviousDeviceNames contains previous names if the device was reset
	PreviousDeviceNames []string `json:"previous_device_names,omitempty"`

	// SecurityPosture provides information about security posture
	SecurityPosture *androidmanagement.SecurityPosture `json:"security_posture,omitempty"`

	// Created timestamp (not from API, set locally)
	CreatedAt time.Time `json:"created_at,omitempty"`

	// Last updated timestamp (not from API, set locally)
	UpdatedAt time.Time `json:"updated_at,omitempty"`
}

// DeviceCommand represents a command that can be issued to a device.
type DeviceCommand struct {
	// Type is the type of command
	Type CommandType `json:"type"`

	// ResetPasswordFlags for password reset commands
	ResetPasswordFlags []string `json:"reset_password_flags,omitempty"`

	// NewPassword for password reset commands
	NewPassword string `json:"new_password,omitempty"`

	// CreateTime when the command was created
	CreateTime string `json:"create_time,omitempty"`

	// Duration for temporary commands (like lock)
	Duration string `json:"duration,omitempty"`

	// UserName for user-specific commands
	UserName string `json:"user_name,omitempty"`

	// ErrorCode if the command failed
	ErrorCode string `json:"error_code,omitempty"`

	// FailureReason if the command failed
	FailureReason string `json:"failure_reason,omitempty"`
}

// DeviceListRequest represents a request to list devices.
type DeviceListRequest struct {
	ListOptions

	// EnterpriseName is the enterprise to list devices for
	EnterpriseName string `json:"enterprise_name"`

	// Filter by device state
	State DeviceState `json:"state,omitempty"`

	// Filter by policy compliance
	PolicyCompliant *bool `json:"policy_compliant,omitempty"`

	// Filter by user name
	UserName string `json:"user_name,omitempty"`
}

// DeviceGetRequest represents a request to get a specific device.
type DeviceGetRequest struct {
	// Name is the device resource name
	Name string `json:"name"`
}

// DeviceCommandRequest represents a request to issue a command to a device.
type DeviceCommandRequest struct {
	// DeviceName is the device resource name
	DeviceName string `json:"device_name"`

	// Command is the command to issue
	Command *DeviceCommand `json:"command"`
}

// DeviceDeleteRequest represents a request to delete a device.
type DeviceDeleteRequest struct {
	// Name is the device resource name
	Name string `json:"name"`

	// WipeDataFlags for wiping device data
	WipeDataFlags []string `json:"wipe_data_flags,omitempty"`

	// WipeExternalStorage indicates whether to wipe external storage
	WipeExternalStorage bool `json:"wipe_external_storage,omitempty"`
}

// Device helper methods

// GetID extracts the device ID from the resource name.
func (d *Device) GetID() string {
	if d.Name == "" {
		return ""
	}

	// Extract ID from name format: enterprises/{enterpriseId}/devices/{deviceId}
	// Find the last '/' and return everything after it
	for i := len(d.Name) - 1; i >= 0; i-- {
		if d.Name[i] == '/' {
			return d.Name[i+1:]
		}
	}

	return d.Name
}

// GetEnterpriseID extracts the enterprise ID from the device resource name.
func (d *Device) GetEnterpriseID() string {
	if d.Name == "" {
		return ""
	}

	// Extract from name format: enterprises/{enterpriseId}/devices/{deviceId}
	const prefix = "enterprises/"
	if len(d.Name) <= len(prefix) || d.Name[:len(prefix)] != prefix {
		return ""
	}

	remaining := d.Name[len(prefix):]
	for i, char := range remaining {
		if char == '/' {
			return remaining[:i]
		}
	}

	return ""
}

// IsOnline checks if the device is currently online based on last status report.
func (d *Device) IsOnline() bool {
	if d.LastStatusReportTime == "" {
		return false
	}

	// Consider device online if last status report was within 5 minutes
	lastReport, err := time.Parse(time.RFC3339, d.LastStatusReportTime)
	if err != nil {
		return false
	}

	return time.Since(lastReport) < 5*time.Minute
}

// IsActive checks if the device is in an active state.
func (d *Device) IsActive() bool {
	return d.State == DeviceStateActive
}

// GetAndroidVersion returns the Android version string if available.
func (d *Device) GetAndroidVersion() string {
	if d.SoftwareInfo != nil {
		return d.SoftwareInfo.AndroidVersion
	}
	return ""
}

// GetDeviceModel returns the device model if available.
func (d *Device) GetDeviceModel() string {
	if d.HardwareInfo != nil {
		return d.HardwareInfo.Model
	}
	return ""
}

// GetDeviceManufacturer returns the device manufacturer if available.
func (d *Device) GetDeviceManufacturer() string {
	if d.HardwareInfo != nil {
		return d.HardwareInfo.Manufacturer
	}
	return ""
}

// GetSecurityPatchLevel returns the security patch level if available.
func (d *Device) GetSecurityPatchLevel() string {
	if d.SoftwareInfo != nil {
		return d.SoftwareInfo.SecurityPatchLevel
	}
	return ""
}

// HasApplication checks if a specific application is installed.
func (d *Device) HasApplication(packageName string) bool {
	for _, report := range d.ApplicationReports {
		if report.PackageName == packageName {
			return true
		}
	}
	return false
}

// GetApplicationReport returns the application report for a specific package.
func (d *Device) GetApplicationReport(packageName string) *androidmanagement.ApplicationReport {
	for _, report := range d.ApplicationReports {
		if report.PackageName == packageName {
			return report
		}
	}
	return nil
}

// Command type constants
const (
	// ResetPasswordFlagRequireEntry requires password entry after reset
	ResetPasswordFlagRequireEntry = "REQUIRE_ENTRY"

	// ResetPasswordFlagDoNotAskCredentialsOnBoot doesn't ask for credentials on boot
	ResetPasswordFlagDoNotAskCredentialsOnBoot = "DO_NOT_ASK_CREDENTIALS_ON_BOOT"

	// WipeDataFlagExternalStorage wipes external storage
	WipeDataFlagExternalStorage = "WIPE_EXTERNAL_STORAGE"

	// WipeDataFlagPreserveResetProtectionData preserves reset protection data
	WipeDataFlagPreserveResetProtectionData = "PRESERVE_RESET_PROTECTION_DATA"
)

// DeviceCommand helper methods

// NewLockCommand creates a new lock command.
func NewLockCommand(duration time.Duration) *DeviceCommand {
	return &DeviceCommand{
		Type:     CommandTypeLock,
		Duration: duration.String(),
	}
}

// NewResetCommand creates a new factory reset command.
func NewResetCommand() *DeviceCommand {
	return &DeviceCommand{
		Type: CommandTypeReset,
	}
}

// NewRebootCommand creates a new reboot command.
func NewRebootCommand() *DeviceCommand {
	return &DeviceCommand{
		Type: CommandTypeReboot,
	}
}

// NewRemovePasswordCommand creates a new remove password command.
func NewRemovePasswordCommand() *DeviceCommand {
	return &DeviceCommand{
		Type: CommandTypeRemovePassword,
	}
}

// NewClearAppDataCommand creates a new clear app data command.
func NewClearAppDataCommand() *DeviceCommand {
	return &DeviceCommand{
		Type: CommandTypeClearAppData,
	}
}

// Convert to/from AMAPI types

// ToAMAPICommand converts to androidmanagement.Command.
func (dc *DeviceCommand) ToAMAPICommand() *androidmanagement.Command {
	if dc == nil {
		return nil
	}

	cmd := &androidmanagement.Command{
		Type:       string(dc.Type),
		CreateTime: dc.CreateTime,
		Duration:   dc.Duration,
		UserName:   dc.UserName,
		ErrorCode:  dc.ErrorCode,
	}

	// Set type-specific fields
	switch dc.Type {
	case CommandTypeRemovePassword:
		if len(dc.ResetPasswordFlags) > 0 {
			cmd.ResetPasswordFlags = dc.ResetPasswordFlags
		}
		if dc.NewPassword != "" {
			cmd.NewPassword = dc.NewPassword
		}
	}

	return cmd
}

// FromAMAPICommand converts from androidmanagement.Command.
func FromAMAPICommand(cmd *androidmanagement.Command) *DeviceCommand {
	if cmd == nil {
		return nil
	}

	return &DeviceCommand{
		Type:               CommandType(cmd.Type),
		ResetPasswordFlags: cmd.ResetPasswordFlags,
		NewPassword:        cmd.NewPassword,
		CreateTime:         cmd.CreateTime,
		Duration:           cmd.Duration,
		UserName:           cmd.UserName,
		ErrorCode:          cmd.ErrorCode,
	}
}