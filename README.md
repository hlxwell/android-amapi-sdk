# Android Management API Go 客户端

这是一个全面、生产就绪的 Google Android Management API Go 客户端库。此包为企业环境中的 Android 设备管理提供了高级接口，包含完整的 SDK 和命令行工具。

## 功能特性

- **完整的 API 覆盖**：100% 覆盖 Google Android Management API 所有功能
  - 企业、策略、设备、注册令牌
  - 迁移令牌、Web 应用、Web 令牌、配置信息
- **功能齐全的 CLI**：11 个命令模块，80+ 个子命令
- **灵活的配置**：支持环境变量、YAML 和 JSON 配置
- **内置可靠性**：自动重试逻辑、速率限制和错误处理
- **类型安全**：完整的类型定义和验证
- **策略预设**：8 种常见用例的预配置模板
- **丰富的文档**：全面的中文文档和示例

## 项目组成

### 1. AMAPI SDK
位于 `pkgs/amapi/` 目录，提供 Go 语言的 Android Management API 客户端库。

### 2. 命令行工具
位于 `cmd/amapi-cli/` 目录，提供完整的命令行界面来演示和使用 SDK 功能。

## 快速体验

### 使用命令行工具

```bash
# 方式 1：使用 Makefile（推荐）
make build                # 构建 CLI 工具到 build/ 目录
./build/amapi-cli --help  # 查看帮助

# 方式 2：使用 go build
go build -o build/amapi-cli ./cmd/amapi-cli

# 使用示例
./build/amapi-cli config show                      # 检查配置
./build/amapi-cli health check                     # 健康检查
./build/amapi-cli enterprise list your-project-id  # 列出企业
./build/amapi-cli enterprise signup-url --project-id your-project-id

# 其他 Makefile 命令
make clean        # 清理构建文件
make build-all    # 跨平台构建
make test         # 运行测试
make install      # 安装到系统
make help         # 查看所有命令
```

### 命令行工具功能

- **企业管理** (`enterprise`)：创建、查看、更新、删除企业，注册 URL，通知管理
- **策略管理** (`policy`)：CRUD 操作，8种预设模板，应用管理，模式切换
- **设备管理** (`device`)：查看、远程控制、删除、操作管理
- **注册令牌** (`enrollment`)：创建、管理令牌，生成 QR 码，批量创建
- **迁移令牌** (`migration`)：管理从其他 EMM 迁移的令牌 ⭐ 新增
- **Web 应用** (`webapp`)：管理企业 Web 应用 ⭐ 新增
- **Web 令牌** (`webtoken`)：管理浏览器访问令牌 ⭐ 新增
- **配置信息** (`provisioning`)：查询设备配置信息 ⭐ 新增
- **配置管理** (`config`)：配置验证、环境变量管理
- **健康检查** (`health`)：API 连接测试、配置验证

详细使用说明请参考：[CLI 使用指南](docs/CLI_USAGE.md)

构建部署指南请参考：[构建指南](docs/BUILD_GUIDE.md)

## ⚠️ 安全提醒

**重要**：本项目需要使用 Google Cloud 服务账号密钥。请务必：

- ✅ **不要**将 `sa-key.json` 提交到版本控制
- ✅ **不要**在代码中硬编码项目 ID 和敏感信息
- ✅ **使用** `.gitignore` 保护敏感文件
- ✅ **参考** `sa-key.json.example` 创建你的密钥文件
- ✅ **阅读** [SECURITY.md](docs/SECURITY.md) 了解安全最佳实践

## 安装

### SDK 安装

```bash
go get github.com/hlxwell/android-api-demo/pkgs/amapi
```

### CLI 工具构建

```bash
# 使用 Makefile（推荐）
make build

# 或使用 go build
go build -o build/amapi-cli ./cmd/amapi-cli
```

## 快速开始

### 基本设置

```go
package main

import (
    "context"
    "log"

    "github.com/hlxwell/android-api-demo/pkgs/amapi/client"
    "github.com/hlxwell/android-api-demo/pkgs/amapi/config"
)

func main() {
    // 加载配置
    cfg, err := config.AutoLoadConfig()
    if err != nil {
        log.Fatal(err)
    }

    // 创建客户端
    c, err := client.New(cfg)
    if err != nil {
        log.Fatal(err)
    }
    defer c.Close()

    // 使用客户端
    enterprises, err := c.Enterprises().List(nil)
    if err != nil {
        log.Fatal(err)
    }

    log.Printf("找到 %d 个企业", len(enterprises.Items))
}
```

### 配置

#### 环境变量

```bash
export GOOGLE_CLOUD_PROJECT="your-project-id"
export GOOGLE_APPLICATION_CREDENTIALS="./service-account-key.json"
export AMAPI_CALLBACK_URL="https://your-app.com/callback"
export AMAPI_LOG_LEVEL="info"
```

#### 配置文件 (config.yaml)

