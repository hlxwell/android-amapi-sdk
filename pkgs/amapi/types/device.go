package types

import (
	"time"

	"google.golang.org/api/androidmanagement/v1"
)

// Device 相关类型和函数
//
// 注意：Device 类型直接使用 androidmanagement.Device，不再定义自定义类型。
// DeviceCommand 类型也直接使用 androidmanagement.Command。
//
// 使用方式：
//
//	import "amapi-pkg/pkgs/amapi/types"
//
//		// 列出设备直接传递参数
//	devices, err := client.Devices().List("enterprises/LC00abc123", 0, "", types.DeviceStateActive, nil, "")
//
//	// 设备命令使用 androidmanagement.Command
//	command := &androidmanagement.Command{
//	    Type: "LOCK",
//	    Duration: "3600s",
//	}

// Device helper functions (for androidmanagement.Device)

// GetDeviceID extracts the device ID from the resource name.
func GetDeviceID(device *androidmanagement.Device) string {
	if device == nil || device.Name == "" {
		return ""
	}

	// Extract ID from name format: enterprises/{enterpriseId}/devices/{deviceId}
	// Find the last '/' and return everything after it
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

	// Extract from name format: enterprises/{enterpriseId}/devices/{deviceId}
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

// IsDeviceOnline checks if the device is currently online based on last status report.
func IsDeviceOnline(device *androidmanagement.Device) bool {
	if device == nil || device.LastStatusReportTime == "" {
		return false
	}

	// Consider device online if last status report was within 5 minutes
	lastReport, err := time.Parse(time.RFC3339, device.LastStatusReportTime)
	if err != nil {
		return false
	}

	return time.Since(lastReport) < 5*time.Minute
}

// IsDeviceActive checks if the device is in an active state.
func IsDeviceActive(device *androidmanagement.Device) bool {
	if device == nil {
		return false
	}
	return device.State == string(DeviceStateActive)
}

// GetDeviceAndroidVersion returns the Android version string if available.
func GetDeviceAndroidVersion(device *androidmanagement.Device) string {
	if device != nil && device.SoftwareInfo != nil {
		return device.SoftwareInfo.AndroidVersion
	}
	return ""
}

// GetDeviceModel returns the device model if available.
func GetDeviceModel(device *androidmanagement.Device) string {
	if device != nil && device.HardwareInfo != nil {
		return device.HardwareInfo.Model
	}
	return ""
}

// GetDeviceManufacturer returns the device manufacturer if available.
func GetDeviceManufacturer(device *androidmanagement.Device) string {
	if device != nil && device.HardwareInfo != nil {
		return device.HardwareInfo.Manufacturer
	}
	return ""
}

// GetDeviceSecurityPatchLevel returns the security patch level if available.
func GetDeviceSecurityPatchLevel(device *androidmanagement.Device) string {
	if device != nil && device.SoftwareInfo != nil {
		return device.SoftwareInfo.SecurityPatchLevel
	}
	return ""
}

// DeviceHasApplication checks if a specific application is installed on the device.
func DeviceHasApplication(device *androidmanagement.Device, packageName string) bool {
	if device == nil || device.ApplicationReports == nil {
		return false
	}
	for _, report := range device.ApplicationReports {
		if report.PackageName == packageName {
			return true
		}
	}
	return false
}

// GetDeviceApplicationReport returns the application report for a specific package.
func GetDeviceApplicationReport(device *androidmanagement.Device, packageName string) *androidmanagement.ApplicationReport {
	if device == nil || device.ApplicationReports == nil {
		return nil
	}
	for _, report := range device.ApplicationReports {
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

// Command helper functions (for androidmanagement.Command)

// NewLockCommand creates a new lock command.
func NewLockCommand(duration time.Duration) *androidmanagement.Command {
	return &androidmanagement.Command{
		Type:     string(CommandTypeLock),
		Duration: duration.String(),
	}
}

// NewResetCommand creates a new factory reset command.
func NewResetCommand() *androidmanagement.Command {
	return &androidmanagement.Command{
		Type: string(CommandTypeReset),
	}
}

// NewRebootCommand creates a new reboot command.
func NewRebootCommand() *androidmanagement.Command {
	return &androidmanagement.Command{
		Type: string(CommandTypeReboot),
	}
}

// NewRemovePasswordCommand creates a new remove password command.
func NewRemovePasswordCommand() *androidmanagement.Command {
	return &androidmanagement.Command{
		Type: string(CommandTypeRemovePassword),
	}
}

// NewClearAppDataCommand creates a new clear app data command.
func NewClearAppDataCommand() *androidmanagement.Command {
	return &androidmanagement.Command{
		Type: string(CommandTypeClearAppData),
	}
}
