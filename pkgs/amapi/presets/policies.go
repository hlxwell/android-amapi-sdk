// Package presets provides pre-configured policy templates for common use cases.
package presets

import (
	"google.golang.org/api/androidmanagement/v1"

	"amapi-pkg/pkgs/amapi/types"
)

// PolicyPreset represents a predefined policy configuration.
type PolicyPreset struct {
	Name        string                      `json:"name"`
	DisplayName string                      `json:"display_name"`
	Description string                      `json:"description"`
	Mode        types.PolicyMode            `json:"mode"`
	Policy      *types.Policy               `json:"policy"`
	Tags        []string                    `json:"tags"`
}

// GetAllPresets returns all available policy presets.
func GetAllPresets() []*PolicyPreset {
	return []*PolicyPreset{
		GetFullyManagedPreset(),
		GetDedicatedDevicePreset(),
		GetWorkProfilePreset(),
		GetKioskModePreset(),
		GetCOPEPreset(),
		GetSecureWorkstationPreset(),
		GetEducationTabletPreset(),
		GetRetailKioskPreset(),
	}
}

// GetPresetByName returns a preset by its name.
func GetPresetByName(name string) *PolicyPreset {
	presets := GetAllPresets()
	for _, preset := range presets {
		if preset.Name == name {
			return preset
		}
	}
	return nil
}

// GetPresetsByMode returns all presets for a specific policy mode.
func GetPresetsByMode(mode types.PolicyMode) []*PolicyPreset {
	var result []*PolicyPreset
	presets := GetAllPresets()
	for _, preset := range presets {
		if preset.Mode == mode {
			result = append(result, preset)
		}
	}
	return result
}

// GetPresetsByTag returns all presets with a specific tag.
func GetPresetsByTag(tag string) []*PolicyPreset {
	var result []*PolicyPreset
	presets := GetAllPresets()
	for _, preset := range presets {
		for _, presetTag := range preset.Tags {
			if presetTag == tag {
				result = append(result, preset)
				break
			}
		}
	}
	return result
}

// GetFullyManagedPreset returns a preset for fully managed devices.
func GetFullyManagedPreset() *PolicyPreset {
	policy := &types.Policy{
		// Device restrictions
		AddUserDisabled:       true,
		UninstallAppsDisabled: true,
		StatusBarDisabled:     false,
		KeyguardDisabled:      false,

		// Security settings
		CameraDisabled:    false,
		BluetoothDisabled: false,

		// System settings
		AutoTimeRequired: true,

		// Status reporting
		StatusReportingSettings: &androidmanagement.StatusReportingSettings{
			ApplicationReportsEnabled:         true,
			DeviceSettingsEnabled:            true,
			DisplayInfoEnabled:               true,
			HardwareStatusEnabled:            true,
			MemoryInfoEnabled:               true,
			NetworkInfoEnabled:              true,
			PowerManagementEventsEnabled:     true,
			SoftwareInfoEnabled:             true,
			SystemPropertiesEnabled:         true,
		},

		// System update policy
		SystemUpdate: &androidmanagement.SystemUpdate{
			Type: "AUTOMATIC",
		},
	}

	return &PolicyPreset{
		Name:        "fully_managed",
		DisplayName: "Fully Managed Device",
		Description: "Standard corporate device policy with full management control. Suitable for company-owned devices.",
		Mode:        types.PolicyModeFullyManaged,
		Policy:      policy,
		Tags:        []string{"corporate", "standard", "secure"},
	}
}

