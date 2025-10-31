package types

import (
	"time"

	"google.golang.org/api/androidmanagement/v1"
)

// MigrationToken is an alias for androidmanagement.MigrationToken.
// Use androidmanagement.MigrationToken directly for all migration token operations.
type MigrationToken = androidmanagement.MigrationToken

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

