package types

import (
	"time"

	"google.golang.org/api/androidmanagement/v1"
)

// Policy represents a policy in the Android Management API.
type Policy struct {
	// Name is the resource name of the policy
	Name string `json:"name"`

	// Version is the version of the policy
	Version int64 `json:"version,omitempty"`

	// Applications lists the applications and their installation policies
	Applications []*androidmanagement.ApplicationPolicy `json:"applications,omitempty"`

	// ComplianceRules defines compliance rules for the policy
	ComplianceRules []*androidmanagement.ComplianceRule `json:"compliance_rules,omitempty"`

	// KeyguardDisabled indicates if the keyguard is disabled
	KeyguardDisabled bool `json:"keyguard_disabled,omitempty"`

	// StatusBarDisabled indicates if the status bar is disabled
	StatusBarDisabled bool `json:"status_bar_disabled,omitempty"`

	// StatusReportingSettings configures status reporting
	StatusReportingSettings *androidmanagement.StatusReportingSettings `json:"status_reporting_settings,omitempty"`

	// SystemUpdate configures system update behavior
	SystemUpdate *androidmanagement.SystemUpdate `json:"system_update,omitempty"`

	// TetheringConfigDisabled indicates if tethering config is disabled
	TetheringConfigDisabled bool `json:"tethering_config_disabled,omitempty"`

	// UninstallAppsDisabled indicates if app uninstallation is disabled
	UninstallAppsDisabled bool `json:"uninstall_apps_disabled,omitempty"`

	// AccountTypesWithManagementDisabled lists account types that can't be managed
	AccountTypesWithManagementDisabled []string `json:"account_types_with_management_disabled,omitempty"`

	// AddUserDisabled indicates if adding users is disabled
	AddUserDisabled bool `json:"add_user_disabled,omitempty"`

	// AdjustVolumeDisabled indicates if volume adjustment is disabled
	AdjustVolumeDisabled bool `json:"adjust_volume_disabled,omitempty"`

	// AlwaysOnVpnPackage specifies the always-on VPN package
	AlwaysOnVpnPackage *androidmanagement.AlwaysOnVpnPackage `json:"always_on_vpn_package,omitempty"`

	// AutoTimeRequired indicates if auto time is required
	AutoTimeRequired bool `json:"auto_time_required,omitempty"`

	// BlockApplicationsEnabled indicates if application blocking is enabled
	BlockApplicationsEnabled bool `json:"block_applications_enabled,omitempty"`

	// BluetoothConfigDisabled indicates if Bluetooth config is disabled
	BluetoothConfigDisabled bool `json:"bluetooth_config_disabled,omitempty"`

	// BluetoothContactSharingDisabled indicates if Bluetooth contact sharing is disabled
	BluetoothContactSharingDisabled bool `json:"bluetooth_contact_sharing_disabled,omitempty"`

	// BluetoothDisabled indicates if Bluetooth is disabled
	BluetoothDisabled bool `json:"bluetooth_disabled,omitempty"`

	// CameraDisabled indicates if camera is disabled
	CameraDisabled bool `json:"camera_disabled,omitempty"`

	// CellBroadcastsConfigDisabled indicates if cell broadcast config is disabled
	CellBroadcastsConfigDisabled bool `json:"cell_broadcasts_config_disabled,omitempty"`

	// ChoosePrivateKeyRules defines rules for private key selection
	ChoosePrivateKeyRules []*androidmanagement.ChoosePrivateKeyRule `json:"choose_private_key_rules,omitempty"`

	// Created timestamp (not from API, set locally)
	CreatedAt time.Time `json:"created_at,omitempty"`

	// Last updated timestamp (not from API, set locally)
	UpdatedAt time.Time `json:"updated_at,omitempty"`
}

// PolicyCreateRequest represents a request to create a policy.
type PolicyCreateRequest struct {
	// EnterpriseName is the enterprise to create the policy for
	EnterpriseName string `json:"enterprise_name"`

	// PolicyID is the ID to assign to the new policy
	PolicyID string `json:"policy_id"`

	// Policy is the policy configuration
	Policy *Policy `json:"policy"`
}

// PolicyUpdateRequest represents a request to update a policy.
type PolicyUpdateRequest struct {
	// Name is the policy resource name
	Name string `json:"name"`

	// Policy is the updated policy configuration
	Policy *Policy `json:"policy"`

	// UpdateMask specifies which fields to update
	UpdateMask []string `json:"update_mask,omitempty"`
}

