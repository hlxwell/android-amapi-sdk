package types

import (
	"time"
)

// MigrationToken 相关类型和函数
//
// 注意：MigrationToken 类型直接使用 androidmanagement.MigrationToken。
// 此文件包含迁移令牌相关的请求类型。
//
// 使用方式：
//
//	import (
//	    "amapi-pkg/pkgs/amapi/types"
//	    "google.golang.org/api/androidmanagement/v1"
//	)
//
//	// 创建迁移令牌请求
//	req := &types.MigrationTokenCreateRequest{
//	    EnterpriseName: "enterprises/LC00abc123",
//	    PolicyName:     "enterprises/LC00abc123/policies/default",
//	    Duration:       24 * time.Hour,
//	}

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
