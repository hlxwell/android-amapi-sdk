package types

import (
	"time"

	"google.golang.org/api/androidmanagement/v1"
)

// Enterprise represents an enterprise in the Android Management API.
type Enterprise struct {
	// Name is the resource name of the enterprise
	Name string `json:"name"`

	// DisplayName is the human-readable name of the enterprise
	DisplayName string `json:"display_name,omitempty"`

	// PrimaryColor is the primary color of the enterprise
	PrimaryColor int64 `json:"primary_color,omitempty"`

	// Logo is the enterprise logo
	Logo *androidmanagement.ExternalData `json:"logo,omitempty"`

	// PubsubTopic is the Cloud Pub/Sub topic for notifications
	PubsubTopic string `json:"pubsub_topic,omitempty"`

	// EnabledNotificationTypes specifies which notifications to enable
	EnabledNotificationTypes []string `json:"enabled_notification_types,omitempty"`

	// AppAutoApprovalEnabled indicates if app auto-approval is enabled
	AppAutoApprovalEnabled bool `json:"app_auto_approval_enabled,omitempty"`

	// ContactInfo provides contact information for the enterprise
	ContactInfo *ContactInfo `json:"contact_info,omitempty"`

	// TermsAndConditions specifies terms and conditions
	TermsAndConditions []*TermsAndConditions `json:"terms_and_conditions,omitempty"`

	// Created timestamp
	CreatedAt time.Time `json:"created_at,omitempty"`

	// Last updated timestamp
	UpdatedAt time.Time `json:"updated_at,omitempty"`
}

// ContactInfo provides contact information for the enterprise.
type ContactInfo struct {
	ContactEmail string `json:"contact_email,omitempty"`
	DataProtectionOfficerName string `json:"data_protection_officer_name,omitempty"`
	DataProtectionOfficerEmail string `json:"data_protection_officer_email,omitempty"`
	DataProtectionOfficerPhone string `json:"data_protection_officer_phone,omitempty"`
	EuRepresentativeName string `json:"eu_representative_name,omitempty"`
	EuRepresentativeEmail string `json:"eu_representative_email,omitempty"`
	EuRepresentativePhone string `json:"eu_representative_phone,omitempty"`
}

// TermsAndConditions represents terms and conditions for the enterprise.
type TermsAndConditions struct {
	Content *androidmanagement.UserFacingMessage `json:"content,omitempty"`
	Header *androidmanagement.UserFacingMessage `json:"header,omitempty"`
}

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
	ContactInfo *ContactInfo `json:"contact_info,omitempty"`
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
	ContactInfo *ContactInfo `json:"contact_info,omitempty"`

	// EnabledNotificationTypes specifies notification types to enable
	EnabledNotificationTypes []string `json:"enabled_notification_types,omitempty"`

	// AppAutoApprovalEnabled controls app auto-approval
	AppAutoApprovalEnabled *bool `json:"app_auto_approval_enabled,omitempty"`

	// TermsAndConditions specifies terms and conditions
	TermsAndConditions []*TermsAndConditions `json:"terms_and_conditions,omitempty"`
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

// Helper methods

// GetID extracts the enterprise ID from the resource name.
func (e *Enterprise) GetID() string {
	if e.Name == "" {
		return ""
	}

	// Extract ID from name format: enterprises/{enterpriseId}
	const prefix = "enterprises/"
	if len(e.Name) > len(prefix) && e.Name[:len(prefix)] == prefix {
		return e.Name[len(prefix):]
	}

	return e.Name
}

// GetResourceName returns the full resource name for the enterprise.
func (e *Enterprise) GetResourceName() string {
	if e.Name != "" {
		return e.Name
	}

	id := e.GetID()
	if id != "" {
		return "enterprises/" + id
	}

	return ""
}

// IsValid checks if the enterprise has required fields.
func (e *Enterprise) IsValid() bool {
	return e.Name != "" || e.GetID() != ""
}

// HasNotificationType checks if a notification type is enabled.
func (e *Enterprise) HasNotificationType(notificationType string) bool {
	for _, nt := range e.EnabledNotificationTypes {
		if nt == notificationType {
			return true
		}
	}
	return false
}

// AddNotificationType adds a notification type if not already present.
func (e *Enterprise) AddNotificationType(notificationType string) {
	if !e.HasNotificationType(notificationType) {
		e.EnabledNotificationTypes = append(e.EnabledNotificationTypes, notificationType)
	}
}

