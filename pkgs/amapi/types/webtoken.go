package types

import (
	"time"

	"google.golang.org/api/androidmanagement/v1"
)

// WebToken represents a web token for accessing Google Play Enterprise Web UI.
type WebToken struct {
	// Name is the resource name of the web token
	Name string `json:"name"`

	// Value is the web token value
	Value string `json:"value"`

	// EnterpriseID is the enterprise this token belongs to
	EnterpriseID string `json:"enterprise_id"`

	// CreatedAt is when the token was created
	CreatedAt time.Time `json:"created_at"`

	// ExpiresAt is when the token expires
	ExpiresAt time.Time `json:"expires_at"`

	// IsActive indicates if the token is still active
	IsActive bool `json:"is_active"`

	// Permissions is the list of permissions granted to this token
	Permissions []string `json:"permissions,omitempty"`

	// LastUsedAt is when the token was last used
	LastUsedAt *time.Time `json:"last_used_at,omitempty"`
}

// WebTokenCreateRequest represents a request to create a web token.
type WebTokenCreateRequest struct {
	// EnterpriseName is the enterprise to create the token for
	EnterpriseName string `json:"enterprise_name"`

	// Duration specifies how long the token should be valid
	Duration time.Duration `json:"duration,omitempty"`

	// EnabledFeatures is the list of features to enable in the embedded UI.
	// If empty, all features are enabled by default.
	// Available features:
	//   - "PLAY_SEARCH": Managed Play search apps page
	//   - "PRIVATE_APPS": Private apps page
	//   - "WEB_APPS": Web apps page
	//   - "STORE_BUILDER": Organize apps page
	//   - "MANAGED_CONFIGURATIONS": Managed configurations page
	//   - "ZERO_TOUCH_CUSTOMER_MANAGEMENT": Zero-touch iframe
	EnabledFeatures []string `json:"enabled_features,omitempty"`

	// ParentFrameUrl is the URL of the parent frame that will host the iframe
	// This is required by the API. Default is "https://localhost" if not specified
	ParentFrameUrl string `json:"parent_frame_url,omitempty"`

	// Description is an optional description for the token
	Description string `json:"description,omitempty"`
}

// WebToken helper methods

// GetID extracts the token ID from the resource name.
func (wt *WebToken) GetID() string {
	if wt.Name == "" {
		return ""
	}

	// Extract ID from name format: enterprises/{enterpriseId}/webTokens/{tokenId}
	for i := len(wt.Name) - 1; i >= 0; i-- {
		if wt.Name[i] == '/' {
			return wt.Name[i+1:]
		}
	}

	return wt.Name
}

// GetEnterpriseID extracts the enterprise ID from the token resource name.
func (wt *WebToken) GetEnterpriseID() string {
	if wt.Name == "" {
		return ""
	}

	// Extract from name format: enterprises/{enterpriseId}/webTokens/{tokenId}
	const prefix = "enterprises/"
	if len(wt.Name) <= len(prefix) || wt.Name[:len(prefix)] != prefix {
		return ""
	}

	remaining := wt.Name[len(prefix):]
	for i, char := range remaining {
		if char == '/' {
			return remaining[:i]
		}
	}

	return ""
}

// IsExpired checks if the web token has expired.
func (wt *WebToken) IsExpired() bool {
	return time.Now().After(wt.ExpiresAt)
}

// Validate validates the web token create request.
func (req *WebTokenCreateRequest) Validate() error {
	if req.EnterpriseName == "" {
		return NewError(ErrCodeInvalidInput, "enterprise name is required")
	}

	return nil
}

// FromAMAPIWebToken converts an Android Management API web token to our type.
func FromAMAPIWebToken(token *androidmanagement.WebToken) *WebToken {
	if token == nil {
		return nil
	}

	webToken := &WebToken{
		Name:      token.Name,
		Value:     token.Value,
		IsActive:  true, // Assume active if not specified
		CreatedAt: time.Now(), // AMAPI doesn't provide creation time
	}

	// Extract enterprise ID from name
	webToken.EnterpriseID = webToken.GetEnterpriseID()

	return webToken
}

// ToAMAPIWebToken converts our web token to Android Management API format.
func (wt *WebToken) ToAMAPIWebToken() *androidmanagement.WebToken {
	if wt == nil {
		return nil
	}

	token := &androidmanagement.WebToken{
		Name:  wt.Name,
		Value: wt.Value,
	}

	return token
}
