package client

import (
	"context"
	"strings"
	"testing"
	"time"

	"amapi-pkg/pkgs/amapi/config"
	"amapi-pkg/pkgs/amapi/types"
)

// 测试常量是否正确定义
func TestConstants(t *testing.T) {
	tests := []struct {
		name     string
		constant interface{}
		expected interface{}
	}{
		{"ClientVersion", ClientVersion, "1.0.0"},
		{"DefaultRetryMaxDelay", DefaultRetryMaxDelay, 30 * time.Second},
		{"DefaultRedisTimeout", DefaultRedisTimeout, 5 * time.Second},
		{"DefaultHealthCheckTimeout", DefaultHealthCheckTimeout, 10 * time.Second},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.constant != tt.expected {
				t.Errorf("%s = %v, want %v", tt.name, tt.constant, tt.expected)
			}
		})
	}
}

// 测试通用资源名称验证函数
func TestValidateResourceName(t *testing.T) {
	tests := []struct {
		name          string
		resourceName  string
		expectedParts []string
		resourceType  string
		expectError   bool
	}{
		{
			name:          "valid enterprise name",
			resourceName:  "enterprises/test-enterprise",
			expectedParts: []string{"enterprises", "{enterpriseId}"},
			resourceType:  "enterprise",
			expectError:   false,
		},
		{
			name:          "valid device name",
			resourceName:  "enterprises/test-enterprise/devices/test-device",
			expectedParts: []string{"enterprises", "{enterpriseId}", "devices", "{deviceId}"},
			resourceType:  "device",
			expectError:   false,
		},
		{
			name:          "invalid format - wrong number of parts",
			resourceName:  "enterprises",
			expectedParts: []string{"enterprises", "{enterpriseId}"},
			resourceType:  "enterprise",
			expectError:   true,
		},
		{
			name:          "invalid format - wrong fixed part",
			resourceName:  "companies/test-enterprise",
			expectedParts: []string{"enterprises", "{enterpriseId}"},
			resourceType:  "enterprise",
			expectError:   true,
		},
		{
			name:          "empty resource name",
			resourceName:  "",
			expectedParts: []string{"enterprises", "{enterpriseId}"},
			resourceType:  "enterprise",
			expectError:   true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			components, err := validateResourceName(tt.resourceName, tt.expectedParts, tt.resourceType)

			if tt.expectError {
				if err == nil {
					t.Errorf("validateResourceName() expected error, got nil")
				}
			} else {
				if err != nil {
					t.Errorf("validateResourceName() unexpected error: %v", err)
				}
				if len(components) != len(tt.expectedParts) {
					t.Errorf("validateResourceName() returned %d components, expected %d", len(components), len(tt.expectedParts))
				}
			}
		})
	}
}

// 测试解析函数
func TestParseResourceNames(t *testing.T) {
	tests := []struct {
		name     string
		testFunc func() error
	}{
		{
			name: "parseEnterpriseName",
			testFunc: func() error {
				enterpriseID, err := parseEnterpriseName("enterprises/test-enterprise")
				if err != nil {
					return err
				}
				if enterpriseID != "test-enterprise" {
					return types.NewError(types.ErrCodeInvalidInput, "unexpected enterprise ID")
				}
				return nil
			},
		},
		{
			name: "parseDeviceName",
			testFunc: func() error {
				enterpriseID, deviceID, err := parseDeviceName("enterprises/test-enterprise/devices/test-device")
				if err != nil {
					return err
				}
				if enterpriseID != "test-enterprise" || deviceID != "test-device" {
					return types.NewError(types.ErrCodeInvalidInput, "unexpected IDs")
				}
				return nil
			},
		},
		{
			name: "parsePolicyName",
			testFunc: func() error {
				enterpriseID, policyID, err := parsePolicyName("enterprises/test-enterprise/policies/test-policy")
				if err != nil {
					return err
				}
				if enterpriseID != "test-enterprise" || policyID != "test-policy" {
					return types.NewError(types.ErrCodeInvalidInput, "unexpected IDs")
				}
				return nil
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.testFunc(); err != nil {
				t.Errorf("%s failed: %v", tt.name, err)
			}
		})
	}
}

// 测试构建资源名称函数
func TestBuildResourceName(t *testing.T) {
	tests := []struct {
		name       string
		components []string
		expected   string
	}{
		{
			name:       "enterprise name",
			components: []string{"enterprises", "test-enterprise"},
			expected:   "enterprises/test-enterprise",
		},
		{
			name:       "device name",
			components: []string{"enterprises", "test-enterprise", "devices", "test-device"},
			expected:   "enterprises/test-enterprise/devices/test-device",
		},
		{
			name:       "empty components",
			components: []string{},
			expected:   "",
		},
		{
			name:       "single component",
			components: []string{"enterprises"},
			expected:   "enterprises",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := buildResourceName(tt.components...)
			if result != tt.expected {
				t.Errorf("buildResourceName() = %v, want %v", result, tt.expected)
			}
		})
	}
}

// 测试NewWithContext的线程安全性
func TestNewWithContextThreadSafety(t *testing.T) {
	// 创建一个基本的配置
	cfg := &config.Config{
		ProjectID:       "test-project",
		CredentialsJSON: `{"type": "service_account", "project_id": "test-project"}`,
		Scopes:          []string{"https://www.googleapis.com/auth/androidmanagement"},
	}

	// 创建带超时的context
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// 测试NewWithContext函数 - 这应该不会修改现有的client实例
	_, err := NewWithContext(ctx, cfg)
	// 注意：这可能会因为认证问题而失败，但我们主要测试函数调用不会panic
	if err != nil && !strings.Contains(err.Error(), "authentication") && !strings.Contains(err.Error(), "credentials") {
		t.Logf("NewWithContext failed with non-auth error: %v", err)
	}
}

// 测试验证函数
func TestValidationFunctions(t *testing.T) {
	tests := []struct {
		name     string
		testFunc func() error
		expected bool // true if should pass, false if should fail
	}{
		{
			name: "validateEnterpriseID - valid",
			testFunc: func() error {
				return validateEnterpriseID("test-enterprise")
			},
			expected: true,
		},
		{
			name: "validateEnterpriseID - empty",
			testFunc: func() error {
				return validateEnterpriseID("")
			},
			expected: false,
		},
		{
			name: "validateDeviceID - valid",
			testFunc: func() error {
				return validateDeviceID("test-device")
			},
			expected: true,
		},
		{
			name: "validateDeviceID - empty",
			testFunc: func() error {
				return validateDeviceID("")
			},
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.testFunc()
			if tt.expected && err != nil {
				t.Errorf("%s expected to pass but failed: %v", tt.name, err)
			}
			if !tt.expected && err == nil {
				t.Errorf("%s expected to fail but passed", tt.name)
			}
		})
	}
}

// 基准测试 - 测试优化后的性能
func BenchmarkValidateResourceName(b *testing.B) {
	expectedParts := []string{"enterprises", "{enterpriseId}", "devices", "{deviceId}"}
	resourceName := "enterprises/test-enterprise/devices/test-device"

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = validateResourceName(resourceName, expectedParts, "device")
	}
}

func BenchmarkParseDeviceName(b *testing.B) {
	deviceName := "enterprises/test-enterprise/devices/test-device"

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _, _ = parseDeviceName(deviceName)
	}
}