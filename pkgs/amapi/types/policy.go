package types

import (
	"google.golang.org/api/androidmanagement/v1"
)

// Policy helper functions (for androidmanagement.Policy)
//
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
