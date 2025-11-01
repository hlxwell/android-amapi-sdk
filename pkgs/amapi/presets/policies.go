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
//
//   - ApplicationReportsEnabled: 应用报告
//
//   - CommonCriteriaModeEnabled: Common Criteria 模式
//
//   - DeviceSettingsEnabled: 设备设置
//
//   - DisplayInfoEnabled: 显示信息
//
//   - HardwareStatusEnabled: 硬件状态
//
//   - MemoryInfoEnabled: 内存信息
//
//   - NetworkInfoEnabled: 网络信息
//
//   - PowerManagementEventsEnabled: 电源管理事件
//
//   - SoftwareInfoEnabled: 软件信息
//
//   - SystemPropertiesEnabled: 系统属性
//
//   - Usage Logs: 启用安全日志
//
//   - SECURITY_LOGS: 安全日志
//
// 这个配置可以让设备上报所有基础数据，便于全面的设备监控和管理。
package presets

import (
	"encoding/json"
	"fmt"

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
	return newBasePolicy()
}

func newBasePolicy() *androidmanagement.Policy {
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

// PolicyPreset describes a reusable policy template with metadata.
type PolicyPreset struct {
	Name        string
	DisplayName string
	Description string
	Tags        []string
	Policy      *androidmanagement.Policy
}

type policyPresetDefinition struct {
	Name        string
	DisplayName string
	Description string
	Tags        []string
	Builder     func() *androidmanagement.Policy
}

var policyPresetDefinitions = []policyPresetDefinition{
	{
		Name:        "fully_managed",
		DisplayName: "Fully Managed Device",
		Description: "全面托管模式，启用完整的状态上报与安全基线。",
		Tags:        []string{"enterprise", "managed"},
		Builder:     func() *androidmanagement.Policy { return clonePolicy(newBasePolicy()) },
	},
	{
		Name:        "dedicated_device",
		DisplayName: "Dedicated Device",
		Description: "适用于专用场景的单应用/多应用固定设备，启用自助亭模式。",
		Tags:        []string{"kiosk", "lock-task"},
		Builder:     buildDedicatedDevicePolicy,
	},
	{
		Name:        "work_profile",
		DisplayName: "Work Profile",
		Description: "在公司拥有设备上启用工作资料分区，允许个人与工作并存。",
		Tags:        []string{"COPE", "BYOD"},
		Builder:     buildWorkProfilePolicy,
	},
	{
		Name:        "kiosk_mode",
		DisplayName: "Kiosk Mode",
		Description: "锁定到单一应用，隐藏系统 UI 元素，适合展台设备。",
		Tags:        []string{"kiosk", "single-app"},
		Builder:     buildRetailKioskPolicy,
	},
	{
		Name:        "cope",
		DisplayName: "Company Owned, Personally Enabled",
		Description: "公司拥有并允许个人使用的模式，强调个人隐私保护。",
		Tags:        []string{"COPE", "privacy"},
		Builder:     buildCopePolicy,
	},
	{
		Name:        "secure_workstation",
		DisplayName: "Secure Workstation",
		Description: "高安全要求的终端策略，禁用截屏与外部存储。",
		Tags:        []string{"security", "compliance"},
		Builder:     buildSecureWorkstationPolicy,
	},
	{
		Name:        "education_tablet",
		DisplayName: "Education Tablet",
		Description: "教育场景推荐策略，保留学习应用同时限制娱乐功能。",
		Tags:        []string{"education", "tablet"},
		Builder:     buildEducationTabletPolicy,
	},
	{
		Name:        "retail_kiosk",
		DisplayName: "Retail Kiosk",
		Description: "零售门店展台策略，固定应用并启用自助任务锁定。",
		Tags:        []string{"retail", "kiosk"},
		Builder:     buildRetailKioskPolicy,
	},
}

// GetAllPresets returns available policy presets.
func GetAllPresets() []*PolicyPreset {
	result := make([]*PolicyPreset, 0, len(policyPresetDefinitions))
	for _, def := range policyPresetDefinitions {
		result = append(result, buildPolicyPreset(def))
	}
	return result
}

// GetPresetByName finds a preset by identifier.
func GetPresetByName(name string) *PolicyPreset {
	for _, def := range policyPresetDefinitions {
		if def.Name == name {
			return buildPolicyPreset(def)
		}
	}
	return nil
}

// CreatePolicyFromPreset returns a cloned policy from the preset and applies an optional customization function.
func CreatePolicyFromPreset(name string, customize func(*androidmanagement.Policy) *androidmanagement.Policy) (*androidmanagement.Policy, error) {
	preset := GetPresetByName(name)
	if preset == nil {
		return nil, fmt.Errorf("unknown policy preset: %s", name)
	}

	policy := clonePolicy(preset.Policy)
	if customize != nil {
		if customized := customize(clonePolicy(preset.Policy)); customized != nil {
			policy = customized
		}
	}

	return policy, nil
}

func buildPolicyPreset(def policyPresetDefinition) *PolicyPreset {
	policy := def.Builder()
	return &PolicyPreset{
		Name:        def.Name,
		DisplayName: def.DisplayName,
		Description: def.Description,
		Tags:        append([]string(nil), def.Tags...),
		Policy:      policy,
	}
}

func buildDedicatedDevicePolicy() *androidmanagement.Policy {
	policy := newBasePolicy()
	policy.KioskCustomLauncherEnabled = true
	policy.StatusBarDisabled = true
	policy.Applications = []*androidmanagement.ApplicationPolicy{
		{
			PackageName:     "com.android.chrome",
			InstallType:     "KIOSK",
			LockTaskAllowed: true,
		},
	}
	return policy
}

func buildWorkProfilePolicy() *androidmanagement.Policy {
	policy := newBasePolicy()
	policy.AddUserDisabled = false
	policy.InstallAppsDisabled = false
	policy.PersonalUsagePolicies = &androidmanagement.PersonalUsagePolicies{
		PersonalPlayStoreMode: "ALLOWLIST",
	}
	return policy
}

func buildCopePolicy() *androidmanagement.Policy {
	policy := newBasePolicy()
	policy.PersonalUsagePolicies = &androidmanagement.PersonalUsagePolicies{
		PersonalPlayStoreMode: "BLACKLIST",
	}
	policy.BluetoothDisabled = false
	policy.CameraDisabled = false
	return policy
}

func buildSecureWorkstationPolicy() *androidmanagement.Policy {
	policy := newBasePolicy()
	policy.ScreenCaptureDisabled = true
	policy.UsbFileTransferDisabled = true
	policy.MicrophoneAccess = "MICROPHONE_ACCESS_ENFORCED"
	policy.CameraAccess = "CAMERA_ACCESS_DISABLED"
	return policy
}

func buildEducationTabletPolicy() *androidmanagement.Policy {
	policy := newBasePolicy()
	policy.FunDisabled = false
	policy.InstallAppsDisabled = false
	policy.Applications = []*androidmanagement.ApplicationPolicy{
		{
			PackageName: "com.android.chrome",
			InstallType: "REQUIRED",
		},
		{
			PackageName: "com.google.android.youtube",
			InstallType: "AVAILABLE",
		},
	}
	return policy
}

func buildRetailKioskPolicy() *androidmanagement.Policy {
	policy := buildDedicatedDevicePolicy()
	policy.PlayStoreMode = "WHITELIST"
	return policy
}

func clonePolicy(p *androidmanagement.Policy) *androidmanagement.Policy {
	if p == nil {
		return nil
	}
	data, err := json.Marshal(p)
	if err != nil {
		return newBasePolicy()
	}
	var clone androidmanagement.Policy
	if err := json.Unmarshal(data, &clone); err != nil {
		return newBasePolicy()
	}
	return &clone
}

// ClonePolicy exposes a safe copy helper for callers needing to duplicate preset policies.
func ClonePolicy(p *androidmanagement.Policy) *androidmanagement.Policy {
	return clonePolicy(p)
}
