# AMAPI SDK 使用文档

Android Management API (AMAPI) Go SDK 是一个功能全面、生产就绪的 Go 客户端库，用于管理企业环境中的 Android 设备。

## 目录

- [特性概览](#特性概览)
- [安装](#安装)
- [快速开始](#快速开始)
- [配置说明](#配置说明)
- [核心功能](#核心功能)
  - [企业管理](#企业管理)
  - [策略管理](#策略管理)
  - [设备管理](#设备管理)
  - [注册令牌](#注册令牌)
  - [迁移令牌](#迁移令牌)
  - [Web 应用](#web-应用)
  - [Web 令牌](#web-令牌)
  - [配置信息](#配置信息)
- [策略预设](#策略预设)
- [高级功能](#高级功能)
- [错误处理](#错误处理)
- [最佳实践](#最佳实践)
- [API 参考](#api-参考)

## 特性概览

### 核心特性

- ✅ **完整的 API 覆盖** - 支持 Android Management API 的所有功能
- ✅ **类型安全** - 完整的类型定义和验证
- ✅ **自动重试** - 内置智能重试逻辑
- ✅ **速率限制** - 自动速率控制和突发容量管理
- ✅ **灵活配置** - 支持环境变量、YAML、JSON 多种配置方式
- ✅ **错误处理** - 详细的错误类型和处理机制
- ✅ **Context 支持** - 完整的 context.Context 支持
- ✅ **策略预设** - 8 种预配置策略模板

### 支持的功能模块

| 功能模块 | 描述 | 状态 |
|---------|------|------|
| 企业管理 | 创建、管理企业，生成注册 URL | ✅ |
| 策略管理 | CRUD 操作，策略预设，应用管理 | ✅ |
| 设备管理 | 设备查询、远程控制、命令执行 | ✅ |
| 注册令牌 | 创建、管理令牌，生成 QR 码 | ✅ |
| 迁移令牌 | EMM 迁移管理 | ✅ |
| Web 应用 | Web 应用配置和管理 | ✅ |
| Web 令牌 | 浏览器访问令牌 | ✅ |
| 配置信息 | 设备配置信息查询 | ✅ |

## 安装

### 使用 go get

```bash
go get amapi-pkg/pkgs/amapi
```

### 导入到项目

```go
import (
    "amapi-pkg/pkgs/amapi"
    "amapi-pkg/pkgs/amapi/client"
    "amapi-pkg/pkgs/amapi/config"
    "amapi-pkg/pkgs/amapi/types"
    "amapi-pkg/pkgs/amapi/presets"
)
```

## 快速开始

### 最简单的示例

```go
package main

import (
    "log"
    "amapi-pkg/pkgs/amapi/client"
    "amapi-pkg/pkgs/amapi/config"
)

func main() {
    // 1. 自动加载配置（从环境变量或配置文件）
    cfg, err := config.AutoLoadConfig()
    if err != nil {
        log.Fatal(err)
    }

    // 2. 创建客户端
    c, err := client.New(cfg)
    if err != nil {
        log.Fatal(err)
    }
    defer c.Close()

    // 3. 测试连接
    if err := c.Health(); err != nil {
        log.Fatal("健康检查失败:", err)
    }

    // 4. 使用客户端
    enterprises, err := c.Enterprises().List(nil)
    if err != nil {
        log.Fatal(err)
    }

    log.Printf("找到 %d 个企业", len(enterprises.Items))
}
```

### 手动配置示例

```go
package main

import (
    "log"
    "time"
    "amapi-pkg/pkgs/amapi/client"
    "amapi-pkg/pkgs/amapi/config"
)

func main() {
    // 手动创建配置
    cfg := &config.Config{
        ProjectID:       "your-project-id",
        CredentialsFile: "./sa-key.json",
        CallbackURL:     "https://your-app.com/callback",

        // 超时和重试配置
        Timeout:        30 * time.Second,
        RetryAttempts:  3,
        RetryDelay:     1 * time.Second,
        EnableRetry:    true,

        // 速率限制
        RateLimit:      100, // 每分钟 100 次请求
        RateBurst:      10,  // 允许突发 10 次

        // 日志配置
        LogLevel:       "info",
        EnableDebugLogging: false,
    }

    // 创建客户端
    c, err := client.New(cfg)
    if err != nil {
        log.Fatal(err)
    }
    defer c.Close()

    // 使用客户端...
}
```

## 配置说明

### 配置方式

AMAPI SDK 支持三种配置方式，按优先级排序：

1. **环境变量** - 最高优先级
2. **配置文件** - 中等优先级（支持 YAML 和 JSON）
3. **程序化配置** - 最低优先级（代码中直接配置）

### 环境变量配置

```bash
# 必需配置
export GOOGLE_CLOUD_PROJECT="your-project-id"
export GOOGLE_APPLICATION_CREDENTIALS="./sa-key.json"

# 可选配置
export AMAPI_CALLBACK_URL="https://your-app.com/callback"
export AMAPI_TIMEOUT="30s"
export AMAPI_RETRY_ATTEMPTS="3"
export AMAPI_RETRY_DELAY="1s"
export AMAPI_ENABLE_RETRY="true"
export AMAPI_RATE_LIMIT="100"
export AMAPI_RATE_BURST="10"
export AMAPI_LOG_LEVEL="info"
export AMAPI_ENABLE_DEBUG_LOGGING="false"
```

### YAML 配置文件

创建 `config.yaml` 文件：

```yaml
# Google Cloud 配置
project_id: "your-project-id"
credentials_file: "./sa-key.json"
# 或使用 credentials_json (二选一)
# credentials_json: '{"type": "service_account", ...}'

# API 配置
service_account_email: "sa@project.iam.gserviceaccount.com"
scopes:
  - "https://www.googleapis.com/auth/androidmanagement"

# 客户端配置
timeout: "30s"
retry_attempts: 3
retry_delay: "1s"
enable_retry: true

# 回调配置
callback_url: "https://your-app.com/callback"

# 缓存配置（可选）
enable_cache: false
cache_ttl: "5m"

# 日志配置
log_level: "info"  # debug, info, warn, error
enable_debug_logging: false

# 速率限制
rate_limit: 100  # 每分钟请求数
rate_burst: 10   # 突发容量
```

### 配置文件搜索路径

SDK 会按以下顺序自动搜索配置文件：

1. `./config.yaml`
2. `./config.yml`
3. `./amapi.yaml`
4. `./amapi.yml`
5. `~/.config/amapi/config.yaml`
6. `~/.config/amapi/config.yml`
7. `/etc/amapi/config.yaml`
8. `/etc/amapi/config.yml`

### 加载配置的方法

```go
// 方法 1: 自动加载（推荐）
// 按优先级尝试环境变量 -> 配置文件 -> 默认值
cfg, err := config.AutoLoadConfig()

// 方法 2: 从环境变量加载
cfg, err := config.LoadFromEnv()

// 方法 3: 从指定配置文件加载
cfg, err := config.LoadFromFile("./config.yaml")

// 方法 4: 使用默认配置
cfg := config.DefaultConfig()
cfg.ProjectID = "your-project-id"
cfg.CredentialsFile = "./sa-key.json"
```

### 配置参数详解

| 参数 | 类型 | 必需 | 默认值 | 说明 |
|------|------|------|--------|------|
| `project_id` | string | ✅ | - | GCP 项目 ID |
| `credentials_file` | string | ✅* | - | 服务账号密钥文件路径 |
| `credentials_json` | string | ✅* | - | 服务账号密钥 JSON 内容 |
| `service_account_email` | string | ❌ | - | 服务账号邮箱 |
| `scopes` | []string | ❌ | androidmanagement | OAuth2 权限范围 |
| `timeout` | duration | ❌ | 30s | API 请求超时时间 |
| `retry_attempts` | int | ❌ | 3 | 重试次数 |
| `retry_delay` | duration | ❌ | 1s | 基础重试延迟 |
| `enable_retry` | bool | ❌ | true | 是否启用重试 |
| `callback_url` | string | ❌ | - | 企业注册回调 URL |
| `enable_cache` | bool | ❌ | false | 是否启用缓存 |
| `cache_ttl` | duration | ❌ | 5m | 缓存有效期 |
| `log_level` | string | ❌ | info | 日志级别 |
| `enable_debug_logging` | bool | ❌ | false | 是否启用调试日志 |
| `rate_limit` | int | ❌ | 100 | 每分钟请求数限制 |
| `rate_burst` | int | ❌ | 10 | 突发请求容量 |

> **注意**: `credentials_file` 和 `credentials_json` 二选一即可

## 核心功能

### 企业管理

企业是 Android Management API 的核心概念，代表一个组织的设备管理实体。

#### 生成企业注册 URL

```go
// 生成注册 URL，让管理员完成企业注册
signupReq := &types.SignupURLRequest{
    ProjectID:             "your-project-id",
    CallbackURL:           "https://your-app.com/callback",
    AdminEmail:            "admin@company.com",
    EnterpriseDisplayName: "Example Company",
}

signupURL, err := c.Enterprises().GenerateSignupURL(signupReq)
if err != nil {
    log.Fatal(err)
}

log.Printf("请访问此 URL 完成注册: %s", signupURL.URL)
```

#### 使用回调令牌创建企业

```go
// 在回调 URL 接收到 token 后创建企业
createReq := &types.EnterpriseCreateRequest{
    SignupToken: "token-from-callback",
    ProjectID:   "your-project-id",
    DisplayName: "My Company",
}

enterprise, err := c.Enterprises().Create(createReq)
if err != nil {
    log.Fatal(err)
}

log.Printf("企业创建成功: %s", enterprise.Name)
```

#### 列出所有企业

```go
// 列出项目下的所有企业
enterprises, err := c.Enterprises().List(nil)
if err != nil {
    log.Fatal(err)
}

for _, enterprise := range enterprises.Items {
    log.Printf("企业: %s (%s)", enterprise.DisplayName, enterprise.Name)
}
```

#### 获取特定企业

```go
// 通过企业 ID 获取详情
enterprise, err := c.Enterprises().GetByID("LC00abc123")
if err != nil {
    log.Fatal(err)
}

log.Printf("企业名称: %s", enterprise.DisplayName)
log.Printf("企业 ID: %s", enterprise.GetID())
log.Printf("联系邮箱: %s", enterprise.ContactInfo.ContactEmail)
```

#### 更新企业信息

```go
updateReq := &types.EnterpriseUpdateRequest{
    DisplayName: "New Company Name",
    ContactInfo: &types.ContactInfo{
        ContactEmail:               "admin@company.com",
        DataProtectionOfficerName:  "John Doe",
        DataProtectionOfficerEmail: "dpo@company.com",
        DataProtectionOfficerPhone: "+1-555-0123",
    },
}

updated, err := c.Enterprises().Update("enterprises/LC00abc123", updateReq)
if err != nil {
    log.Fatal(err)
}

log.Printf("企业信息已更新")
```

#### 配置通知

```go
// 启用企业通知
notificationTypes := []string{
    types.NotificationTypeEnrollment,      // 设备注册通知
    types.NotificationTypeComplianceReport, // 合规报告
    types.NotificationTypeStatusReport,     // 状态报告
    types.NotificationTypeCommand,          // 命令执行通知
}

enterprise, err := c.Enterprises().EnableNotifications(
    "enterprises/LC00abc123",
    notificationTypes,
)
if err != nil {
    log.Fatal(err)
}

log.Printf("已启用 %d 种通知类型", len(enterprise.EnabledNotificationTypes))
```

#### 配置 Pub/Sub 主题

```go
// 设置接收通知的 Pub/Sub 主题
topicName := "projects/your-project-id/topics/amapi-events"

enterprise, err := c.Enterprises().SetPubSubTopic(
    "enterprises/LC00abc123",
    topicName,
)
if err != nil {
    log.Fatal(err)
}

log.Printf("Pub/Sub 主题已设置: %s", enterprise.PubsubTopic)
```

### 策略管理

策略定义了设备的行为和限制。

#### 创建基础策略

```go
// 创建一个简单的策略
policy := &types.Policy{
    StatusBarDisabled:     false,
    KeyguardDisabled:      false,
    AddUserDisabled:       true,
    UninstallAppsDisabled: true,
    CameraDisabled:        false,
    BluetoothDisabled:     false,
    AutoTimeRequired:      true,
}

created, err := c.Policies().CreateByEnterpriseID(
    "LC00abc123",
    "basic-policy",
    policy,
)
if err != nil {
    log.Fatal(err)
}

log.Printf("策略已创建: %s", created.Name)
```

#### 使用预设创建策略

```go
// 使用完全托管设备预设
preset := presets.GetFullyManagedPreset()

policy, err := c.Policies().CreateByEnterpriseID(
    "LC00abc123",
    "fully-managed-policy",
    preset.Policy,
)
if err != nil {
    log.Fatal(err)
}

log.Printf("使用预设创建策略: %s", preset.Name)
```

#### 列出策略

```go
// 列出企业下的所有策略
policies, err := c.Policies().ListByEnterpriseID("LC00abc123", nil)
if err != nil {
    log.Fatal(err)
}

for _, policy := range policies.Items {
    log.Printf("策略: %s (模式: %s)", policy.GetID(), policy.GetPolicyMode())
}
```

#### 更新策略

```go
// 获取现有策略
policy, err := c.Policies().GetByID("LC00abc123", "my-policy")
if err != nil {
    log.Fatal(err)
}

// 修改策略设置
policy.CameraDisabled = true
policy.BluetoothDisabled = true

// 更新策略
updated, err := c.Policies().UpdateByID("LC00abc123", "my-policy", policy)
if err != nil {
    log.Fatal(err)
}

log.Printf("策略已更新")
```

#### 为策略添加应用

```go
// 添加必需应用
policy.AddApplication(types.NewRequiredApp("com.company.vpn"))
policy.AddApplication(types.NewRequiredApp("com.company.security"))

// 添加 Kiosk 模式应用
policy.AddApplication(types.NewKioskApp("com.company.kioskapp"))

// 阻止特定应用
policy.AddApplication(types.NewBlockedApp("com.facebook.katana"))
policy.AddApplication(types.NewBlockedApp("com.instagram.android"))

// 更新策略
updated, err := c.Policies().UpdateByID("LC00abc123", "my-policy", policy)
```

#### 删除策略

```go
err := c.Policies().DeleteByID("LC00abc123", "old-policy")
if err != nil {
    log.Fatal(err)
}

log.Printf("策略已删除")
```

#### 获取使用策略的设备

```go
// 查找使用特定策略的设备
devices, err := c.Policies().GetDevicesUsingPolicy(
    "enterprises/LC00abc123/policies/my-policy",
)
if err != nil {
    log.Fatal(err)
}

log.Printf("使用此策略的设备数: %d", len(devices.Items))
```

### 设备管理

设备管理功能允许您查询、监控和控制已注册的设备。

#### 列出设备

```go
// 列出企业下的所有设备
devices, err := c.Devices().ListByEnterpriseID("LC00abc123", nil)
if err != nil {
    log.Fatal(err)
}

for _, device := range devices.Items {
    log.Printf("设备: %s", device.GetID())
    log.Printf("  状态: %s", device.State)
    log.Printf("  合规: %t", device.PolicyCompliant)
    log.Printf("  用户: %s", device.UserName)
    log.Printf("  型号: %s", device.GetDeviceModel())
    log.Printf("  Android 版本: %s", device.GetAndroidVersion())
}
```

#### 获取设备详情

```go
// 获取特定设备的详细信息
device, err := c.Devices().GetByID("LC00abc123", "A1B2C3D4E5F6")
if err != nil {
    log.Fatal(err)
}

log.Printf("设备名称: %s", device.Name)
log.Printf("IMEI: %s", device.GetIMEI())
log.Printf("序列号: %s", device.GetSerialNumber())
log.Printf("是否在线: %t", device.IsOnline())
log.Printf("最后同步: %s", device.GetLastSyncTime())
```

#### 设备命令

```go
// 锁定设备（10 分钟）
err := c.Devices().LockByID("LC00abc123", "device-id", "PT10M")
if err != nil {
    log.Fatal(err)
}

// 重启设备
err = c.Devices().RebootByID("LC00abc123", "device-id")
if err != nil {
    log.Fatal(err)
}

// 重置设备（恢复出厂设置）
err = c.Devices().ResetByID("LC00abc123", "device-id")
if err != nil {
    log.Fatal(err)
}

// 启动丢失模式
lostModeReq := &types.StartLostModeRequest{
    Message:     "此设备已丢失，请联系 IT 部门",
    PhoneNumber: "+1-555-0123",
}
err = c.Devices().StartLostMode("LC00abc123", "device-id", lostModeReq)
if err != nil {
    log.Fatal(err)
}

// 停止丢失模式
err = c.Devices().StopLostMode("LC00abc123", "device-id")
```

#### 查询合规设备

```go
// 获取合规设备
compliantDevices, err := c.Devices().GetCompliantDevices("LC00abc123")
if err != nil {
    log.Fatal(err)
}

log.Printf("合规设备数: %d", len(compliantDevices.Items))

// 获取不合规设备
nonCompliantDevices, err := c.Devices().GetNonCompliantDevices("LC00abc123")
if err != nil {
    log.Fatal(err)
}

log.Printf("不合规设备数: %d", len(nonCompliantDevices.Items))

// 显示不合规原因
for _, device := range nonCompliantDevices.Items {
    log.Printf("设备 %s 不合规:", device.GetID())
    for _, violation := range device.NonComplianceDetails {
        log.Printf("  - %s", violation.SettingName)
    }
}
```

#### 删除设备

```go
// 从企业中移除设备
err := c.Devices().DeleteByID("LC00abc123", "device-id")
if err != nil {
    log.Fatal(err)
}

log.Printf("设备已移除")
```

### 注册令牌

注册令牌用于将新设备注册到企业。

#### 创建注册令牌

```go
// 创建一个有效期为 24 小时的注册令牌
token, err := c.EnrollmentTokens().CreateByEnterpriseID(
    "LC00abc123",
    "policy-id",
    24 * time.Hour,
)
if err != nil {
    log.Fatal(err)
}

log.Printf("令牌已创建: %s", token.Value)
log.Printf("过期时间: %s", token.ExpirationTimestamp)
```

#### 高级令牌创建

```go
// 创建具有自定义选项的令牌
req := &types.EnrollmentTokenCreateRequest{
    EnterpriseName:     "enterprises/LC00abc123",
    PolicyName:         "enterprises/LC00abc123/policies/my-policy",
    Duration:           7 * 24 * time.Hour, // 7 天
    AllowPersonalUsage: true,                // BYOD 模式
    OneTimeOnly:        true,                // 一次性使用
}

token, err := c.EnrollmentTokens().Create(req)
if err != nil {
    log.Fatal(err)
}

log.Printf("高级令牌已创建: %s", token.GetID())
```

#### 生成 QR 码

```go
// 生成包含 WiFi 配置的 QR 码
qrOptions := &types.QRCodeOptions{
    WiFiSSID:         "CompanyWiFi",
    WiFiPassword:     "password123",
    WiFiSecurityType: types.WiFiSecurityTypeWPA2,
    SkipSetupWizard:  true,
    Locale:           "zh_CN",
}

qrData, err := c.EnrollmentTokens().GenerateQRCodeByID(
    "LC00abc123",
    "token-id",
    qrOptions,
)
if err != nil {
    log.Fatal(err)
}

// 获取 QR 码 JSON 数据
qrJSON, _ := qrData.ToJSON()
log.Printf("QR 码数据: %s", qrJSON)

// 或者获取可扫描的字符串
qrString := qrData.ToQRString()
log.Printf("QR 码字符串: %s", qrString)
```

#### 列出活动令牌

```go
// 获取所有活动的注册令牌
activeTokens, err := c.EnrollmentTokens().GetActiveTokens("LC00abc123")
if err != nil {
    log.Fatal(err)
}

for _, token := range activeTokens.Items {
    log.Printf("令牌: %s", token.GetID())
    log.Printf("  策略: %s", token.GetPolicyID())
    log.Printf("  剩余时间: %v", token.TimeUntilExpiration())
}
```

#### 获取令牌统计

```go
// 获取令牌统计信息
stats, err := c.EnrollmentTokens().GetTokenStatistics("LC00abc123")
if err != nil {
    log.Fatal(err)
}

log.Printf("总令牌数: %d", stats.TotalTokens)
log.Printf("活动令牌: %d", stats.ActiveTokens)
log.Printf("已过期令牌: %d", stats.ExpiredTokens)
```

#### 撤销令牌

```go
// 删除/撤销注册令牌
err := c.EnrollmentTokens().DeleteByID("LC00abc123", "token-id")
if err != nil {
    log.Fatal(err)
}

log.Printf("令牌已撤销")
```

### 迁移令牌

迁移令牌用于从其他 EMM（企业移动管理）系统迁移设备。

#### 创建迁移令牌

```go
// 创建迁移令牌
req := &types.MigrationTokenCreateRequest{
    EnterpriseName: "enterprises/LC00abc123",
    PolicyName:     "enterprises/LC00abc123/policies/migration-policy",
    Duration:       30 * 24 * time.Hour, // 30 天
}

token, err := c.MigrationTokens().Create(req)
if err != nil {
    log.Fatal(err)
}

log.Printf("迁移令牌: %s", token.Value)
```

#### 列出迁移令牌

```go
tokens, err := c.MigrationTokens().ListByEnterpriseID("LC00abc123", nil)
if err != nil {
    log.Fatal(err)
}

for _, token := range tokens.Items {
    log.Printf("迁移令牌: %s (过期: %s)",
        token.GetID(),
        token.ExpirationTimestamp)
}
```

### Web 应用

管理企业 Web 应用。

#### 创建 Web 应用

```go
webAppReq := &types.WebAppCreateRequest{
    EnterpriseName: "enterprises/LC00abc123",
    DisplayName:    "Company Portal",
    StartURL:       "https://portal.company.com",
    DisplayMode:    types.WebAppDisplayModeStandalone,
    Icons: []types.WebAppIcon{
        {
            ImageData: "base64-encoded-image",
        },
    },
}

webApp, err := c.WebApps().Create(webAppReq)
if err != nil {
    log.Fatal(err)
}

log.Printf("Web 应用已创建: %s", webApp.Name)
```

#### 列出 Web 应用

```go
webApps, err := c.WebApps().ListByEnterpriseID("LC00abc123", nil)
if err != nil {
    log.Fatal(err)
}

for _, app := range webApps.Items {
    log.Printf("Web 应用: %s - %s", app.DisplayName, app.StartURL)
}
```

### Web 令牌

用于生成管理员访问企业管理界面的临时令牌。

#### 创建 Web 令牌

```go
req := &types.WebTokenCreateRequest{
    EnterpriseName:  "enterprises/LC00abc123",
    ParentFrameURL:  "https://admin.company.com",
    Permissions:     []string{"MANAGE_POLICIES", "MANAGE_DEVICES"},
}

webToken, err := c.WebTokens().Create(req)
if err != nil {
    log.Fatal(err)
}

// 重定向管理员到此 URL
log.Printf("管理员访问 URL: %s", webToken.Value)
```

### 配置信息

查询设备配置和预配置信息。

#### 获取配置信息

```go
provisioningInfo, err := c.ProvisioningInfo().Get("your-project-id")
if err != nil {
    log.Fatal(err)
}

log.Printf("配置信息名称: %s", provisioningInfo.Name)
```

## 策略预设

SDK 提供 8 种预配置的策略模板，适用于不同的使用场景。

### 可用的预设

```go
// 获取所有预设
allPresets := presets.GetAllPresets()

for _, preset := range allPresets {
    log.Printf("预设: %s - %s", preset.Name, preset.Description)
}
```

### 预设列表

| 预设名称 | 适用场景 | 特点 |
|----------|----------|------|
| `fully_managed` | 企业完全托管设备 | 标准企业策略，平衡安全性和功能性 |
| `dedicated_device` | 专用设备/信息亭 | 锁定模式，限制用户操作 |
| `work_profile` | BYOD 工作配置文件 | 分离工作和个人数据 |
| `kiosk_mode` | 单应用信息亭 | 仅允许运行一个应用 |
| `cope` | 企业拥有，个人使用 | 允许个人使用的企业设备 |
| `secure_workstation` | 高安全性工作站 | 严格的安全限制 |
| `education_tablet` | 教育平板 | 适合学校和教育机构 |
| `retail_kiosk` | 零售终端 | 销售和客户服务终端 |

### 使用预设

```go
// 1. 获取完全托管设备预设
preset := presets.GetFullyManagedPreset()

// 2. 直接使用预设创建策略
policy, err := c.Policies().CreateByEnterpriseID(
    "LC00abc123",
    "my-policy",
    preset.Policy,
)

// 3. 或者基于预设进行自定义
customPolicy := preset.Policy.Clone()
customPolicy.CameraDisabled = true
customPolicy.AddApplication(types.NewRequiredApp("com.company.app"))

policy, err = c.Policies().CreateByEnterpriseID(
    "LC00abc123",
    "custom-policy",
    customPolicy,
)
```

### 预设详情

#### 1. 完全托管设备（fully_managed）

```go
preset := presets.GetFullyManagedPreset()

// 特性：
// - 企业完全控制设备
// - 允许必要的设备功能
// - 强制自动更新
// - 支持应用管理
```

#### 2. 专用设备（dedicated_device）

```go
preset := presets.GetDedicatedDevicePreset()

// 特性：
// - 锁定模式
// - 禁用大部分用户设置
// - 适合固定用途设备
// - 自动启动指定应用
```

#### 3. 工作配置文件（work_profile）

```go
preset := presets.GetWorkProfilePreset()

// 特性：
// - BYOD 模式
// - 工作和个人数据分离
// - 用户保留设备控制权
// - 企业只管理工作配置文件
```

#### 4. Kiosk 模式（kiosk_mode）

```go
preset := presets.GetKioskModePreset()

// 特性：
// - 单应用模式
// - 禁用所有系统 UI
// - 防止退出应用
// - 适合公共展示
```

#### 5. 企业拥有个人使用（cope）

```go
preset := presets.GetCOPEPreset()

// 特性：
// - 允许个人使用
// - 企业保留设备所有权
// - 平衡工作和生活
// - 可配置使用限制
```

#### 6. 安全工作站（secure_workstation）

```go
preset := presets.GetSecureWorkstationPreset()

// 特性：
// - 最高安全级别
// - 禁用大部分功能
// - 强制加密
// - 严格的应用白名单
```

#### 7. 教育平板（education_tablet）

```go
preset := presets.GetEducationTabletPreset()

// 特性：
// - 适合学校使用
// - 允许教育应用
// - 限制社交媒体
// - 家长控制选项
```

#### 8. 零售终端（retail_kiosk）

```go
preset := presets.GetRetailKioskPreset()

// 特性：
// - 销售终端优化
// - 客户交互友好
// - 支付应用支持
// - 限制非业务功能
```

### 自定义预设策略

```go
// 从预设开始，添加自定义配置
preset := presets.GetFullyManagedPreset()
policy := preset.Policy.Clone()

// 添加公司应用
policy.AddApplication(types.NewRequiredApp("com.company.vpn"))
policy.AddApplication(types.NewRequiredApp("com.company.mail"))
policy.AddApplication(types.NewRequiredApp("com.company.chat"))

// 阻止特定应用
policy.AddApplication(types.NewBlockedApp("com.facebook.katana"))
policy.AddApplication(types.NewBlockedApp("com.instagram.android"))

// 调整设置
policy.CameraDisabled = true
policy.BluetoothDisabled = false
policy.WifiDisabled = false

// 创建策略
created, err := c.Policies().CreateByEnterpriseID(
    "LC00abc123",
    "custom-company-policy",
    policy,
)
```

## 高级功能

### Context 支持

所有 API 调用都支持 context.Context，可用于超时控制和取消操作。

```go
// 创建带超时的 context
ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
defer cancel()

// 使用 context 创建客户端
c, err := client.NewWithContext(ctx, cfg)
if err != nil {
    log.Fatal(err)
}

// 所有操作都会遵守 context 的超时设置
enterprises, err := c.Enterprises().List(nil)
if err != nil {
    if ctx.Err() == context.DeadlineExceeded {
        log.Println("操作超时")
    }
}
```

### 自动重试

SDK 内置智能重试机制，自动处理临时性错误。

```go
cfg := &config.Config{
    ProjectID:       "your-project-id",
    CredentialsFile: "./sa-key.json",

    // 重试配置
    EnableRetry:   true,
    RetryAttempts: 5,                  // 最多重试 5 次
    RetryDelay:    2 * time.Second,    // 基础延迟 2 秒
    // SDK 会使用指数退避算法自动调整延迟
}

c, err := client.New(cfg)

// API 调用会自动重试（如果遇到可重试的错误）
devices, err := c.Devices().ListByEnterpriseID("LC00abc123", nil)
// 如果失败，SDK 会自动重试最多 5 次
```

### 速率限制

SDK 自动处理 API 速率限制，防止超出配额。

```go
cfg := &config.Config{
    ProjectID:       "your-project-id",
    CredentialsFile: "./sa-key.json",

    // 速率限制配置
    RateLimit: 100,  // 每分钟 100 个请求
    RateBurst: 20,   // 允许突发 20 个请求
}

c, err := client.New(cfg)

// SDK 会自动控制请求速率
for i := 0; i < 200; i++ {
    // 即使循环 200 次，SDK 也会自动限速
    // 确保不超过每分钟 100 次的限制
    devices, err := c.Devices().ListByEnterpriseID("LC00abc123", nil)
    if err != nil {
        log.Printf("错误: %v", err)
    }
}
```

### 健康检查

```go
// 检查 API 连接状态
if err := c.Health(); err != nil {
    log.Printf("健康检查失败: %v", err)
    // 采取补救措施...
} else {
    log.Println("API 连接正常")
}
```

### 客户端信息

```go
// 获取客户端信息
info := c.GetInfo()

log.Printf("SDK 版本: %s", info.Version)
log.Printf("项目 ID: %s", info.ProjectID)
log.Printf("用户代理: %s", info.UserAgent)
log.Printf("支持的功能: %v", info.Capabilities)
log.Printf("创建时间: %s", info.CreatedAt)
```

### 配置克隆

```go
// 获取客户端配置的副本
cfg := c.GetConfig()

// 修改配置创建新客户端
cfg.RetryAttempts = 10
newClient, err := client.New(cfg)
```

## 错误处理

### 错误类型

SDK 提供了详细的错误类型系统。

```go
devices, err := c.Devices().GetByID("LC00abc123", "invalid-device-id")
if err != nil {
    // 类型断言为 API 错误
    if apiErr, ok := err.(*types.Error); ok {
        log.Printf("错误代码: %d", apiErr.Code)
        log.Printf("错误消息: %s", apiErr.Message)
        log.Printf("错误类型: %s", apiErr.GetErrorType())

        // 检查是否可重试
        if apiErr.IsRetryable() {
            delay := apiErr.RetryDelay(1, time.Second)
            log.Printf("可重试，建议延迟: %v", delay)
        }

        // 获取原始错误
        if apiErr.Cause != nil {
            log.Printf("原始错误: %v", apiErr.Cause)
        }
    }
}
```

### 常见错误代码

```go
switch apiErr.Code {
case types.ErrCodeNotFound:
    log.Println("资源未找到")

case types.ErrCodeInvalidInput:
    log.Println("输入参数无效")

case types.ErrCodeUnauthorized:
    log.Println("认证失败")

case types.ErrCodeForbidden:
    log.Println("权限不足")

case types.ErrCodeTooManyRequests:
    log.Println("请求过于频繁")
    delay := apiErr.RetryDelay(1, time.Second)
    time.Sleep(delay)
    // 重试...

case types.ErrCodeInternalServerError:
    log.Println("服务器内部错误")

case types.ErrCodeServiceUnavailable:
    log.Println("服务暂时不可用")
}
```

### 错误处理最佳实践

```go
func handleDeviceOperation(c *client.Client, enterpriseID, deviceID string) error {
    device, err := c.Devices().GetByID(enterpriseID, deviceID)
    if err != nil {
        // 1. 检查是否为 API 错误
        apiErr, ok := err.(*types.Error)
        if !ok {
            return fmt.Errorf("未知错误: %w", err)
        }

        // 2. 根据错误类型处理
        switch {
        case apiErr.IsNotFound():
            return fmt.Errorf("设备不存在: %s", deviceID)

        case apiErr.IsAuthenticationError():
            return fmt.Errorf("认证失败，请检查凭证")

        case apiErr.IsPermissionError():
            return fmt.Errorf("权限不足，需要更高级别的权限")

        case apiErr.IsRateLimitError():
            delay := apiErr.RetryDelay(1, time.Second)
            log.Printf("触发速率限制，等待 %v 后重试", delay)
            time.Sleep(delay)
            return handleDeviceOperation(c, enterpriseID, deviceID)

        case apiErr.IsRetryable():
            log.Printf("临时错误，将自动重试")
            return err

        default:
            return fmt.Errorf("API 错误: %w", err)
        }
    }

    // 处理设备...
    log.Printf("设备信息: %+v", device)
    return nil
}
```

## 最佳实践

### 1. 资源管理

```go
// 始终关闭客户端以释放资源
c, err := client.New(cfg)
if err != nil {
    log.Fatal(err)
}
defer c.Close()  // 确保资源被释放
```

### 2. 配置验证

```go
// 在使用配置前进行验证
cfg := &config.Config{
    ProjectID:       "your-project-id",
    CredentialsFile: "./sa-key.json",
}

// 验证配置
if err := cfg.Validate(); err != nil {
    log.Fatalf("配置无效: %v", err)
}

c, err := client.New(cfg)
```

### 3. 错误处理

```go
// 使用类型断言处理特定错误
devices, err := c.Devices().ListByEnterpriseID(enterpriseID, nil)
if err != nil {
    if apiErr, ok := err.(*types.Error); ok {
        // 处理 API 特定错误
        if apiErr.IsRetryable() {
            // 重试逻辑
        }
    }
    return err
}
```

### 4. Context 使用

```go
// 使用 context 控制超时
ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
defer cancel()

c, err := client.NewWithContext(ctx, cfg)
if err != nil {
    log.Fatal(err)
}
defer c.Close()

// 所有操作都会遵守 context 的超时
```

### 5. 策略验证

```go
// 在创建策略前验证
policy := &types.Policy{
    // ... 策略配置
}

if err := policy.Validate(); err != nil {
    log.Fatalf("策略无效: %v", err)
}

created, err := c.Policies().CreateByEnterpriseID(
    enterpriseID,
    "my-policy",
    policy,
)
```

### 6. 批量操作

```go
// 批量创建令牌
func createMultipleTokens(c *client.Client, enterpriseID, policyID string, count int) error {
    for i := 0; i < count; i++ {
        token, err := c.EnrollmentTokens().CreateByEnterpriseID(
            enterpriseID,
            policyID,
            24*time.Hour,
        )
        if err != nil {
            return fmt.Errorf("创建第 %d 个令牌失败: %w", i+1, err)
        }
        log.Printf("令牌 %d 已创建: %s", i+1, token.Value)

        // 添加小延迟避免速率限制
        time.Sleep(100 * time.Millisecond)
    }
    return nil
}
```

### 7. 日志记录

```go
// 使用适当的日志级别
cfg := &config.Config{
    ProjectID:       "your-project-id",
    CredentialsFile: "./sa-key.json",
    LogLevel:        "info",  // 生产环境使用 info
    EnableDebugLogging: false, // 调试时设为 true
}
```

### 8. 安全实践

```go
// 不要在代码中硬编码敏感信息
// ❌ 错误示例
cfg := &config.Config{
    ProjectID:       "my-project-123",  // 不要硬编码
    CredentialsJSON: `{"type": "service_account", ...}`,  // 不要硬编码
}

// ✅ 正确示例
cfg, err := config.AutoLoadConfig()  // 从环境变量或文件加载
if err != nil {
    log.Fatal(err)
}
```

## API 参考

### Client 方法

```go
// 创建客户端
c, err := client.New(cfg)
c, err := client.NewWithContext(ctx, cfg)

// 客户端管理
err := c.Health()           // 健康检查
info := c.GetInfo()         // 获取客户端信息
cfg := c.GetConfig()        // 获取配置副本
err := c.Close()            // 关闭客户端
```

### 企业 API

```go
// 企业服务
svc := c.Enterprises()

// 方法
signupURL, err := svc.GenerateSignupURL(req)
enterprise, err := svc.Create(req)
enterprise, err := svc.Get(name)
enterprise, err := svc.GetByID(id)
enterprises, err := svc.List(req)
enterprise, err := svc.Update(name, req)
err := svc.Delete(name)
err := svc.DeleteByID(id)
enterprise, err := svc.EnableNotifications(name, types)
enterprise, err := svc.SetPubSubTopic(name, topic)
```

### 策略 API

```go
// 策略服务
svc := c.Policies()

// 方法
policy, err := svc.Create(name, policy)
policy, err := svc.CreateByEnterpriseID(enterpriseID, policyID, policy)
policy, err := svc.Get(name)
policy, err := svc.GetByID(enterpriseID, policyID)
policies, err := svc.List(parent, req)
policies, err := svc.ListByEnterpriseID(enterpriseID, req)
policy, err := svc.Update(name, policy)
policy, err := svc.UpdateByID(enterpriseID, policyID, policy)
err := svc.Delete(name)
err := svc.DeleteByID(enterpriseID, policyID)
devices, err := svc.GetDevicesUsingPolicy(policyName)
```

### 设备 API

```go
// 设备服务
svc := c.Devices()

// 方法
device, err := svc.Get(name)
device, err := svc.GetByID(enterpriseID, deviceID)
devices, err := svc.List(parent, req)
devices, err := svc.ListByEnterpriseID(enterpriseID, req)
err := svc.Delete(name)
err := svc.DeleteByID(enterpriseID, deviceID)
err := svc.Lock(name, duration)
err := svc.LockByID(enterpriseID, deviceID, duration)
err := svc.Reboot(name)
err := svc.RebootByID(enterpriseID, deviceID)
err := svc.Reset(name)
err := svc.ResetByID(enterpriseID, deviceID)
devices, err := svc.GetCompliantDevices(enterpriseID)
devices, err := svc.GetNonCompliantDevices(enterpriseID)
err := svc.StartLostMode(enterpriseID, deviceID, req)
err := svc.StopLostMode(enterpriseID, deviceID)
```

### 注册令牌 API

```go
// 注册令牌服务
svc := c.EnrollmentTokens()

// 方法
token, err := svc.Create(req)
token, err := svc.CreateByEnterpriseID(enterpriseID, policyID, duration)
token, err := svc.Get(name)
token, err := svc.GetByID(enterpriseID, tokenID)
tokens, err := svc.List(parent, req)
tokens, err := svc.ListByEnterpriseID(enterpriseID, req)
err := svc.Delete(name)
err := svc.DeleteByID(enterpriseID, tokenID)
qrData, err := svc.GenerateQRCode(tokenName, options)
qrData, err := svc.GenerateQRCodeByID(enterpriseID, tokenID, options)
tokens, err := svc.GetActiveTokens(enterpriseID)
stats, err := svc.GetTokenStatistics(enterpriseID)
```

### 迁移令牌 API

```go
// 迁移令牌服务
svc := c.MigrationTokens()

// 方法
token, err := svc.Create(req)
token, err := svc.Get(name)
tokens, err := svc.List(parent, req)
tokens, err := svc.ListByEnterpriseID(enterpriseID, req)
```

### Web 应用 API

```go
// Web 应用服务
svc := c.WebApps()

// 方法
webApp, err := svc.Create(req)
webApp, err := svc.Get(name)
webApps, err := svc.List(parent, req)
webApps, err := svc.ListByEnterpriseID(enterpriseID, req)
webApp, err := svc.Update(name, webApp)
err := svc.Delete(name)
```

### Web 令牌 API

```go
// Web 令牌服务
svc := c.WebTokens()

// 方法
webToken, err := svc.Create(req)
```

### 配置信息 API

```go
// 配置信息服务
svc := c.ProvisioningInfo()

// 方法
info, err := svc.Get(name)
```

## 完整示例

这里提供一个完整的企业设置工作流示例：

```go
package main

import (
    "log"
    "time"

    "amapi-pkg/pkgs/amapi/client"
    "amapi-pkg/pkgs/amapi/config"
    "amapi-pkg/pkgs/amapi/presets"
    "amapi-pkg/pkgs/amapi/types"
)

func main() {
    // 1. 加载配置
    cfg, err := config.AutoLoadConfig()
    if err != nil {
        log.Fatal(err)
    }

    // 2. 创建客户端
    c, err := client.New(cfg)
    if err != nil {
        log.Fatal(err)
    }
    defer c.Close()

    // 3. 检查连接
    if err := c.Health(); err != nil {
        log.Fatal("健康检查失败:", err)
    }

    // 4. 生成企业注册 URL
    signupReq := &types.SignupURLRequest{
        ProjectID:             cfg.ProjectID,
        CallbackURL:           cfg.CallbackURL,
        AdminEmail:            "admin@company.com",
        EnterpriseDisplayName: "Example Company",
    }

    signupURL, err := c.Enterprises().GenerateSignupURL(signupReq)
    if err != nil {
        log.Fatal(err)
    }
    log.Printf("注册 URL: %s", signupURL.URL)

    // 5. 假设企业已创建，获取企业列表
    enterprises, err := c.Enterprises().List(nil)
    if err != nil {
        log.Fatal(err)
    }

    if len(enterprises.Items) == 0 {
        log.Println("请先完成企业注册")
        return
    }

    enterprise := enterprises.Items[0]
    enterpriseID := enterprise.GetID()
    log.Printf("使用企业: %s", enterprise.DisplayName)

    // 6. 创建策略
    preset := presets.GetFullyManagedPreset()
    policy := preset.Policy.Clone()

    // 自定义策略
    policy.AddApplication(types.NewRequiredApp("com.company.vpn"))
    policy.CameraDisabled = false
    policy.BluetoothDisabled = false

    createdPolicy, err := c.Policies().CreateByEnterpriseID(
        enterpriseID,
        "company-policy",
        policy,
    )
    if err != nil {
        log.Fatal(err)
    }
    log.Printf("策略已创建: %s", createdPolicy.GetID())

    // 7. 创建注册令牌
    token, err := c.EnrollmentTokens().CreateByEnterpriseID(
        enterpriseID,
        createdPolicy.GetID(),
        24*time.Hour,
    )
    if err != nil {
        log.Fatal(err)
    }
    log.Printf("注册令牌: %s", token.Value)

    // 8. 生成 QR 码
    qrOptions := &types.QRCodeOptions{
        WiFiSSID:         "CompanyWiFi",
        WiFiPassword:     "password123",
        WiFiSecurityType: types.WiFiSecurityTypeWPA2,
        SkipSetupWizard:  true,
        Locale:           "zh_CN",
    }

    qrData, err := c.EnrollmentTokens().GenerateQRCodeByID(
        enterpriseID,
        token.GetID(),
        qrOptions,
    )
    if err != nil {
        log.Fatal(err)
    }

    qrJSON, _ := qrData.ToJSON()
    log.Printf("QR 码数据: %s", qrJSON)

    // 9. 列出设备（如果有）
    devices, err := c.Devices().ListByEnterpriseID(enterpriseID, nil)
    if err != nil {
        log.Fatal(err)
    }

    log.Printf("已注册设备数: %d", len(devices.Items))
    for _, device := range devices.Items {
        log.Printf("  设备: %s (合规: %t)",
            device.GetID(),
            device.PolicyCompliant)
    }

    log.Println("设置完成！")
}
```

## GoDoc 文档

本 SDK 包含完整的 godoc 格式注释文档。

### 查看文档

```bash
# 方式 1: 使用便捷脚本（推荐）
cd pkgs/amapi
./docs.sh serve
# 然后访问 http://localhost:6060/pkg/amapi-pkg/pkgs/amapi/

# 方式 2: 查看包概览
./docs.sh

# 方式 3: 查看完整文档
./docs.sh all

# 方式 4: 查看特定类型
./docs.sh type Client
./docs.sh type Config

# 方式 5: 使用 go doc 命令
go doc
go doc -all
go doc Client
go doc config.Config
```

详细说明请查看 [GODOC.md](./GODOC.md)。

### 文档特点

✅ **全中文注释** - 所有导出的类型、函数都有中文说明
✅ **详细示例** - 包含完整的代码示例
✅ **多种查看方式** - 命令行、Web 界面、编辑器集成
✅ **完整覆盖** - 涵盖所有主要包和功能

## 更多资源

### 官方文档

- [Android Management API 官方文档](https://developers.google.com/android/management)
- [Google Cloud 文档](https://cloud.google.com/docs)
- [服务账号最佳实践](https://cloud.google.com/iam/docs/best-practices-service-accounts)

### 项目文档

- [快速开始指南](../../docs/QUICKSTART.md)
- [CLI 使用手册](../../docs/CLI_USAGE.md)
- [构建指南](../../docs/BUILD_GUIDE.md)
- [安全指南](../../docs/SECURITY.md)
- [GoDoc 使用指南](./GODOC.md) ⭐ 新增

### 示例代码

- [基础使用示例](./examples/basic_usage.go)
- [企业设置示例](./examples/enterprise_setup.go)

## 文档文件

本 SDK 包含以下文档：

- `README.md` - SDK 完整使用文档（本文件）
- `GODOC.md` - GoDoc 文档查看指南
- `docs.sh` - 文档查看便捷脚本
- `generate_docs.sh` - 文档生成脚本

## 支持与反馈

如有问题或建议：

- 在 GitHub 上创建 Issue
- 查看 GoDoc 文档（运行 `./docs.sh serve`）
- 查看示例代码
- 阅读官方 API 文档

## 许可证

本项目采用 MIT 许可证。详见项目根目录的 LICENSE 文件。

