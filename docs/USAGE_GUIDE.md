# AMAPI Package Usage Guide

这是一个完整的指南，介绍如何使用 `pkgs/amapi` 包来管理 Android Management API。

## 概述

`pkgs/amapi` 包提供了一个高级的、类型安全的 Go 客户端，用于与 Google Android Management API 进行交互。它支持企业管理、策略配置、设备管理和注册令牌等核心功能。

## 快速开始

### 1. 安装和配置

```go
import "github.com/hlxwell/android-api-demo/pkgs/amapi/client"
import "github.com/hlxwell/android-api-demo/pkgs/amapi/config"
```

### 2. 基本设置

```go
// 方法 1: 自动加载配置
cfg, err := config.AutoLoadConfig()
if err != nil {
    log.Fatal(err)
}

// 方法 2: 手动配置
cfg := &config.Config{
    ProjectID:       "your-gcp-project",
    CredentialsFile: "./service-account-key.json",
    CallbackURL:     "https://your-app.com/callback",
    Timeout:         30 * time.Second,
    EnableRetry:     true,
    RetryAttempts:   3,
}

// 创建客户端
client, err := client.New(cfg)
if err != nil {
    log.Fatal(err)
}
defer client.Close()
```

## 核心功能

### 企业管理

#### 生成注册URL
```go
signupReq := &types.SignupURLRequest{
    ProjectID:   "your-project-id",
    CallbackURL: "https://your-app.com/callback",
    AdminEmail:  "admin@company.com",
}

signupURL, err := client.Enterprises().GenerateSignupURL(signupReq)
```

#### 创建企业
```go
createReq := &types.EnterpriseCreateRequest{
    SignupToken: "token-from-callback",
    ProjectID:   "your-project-id",
    DisplayName: "My Company",
}

enterprise, err := client.Enterprises().Create(createReq)
```

#### 企业配置
```go
// 启用通知
notificationTypes := []string{
    types.NotificationTypeEnrollment,
    types.NotificationTypeComplianceReport,
}

updated, err := client.Enterprises().EnableNotifications(
    enterpriseName,
    notificationTypes,
)

// 设置 Pub/Sub 主题
topicName := "projects/your-project/topics/amapi-events"
enterprise, err := client.Enterprises().SetPubSubTopic(
    enterpriseName,
    topicName,
)
```

### 策略管理

#### 使用预设策略
```go
// 获取全面管理设备的预设
preset := presets.GetFullyManagedPreset()

// 创建策略
policy, err := client.Policies().CreateByEnterpriseID(
    "enterprise-id",
    "my-policy",
    preset.Policy,
)
```

#### 自定义策略
```go
// 创建自定义策略
policy := &types.Policy{
    StatusBarDisabled:     false,
    KeyguardDisabled:      false,
    AddUserDisabled:       true,
    UninstallAppsDisabled: true,
    CameraDisabled:        false,
    AutoTimeRequired:      true,
}

created, err := client.Policies().CreateByEnterpriseID(
    "enterprise-id",
    "custom-policy",
    policy,
)
```

#### 管理应用程序
```go
// 添加必需应用
err := client.Policies().RequireApplication(
    policyName,
    "com.company.app",
)

// 阻止应用
err := client.Policies().BlockApplication(
    policyName,
    "com.social.app",
)

// 设置 Kiosk 模式
policy, err := client.Policies().SetKioskMode(
    policyName,
    "com.company.kioskapp",
)
```

### 设备管理

#### 列出设备
```go
// 列出所有设备
devices, err := client.Devices().ListByEnterpriseID("enterprise-id", nil)

// 获取合规设备
compliantDevices, err := client.Devices().GetCompliantDevices("enterprise-id")

// 获取非合规设备
nonCompliantDevices, err := client.Devices().GetNonCompliantDevices("enterprise-id")
```

#### 设备操作
```go
// 锁定设备 10 分钟
err := client.Devices().LockByID("enterprise-id", "device-id", "PT10M")

// 重启设备
err := client.Devices().RebootByID("enterprise-id", "device-id")

// 恢复出厂设置
err := client.Devices().ResetByID("enterprise-id", "device-id")

// 删除密码
err := client.Devices().RemovePasswordByID("enterprise-id", "device-id")
```

#### 获取设备信息
```go
device, err := client.Devices().GetByID("enterprise-id", "device-id")

fmt.Printf("设备状态: %s\n", device.State)
fmt.Printf("合规状态: %t\n", device.PolicyCompliant)
fmt.Printf("Android 版本: %s\n", device.GetAndroidVersion())
fmt.Printf("设备型号: %s\n", device.GetDeviceModel())
fmt.Printf("在线状态: %t\n", device.IsOnline())
```

