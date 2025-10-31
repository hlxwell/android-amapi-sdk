package types

import (
	"google.golang.org/api/androidmanagement/v1"
)

// WebToken 相关类型和函数
//
// 注意：WebToken 类型直接使用 androidmanagement.WebToken。
// 此文件包含 Web 令牌相关的辅助函数。
//
// 使用方式：
//
//	import (
//	    "amapi-pkg/pkgs/amapi/types"
//	    "google.golang.org/api/androidmanagement/v1"
//	)
//
//		// 创建 Web 令牌直接传递参数
//	token, err := client.WebTokens().Create("enterprises/LC00abc123", "", nil)

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

// Note: Type conversion functions removed
// Use androidmanagement.WebToken directly instead of custom WebToken type
