package types

import (
	"google.golang.org/api/androidmanagement/v1"
)

// MigrationToken 相关类型和函数
//
// 注意：MigrationToken 类型直接使用 androidmanagement.MigrationToken。
// 此文件包含迁移令牌相关的辅助函数。
//
// 使用方式：
//
//	import "amapi-pkg/pkgs/amapi/types"
//
//	// 提取迁移令牌 ID
//	tokenID := types.GetMigrationTokenID(token)
//	enterpriseID := types.GetMigrationTokenEnterpriseID(token)

// GetMigrationTokenID extracts the migration token ID from the resource name.
//
// This is a convenience wrapper around ExtractResourceField.
func GetMigrationTokenID(token *androidmanagement.MigrationToken) string {
	if token == nil || token.Name == "" {
		return ""
	}
	return ExtractResourceField(token.Name, "MigrationTokenID")
}

// GetMigrationTokenEnterpriseID extracts the enterprise ID from the migration token resource name.
//
// This is a convenience wrapper around ExtractResourceField.
func GetMigrationTokenEnterpriseID(token *androidmanagement.MigrationToken) string {
	if token == nil || token.Name == "" {
		return ""
	}
	return ExtractResourceField(token.Name, "EnterpriseID")
}
