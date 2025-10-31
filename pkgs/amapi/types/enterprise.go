package types

import (
	"time"

	"google.golang.org/api/androidmanagement/v1"
)

// Enterprise is an alias for androidmanagement.Enterprise.
// Use androidmanagement.Enterprise directly for all enterprise operations.
type Enterprise = androidmanagement.Enterprise

// EnterpriseSignupURL represents a signup URL for enterprise creation.
type EnterpriseSignupURL struct {
	// URL is the signup URL
	URL string `json:"url"`

	// CallbackURL is the URL to redirect to after signup
	CallbackURL string `json:"callback_url,omitempty"`

	// ProjectID is the Google Cloud project ID
	ProjectID string `json:"project_id"`

	// CompletionToken is used to complete the signup process
	CompletionToken string `json:"completion_token,omitempty"`

	// ExpiresAt indicates when the signup URL expires
	ExpiresAt time.Time `json:"expires_at,omitempty"`

	// CreatedAt timestamp
	CreatedAt time.Time `json:"created_at"`
}

// SignupURLRequest represents a request to generate a signup URL.
type SignupURLRequest struct {
	// ProjectID is the Google Cloud project ID
	ProjectID string `json:"project_id"`

	// CallbackURL is the URL to redirect to after signup
	CallbackURL string `json:"callback_url,omitempty"`

	// AdminEmail is the email of the admin user
	AdminEmail string `json:"admin_email,omitempty"`

	// EnterpriseDisplayName is the display name for the enterprise
	EnterpriseDisplayName string `json:"enterprise_display_name,omitempty"`

	// Locale for the signup process
	Locale string `json:"locale,omitempty"`
}

// EnterpriseUpgradeURL represents an enterprise upgrade URL.
type EnterpriseUpgradeURL struct {
	// URL is the upgrade URL
	URL string `json:"url"`

	// EnterpriseName is the enterprise this URL is for
	EnterpriseName string `json:"enterprise_name"`

	// ProjectID is the Google Cloud project ID
	ProjectID string `json:"project_id"`

	// CreatedAt timestamp
	CreatedAt time.Time `json:"created_at"`

	// ExpiresAt indicates when the upgrade URL expires
	ExpiresAt time.Time `json:"expires_at,omitempty"`
}

// EnterpriseUpgradeURLRequest represents a request to generate an enterprise upgrade URL.
type EnterpriseUpgradeURLRequest struct {
	// EnterpriseName is the enterprise to generate the upgrade URL for
	EnterpriseName string `json:"enterprise_name"`

	// ProjectID is the Google Cloud project ID
	ProjectID string `json:"project_id"`

	// CallbackURL is the URL to redirect to after upgrade
	CallbackURL string `json:"callback_url,omitempty"`

	// AdminEmail is the email of the admin user
	AdminEmail string `json:"admin_email,omitempty"`

	// Locale for the upgrade process
	Locale string `json:"locale,omitempty"`
}

// EnterpriseCreateRequest represents a request to create an enterprise.
type EnterpriseCreateRequest struct {
	// SignupToken from the signup completion
	SignupToken string `json:"signup_token"`

	// ProjectID is the Google Cloud project ID
	ProjectID string `json:"project_id"`

	// EnterpriseToken from the signup process
	EnterpriseToken string `json:"enterprise_token,omitempty"`

	// DisplayName for the enterprise
	DisplayName string `json:"display_name,omitempty"`

	// ContactInfo for the enterprise
	ContactInfo *androidmanagement.ContactInfo `json:"contact_info,omitempty"`
}

// EnterpriseUpdateRequest represents a request to update an enterprise.
type EnterpriseUpdateRequest struct {
	// DisplayName is the new display name
	DisplayName string `json:"display_name,omitempty"`

	// PrimaryColor is the new primary color
	PrimaryColor *int64 `json:"primary_color,omitempty"`

	// Logo is the new logo
	Logo *androidmanagement.ExternalData `json:"logo,omitempty"`

	// ContactInfo is the new contact information
	ContactInfo *androidmanagement.ContactInfo `json:"contact_info,omitempty"`

	// EnabledNotificationTypes specifies notification types to enable
	EnabledNotificationTypes []string `json:"enabled_notification_types,omitempty"`

	// AppAutoApprovalEnabled controls app auto-approval
	AppAutoApprovalEnabled *bool `json:"app_auto_approval_enabled,omitempty"`

	// TermsAndConditions specifies terms and conditions
	TermsAndConditions []*androidmanagement.TermsAndConditions `json:"terms_and_conditions,omitempty"`
}

// EnterpriseListRequest represents a request to list enterprises.
type EnterpriseListRequest struct {
	ListOptions

	// ProjectID to filter by
	ProjectID string `json:"project_id,omitempty"`
}

// EnterpriseDeleteRequest represents a request to delete an enterprise.
type EnterpriseDeleteRequest struct {
	// Name is the enterprise resource name
	Name string `json:"name"`

	// Force deletion even if the enterprise has devices
	Force bool `json:"force,omitempty"`
}

// Notification type constants
const (
	NotificationTypeEnrollment          = "ENROLLMENT"
	NotificationTypeComplianceReport    = "COMPLIANCE_REPORT"
	NotificationTypeStatusReport        = "STATUS_REPORT"
	NotificationTypeCommand             = "COMMAND"
	NotificationTypeUsageLog            = "USAGE_LOG_ENABLED"
)