// RemoveNotificationType removes a notification type.
func (e *Enterprise) RemoveNotificationType(notificationType string) {
	for i, nt := range e.EnabledNotificationTypes {
		if nt == notificationType {
			e.EnabledNotificationTypes = append(
				e.EnabledNotificationTypes[:i],
				e.EnabledNotificationTypes[i+1:]...,
			)
			break
		}
	}
}

// Convert between types

// ToAMAPIEnterprise converts to androidmanagement.Enterprise.
func (e *Enterprise) ToAMAPIEnterprise() *androidmanagement.Enterprise {
	if e == nil {
		return nil
	}

	return &androidmanagement.Enterprise{
		Name:                     e.Name,
		// Note: DisplayName is not supported in the Google API struct
		PrimaryColor:             e.PrimaryColor,
		Logo:                     e.Logo,
		PubsubTopic:              e.PubsubTopic,
		EnabledNotificationTypes: e.EnabledNotificationTypes,
		AppAutoApprovalEnabled:   e.AppAutoApprovalEnabled,
		ContactInfo:              e.ContactInfo.ToAMAPIContactInfo(),
		TermsAndConditions:       e.termsToAMAPI(),
	}
}

// FromAMAPIEnterprise converts from androidmanagement.Enterprise.
func FromAMAPIEnterprise(e *androidmanagement.Enterprise) *Enterprise {
	if e == nil {
		return nil
	}

	return &Enterprise{
		Name:                     e.Name,
		DisplayName:              "", // DisplayName is not available in the Google API struct
		PrimaryColor:             e.PrimaryColor,
		Logo:                     e.Logo,
		PubsubTopic:              e.PubsubTopic,
		EnabledNotificationTypes: e.EnabledNotificationTypes,
		AppAutoApprovalEnabled:   e.AppAutoApprovalEnabled,
		ContactInfo:              FromAMAPIContactInfo(e.ContactInfo),
		TermsAndConditions:       fromAMAPITerms(e.TermsAndConditions),
		CreatedAt:                time.Now(), // API doesn't provide this
		UpdatedAt:                time.Now(), // API doesn't provide this
	}
}

// ToAMAPIContactInfo converts ContactInfo to androidmanagement format.
func (c *ContactInfo) ToAMAPIContactInfo() *androidmanagement.ContactInfo {
	if c == nil {
		return nil
	}

	return &androidmanagement.ContactInfo{
		ContactEmail:                   c.ContactEmail,
		DataProtectionOfficerEmail:     c.DataProtectionOfficerEmail,
		DataProtectionOfficerName:      c.DataProtectionOfficerName,
		DataProtectionOfficerPhone:     c.DataProtectionOfficerPhone,
		EuRepresentativeEmail:          c.EuRepresentativeEmail,
		EuRepresentativeName:           c.EuRepresentativeName,
		EuRepresentativePhone:          c.EuRepresentativePhone,
	}
}

// FromAMAPIContactInfo converts from androidmanagement ContactInfo.
func FromAMAPIContactInfo(c *androidmanagement.ContactInfo) *ContactInfo {
	if c == nil {
		return nil
	}

	return &ContactInfo{
		ContactEmail:                   c.ContactEmail,
		DataProtectionOfficerEmail:     c.DataProtectionOfficerEmail,
		DataProtectionOfficerName:      c.DataProtectionOfficerName,
		DataProtectionOfficerPhone:     c.DataProtectionOfficerPhone,
		EuRepresentativeEmail:          c.EuRepresentativeEmail,
		EuRepresentativeName:           c.EuRepresentativeName,
		EuRepresentativePhone:          c.EuRepresentativePhone,
	}
}

// Helper functions for terms and conditions conversion
func (e *Enterprise) termsToAMAPI() []*androidmanagement.TermsAndConditions {
	if e.TermsAndConditions == nil {
		return nil
	}

	result := make([]*androidmanagement.TermsAndConditions, len(e.TermsAndConditions))
	for i, tc := range e.TermsAndConditions {
		result[i] = &androidmanagement.TermsAndConditions{
			Content: tc.Content,
			Header:  tc.Header,
		}
	}
	return result
}

func fromAMAPITerms(terms []*androidmanagement.TermsAndConditions) []*TermsAndConditions {
	if terms == nil {
		return nil
	}

	result := make([]*TermsAndConditions, len(terms))
	for i, tc := range terms {
		result[i] = &TermsAndConditions{
			Content: tc.Content,
			Header:  tc.Header,
		}
	}
	return result
}