// GetDedicatedDevicePreset returns a preset for dedicated devices (kiosk mode).
func GetDedicatedDevicePreset() *PolicyPreset {
	policy := &types.Policy{
		// Kiosk restrictions
		StatusBarDisabled:          true,
		KeyguardDisabled:          true,
		AddUserDisabled:           true,
		UninstallAppsDisabled:     true,
		AdjustVolumeDisabled:      true,
		TetheringConfigDisabled:   true,

		// Disable various features for kiosk
		BluetoothDisabled:                false, // May need Bluetooth for some kiosk apps
		CameraDisabled:                   false, // May need camera for some kiosk apps
		BluetoothConfigDisabled:          true,
		CellBroadcastsConfigDisabled:     true,

		// Block all apps by default - specific apps should be added
		BlockApplicationsEnabled: true,

		// System settings
		AutoTimeRequired: true,

		// Minimal status reporting for dedicated devices
		StatusReportingSettings: &androidmanagement.StatusReportingSettings{
			ApplicationReportsEnabled:    true,
			DeviceSettingsEnabled:       true,
			HardwareStatusEnabled:       true,
			SoftwareInfoEnabled:         true,
		},

		// Automatic system updates for stability
		SystemUpdate: &androidmanagement.SystemUpdate{
			Type: "AUTOMATIC",
		},
	}

	return &PolicyPreset{
		Name:        "dedicated_device",
		DisplayName: "Dedicated Device (Kiosk)",
		Description: "Locked-down device for single-purpose use. Ideal for digital signage, point-of-sale, or information kiosks.",
		Mode:        types.PolicyModeDedicated,
		Policy:      policy,
		Tags:        []string{"kiosk", "lockdown", "single-purpose", "retail"},
	}
}

// GetWorkProfilePreset returns a preset for work profile devices.
func GetWorkProfilePreset() *PolicyPreset {
	policy := &types.Policy{
		// Less restrictive for work profile
		AddUserDisabled:       false,
		UninstallAppsDisabled: false,
		StatusBarDisabled:     false,
		KeyguardDisabled:      false,

		// Allow more user control
		CameraDisabled:          false,
		BluetoothDisabled:       false,
		AdjustVolumeDisabled:    false,
		TetheringConfigDisabled: false,

		// Basic security
		AutoTimeRequired: true,

		// Status reporting for work profile
		StatusReportingSettings: &androidmanagement.StatusReportingSettings{
			ApplicationReportsEnabled: true,
			DeviceSettingsEnabled:    true,
			SoftwareInfoEnabled:      true,
		},

		// User-controlled system updates
		SystemUpdate: &androidmanagement.SystemUpdate{
			Type: "WINDOWED",
		},
	}

	return &PolicyPreset{
		Name:        "work_profile",
		DisplayName: "Work Profile (BYOD)",
		Description: "Minimal restrictions for personal devices with work profile. Balances security with user privacy.",
		Mode:        types.PolicyModeWorkProfile,
		Policy:      policy,
		Tags:        []string{"byod", "personal", "flexible"},
	}
}

// GetKioskModePreset returns a preset for kiosk mode with a single app.
func GetKioskModePreset() *PolicyPreset {
	policy := &types.Policy{
		// Maximum restrictions for kiosk
		StatusBarDisabled:               true,
		KeyguardDisabled:               true,
		AddUserDisabled:                true,
		UninstallAppsDisabled:          true,
		AdjustVolumeDisabled:           true,
		TetheringConfigDisabled:        true,
		BluetoothConfigDisabled:        true,
		CellBroadcastsConfigDisabled:   true,

		// Block everything by default
		BlockApplicationsEnabled: true,

		// Disable system features
		CameraDisabled:    true,
		BluetoothDisabled: true,

		// System settings
		AutoTimeRequired: true,

		// Minimal reporting
		StatusReportingSettings: &androidmanagement.StatusReportingSettings{
			ApplicationReportsEnabled: true,
			HardwareStatusEnabled:    true,
		},

		// Automatic updates
		SystemUpdate: &androidmanagement.SystemUpdate{
			Type: "AUTOMATIC",
		},

		// Applications will be added dynamically based on use case
		Applications: []*androidmanagement.ApplicationPolicy{},
	}

	return &PolicyPreset{
		Name:        "kiosk_mode",
		DisplayName: "Single App Kiosk",
		Description: "Extreme lockdown for single application use. Perfect for digital signage or interactive displays.",
		Mode:        types.PolicyModeDedicated,
		Policy:      policy,
		Tags:        []string{"kiosk", "single-app", "lockdown", "signage"},
	}
}

