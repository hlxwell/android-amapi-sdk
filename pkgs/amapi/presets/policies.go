// Package presets provides default policy configuration.
package presets

import (
	"google.golang.org/api/androidmanagement/v1"
)

// GetDefaultPolicy returns a default policy with all status reporting enabled.
// This policy enables comprehensive device monitoring and data collection.
func GetDefaultPolicy() *androidmanagement.Policy {
	return &androidmanagement.Policy{
		// Enable all status reporting settings for comprehensive device monitoring
		StatusReportingSettings: &androidmanagement.StatusReportingSettings{
			ApplicationReportsEnabled:    true,
			CommonCriteriaModeEnabled:    true,
			DeviceSettingsEnabled:        true,
			DisplayInfoEnabled:           true,
			HardwareStatusEnabled:        true,
			MemoryInfoEnabled:            true,
			NetworkInfoEnabled:           true,
			PowerManagementEventsEnabled: true,
			SoftwareInfoEnabled:          true,
			SystemPropertiesEnabled:      true,
		},
		// Enable security logs in usage log
		UsageLog: &androidmanagement.UsageLog{
			EnabledLogTypes: []string{"SECURITY_LOGS"},
		},
	}
}
