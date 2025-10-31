package types

import (
	"google.golang.org/api/androidmanagement/v1"
)

// Policy 相关类型和函数
//
// 注意：Policy 类型直接使用 androidmanagement.Policy，不再定义自定义类型。
// 此文件包含请求类型和操作 Policy 的辅助函数。
//
// 使用方式：
//
//	import (
//	    "amapi-pkg/pkgs/amapi/types"
//	    "google.golang.org/api/androidmanagement/v1"
//	)
//
//	// 创建策略请求
//	req := &types.PolicyCreateRequest{
//	    EnterpriseName: "enterprises/LC00abc123",
//	    PolicyID:       "default-policy",
//	    Policy:         &androidmanagement.Policy{
//	        // ... 策略配置
//	    },
//	}
//
//	// 使用辅助函数操作策略
//	types.AddApplication(policy, app)
//	types.RemoveApplication(policy, "com.example.app")
//	err := types.ValidatePolicy(policy)

// PolicyCreateRequest represents a request to create a policy.
type PolicyCreateRequest struct {
	// EnterpriseName is the enterprise to create the policy for
	EnterpriseName string `json:"enterprise_name"`

	// PolicyID is the ID to assign to the new policy
	PolicyID string `json:"policy_id"`

	// Policy is the policy configuration
	Policy *androidmanagement.Policy `json:"policy"`
}

// PolicyUpdateRequest represents a request to update a policy.
type PolicyUpdateRequest struct {
	// Name is the policy resource name
	Name string `json:"name"`

	// Policy is the updated policy configuration
	Policy *androidmanagement.Policy `json:"policy"`

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
	Policy *androidmanagement.Policy `json:"policy"`

	// Tags for categorizing templates
	Tags []string `json:"tags,omitempty"`

	// Version of the template
	Version string `json:"version,omitempty"`
}

// Policy helper functions (for androidmanagement.Policy)

// HasApplication checks if a specific application is configured in the policy.
func HasApplication(p *androidmanagement.Policy, packageName string) bool {
	if p == nil || p.Applications == nil {
		return false
	}
	for _, app := range p.Applications {
		if app.PackageName == packageName {
			return true
		}
	}
	return false
}

// GetApplication returns the application policy for a specific package.
func GetApplication(p *androidmanagement.Policy, packageName string) *androidmanagement.ApplicationPolicy {
	if p == nil || p.Applications == nil {
		return nil
	}
	for _, app := range p.Applications {
		if app.PackageName == packageName {
			return app
		}
	}
	return nil
}

// AddApplication adds an application policy to the policy.
func AddApplication(p *androidmanagement.Policy, app *androidmanagement.ApplicationPolicy) {
	if p == nil || app == nil || app.PackageName == "" {
		return
	}

	// Remove existing application with same package name
	RemoveApplication(p, app.PackageName)

	// Add the new application
	if p.Applications == nil {
		p.Applications = []*androidmanagement.ApplicationPolicy{}
	}
	p.Applications = append(p.Applications, app)
}

// RemoveApplication removes an application policy from the policy.
func RemoveApplication(p *androidmanagement.Policy, packageName string) {
	if p == nil || p.Applications == nil {
		return
	}
	for i, app := range p.Applications {
		if app.PackageName == packageName {
			p.Applications = append(p.Applications[:i], p.Applications[i+1:]...)
			break
		}
	}
}

// ValidatePolicy validates a policy configuration.
func ValidatePolicy(p *androidmanagement.Policy) error {
	if p == nil {
		return NewError(ErrCodeInvalidInput, "policy is required")
	}

	// Validate applications
	if p.Applications != nil {
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
	}

	return nil
}

// Application installation type helpers

// NewRequiredApp creates an application policy for a required app.
func NewRequiredApp(packageName string) *androidmanagement.ApplicationPolicy {
	return &androidmanagement.ApplicationPolicy{
		PackageName:     packageName,
		InstallType:     string(InstallTypeRequired),
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
		PackageName:             packageName,
		InstallType:             string(InstallTypeKiosk),
		LockTaskAllowed:         true,
		DefaultPermissionPolicy: "GRANT",
	}
}