// PolicyGetRequest represents a request to get a specific policy.
type PolicyGetRequest struct {
	// Name is the policy resource name
	Name string `json:"name"`
}

// PolicyListRequest represents a request to list policies.
type PolicyListRequest struct {
	ListOptions

	// EnterpriseName is the enterprise to list policies for
	EnterpriseName string `json:"enterprise_name"`
}

// PolicyDeleteRequest represents a request to delete a policy.
type PolicyDeleteRequest struct {
	// Name is the policy resource name
	Name string `json:"name"`
}

// PolicyTemplate represents a policy template for common configurations.
type PolicyTemplate struct {
	// Name is the template name
	Name string `json:"name"`

	// DisplayName is the human-readable name
	DisplayName string `json:"display_name"`

	// Description describes the template
	Description string `json:"description"`

	// Mode is the policy mode this template is for
	Mode PolicyMode `json:"mode"`

	// Policy is the template policy configuration
	Policy *Policy `json:"policy"`

	// Tags for categorizing templates
	Tags []string `json:"tags,omitempty"`

	// Version of the template
	Version string `json:"version,omitempty"`
}

// Policy helper methods

// GetID extracts the policy ID from the resource name.
func (p *Policy) GetID() string {
	if p.Name == "" {
		return ""
	}

	// Extract ID from name format: enterprises/{enterpriseId}/policies/{policyId}
	for i := len(p.Name) - 1; i >= 0; i-- {
		if p.Name[i] == '/' {
			return p.Name[i+1:]
		}
	}

	return p.Name
}

// GetEnterpriseID extracts the enterprise ID from the policy resource name.
func (p *Policy) GetEnterpriseID() string {
	if p.Name == "" {
		return ""
	}

	// Extract from name format: enterprises/{enterpriseId}/policies/{policyId}
	const prefix = "enterprises/"
	if len(p.Name) <= len(prefix) || p.Name[:len(prefix)] != prefix {
		return ""
	}

	remaining := p.Name[len(prefix):]
	for i, char := range remaining {
		if char == '/' {
			return remaining[:i]
		}
	}

	return ""
}

// HasApplication checks if a specific application is configured in the policy.
func (p *Policy) HasApplication(packageName string) bool {
	for _, app := range p.Applications {
		if app.PackageName == packageName {
			return true
		}
	}
	return false
}

// GetApplication returns the application policy for a specific package.
func (p *Policy) GetApplication(packageName string) *androidmanagement.ApplicationPolicy {
	for _, app := range p.Applications {
		if app.PackageName == packageName {
			return app
		}
	}
	return nil
}

// AddApplication adds an application policy to the policy.
func (p *Policy) AddApplication(app *androidmanagement.ApplicationPolicy) {
	if app == nil || app.PackageName == "" {
		return
	}

	// Remove existing application with same package name
	p.RemoveApplication(app.PackageName)

	// Add the new application
	p.Applications = append(p.Applications, app)
}

// RemoveApplication removes an application policy from the policy.
func (p *Policy) RemoveApplication(packageName string) {
	for i, app := range p.Applications {
		if app.PackageName == packageName {
			p.Applications = append(p.Applications[:i], p.Applications[i+1:]...)
			break
		}
	}
}

// IsFullyManaged checks if this is a fully managed device policy.
func (p *Policy) IsFullyManaged() bool {
	// A fully managed policy typically has certain restrictions enabled
	return p.AddUserDisabled && p.UninstallAppsDisabled
}

// IsDedicated checks if this is a dedicated device policy.
func (p *Policy) IsDedicated() bool {
	// A dedicated device policy typically has strict kiosk-like restrictions
	return p.StatusBarDisabled && p.KeyguardDisabled
}

// IsWorkProfile checks if this is a work profile policy.
func (p *Policy) IsWorkProfile() bool {
	// Work profile policies are typically less restrictive
	return !p.AddUserDisabled && !p.StatusBarDisabled && !p.KeyguardDisabled
}

// GetPolicyMode determines the policy mode based on configuration.
func (p *Policy) GetPolicyMode() PolicyMode {
	if p.IsDedicated() {
		return PolicyModeDedicated
	}
	if p.IsWorkProfile() {
		return PolicyModeWorkProfile
	}
	return PolicyModeFullyManaged
}