```bash
# 复制示例配置文件
cp config.yaml.example config.yaml

# 编辑配置文件
vi config.yaml
```

```yaml
project_id: "your-project-id"
credentials_file: "./sa-key.json"
callback_url: "https://your-app.com/callback"
timeout: "30s"
retry_attempts: 3
enable_retry: true
log_level: "info"
```

#### 程序化配置

```go
cfg := &config.Config{
    ProjectID:       "your-project-id",
    CredentialsFile: "./service-account-key.json",
    CallbackURL:     "https://your-app.com/callback",
    Timeout:         30 * time.Second,
    RetryAttempts:   3,
    EnableRetry:     true,
    LogLevel:        "info",
}
```

## 核心组件

### 企业管理

```go
// 生成注册 URL
signupReq := &types.SignupURLRequest{
    ProjectID:   "your-project-id",
    CallbackURL: "https://your-app.com/callback",
}
signupURL, err := client.Enterprises().GenerateSignupURL(signupReq)

// 注册后创建企业
createReq := &types.EnterpriseCreateRequest{
    SignupToken: "token-from-callback",
    ProjectID:   "your-project-id",
    DisplayName: "My Company",
}
enterprise, err := client.Enterprises().Create(createReq)

// 列出企业
enterprises, err := client.Enterprises().List(nil)

// 获取特定企业
enterprise, err := client.Enterprises().GetByID("enterprise-id")
```

### 策略管理

```go
// 从预设创建策略
preset := presets.GetFullyManagedPreset()
policy, err := client.Policies().CreateByEnterpriseID(
    "enterprise-id",
    "my-policy",
    preset.Policy,
)

// 更新策略
policy.CameraDisabled = true
updated, err := client.Policies().UpdateByID(
    "enterprise-id",
    "my-policy",
    policy,
)

// 列出策略
policies, err := client.Policies().ListByEnterpriseID("enterprise-id", nil)
```

### 设备管理

```go
// 列出设备
devices, err := client.Devices().ListByEnterpriseID("enterprise-id", nil)

// 获取设备详情
device, err := client.Devices().GetByID("enterprise-id", "device-id")

// 发送命令
err = client.Devices().LockByID("enterprise-id", "device-id", "PT10M") // 10 分钟
err = client.Devices().RebootByID("enterprise-id", "device-id")
err = client.Devices().ResetByID("enterprise-id", "device-id")

// 获取合规状态
compliantDevices, err := client.Devices().GetCompliantDevices("enterprise-id")
nonCompliantDevices, err := client.Devices().GetNonCompliantDevices("enterprise-id")
```

### 注册令牌管理

```go
// 创建注册令牌
token, err := client.EnrollmentTokens().CreateByEnterpriseID(
    "enterprise-id",
    "policy-id",
    24*time.Hour, // 有效期 24 小时
)

// 生成二维码
qrOptions := &types.QRCodeOptions{
    WiFiSSID:        "CompanyWiFi",
    WiFiPassword:    "password123",
    WiFiSecurityType: types.WiFiSecurityTypeWPA2,
    SkipSetupWizard: true,
}
qrData, err := client.EnrollmentTokens().GenerateQRCodeByID(
    "enterprise-id",
    "token-id",
    qrOptions,
)

// 列出活动令牌
tokens, err := client.EnrollmentTokens().GetActiveTokens("enterprise-id")
```

## 策略预设

本库包含常见场景的预配置策略模板：

```go
// 可用的预设
presets := presets.GetAllPresets()

// 特定预设
fullyManaged := presets.GetFullyManagedPreset()
dedicatedDevice := presets.GetDedicatedDevicePreset()
workProfile := presets.GetWorkProfilePreset()
kioskMode := presets.GetKioskModePreset()

// 从预设创建策略并自定义
customizations := map[string]interface{}{
    "camera_disabled": true,
    "bluetooth_disabled": false,
}
policy, err := presets.CreatePolicyFromPreset("fully_managed", customizations)
```

### 可用的预设

- **fully_managed**: 标准企业设备策略
- **dedicated_device**: 锁定的信息亭模式
- **work_profile**: BYOD（自带设备办公）工作配置文件
- **kiosk_mode**: 单应用信息亭
- **cope**: 企业拥有，个人使用
- **secure_workstation**: 高安全性配置
- **education_tablet**: 针对教育场景优化
- **retail_kiosk**: 销售终端和客户交互

## 高级功能

### 错误处理

```go
devices, err := client.Devices().ListByEnterpriseID("enterprise-id", nil)
if err != nil {
    if apiErr, ok := err.(*types.Error); ok {
        switch apiErr.Code {
        case types.ErrCodeNotFound:
            log.Println("企业未找到")
        case types.ErrCodeTooManyRequests:
            log.Println("速率受限，重试时间:", apiErr.RetryDelay(1, time.Second))
        default:
            log.Printf("API 错误: %s", apiErr.Error())
        }
    }
}
```

### 重试和速率限制

