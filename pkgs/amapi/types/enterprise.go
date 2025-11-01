package types

import (
	"time"
)

// Enterprise 相关类型和函数
//
// 注意：Enterprise 类型直接使用 androidmanagement.Enterprise。
// 此文件包含企业相关的辅助类型和辅助函数。
//
// 使用方式：
//
//	import (
//	    "amapi-pkg/pkgs/amapi/types"
//	    "google.golang.org/api/androidmanagement/v1"
//	)
//
//		// 创建企业直接传递参数
//	//	enterprise, err := client.Enterprises().Create(signupToken, projectID, enterpriseToken, contactInfo)
//
//	// Enterprise 类型直接使用 androidmanagement.Enterprise
//	// var enterprise *androidmanagement.Enterprise

// EnterpriseSignupURL represents a signup URL for enterprise creation.
type EnterpriseSignupURL struct {
	// URL is the signup URL
	URL string `json:"url"`

	// CallbackURL is the URL to redirect to after signup
	CallbackURL string `json:"callback_url,omitempty"`

	// ProjectID is the Google Cloud project ID
	ProjectID string `json:"project_id"`

	// CompletionToken is used to complete the signup process
	CompletionToken string `json:"completion_token,omitempty"`

	// ExpiresAt indicates when the signup URL expires
	ExpiresAt time.Time `json:"expires_at,omitempty"`

	// CreatedAt timestamp
	CreatedAt time.Time `json:"created_at"`
}

// EnterpriseUpgradeURL represents an enterprise upgrade URL.
type EnterpriseUpgradeURL struct {
	// URL is the upgrade URL
	URL string `json:"url"`

	// EnterpriseName is the enterprise this URL is for
	EnterpriseName string `json:"enterprise_name"`

	// ProjectID is the Google Cloud project ID
	ProjectID string `json:"project_id"`

	// CreatedAt timestamp
	CreatedAt time.Time `json:"created_at"`

	// ExpiresAt indicates when the upgrade URL expires
	ExpiresAt time.Time `json:"expires_at,omitempty"`
}

// Notification type constants
const (
	NotificationTypeEnrollment       = "ENROLLMENT"
	NotificationTypeComplianceReport = "COMPLIANCE_REPORT"
	NotificationTypeStatusReport     = "STATUS_REPORT"
	NotificationTypeCommand          = "COMMAND"
	NotificationTypeUsageLog         = "USAGE_LOG_ENABLED"
)
