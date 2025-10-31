// Package types provides unified resource name parsing utilities.
//
// 提供统一的资源名称解析函数，用于从 Android Management API 的资源名称中提取 ID。
//
// Android Management API 的资源名称格式：
//   - Enterprise: "enterprises/{enterpriseId}"
//   - Policy: "enterprises/{enterpriseId}/policies/{policyId}"
//   - Device: "enterprises/{enterpriseId}/devices/{deviceId}"
//   - EnrollmentToken: "enterprises/{enterpriseId}/enrollmentTokens/{tokenId}"
//   - MigrationToken: "enterprises/{enterpriseId}/migrationTokens/{tokenId}"
//   - WebApp: "enterprises/{enterpriseId}/webApps/{webAppId}"
//   - WebToken: "enterprises/{enterpriseId}/webTokens/{tokenId}"
//   - ProvisioningInfo: "signupUrls/{signupUrlId}" 或 "enterprises/{enterpriseId}/devices/{deviceId}/provisioningInfo"
//
// 使用方式：
//
//	import "amapi-pkg/pkgs/amapi/types"
//
//	// 方式 1: 使用结构体（推荐，更灵活）
//	resourceName := types.ParseResourceNameStruct("enterprises/LC00abc123/policies/default")
//	enterpriseID := resourceName.EnterpriseID    // "LC00abc123"
//	policyID := resourceName.PolicyID          // "default"
//
//	// 方式 2: 使用统一字段提取函数
//	enterpriseID := types.ExtractResourceField("enterprises/LC00abc123/policies/default", "EnterpriseID")
package types

import "strings"

// ResourceName represents a parsed Android Management API resource name.
//
// This struct provides type-safe access to all segments of a resource name.
// Fields are populated based on the resource type detected during parsing.
//
// Example:
//
//	name := ParseResourceNameStruct("enterprises/LC00abc123/policies/default")
//	// name.ResourceType = "policy"
//	// name.EnterpriseID = "LC00abc123"
//	// name.PolicyID = "default"
//
//	name := ParseResourceNameStruct("enterprises/LC00abc123/devices/device123")
//	// name.ResourceType = "device"
//	// name.EnterpriseID = "LC00abc123"
//	// name.DeviceID = "device123"
type ResourceName struct {
	// ResourceType indicates the type of resource (e.g., "enterprise", "policy", "device")
	ResourceType string

	// Common fields
	EnterpriseID string

	// Specific resource IDs (populated based on ResourceType)
	PolicyID          string
	DeviceID          string
	EnrollmentTokenID string
	MigrationTokenID  string
	WebAppID          string
	WebTokenID        string
	SignupURLID       string

	// Original resource name
	Original string

	// All segments for advanced usage
	Segments []string
}

// ParseResourceNameStruct parses a resource name into a structured ResourceName object.
//
// This is the recommended way to extract IDs from resource names as it provides
// type-safe access to all fields and automatically detects the resource type.
//
// Example:
//
//	name := ParseResourceNameStruct("enterprises/LC00abc123/policies/default")
//	if name.ResourceType == "policy" {
//	    fmt.Printf("Enterprise: %s, Policy: %s\n", name.EnterpriseID, name.PolicyID)
//	}
func ParseResourceNameStruct(resourceName string) *ResourceName {
	if resourceName == "" {
		return &ResourceName{Original: resourceName}
	}

	segments := strings.Split(resourceName, "/")
	// Filter out empty segments
	filteredSegments := make([]string, 0, len(segments))
	for _, seg := range segments {
		if seg != "" {
			filteredSegments = append(filteredSegments, seg)
		}
	}

	rn := &ResourceName{
		Original:     resourceName,
		Segments:     filteredSegments,
		ResourceType: detectResourceType(filteredSegments),
	}

	// Parse based on resource type
	if len(filteredSegments) >= 2 && filteredSegments[0] == "enterprises" {
		rn.EnterpriseID = filteredSegments[1]

		// Parse nested resources
		if len(filteredSegments) >= 4 {
			switch filteredSegments[2] {
			case "policies":
				rn.PolicyID = filteredSegments[3]
			case "devices":
				rn.DeviceID = filteredSegments[3]
				// Check for provisioning info
				if len(filteredSegments) >= 5 && filteredSegments[4] == "provisioningInfo" {
					rn.ResourceType = "provisioningInfo"
				}
			case "enrollmentTokens":
				rn.EnrollmentTokenID = filteredSegments[3]
			case "migrationTokens":
				rn.MigrationTokenID = filteredSegments[3]
			case "webApps":
				rn.WebAppID = filteredSegments[3]
			case "webTokens":
				rn.WebTokenID = filteredSegments[3]
			}
		}
	} else if len(filteredSegments) >= 2 && filteredSegments[0] == "signupUrls" {
		rn.SignupURLID = filteredSegments[1]
	}

	return rn
}