```go
cfg := &config.Config{
    // ... 其他配置
    RetryAttempts: 5,
    RetryDelay:    2 * time.Second,
    EnableRetry:   true,
    RateLimit:     200, // 每分钟请求数
    RateBurst:     20,  // 突发容量
}
```

### Context 支持

```go
ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
defer cancel()

client, err := client.NewWithContext(ctx, cfg)
```

## 配置参考

### 环境变量

| 变量 | 描述 | 默认值 |
|----------|-------------|---------|
| `GOOGLE_CLOUD_PROJECT` | GCP 项目 ID | 必需 |
| `GOOGLE_APPLICATION_CREDENTIALS` | 服务账号密钥路径 | 必需 |
| `AMAPI_CALLBACK_URL` | 企业注册回调 URL | "" |
| `AMAPI_TIMEOUT` | API 请求超时时间 | "30s" |
| `AMAPI_RETRY_ATTEMPTS` | 重试次数 | 3 |
| `AMAPI_RETRY_DELAY` | 基础重试延迟 | "1s" |
| `AMAPI_ENABLE_RETRY` | 启用重试逻辑 | true |
| `AMAPI_RATE_LIMIT` | 每分钟请求数 | 100 |
| `AMAPI_RATE_BURST` | 突发容量 | 10 |
| `AMAPI_LOG_LEVEL` | 日志级别 (debug,info,warn,error) | "info" |
| `AMAPI_ENABLE_DEBUG_LOGGING` | 启用调试日志 | false |

### 配置文件

库会按以下顺序自动搜索配置文件：

1. `./config.yaml`
2. `./config.yml`
3. `./amapi.yaml`
4. `./amapi.yml`
5. `~/.config/amapi/config.yaml`
6. `~/.config/amapi/config.yml`
7. `/etc/amapi/config.yaml`
8. `/etc/amapi/config.yml`

## 📚 文档

- [📖 快速开始](docs/QUICKSTART.md) - 5分钟快速设置指南
- [🔨 构建指南](docs/BUILD_GUIDE.md) - 如何构建和部署
- [📋 CLI 使用手册](docs/CLI_USAGE.md) - 命令行工具完整使用文档
- [📖 使用指南](docs/USAGE_GUIDE.md) - SDK 详细使用说明
- [🔐 安全指南](docs/SECURITY.md) - 安全最佳实践
- [📝 脱敏总结](docs/DESENSITIZATION_SUMMARY.md) - 代码脱敏报告

## 配置示例

项目根目录提供了配置文件示例：

- `config.yaml.example` - YAML 配置文件模板
- `.env.example` - 环境变量配置模板
- `sa-key.json.example` - 服务账号密钥文件模板

复制并修改这些示例文件来配置你的环境。

## 测试

```bash
# 运行测试
go test ./...

# 运行测试并显示覆盖率
go test -cover ./...

# 运行集成测试（需要有效的 GCP 凭证）
go test -tags=integration ./...
```

## 系统要求

- Go 1.19 或更高版本
- 有效的 Google Cloud Platform 项目
- 具有 Android Management API 权限的服务账号

## 认证设置

### 1. 创建 GCP 项目

```bash
gcloud projects create your-project-id
gcloud config set project your-project-id
```

### 2. 启用 Android Management API

```bash
gcloud services enable androidmanagement.googleapis.com
```

### 3. 创建服务账号

```bash
gcloud iam service-accounts create amapi-service-account \
  --display-name="Android Management API Service Account"
```

### 4. 分配权限

```bash
gcloud projects add-iam-policy-binding your-project-id \
  --member="serviceAccount:amapi-service-account@your-project-id.iam.gserviceaccount.com" \
  --role="roles/androidmanagement.user"
```

### 5. 下载服务账号密钥

```bash
gcloud iam service-accounts keys create sa-key.json \
  --iam-account=amapi-service-account@your-project-id.iam.gserviceaccount.com
```

**⚠️ 注意**：密钥文件 `sa-key.json` 已在 `.gitignore` 中，不会被提交到 Git。

### 6. 设置环境变量

```bash
export GOOGLE_APPLICATION_CREDENTIALS="./sa-key.json"
export GOOGLE_CLOUD_PROJECT="your-project-id"
```

## 贡献

1. Fork 本仓库
2. 创建功能分支
3. 为新功能添加测试
4. 确保所有测试通过
5. 提交 Pull Request

## 许可证

本项目采用 MIT 许可证 - 详见 [LICENSE](LICENSE) 文件。

## 支持

如有问题或疑问：

- 在 GitHub 上创建 Issue
- 查看 [API 文档](https://developers.google.com/android/management)
- 查看 [examples](examples/) 目录
- 阅读 [文档](docs/)

## 更多资源

- [Google Cloud 文档](https://cloud.google.com/docs)
- [Android Management API 参考](https://developers.google.com/android/management/reference/rest)
- [服务账号最佳实践](https://cloud.google.com/iam/docs/best-practices-service-accounts)
