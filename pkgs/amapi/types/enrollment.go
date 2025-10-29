package types

import (
	"encoding/json"
	"time"

	"google.golang.org/api/androidmanagement/v1"
)

// EnrollmentToken represents an enrollment token in the Android Management API.
type EnrollmentToken struct {
	// Name is the resource name of the enrollment token
	Name string `json:"name"`

	// PolicyName is the name of the policy associated with this token
	PolicyName string `json:"policy_name,omitempty"`

	// Value is the actual token value
	Value string `json:"value,omitempty"`

	// QrCode is the QR code data for enrollment
	QrCode string `json:"qr_code,omitempty"`

	// User is the user information for enrollment
	User *androidmanagement.User `json:"user,omitempty"`

	// AdditionalData contains additional enrollment data
	AdditionalData map[string]interface{} `json:"additional_data,omitempty"`

	// Duration specifies how long the token is valid
	Duration string `json:"duration,omitempty"`

	// AllowPersonalUsage indicates if personal usage is allowed (for work profile)
	AllowPersonalUsage bool `json:"allow_personal_usage,omitempty"`

	// OneTimeOnly indicates if the token can only be used once
	OneTimeOnly bool `json:"one_time_only,omitempty"`

	// ExpirationTimestamp when the token expires
	ExpirationTimestamp string `json:"expiration_timestamp,omitempty"`

	// Created timestamp (not from API, set locally)
	CreatedAt time.Time `json:"created_at,omitempty"`

	// Last updated timestamp (not from API, set locally)
	UpdatedAt time.Time `json:"updated_at,omitempty"`
}

// EnrollmentTokenCreateRequest represents a request to create an enrollment token.
type EnrollmentTokenCreateRequest struct {
	// EnterpriseName is the enterprise to create the token for
	EnterpriseName string `json:"enterprise_name"`

	// PolicyName is the policy to associate with the token
	PolicyName string `json:"policy_name"`

	// Duration specifies how long the token should be valid
	Duration time.Duration `json:"duration,omitempty"`

	// AllowPersonalUsage indicates if personal usage is allowed
	AllowPersonalUsage bool `json:"allow_personal_usage,omitempty"`

	// OneTimeOnly indicates if the token can only be used once
	OneTimeOnly bool `json:"one_time_only,omitempty"`

	// User information for the enrollment
	User *EnrollmentUser `json:"user,omitempty"`

	// AdditionalData contains additional enrollment data
	AdditionalData map[string]interface{} `json:"additional_data,omitempty"`
}

// EnrollmentTokenListRequest represents a request to list enrollment tokens.
type EnrollmentTokenListRequest struct {
	ListOptions

	// EnterpriseName is the enterprise to list tokens for
	EnterpriseName string `json:"enterprise_name"`

	// PolicyName filters tokens by policy
	PolicyName string `json:"policy_name,omitempty"`

	// IncludeExpired indicates whether to include expired tokens
	IncludeExpired bool `json:"include_expired,omitempty"`
}

// EnrollmentTokenGetRequest represents a request to get a specific enrollment token.
type EnrollmentTokenGetRequest struct {
	// Name is the enrollment token resource name
	Name string `json:"name"`
}

// EnrollmentTokenDeleteRequest represents a request to delete an enrollment token.
type EnrollmentTokenDeleteRequest struct {
	// Name is the enrollment token resource name
	Name string `json:"name"`
}

// EnrollmentUser represents user information for enrollment.
type EnrollmentUser struct {
	// AccountIdentifier is the account identifier
	AccountIdentifier string `json:"account_identifier,omitempty"`
}

// QRCodeData represents the data encoded in enrollment QR codes.
type QRCodeData struct {
	// EnrollmentToken is the token value
	EnrollmentToken string `json:"android.app.extra.PROVISIONING_DEVICE_ADMIN_COMPONENT_NAME,omitempty"`

	// WiFiSSID for WiFi configuration during enrollment
	WiFiSSID string `json:"android.app.extra.PROVISIONING_WIFI_SSID,omitempty"`

	// WiFiPassword for WiFi configuration during enrollment
	WiFiPassword string `json:"android.app.extra.PROVISIONING_WIFI_PASSWORD,omitempty"`

	// WiFiSecurityType for WiFi configuration
	WiFiSecurityType string `json:"android.app.extra.PROVISIONING_WIFI_SECURITY_TYPE,omitempty"`

	// WiFiHidden indicates if WiFi network is hidden
	WiFiHidden bool `json:"android.app.extra.PROVISIONING_WIFI_HIDDEN,omitempty"`

	// TimeZone for device configuration
	TimeZone string `json:"android.app.extra.PROVISIONING_TIME_ZONE,omitempty"`

	// Locale for device configuration
	Locale string `json:"android.app.extra.PROVISIONING_LOCALE,omitempty"`

	// SkipSetupWizard indicates whether to skip setup wizard
	SkipSetupWizard bool `json:"android.app.extra.PROVISIONING_SKIP_SETUP_WIZARD,omitempty"`

	// LeaveAllSystemAppsEnabled indicates whether to leave system apps enabled
	LeaveAllSystemAppsEnabled bool `json:"android.app.extra.PROVISIONING_LEAVE_ALL_SYSTEM_APPS_ENABLED,omitempty"`

	// AdminExtrasBundle contains additional admin configuration
	AdminExtrasBundle map[string]interface{} `json:"android.app.extra.PROVISIONING_ADMIN_EXTRAS_BUNDLE,omitempty"`
}

