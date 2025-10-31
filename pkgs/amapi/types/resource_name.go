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
//	// 方式 2: 使用便捷方法（向后兼容）
//	deviceID := types.ExtractLastSegment(device.Name)
//	enterpriseID := types.ExtractEnterpriseID(device.Name)
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

// GetID returns the primary resource ID based on the resource type.
//
// This is a convenience method that returns the appropriate ID field
// based on the detected resource type.
func (rn *ResourceName) GetID() string {
	if rn == nil {
		return ""
	}

	switch rn.ResourceType {
	case "enterprise":
		return rn.EnterpriseID
	case "policy":
		return rn.PolicyID
	case "device":
		return rn.DeviceID
	case "enrollmentToken":
		return rn.EnrollmentTokenID
	case "migrationToken":
		return rn.MigrationTokenID
	case "webApp":
		return rn.WebAppID
	case "webToken":
		return rn.WebTokenID
	case "signupUrl":
		return rn.SignupURLID
	default:
		// Fallback: return last segment
		if len(rn.Segments) > 0 {
			return rn.Segments[len(rn.Segments)-1]
		}
		return ""
	}
}

// ExtractSegment extracts a segment from a resource name by index (0-based).
//
// Resource name format: "enterprises/{enterpriseId}/policies/{policyId}"
//   - ExtractSegment(name, 0) returns "enterprises"
//   - ExtractSegment(name, 1) returns "{enterpriseId}"
//   - ExtractSegment(name, 2) returns "policies"
//   - ExtractSegment(name, 3) returns "{policyId}"
//
// Returns empty string if index is out of range or name is empty.
//
// Deprecated: Use ParseResourceNameStruct() instead for type-safe access.
func ExtractSegment(resourceName string, index int) string {
	if resourceName == "" {
		return ""
	}

	segments := strings.Split(resourceName, "/")
	if index < 0 || index >= len(segments) {
		return ""
	}

	return segments[index]
}

// ExtractLastSegment extracts the last segment from a resource name.
//
// This is commonly used to extract resource IDs:
//   - "enterprises/LC00abc123" -> "LC00abc123"
//   - "enterprises/LC00abc123/policies/default" -> "default"
//   - "enterprises/LC00abc123/devices/device123" -> "device123"
//
// Returns empty string if name is empty or has no segments.
//
// Deprecated: Use ParseResourceNameStruct() instead for type-safe access.
func ExtractLastSegment(resourceName string) string {
	if resourceName == "" {
		return ""
	}

	// Find the last '/' and return everything after it
	lastSlash := strings.LastIndex(resourceName, "/")
	if lastSlash == -1 {
		// No slash found, return the whole string (e.g., just an ID)
		return resourceName
	}

	if lastSlash == len(resourceName)-1 {
		// Last character is '/', invalid format
		return ""
	}

	return resourceName[lastSlash+1:]
}

// ExtractEnterpriseID extracts the enterprise ID from any resource name.
//
// Works with any resource name format:
//   - "enterprises/LC00abc123" -> "LC00abc123"
//   - "enterprises/LC00abc123/policies/default" -> "LC00abc123"
//   - "enterprises/LC00abc123/devices/device123" -> "LC00abc123"
//
// Returns empty string if:
//   - name is empty
//   - name doesn't start with "enterprises/"
//   - enterprise ID segment is missing
//
// Deprecated: Use ParseResourceNameStruct() instead for type-safe access.
func ExtractEnterpriseID(resourceName string) string {
	if resourceName == "" {
		return ""
	}

	const prefix = "enterprises/"
	if !strings.HasPrefix(resourceName, prefix) {
		return ""
	}

	// Skip the prefix and extract the next segment
	remaining := resourceName[len(prefix):]
	if remaining == "" {
		return ""
	}

	// Find the next '/' (if any) to get just the enterprise ID
	nextSlash := strings.Index(remaining, "/")
	if nextSlash == -1 {
		// No more segments, the whole remaining part is the enterprise ID
		return remaining
	}

	// Return the segment before the next '/'
	return remaining[:nextSlash]
}

// ParseResourceName parses a resource name and returns all segments.
//
// Example:
//   - "enterprises/LC00abc123/policies/default" -> ["enterprises", "LC00abc123", "policies", "default"]
//
// Returns nil if name is empty.
//
// Deprecated: Use ParseResourceNameStruct() instead for structured access.
func ParseResourceName(resourceName string) []string {
	rn := ParseResourceNameStruct(resourceName)
	if rn == nil || len(rn.Segments) == 0 {
		return nil
	}
	return rn.Segments
}
