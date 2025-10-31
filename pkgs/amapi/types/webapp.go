package types

import (
	"google.golang.org/api/androidmanagement/v1"
)

// WebApp 相关类型和函数
//
// 注意：WebApp 类型直接使用 androidmanagement.WebApp。
// 此文件包含 Web 应用相关的辅助函数。
//
// 使用方式：
//
//	import (
//	    "amapi-pkg/pkgs/amapi/types"
//	    "google.golang.org/api/androidmanagement/v1"
//	)
//
//		// 创建 Web 应用直接传递参数
//	webApp, err := client.WebApps().Create(
//	    "enterprises/LC00abc123",
//	    "https://example.com",
//	    nil, // icons
//	    0,   // versionCode
//	)

// WebApp helper functions (for androidmanagement.WebApp)

// GetWebAppID extracts the web app ID from the resource name.
//
// This is a convenience wrapper around ParseResourceNameStruct.
func GetWebAppID(webApp *androidmanagement.WebApp) string {
	if webApp == nil || webApp.Name == "" {
		return ""
	}
	rn := ParseResourceNameStruct(webApp.Name)
	if rn == nil {
		return ""
	}
	return rn.WebAppID
}

// GetWebAppEnterpriseID extracts the enterprise ID from the web app resource name.
//
// This is a convenience wrapper around ParseResourceNameStruct.
func GetWebAppEnterpriseID(webApp *androidmanagement.WebApp) string {
	if webApp == nil || webApp.Name == "" {
		return ""
	}
	rn := ParseResourceNameStruct(webApp.Name)
	if rn == nil {
		return ""
	}
	return rn.EnterpriseID
}

// Note: Type conversion functions removed
// Use androidmanagement.WebApp directly instead of custom WebApp type
