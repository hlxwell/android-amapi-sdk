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
func GetWebAppID(webApp *androidmanagement.WebApp) string {
	if webApp == nil || webApp.Name == "" {
		return ""
	}

	// Extract ID from name format: enterprises/{enterpriseId}/webApps/{webAppId}
	for i := len(webApp.Name) - 1; i >= 0; i-- {
		if webApp.Name[i] == '/' {
			return webApp.Name[i+1:]
		}
	}

	return webApp.Name
}

// GetWebAppEnterpriseID extracts the enterprise ID from the web app resource name.
func GetWebAppEnterpriseID(webApp *androidmanagement.WebApp) string {
	if webApp == nil || webApp.Name == "" {
		return ""
	}

	// Extract from name format: enterprises/{enterpriseId}/webApps/{webAppId}
	const prefix = "enterprises/"
	if len(webApp.Name) <= len(prefix) || webApp.Name[:len(prefix)] != prefix {
		return ""
	}

	remaining := webApp.Name[len(prefix):]
	for i, char := range remaining {
		if char == '/' {
			return remaining[:i]
		}
	}

	return ""
}

// Note: Type conversion functions removed
// Use androidmanagement.WebApp directly instead of custom WebApp type
