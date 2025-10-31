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

// GetWebTokenID extracts the web token ID from the resource name.
//
// This is a convenience wrapper around ParseResourceNameStruct.
func GetWebTokenID(token *androidmanagement.WebToken) string {
	if token == nil || token.Name == "" {
		return ""
	}
	rn := ParseResourceNameStruct(token.Name)
	if rn == nil {
		return ""
	}
	return rn.WebTokenID
}

// GetWebTokenEnterpriseID extracts the enterprise ID from the token resource name.
//
// This is a convenience wrapper around ParseResourceNameStruct.
func GetWebTokenEnterpriseID(token *androidmanagement.WebToken) string {
	if token == nil || token.Name == "" {
		return ""
	}
	rn := ParseResourceNameStruct(token.Name)
	if rn == nil {
		return ""
	}
	return rn.EnterpriseID
}

// Note: Type conversion functions removed
// Use androidmanagement.WebToken directly instead of custom WebToken type