// GetCOPEPreset returns a preset for Corporate Owned, Personally Enabled devices.
func GetCOPEPreset() *PolicyPreset {
	policy := &types.Policy{
		// Balanced restrictions for COPE
		AddUserDisabled:       false,
		UninstallAppsDisabled: true, // Prevent removal of corporate apps
		StatusBarDisabled:     false,
		KeyguardDisabled:      false,

		// Allow personal use features
		CameraDisabled:          false,
		BluetoothDisabled:       false,
		AdjustVolumeDisabled:    false,
		TetheringConfigDisabled: false,

		// Security requirements
		AutoTimeRequired: true,

		// Comprehensive reporting for corporate compliance
		StatusReportingSettings: &androidmanagement.StatusReportingSettings{
			ApplicationReportsEnabled:     true,
			DeviceSettingsEnabled:        true,
			DisplayInfoEnabled:           true,
			HardwareStatusEnabled:        true,
			MemoryInfoEnabled:           true,
			NetworkInfoEnabled:          true,
			SoftwareInfoEnabled:         true,
		},

		// Controlled system updates
		SystemUpdate: &androidmanagement.SystemUpdate{
			Type: "WINDOWED",
		},
	}

	return &PolicyPreset{
		Name:        "cope",
		DisplayName: "Corporate Owned, Personally Enabled",
		Description: "Corporate-owned device with personal use allowed. Balances corporate control with user flexibility.",
		Mode:        types.PolicyModeFullyManaged,
		Policy:      policy,
		Tags:        []string{"cope", "corporate", "personal-use", "balanced"},
	}
}

// GetSecureWorkstationPreset returns a preset for secure workstation devices.
func GetSecureWorkstationPreset() *PolicyPreset {
	policy := &types.Policy{
		// Security-focused restrictions
		AddUserDisabled:                    true,
		UninstallAppsDisabled:              true,
		StatusBarDisabled:                  false,
		KeyguardDisabled:                   false,
		CameraDisabled:                     true, // Disable camera for security
		BluetoothDisabled:                  true, // Disable Bluetooth for security
		TetheringConfigDisabled:            true,
		BluetoothConfigDisabled:            true,
		BluetoothContactSharingDisabled:    true,
		CellBroadcastsConfigDisabled:       true,

		// System security
		AutoTimeRequired: true,

		// Enhanced reporting for security monitoring
		StatusReportingSettings: &androidmanagement.StatusReportingSettings{
			ApplicationReportsEnabled:         true,
			DeviceSettingsEnabled:            true,
			DisplayInfoEnabled:               true,
			HardwareStatusEnabled:            true,
			MemoryInfoEnabled:               true,
			NetworkInfoEnabled:              true,
			PowerManagementEventsEnabled:     true,
			SoftwareInfoEnabled:             true,
			SystemPropertiesEnabled:         true,
		},

		// Immediate security updates
		SystemUpdate: &androidmanagement.SystemUpdate{
			Type: "AUTOMATIC",
		},
	}

	return &PolicyPreset{
		Name:        "secure_workstation",
		DisplayName: "Secure Workstation",
		Description: "High-security configuration for sensitive work environments. Disables potential attack vectors.",
		Mode:        types.PolicyModeFullyManaged,
		Policy:      policy,
		Tags:        []string{"security", "enterprise", "restricted", "sensitive"},
	}
}

// GetEducationTabletPreset returns a preset for education tablets.
func GetEducationTabletPreset() *PolicyPreset {
	policy := &types.Policy{
		// Educational environment restrictions
		AddUserDisabled:       true,  // Prevent students from adding users
		UninstallAppsDisabled: true,  // Prevent removal of educational apps
		StatusBarDisabled:     false, // Allow access to notifications
		KeyguardDisabled:      false, // Keep lock screen for security

		// Allow educational features
		CameraDisabled:          false, // Camera may be needed for projects
		BluetoothDisabled:       false, // Bluetooth for accessories
		AdjustVolumeDisabled:    false, // Students need volume control
		TetheringConfigDisabled: true,  // Prevent network sharing

		// System settings
		AutoTimeRequired: true,

		// Educational reporting
		StatusReportingSettings: &androidmanagement.StatusReportingSettings{
			ApplicationReportsEnabled: true,
			DeviceSettingsEnabled:    true,
			HardwareStatusEnabled:    true,
			SoftwareInfoEnabled:      true,
		},

		// Controlled updates to maintain stability
		SystemUpdate: &androidmanagement.SystemUpdate{
			Type: "WINDOWED",
		},
	}

	return &PolicyPreset{
		Name:        "education_tablet",
		DisplayName: "Education Tablet",
		Description: "Optimized for educational environments with appropriate restrictions and feature access.",
		Mode:        types.PolicyModeFullyManaged,
		Policy:      policy,
		Tags:        []string{"education", "student", "tablet", "learning"},
	}
}