// EnrollmentToken helper methods

// GetID extracts the token ID from the resource name.
func (et *EnrollmentToken) GetID() string {
	if et.Name == "" {
		return ""
	}

	// Extract ID from name format: enterprises/{enterpriseId}/enrollmentTokens/{tokenId}
	for i := len(et.Name) - 1; i >= 0; i-- {
		if et.Name[i] == '/' {
			return et.Name[i+1:]
		}
	}

	return et.Name
}

// GetEnterpriseID extracts the enterprise ID from the token resource name.
func (et *EnrollmentToken) GetEnterpriseID() string {
	if et.Name == "" {
		return ""
	}

	// Extract from name format: enterprises/{enterpriseId}/enrollmentTokens/{tokenId}
	const prefix = "enterprises/"
	if len(et.Name) <= len(prefix) || et.Name[:len(prefix)] != prefix {
		return ""
	}

	remaining := et.Name[len(prefix):]
	for i, char := range remaining {
		if char == '/' {
			return remaining[:i]
		}
	}

	return ""
}

// IsExpired checks if the enrollment token has expired.
func (et *EnrollmentToken) IsExpired() bool {
	if et.ExpirationTimestamp == "" {
		return false
	}

	expiration, err := time.Parse(time.RFC3339, et.ExpirationTimestamp)
	if err != nil {
		return false
	}

	return time.Now().After(expiration)
}

// TimeUntilExpiration returns the duration until the token expires.
func (et *EnrollmentToken) TimeUntilExpiration() time.Duration {
	if et.ExpirationTimestamp == "" {
		return 0
	}

	expiration, err := time.Parse(time.RFC3339, et.ExpirationTimestamp)
	if err != nil {
		return 0
	}

	if time.Now().After(expiration) {
		return 0
	}

	return time.Until(expiration)
}

// GetPolicyID extracts the policy ID from the policy name.
func (et *EnrollmentToken) GetPolicyID() string {
	if et.PolicyName == "" {
		return ""
	}

	// Extract ID from policy name format: enterprises/{enterpriseId}/policies/{policyId}
	for i := len(et.PolicyName) - 1; i >= 0; i-- {
		if et.PolicyName[i] == '/' {
			return et.PolicyName[i+1:]
		}
	}

	return et.PolicyName
}

// GenerateQRCodeData generates QR code data for enrollment.
func (et *EnrollmentToken) GenerateQRCodeData(options *QRCodeOptions) *QRCodeData {
	data := &QRCodeData{
		EnrollmentToken: et.Value,
	}

	if options != nil {
		data.WiFiSSID = options.WiFiSSID
		data.WiFiPassword = options.WiFiPassword
		data.WiFiSecurityType = options.WiFiSecurityType
		data.WiFiHidden = options.WiFiHidden
		data.TimeZone = options.TimeZone
		data.Locale = options.Locale
		data.SkipSetupWizard = options.SkipSetupWizard
		data.LeaveAllSystemAppsEnabled = options.LeaveAllSystemAppsEnabled
		data.AdminExtrasBundle = options.AdminExtrasBundle
	}

	return data
}

// QRCodeOptions provides options for QR code generation.
type QRCodeOptions struct {
	WiFiSSID                      string                 `json:"wifi_ssid,omitempty"`
	WiFiPassword                  string                 `json:"wifi_password,omitempty"`
	WiFiSecurityType              string                 `json:"wifi_security_type,omitempty"`
	WiFiHidden                    bool                   `json:"wifi_hidden,omitempty"`
	TimeZone                      string                 `json:"time_zone,omitempty"`
	Locale                        string                 `json:"locale,omitempty"`
	SkipSetupWizard               bool                   `json:"skip_setup_wizard,omitempty"`
	LeaveAllSystemAppsEnabled     bool                   `json:"leave_all_system_apps_enabled,omitempty"`
	AdminExtrasBundle             map[string]interface{} `json:"admin_extras_bundle,omitempty"`
}

// WiFi security type constants
const (
	WiFiSecurityTypeNone = "NONE"
	WiFiSecurityTypeWEP  = "WEP"
	WiFiSecurityTypeWPA  = "WPA"
	WiFiSecurityTypeWPA2 = "WPA2"
)

// Default QR code options