### 注册令牌管理

#### 创建注册令牌
```go
// 基本令牌创建
token, err := client.EnrollmentTokens().CreateByEnterpriseID(
    "enterprise-id",
    "policy-id",
    24*time.Hour, // 有效期 24 小时
)

// 工作配置文件令牌
token, err := client.EnrollmentTokens().CreateForWorkProfile(
    "enterprise-id",
    "policy-id",
    7*24*time.Hour, // 有效期 1 周
)
```

#### 生成 QR 码
```go
qrOptions := &types.QRCodeOptions{
    WiFiSSID:         "CompanyWiFi",
    WiFiPassword:     "password123",
    WiFiSecurityType: types.WiFiSecurityTypeWPA2,
    SkipSetupWizard:  true,
    Locale:           "en_US",
}

qrData, err := client.EnrollmentTokens().GenerateQRCodeByID(
    "enterprise-id",
    "token-id",
    qrOptions,
)

// 转换为 JSON 用于 QR 码生成
qrJSON, err := qrData.ToJSON()
```

#### 令牌管理
```go
// 获取活跃令牌
activeTokens, err := client.EnrollmentTokens().GetActiveTokens("enterprise-id")

// 获取特定策略的令牌
policyTokens, err := client.EnrollmentTokens().GetTokensForPolicy(
    "enterprise-id",
    "policy-id",
)

// 撤销令牌
err := client.EnrollmentTokens().RevokeTokenByID("enterprise-id", "token-id")

// 批量创建令牌
tokens, err := client.EnrollmentTokens().CreateBulkTokens(
    "enterprise-id",
    "policy-id",
    10,                // 创建 10 个令牌
    24*time.Hour,      // 每个有效期 24 小时
)
```

## 预设策略模板

### 可用预设

```go
// 全面管理设备
fullyManaged := presets.GetFullyManagedPreset()

// 专用设备（Kiosk）
dedicated := presets.GetDedicatedDevicePreset()

// 工作配置文件（BYOD）
workProfile := presets.GetWorkProfilePreset()

// 企业自有、个人启用（COPE）
cope := presets.GetCOPEPreset()

// 安全工作站
secure := presets.GetSecureWorkstationPreset()

// 教育平板
education := presets.GetEducationTabletPreset()

// 零售 Kiosk
retail := presets.GetRetailKioskPreset()
```

### 自定义预设

```go
// 从预设创建自定义策略
customizations := map[string]interface{}{
    "camera_disabled": true,
    "bluetooth_disabled": false,
    "status_bar_disabled": false,
}

policy, err := presets.CreatePolicyFromPreset(
    "fully_managed",
    customizations,
)
```

## 错误处理

### 结构化错误处理

```go
devices, err := client.Devices().ListByEnterpriseID("enterprise-id", nil)
if err != nil {
    if apiErr, ok := err.(*types.Error); ok {
        switch apiErr.Code {
        case types.ErrCodeNotFound:
            fmt.Println("企业未找到")
        case types.ErrCodeTooManyRequests:
            fmt.Printf("请求频率过高，请在 %v 后重试\n",
                apiErr.RetryDelay(1, time.Second))
        case types.ErrCodeUnauthorized:
            fmt.Println("认证失败，请检查凭据")
        default:
            fmt.Printf("API 错误: %s\n", apiErr.Error())
        }
    } else {
        fmt.Printf("其他错误: %v\n", err)
    }
}
```

### 重试机制

```go
// 配置重试
cfg := &config.Config{
    // ... 其他配置
    RetryAttempts: 5,
    RetryDelay:    2 * time.Second,
    EnableRetry:   true,
}

// 检查错误是否可重试
if types.IsRetryableError(err) {
    delay := types.GetRetryDelay(err, attempt, baseDelay)
    time.Sleep(delay)
    // 重试操作
}
```

## 配置选项

### 环境变量配置

```bash
# 必需配置
export GOOGLE_CLOUD_PROJECT="your-project-id"
export GOOGLE_APPLICATION_CREDENTIALS="./service-account-key.json"

# 可选配置
export AMAPI_CALLBACK_URL="https://your-app.com/callback"
export AMAPI_TIMEOUT="30s"
export AMAPI_RETRY_ATTEMPTS="3"
export AMAPI_ENABLE_RETRY="true"
export AMAPI_RATE_LIMIT="100"
export AMAPI_LOG_LEVEL="info"
```

### YAML 配置文件

```yaml
# amapi.yaml
project_id: "your-project-id"
credentials_file: "./service-account-key.json"
callback_url: "https://your-app.com/callback"
timeout: "30s"
retry_attempts: 3
retry_delay: "1s"
enable_retry: true
rate_limit: 100
rate_burst: 10
log_level: "info"
enable_debug_logging: false
```

