package types

import (
	"google.golang.org/api/androidmanagement/v1"
)

// WebApp 相关类型和函数
//
// 注意：WebApp 类型直接使用 androidmanagement.WebApp。
// 此文件包含 Web 应用相关的请求类型。
//
// 使用方式：
//
//	import (
//	    "amapi-pkg/pkgs/amapi/types"
//	    "google.golang.org/api/androidmanagement/v1"
//	)
//
//	// 创建 Web 应用请求
//	req := &types.WebAppCreateRequest{
//	    EnterpriseName: "enterprises/LC00abc123",
//	    DisplayName:    "My Web App",
//	    StartURL:       "https://example.com",
//	}

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

// Note: Type conversion functions removed
// Use androidmanagement.WebApp directly instead of custom WebApp type
