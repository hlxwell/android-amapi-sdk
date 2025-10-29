package types

import (
	"time"

	"google.golang.org/api/androidmanagement/v1"
)

// MigrationToken represents a migration token for device migration.
type MigrationToken struct {
	// Name is the resource name of the migration token
	Name string `json:"name"`

	// Value is the migration token value
	Value string `json:"value"`

	// EnterpriseID is the enterprise this token belongs to
	EnterpriseID string `json:"enterprise_id"`

	// PolicyName is the policy to apply to migrated devices
	PolicyName string `json:"policy_name"`

	// CreatedAt is when the token was created
	CreatedAt time.Time `json:"created_at"`

	// ExpiresAt is when the token expires
	ExpiresAt time.Time `json:"expires_at"`

	// IsActive indicates if the token is still active
	IsActive bool `json:"is_active"`

	// DeviceCount is the number of devices that have used this token
	DeviceCount int `json:"device_count"`

	// LastUsedAt is when the token was last used
	LastUsedAt *time.Time `json:"last_used_at,omitempty"`
}

// MigrationTokenCreateRequest represents a request to create a migration token.
type MigrationTokenCreateRequest struct {
	// EnterpriseName is the enterprise to create the token for
	EnterpriseName string `json:"enterprise_name"`

	// PolicyName is the policy to apply to migrated devices
	PolicyName string `json:"policy_name"`

	// Duration specifies how long the token should be valid
	Duration time.Duration `json:"duration,omitempty"`

	// Description is an optional description for the token
	Description string `json:"description,omitempty"`
}

// MigrationTokenListRequest represents a request to list migration tokens.
type MigrationTokenListRequest struct {
	ListOptions

	// EnterpriseName is the enterprise to list tokens for
	EnterpriseName string `json:"enterprise_name"`

	// IncludeExpired indicates whether to include expired tokens
	IncludeExpired bool `json:"include_expired,omitempty"`

	// ActiveOnly indicates whether to include only active tokens
	ActiveOnly bool `json:"active_only,omitempty"`
}

// MigrationTokenDeleteRequest represents a request to delete a migration token.
type MigrationTokenDeleteRequest struct {
	// Name is the migration token resource name
	Name string `json:"name"`
}

// MigrationToken helper methods

// GetID extracts the token ID from the resource name.
func (mt *MigrationToken) GetID() string {
	if mt.Name == "" {
		return ""
	}

	// Extract ID from name format: enterprises/{enterpriseId}/migrationTokens/{tokenId}
	for i := len(mt.Name) - 1; i >= 0; i-- {
		if mt.Name[i] == '/' {
			return mt.Name[i+1:]
		}
	}

	return mt.Name
}

// GetEnterpriseID extracts the enterprise ID from the token resource name.
func (mt *MigrationToken) GetEnterpriseID() string {
	if mt.Name == "" {
		return ""
	}

	// Extract from name format: enterprises/{enterpriseId}/migrationTokens/{tokenId}
	const prefix = "enterprises/"
	if len(mt.Name) <= len(prefix) || mt.Name[:len(prefix)] != prefix {
		return ""
	}

	remaining := mt.Name[len(prefix):]
	for i, char := range remaining {
		if char == '/' {
			return remaining[:i]
		}
	}

	return ""
}

// IsExpired checks if the migration token has expired.
func (mt *MigrationToken) IsExpired() bool {
	return time.Now().After(mt.ExpiresAt)
}

// Validate validates the migration token create request.
func (req *MigrationTokenCreateRequest) Validate() error {
	if req.EnterpriseName == "" {
		return NewError(ErrCodeInvalidInput, "enterprise name is required")
	}

	if req.PolicyName == "" {
		return NewError(ErrCodeInvalidInput, "policy name is required")
	}

	return nil
}

// FromAMAPIMigrationToken converts an Android Management API migration token to our type.
func FromAMAPIMigrationToken(token *androidmanagement.MigrationToken) *MigrationToken {
	if token == nil {
		return nil
	}

	migrationToken := &MigrationToken{
		Name:      token.Name,
		Value:     token.Value,
		CreatedAt: time.Now(), // AMAPI doesn't provide creation time
		IsActive:  true,       // Assume active if not specified
	}

	// Extract enterprise ID from name
	migrationToken.EnterpriseID = migrationToken.GetEnterpriseID()

	return migrationToken
}

// ToAMAPIMigrationToken converts our migration token to Android Management API format.
func (mt *MigrationToken) ToAMAPIMigrationToken() *androidmanagement.MigrationToken {
	if mt == nil {
		return nil
	}

	token := &androidmanagement.MigrationToken{
		Name:  mt.Name,
		Value: mt.Value,
	}

	return token
}