// NewBasicQRCodeOptions creates basic QR code options.
func NewBasicQRCodeOptions() *QRCodeOptions {
	return &QRCodeOptions{
		SkipSetupWizard: true,
		Locale:          "en_US",
	}
}

// NewQRCodeOptionsWithWiFi creates QR code options with WiFi configuration.
func NewQRCodeOptionsWithWiFi(ssid, password, securityType string) *QRCodeOptions {
	options := NewBasicQRCodeOptions()
	options.WiFiSSID = ssid
	options.WiFiPassword = password
	options.WiFiSecurityType = securityType
	return options
}

// ToJSON converts QR code data to JSON string for encoding.
func (qr *QRCodeData) ToJSON() (string, error) {
	data, err := json.Marshal(qr)
	if err != nil {
		return "", err
	}
	return string(data), nil
}

// Convert to/from AMAPI types

// ToAMAPIEnrollmentToken converts to androidmanagement.EnrollmentToken.
func (et *EnrollmentToken) ToAMAPIEnrollmentToken() *androidmanagement.EnrollmentToken {
	if et == nil {
		return nil
	}

	// Convert bool to string for AllowPersonalUsage
	var allowPersonalUsage string
	if et.AllowPersonalUsage {
		allowPersonalUsage = "PERSONAL_USAGE_ALLOWED"
	} else {
		allowPersonalUsage = "PERSONAL_USAGE_DISALLOWED"
	}

	token := &androidmanagement.EnrollmentToken{
		Name:                et.Name,
		PolicyName:          et.PolicyName,
		Value:               et.Value,
		QrCode:              et.QrCode,
		User:                et.User,
		Duration:            et.Duration,
		AllowPersonalUsage:  allowPersonalUsage,
		OneTimeOnly:         et.OneTimeOnly,
		ExpirationTimestamp: et.ExpirationTimestamp,
	}

	// Convert additional data
	if et.AdditionalData != nil {
		// Note: AMAPI doesn't have AdditionalData field directly,
		// this would need to be handled in the admin extras bundle
	}

	return token
}

// FromAMAPIEnrollmentToken converts from androidmanagement.EnrollmentToken.
func FromAMAPIEnrollmentToken(token *androidmanagement.EnrollmentToken) *EnrollmentToken {
	if token == nil {
		return nil
	}

	// Convert string to bool for AllowPersonalUsage
	allowPersonalUsage := token.AllowPersonalUsage == "PERSONAL_USAGE_ALLOWED"

	return &EnrollmentToken{
		Name:                token.Name,
		PolicyName:          token.PolicyName,
		Value:               token.Value,
		QrCode:              token.QrCode,
		User:                token.User,
		Duration:            token.Duration,
		AllowPersonalUsage:  allowPersonalUsage,
		OneTimeOnly:         token.OneTimeOnly,
		ExpirationTimestamp: token.ExpirationTimestamp,
		CreatedAt:           time.Now(), // API doesn't provide this
		UpdatedAt:           time.Now(), // API doesn't provide this
	}
}

// ToAMAPIUser converts EnrollmentUser to androidmanagement.User.
func (eu *EnrollmentUser) ToAMAPIUser() *androidmanagement.User {
	if eu == nil {
		return nil
	}

	return &androidmanagement.User{
		AccountIdentifier: eu.AccountIdentifier,
	}
}

// FromAMAPIUser converts from androidmanagement.User.
func FromAMAPIUser(user *androidmanagement.User) *EnrollmentUser {
	if user == nil {
		return nil
	}

	return &EnrollmentUser{
		AccountIdentifier: user.AccountIdentifier,
	}
}

// Validation

// Validate checks if the enrollment token is valid.
func (et *EnrollmentToken) Validate() error {
	if et.Name == "" {
		return NewError(ErrCodeInvalidInput, "enrollment token name is required")
	}

	if et.PolicyName == "" {
		return NewError(ErrCodeInvalidInput, "policy name is required")
	}

	if et.IsExpired() {
		return NewError(ErrCodeInvalidInput, "enrollment token has expired")
	}

	return nil
}

// Validate checks if the enrollment token create request is valid.
func (req *EnrollmentTokenCreateRequest) Validate() error {
	if req.EnterpriseName == "" {
		return NewError(ErrCodeInvalidInput, "enterprise name is required")
	}

	if req.PolicyName == "" {
		return NewError(ErrCodeInvalidInput, "policy name is required")
	}

	if req.Duration < 0 {
		return NewError(ErrCodeInvalidInput, "duration cannot be negative")
	}

	// Default duration if not specified
	if req.Duration == 0 {
		req.Duration = 24 * time.Hour // Default to 24 hours
	}

	// Maximum duration check (API limit is typically 30 days)
	if req.Duration > 30*24*time.Hour {
		return NewError(ErrCodeInvalidInput, "duration cannot exceed 30 days")
	}

	return nil
}