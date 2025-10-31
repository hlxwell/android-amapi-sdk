// Package presets provides default policy configuration.
//
// 这个包提供了一个默认策略配置，用于快速开始设备管理。
// 默认策略启用了所有基础数据上报功能，便于监控和管理设备。
//
// # 使用方式
//
//	import "amapi-pkg/pkgs/amapi/presets"
//
//	// 获取默认策略
//	policy := presets.GetDefaultPolicy()
//
//	// 使用策略创建设备管理策略
//	createdPolicy, err := client.Policies().CreateByEnterpriseID(
//	    enterpriseID,
//	    "default-policy",
//	    policy,
//	)
//
// # 默认策略功能
//
// 默认策略包含以下配置：
//
//   - Status Reporting Settings: 启用所有状态上报选项
//     - ApplicationReportsEnabled: 应用报告
//     - CommonCriteriaModeEnabled: Common Criteria 模式
//     - DeviceSettingsEnabled: 设备设置
//     - DisplayInfoEnabled: 显示信息
//     - HardwareStatusEnabled: 硬件状态
//     - MemoryInfoEnabled: 内存信息
//     - NetworkInfoEnabled: 网络信息
//     - PowerManagementEventsEnabled: 电源管理事件
//     - SoftwareInfoEnabled: 软件信息
//     - SystemPropertiesEnabled: 系统属性
//
//   - Usage Logs: 启用安全日志
//     - SECURITY_LOGS: 安全日志
//
// 这个配置可以让设备上报所有基础数据，便于全面的设备监控和管理。
package presets

import (
	"google.golang.org/api/androidmanagement/v1"
)

// GetDefaultPolicy returns a default policy with all status reporting enabled.
//
// 此策略启用了全面的设备监控和数据收集功能。
// 使用此策略可以确保设备上报所有基础数据，包括：
//   - 应用安装和使用情况
//   - 设备硬件和软件状态
//   - 网络连接信息
//   - 电源管理事件
//   - 安全日志
//
// 适用于需要全面了解设备状态的场景。
//
// 返回的 Policy 可以直接用于创建或更新策略：
//
//	policy := presets.GetDefaultPolicy()
//	createdPolicy, err := client.Policies().CreateByEnterpriseID(
//	    enterpriseID,
//	    "default-policy",
//	    policy,
//	)
func GetDefaultPolicy() *androidmanagement.Policy {
	return &androidmanagement.Policy{
		// Enable all status reporting settings for comprehensive device monitoring
		StatusReportingSettings: &androidmanagement.StatusReportingSettings{
			ApplicationReportsEnabled:    true,
			CommonCriteriaModeEnabled:    true,
			DeviceSettingsEnabled:        true,
			DisplayInfoEnabled:           true,
			HardwareStatusEnabled:        true,
			MemoryInfoEnabled:            true,
			NetworkInfoEnabled:           true,
			PowerManagementEventsEnabled: true,
			SoftwareInfoEnabled:          true,
			SystemPropertiesEnabled:      true,
		},
		// Enable security logs in usage log
		UsageLog: &androidmanagement.UsageLog{
			EnabledLogTypes: []string{"SECURITY_LOGS"},
		},
	}
}