// Clone creates a deep copy of the policy.
func (p *Policy) Clone() *Policy {
	if p == nil {
		return nil
	}

	// Create a new policy and copy basic fields
	clone := &Policy{
		Name:                               p.Name,
		Version:                            p.Version,
		KeyguardDisabled:                   p.KeyguardDisabled,
		StatusBarDisabled:                  p.StatusBarDisabled,
		TetheringConfigDisabled:            p.TetheringConfigDisabled,
		UninstallAppsDisabled:              p.UninstallAppsDisabled,
		AddUserDisabled:                    p.AddUserDisabled,
		AdjustVolumeDisabled:               p.AdjustVolumeDisabled,
		AutoTimeRequired:                   p.AutoTimeRequired,
		BlockApplicationsEnabled:           p.BlockApplicationsEnabled,
		BluetoothConfigDisabled:            p.BluetoothConfigDisabled,
		BluetoothContactSharingDisabled:    p.BluetoothContactSharingDisabled,
		BluetoothDisabled:                  p.BluetoothDisabled,
		CameraDisabled:                     p.CameraDisabled,
		CellBroadcastsConfigDisabled:       p.CellBroadcastsConfigDisabled,
		CreatedAt:                          p.CreatedAt,
		UpdatedAt:                          p.UpdatedAt,
	}

	// Deep copy slices
	if p.Applications != nil {
		clone.Applications = make([]*androidmanagement.ApplicationPolicy, len(p.Applications))
		copy(clone.Applications, p.Applications)
	}

	if p.ComplianceRules != nil {
		clone.ComplianceRules = make([]*androidmanagement.ComplianceRule, len(p.ComplianceRules))
		copy(clone.ComplianceRules, p.ComplianceRules)
	}

	if p.AccountTypesWithManagementDisabled != nil {
		clone.AccountTypesWithManagementDisabled = make([]string, len(p.AccountTypesWithManagementDisabled))
		copy(clone.AccountTypesWithManagementDisabled, p.AccountTypesWithManagementDisabled)
	}

	if p.ChoosePrivateKeyRules != nil {
		clone.ChoosePrivateKeyRules = make([]*androidmanagement.ChoosePrivateKeyRule, len(p.ChoosePrivateKeyRules))
		copy(clone.ChoosePrivateKeyRules, p.ChoosePrivateKeyRules)
	}

	// Copy pointers (these would need deep copying for complete isolation)
	clone.StatusReportingSettings = p.StatusReportingSettings
	clone.SystemUpdate = p.SystemUpdate
	clone.AlwaysOnVpnPackage = p.AlwaysOnVpnPackage

	return clone
}

// Validate checks if the policy configuration is valid.
func (p *Policy) Validate() error {
	if p.Name == "" {
		return NewError(ErrCodeInvalidInput, "policy name is required")
	}

	// Validate applications
	packageNames := make(map[string]bool)
	for _, app := range p.Applications {
		if app.PackageName == "" {
			return NewError(ErrCodeInvalidInput, "application package name cannot be empty")
		}

		if packageNames[app.PackageName] {
			return NewErrorWithDetails(ErrCodeInvalidInput, "duplicate application",
				"package name "+app.PackageName+" appears multiple times")
		}
		packageNames[app.PackageName] = true
	}

	return nil
}

// Application installation type helpers

// NewRequiredApp creates an application policy for a required app.
func NewRequiredApp(packageName string) *androidmanagement.ApplicationPolicy {
	return &androidmanagement.ApplicationPolicy{
		PackageName:   packageName,
		InstallType:   string(InstallTypeRequired),
		LockTaskAllowed: true,
	}
}

// NewPreinstalledApp creates an application policy for a preinstalled app.
func NewPreinstalledApp(packageName string) *androidmanagement.ApplicationPolicy {
	return &androidmanagement.ApplicationPolicy{
		PackageName: packageName,
		InstallType: string(InstallTypePreinstalled),
	}
}

// NewBlockedApp creates an application policy for a blocked app.
func NewBlockedApp(packageName string) *androidmanagement.ApplicationPolicy {
	return &androidmanagement.ApplicationPolicy{
		PackageName: packageName,
		InstallType: string(InstallTypeBlocked),
	}
}

// NewKioskApp creates an application policy for a kiosk app.
func NewKioskApp(packageName string) *androidmanagement.ApplicationPolicy {
	return &androidmanagement.ApplicationPolicy{
		PackageName:        packageName,
		InstallType:        string(InstallTypeKiosk),
		LockTaskAllowed:    true,
		DefaultPermissionPolicy: "GRANT",
	}
}