// GetRetailKioskPreset returns a preset for retail kiosk devices.
func GetRetailKioskPreset() *PolicyPreset {
	policy := &types.Policy{
		// Retail kiosk restrictions
		StatusBarDisabled:               true,
		KeyguardDisabled:               true,
		AddUserDisabled:                true,
		UninstallAppsDisabled:          true,
		AdjustVolumeDisabled:           false, // May need volume for customer interaction
		TetheringConfigDisabled:        true,
		BluetoothConfigDisabled:        true,
		CellBroadcastsConfigDisabled:   true,

		// Allow some features for retail use
		CameraDisabled:    false, // May need camera for QR codes
		BluetoothDisabled: false, // May need Bluetooth for payments

		// Block apps by default
		BlockApplicationsEnabled: true,

		// System settings
		AutoTimeRequired: true,

		// Retail-focused reporting
		StatusReportingSettings: &androidmanagement.StatusReportingSettings{
			ApplicationReportsEnabled: true,
			DeviceSettingsEnabled:    true,
			HardwareStatusEnabled:    true,
			SoftwareInfoEnabled:      true,
		},

		// Automatic updates for stability
		SystemUpdate: &androidmanagement.SystemUpdate{
			Type: "AUTOMATIC",
		},
	}

	return &PolicyPreset{
		Name:        "retail_kiosk",
		DisplayName: "Retail Kiosk",
		Description: "Optimized for retail point-of-sale and customer interaction with necessary features enabled.",
		Mode:        types.PolicyModeDedicated,
		Policy:      policy,
		Tags:        []string{"retail", "pos", "customer", "commerce"},
	}
}

// CreatePolicyFromPreset creates a new policy based on a preset.
func CreatePolicyFromPreset(presetName string, customizations map[string]interface{}) (*types.Policy, error) {
	preset := GetPresetByName(presetName)
	if preset == nil {
		return nil, types.NewError(types.ErrCodeNotFound, "preset not found: "+presetName)
	}

	// Clone the preset policy
	policy := preset.Policy.Clone()

	// Apply customizations
	if customizations != nil {
		if err := applyCustomizations(policy, customizations); err != nil {
			return nil, err
		}
	}

	return policy, nil
}

// applyCustomizations applies custom settings to a policy.
func applyCustomizations(policy *types.Policy, customizations map[string]interface{}) error {
	for key, value := range customizations {
		switch key {
		case "camera_disabled":
			if v, ok := value.(bool); ok {
				policy.CameraDisabled = v
			}
		case "bluetooth_disabled":
			if v, ok := value.(bool); ok {
				policy.BluetoothDisabled = v
			}
		case "status_bar_disabled":
			if v, ok := value.(bool); ok {
				policy.StatusBarDisabled = v
			}
		case "keyguard_disabled":
			if v, ok := value.(bool); ok {
				policy.KeyguardDisabled = v
			}
		case "add_user_disabled":
			if v, ok := value.(bool); ok {
				policy.AddUserDisabled = v
			}
		case "uninstall_apps_disabled":
			if v, ok := value.(bool); ok {
				policy.UninstallAppsDisabled = v
			}
		case "adjust_volume_disabled":
			if v, ok := value.(bool); ok {
				policy.AdjustVolumeDisabled = v
			}
		case "auto_time_required":
			if v, ok := value.(bool); ok {
				policy.AutoTimeRequired = v
			}
		case "applications":
			if v, ok := value.([]*androidmanagement.ApplicationPolicy); ok {
				policy.Applications = v
			}
		default:
			return types.NewErrorWithDetails(types.ErrCodeInvalidInput,
				"unknown customization key", "key: "+key)
		}
	}

	return nil
}

// GetPresetNames returns a list of all available preset names.
func GetPresetNames() []string {
	presets := GetAllPresets()
	names := make([]string, len(presets))
	for i, preset := range presets {
		names[i] = preset.Name
	}
	return names
}

// ValidatePreset validates a preset configuration.
func ValidatePreset(preset *PolicyPreset) error {
	if preset == nil {
		return types.NewError(types.ErrCodeInvalidInput, "preset is required")
	}

	if preset.Name == "" {
		return types.NewError(types.ErrCodeInvalidInput, "preset name is required")
	}

	if preset.Policy == nil {
		return types.NewError(types.ErrCodeInvalidInput, "preset policy is required")
	}

	return preset.Policy.Validate()
}