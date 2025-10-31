package types

import (
	"google.golang.org/api/androidmanagement/v1"
)

// ProvisioningInfo 相关类型和函数
//
// 注意：ProvisioningInfo 类型直接使用 androidmanagement.ProvisioningInfo。
// 此文件包含配置信息相关的辅助函数。
//
// 使用方式：
//
//	import "amapi-pkg/pkgs/amapi/types"
//
//	// ProvisioningInfo 通常没有资源名称格式，如果需要提取信息，请直接访问字段

// GetProvisioningInfoName extracts the provisioning info name or ID.
//
// ProvisioningInfo 的资源名称格式通常是：
//   - "signupUrls/{signupUrlId}"
//   - "enterprises/{enterpriseId}/devices/{deviceId}/provisioningInfo"
//
// This is a convenience wrapper around ParseResourceNameStruct.
func GetProvisioningInfoName(info *androidmanagement.ProvisioningInfo) string {
	if info == nil || info.Name == "" {
		return ""
	}
	rn := ParseResourceNameStruct(info.Name)
	if rn == nil {
		return ""
	}
	// Try to get the appropriate ID based on resource type
	if rn.ResourceType == "signupUrl" {
		return rn.SignupURLID
	}
	// For device provisioning info, return the device ID
	if rn.ResourceType == "provisioningInfo" {
		return rn.DeviceID
	}
	// Fallback to GetID
	return rn.GetID()
}
