package types

import (
	"time"

	"google.golang.org/api/androidmanagement/v1"
)

// WebApp represents a web application in the Android Management API.
type WebApp struct {
	// Name is the resource name of the web app
	Name string `json:"name"`

	// DisplayName is the human-readable name of the web app
	DisplayName string `json:"display_name"`

	// StartURL is the URL where the web app starts
	StartURL string `json:"start_url"`

	// Icons is the list of icons for the web app
	Icons []*androidmanagement.WebAppIcon `json:"icons,omitempty"`

	// VersionCode is the version code of the web app
	VersionCode int64 `json:"version_code,omitempty"`

	// EnterpriseID is the enterprise this web app belongs to
	EnterpriseID string `json:"enterprise_id"`

	// CreatedAt is when the web app was created
	CreatedAt time.Time `json:"created_at"`

	// UpdatedAt is when the web app was last updated
	UpdatedAt time.Time `json:"updated_at"`

	// IsActive indicates if the web app is active
	IsActive bool `json:"is_active"`
}

// WebAppCreateRequest represents a request to create a web app.
type WebAppCreateRequest struct {
	// EnterpriseName is the enterprise to create the web app for
	EnterpriseName string `json:"enterprise_name"`

	// DisplayName is the human-readable name of the web app
	DisplayName string `json:"display_name"`

	// StartURL is the URL where the web app starts
	StartURL string `json:"start_url"`

	// Icons is the list of icons for the web app
	Icons []*androidmanagement.WebAppIcon `json:"icons,omitempty"`

	// VersionCode is the version code of the web app
	VersionCode int64 `json:"version_code,omitempty"`
}

// WebAppUpdateRequest represents a request to update a web app.
type WebAppUpdateRequest struct {
	// Name is the web app resource name
	Name string `json:"name"`

	// DisplayName is the human-readable name of the web app
	DisplayName string `json:"display_name,omitempty"`

	// StartURL is the URL where the web app starts
	StartURL string `json:"start_url,omitempty"`

	// Icons is the list of icons for the web app
	Icons []*androidmanagement.WebAppIcon `json:"icons,omitempty"`

	// VersionCode is the version code of the web app
	VersionCode int64 `json:"version_code,omitempty"`

	// UpdateMask specifies which fields to update
	UpdateMask []string `json:"update_mask,omitempty"`
}

// WebAppListRequest represents a request to list web apps.
type WebAppListRequest struct {
	ListOptions

	// EnterpriseName is the enterprise to list web apps for
	EnterpriseName string `json:"enterprise_name"`

	// ActiveOnly indicates whether to include only active web apps
	ActiveOnly bool `json:"active_only,omitempty"`
}

// WebAppDeleteRequest represents a request to delete a web app.
type WebAppDeleteRequest struct {
	// Name is the web app resource name
	Name string `json:"name"`
}

// WebApp helper methods

// GetID extracts the web app ID from the resource name.
func (wa *WebApp) GetID() string {
	if wa.Name == "" {
		return ""
	}

	// Extract ID from name format: enterprises/{enterpriseId}/webApps/{webAppId}
	for i := len(wa.Name) - 1; i >= 0; i-- {
		if wa.Name[i] == '/' {
			return wa.Name[i+1:]
		}
	}

	return wa.Name
}

// GetEnterpriseID extracts the enterprise ID from the web app resource name.
func (wa *WebApp) GetEnterpriseID() string {
	if wa.Name == "" {
		return ""
	}

	// Extract from name format: enterprises/{enterpriseId}/webApps/{webAppId}
	const prefix = "enterprises/"
	if len(wa.Name) <= len(prefix) || wa.Name[:len(prefix)] != prefix {
		return ""
	}

	remaining := wa.Name[len(prefix):]
	for i, char := range remaining {
		if char == '/' {
			return remaining[:i]
		}
	}

	return ""
}

// Validate validates the web app create request.
func (req *WebAppCreateRequest) Validate() error {
	if req.EnterpriseName == "" {
		return NewError(ErrCodeInvalidInput, "enterprise name is required")
	}

	if req.DisplayName == "" {
		return NewError(ErrCodeInvalidInput, "display name is required")
	}

	if req.StartURL == "" {
		return NewError(ErrCodeInvalidInput, "start URL is required")
	}

	return nil
}

// FromAMAPIWebApp converts an Android Management API web app to our type.
func FromAMAPIWebApp(webApp *androidmanagement.WebApp) *WebApp {
	if webApp == nil {
		return nil
	}

	wa := &WebApp{
		Name:        webApp.Name,
		StartURL:    webApp.StartUrl,
		Icons:       webApp.Icons,
		VersionCode: webApp.VersionCode,
		IsActive:    true, // Assume active if not specified
		CreatedAt:   time.Now(), // AMAPI doesn't provide creation time
		UpdatedAt:   time.Now(),
	}

	// Extract enterprise ID from name
	wa.EnterpriseID = wa.GetEnterpriseID()

	return wa
}

// ToAMAPIWebApp converts our web app to Android Management API format.
func (wa *WebApp) ToAMAPIWebApp() *androidmanagement.WebApp {
	if wa == nil {
		return nil
	}

	return &androidmanagement.WebApp{
		Name:        wa.Name,
		StartUrl:    wa.StartURL,
		Icons:       wa.Icons,
		VersionCode: wa.VersionCode,
	}
}