// detectResourceType detects the resource type from segments.
func detectResourceType(segments []string) string {
	if len(segments) == 0 {
		return "unknown"
	}

	// Check top-level resource type
	switch segments[0] {
	case "enterprises":
		if len(segments) == 2 {
			return "enterprise"
		}
		if len(segments) >= 4 {
			switch segments[2] {
			case "policies":
				return "policy"
			case "devices":
				if len(segments) >= 5 && segments[4] == "provisioningInfo" {
					return "provisioningInfo"
				}
				return "device"
			case "enrollmentTokens":
				return "enrollmentToken"
			case "migrationTokens":
				return "migrationToken"
			case "webApps":
				return "webApp"
			case "webTokens":
				return "webToken"
			}
		}
	case "signupUrls":
		return "signupUrl"
	}

	return "unknown"
}

// GetField extracts a field value by field name from the ResourceName struct.
//
// Supported field names:
//   - "EnterpriseID"
//   - "PolicyID"
//   - "DeviceID"
//   - "EnrollmentTokenID"
//   - "MigrationTokenID"
//   - "WebAppID"
//   - "WebTokenID"
//   - "SignupURLID"
//   - "ResourceType"
//
// Returns empty string if field name is invalid or field is empty.
//
// This is used internally by convenience wrapper functions to extract specific IDs.
func (rn *ResourceName) GetField(fieldName string) string {
	if rn == nil {
		return ""
	}

	switch fieldName {
	case "EnterpriseID":
		return rn.EnterpriseID
	case "PolicyID":
		return rn.PolicyID
	case "DeviceID":
		return rn.DeviceID
	case "EnrollmentTokenID":
		return rn.EnrollmentTokenID
	case "MigrationTokenID":
		return rn.MigrationTokenID
	case "WebAppID":
		return rn.WebAppID
	case "WebTokenID":
		return rn.WebTokenID
	case "SignupURLID":
		return rn.SignupURLID
	case "ResourceType":
		return rn.ResourceType
	default:
		return ""
	}
}

// ExtractResourceField is a generic helper function that extracts a field from a resource name.
//
// This function:
// 1. Checks if the resource name is valid
// 2. Parses it into a ResourceName struct
// 3. Extracts the specified field by field name
//
// Supported field names:
//   - "EnterpriseID"
//   - "PolicyID"
//   - "DeviceID"
//   - "EnrollmentTokenID"
//   - "MigrationTokenID"
//   - "WebAppID"
//   - "WebTokenID"
//   - "SignupURLID"
//
// This is used by all Get*ID and Get*EnterpriseID convenience wrapper functions.
//
// Example:
//
//	enterpriseID := ExtractResourceField("enterprises/LC00abc123/policies/default", "EnterpriseID")
//	policyID := ExtractResourceField("enterprises/LC00abc123/policies/default", "PolicyID")
func ExtractResourceField(resourceName string, fieldName string) string {
	if resourceName == "" {
		return ""
	}

	rn := ParseResourceNameStruct(resourceName)
	if rn == nil {
		return ""
	}

	return rn.GetField(fieldName)
}