### 编程配置

```go
cfg := &config.Config{
    ProjectID:               "your-project-id",
    CredentialsFile:         "./service-account-key.json",
    CallbackURL:             "https://your-app.com/callback",
    Timeout:                 30 * time.Second,
    RetryAttempts:          3,
    RetryDelay:             1 * time.Second,
    EnableRetry:            true,
    RateLimit:              100,
    RateBurst:              10,
    LogLevel:               "info",
    EnableDebugLogging:     false,
}
```

## 最佳实践

### 1. 配置管理
```go
// 使用自动配置加载
cfg, err := config.AutoLoadConfig()
if err != nil {
    // 回退到环境变量
    cfg, err = config.LoadFromEnv()
}
```

### 2. 错误处理
```go
// 始终检查和处理错误
devices, err := client.Devices().ListByEnterpriseID(enterpriseID, nil)
if err != nil {
    log.Printf("Failed to list devices: %v", err)
    return
}
```

### 3. 资源清理
```go
client, err := client.New(cfg)
if err != nil {
    log.Fatal(err)
}
defer client.Close() // 确保清理资源
```

### 4. 上下文使用
```go
ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
defer cancel()

client, err := client.NewWithContext(ctx, cfg)
```

### 5. 策略验证
```go
// 创建策略前验证
if err := policy.Validate(); err != nil {
    log.Printf("Policy validation failed: %v", err)
    return
}
```

## 高级用例

### 1. 批量设备管理
```go
// 获取所有非合规设备并发送重启命令
nonCompliantDevices, err := client.Devices().GetNonCompliantDevices(enterpriseID)
for _, device := range nonCompliantDevices.Items {
    err := client.Devices().RebootByID(enterpriseID, device.GetID())
    if err != nil {
        log.Printf("Failed to reboot device %s: %v", device.GetID(), err)
    }
}
```

### 2. 动态策略更新
```go
// 根据合规状态调整策略
policies, _ := client.Policies().ListByEnterpriseID(enterpriseID, nil)
for _, policy := range policies.Items {
    devicesUsingPolicy, _ := client.Policies().GetDevicesUsingPolicy(policy.Name)

    // 计算合规率
    compliantCount := 0
    for _, device := range devicesUsingPolicy.Items {
        if device.PolicyCompliant {
            compliantCount++
        }
    }

    complianceRate := float64(compliantCount) / float64(len(devicesUsingPolicy.Items))

    // 如果合规率低于 80%，加强策略
    if complianceRate < 0.8 {
        policy.CameraDisabled = true
        policy.BluetoothDisabled = true

        updated, err := client.Policies().UpdateByID(
            enterpriseID,
            policy.GetID(),
            &policy,
        )
        if err == nil {
            log.Printf("Tightened policy %s due to low compliance", policy.GetID())
        }
    }
}
```

### 3. 令牌生命周期管理
```go
// 定期清理过期令牌
tokens, _ := client.EnrollmentTokens().ListByEnterpriseID(enterpriseID, nil)
for _, token := range tokens.Items {
    if token.IsExpired() {
        err := client.EnrollmentTokens().RevokeTokenByID(
            enterpriseID,
            token.GetID(),
        )
        if err == nil {
            log.Printf("Cleaned up expired token %s", token.GetID())
        }
    }
}
```

## 故障排除

### 常见问题

1. **认证失败**
   - 检查服务账户密钥文件路径
   - 确认服务账户具有 Android Management API 权限
   - 验证项目 ID 正确

2. **请求超时**
   - 增加超时配置
   - 检查网络连接
   - 启用重试机制

3. **API 限制**
   - 配置合适的速率限制
   - 实现退避重试
   - 监控 API 配额使用

4. **设备不合规**
   - 检查策略配置
   - 验证设备状态报告
   - 确认策略推送成功

### 调试技巧

```go
// 启用调试日志
cfg.EnableDebugLogging = true
cfg.LogLevel = "debug"

// 使用客户端健康检查
if err := client.Health(); err != nil {
    log.Printf("Client health check failed: %v", err)
}

// 获取客户端信息
info := client.GetInfo()
log.Printf("Client info: %+v", info)
```

## 集成示例

参考 `pkgs/amapi/examples/` 目录中的完整示例：

- `basic_usage.go` - 基本用法演示
- `enterprise_setup.go` - 完整企业设置流程
- `device_management.go` - 设备管理操作
- `policy_templates.go` - 策略模板使用

这些示例展示了如何在实际应用中使用 amapi 包的各种功能。