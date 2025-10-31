package types

import (
	"time"

	"google.golang.org/api/androidmanagement/v1"
)

// WebToken 相关类型和函数
//
// 注意：WebToken 类型直接使用 androidmanagement.WebToken。
// 此文件包含 Web 令牌相关的请求类型。
//
// 使用方式：
//
//	import (
//	    "amapi-pkg/pkgs/amapi/types"
//	    "google.golang.org/api/androidmanagement/v1"
//	)
//
//	// 创建 Web 令牌请求
//	req := &types.WebTokenCreateRequest{
//	    EnterpriseName: "enterprises/LC00abc123",
//	    Duration:       1 * time.Hour,
//	}

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

// WebToken helper functions (for androidmanagement.WebToken)

// GetWebTokenID extracts the token ID from the resource name.
func GetWebTokenID(token *androidmanagement.WebToken) string {
	if token == nil || token.Name == "" {
		return ""
	}

	// Extract ID from name format: enterprises/{enterpriseId}/webTokens/{tokenId}
	for i := len(token.Name) - 1; i >= 0; i-- {
		if token.Name[i] == '/' {
			return token.Name[i+1:]
		}
	}

	return token.Name
}

// GetWebTokenEnterpriseID extracts the enterprise ID from the token resource name.
func GetWebTokenEnterpriseID(token *androidmanagement.WebToken) string {
	if token == nil || token.Name == "" {
		return ""
	}

	// Extract from name format: enterprises/{enterpriseId}/webTokens/{tokenId}
	const prefix = "enterprises/"
	if len(token.Name) <= len(prefix) || token.Name[:len(prefix)] != prefix {
		return ""
	}

	remaining := token.Name[len(prefix):]
	for i, char := range remaining {
		if char == '/' {
			return remaining[:i]
		}
	}

	return ""
}

// Validate validates the web token create request.
func (req *WebTokenCreateRequest) Validate() error {
	if req.EnterpriseName == "" {
		return NewError(ErrCodeInvalidInput, "enterprise name is required")
	}

	return nil
}

// Note: Type conversion functions removed
// Use androidmanagement.WebToken directly instead of custom WebToken type
