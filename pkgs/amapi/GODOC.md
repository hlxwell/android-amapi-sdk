# AMAPI SDK GoDoc 文档指南

本文档说明如何查看和使用 AMAPI SDK 的 GoDoc 文档。

## 文档已完成

✅ 所有主要包和类型都已添加完整的 godoc 格式注释：

- ✅ `amapi` - 主包
- ✅ `amapi/config` - 配置管理
- ✅ `amapi/client` - API 客户端
- ✅ `amapi/types` - 类型定义
- ✅ `amapi/presets` - 策略预设

## 查看文档的方式

### 方式 1: 使用快捷脚本（推荐）

我们提供了便捷的脚本来查看文档：

```bash
cd pkgs/amapi

# 查看包概览
./docs.sh

# 查看完整文档
./docs.sh all

# 查看特定类型
./docs.sh type Client
./docs.sh type Config
./docs.sh type NewClient

# 查看子包文档
./docs.sh config
./docs.sh client

# 启动 Web 界面（推荐！）
./docs.sh serve
# 然后访问 http://localhost:6060/pkg/amapi-pkg/pkgs/amapi/
```

### 方式 2: 使用 go doc 命令

```bash
cd pkgs/amapi

# 查看包文档
go doc

# 查看完整文档
go doc -all

# 查看特定类型或函数
go doc Client
go doc NewClient
go doc Config
go doc config.LoadFromFile
go doc config.DefaultConfig

# 查看子包
go doc config
go doc client
go doc types
go doc presets
```

### 方式 3: 启动 godoc Web 服务器

这是最佳的文档浏览方式，提供完整的网页界面：

```bash
# 方式 A: 使用脚本
cd pkgs/amapi
./docs.sh serve

# 方式 B: 手动启动
# 首先安装 godoc（如果未安装）
go install golang.org/x/tools/cmd/godoc@latest

# 从项目根目录启动
cd /path/to/amapi-pkg
godoc -http=:6060
```

然后在浏览器中访问：

- **主包**: http://localhost:6060/pkg/amapi-pkg/pkgs/amapi/
- **配置包**: http://localhost:6060/pkg/amapi-pkg/pkgs/amapi/config/
- **客户端包**: http://localhost:6060/pkg/amapi-pkg/pkgs/amapi/client/
- **类型包**: http://localhost:6060/pkg/amapi-pkg/pkgs/amapi/types/
- **预设包**: http://localhost:6060/pkg/amapi-pkg/pkgs/amapi/presets/

### 方式 4: 在编辑器中查看

大多数现代 Go 编辑器都支持查看 godoc：

**VS Code**:
- 鼠标悬停在类型或函数上
- 使用 `Go: Show Documentation` 命令

**GoLand**:
- Ctrl+Q (Windows/Linux) 或 Cmd+J (Mac)
- 查看 Quick Documentation

**Vim/Neovim (with gopls)**:
- 使用 LSP 的 hover 功能
- `:GoDoc` 命令

## 文档结构

### 主包 (amapi)

主包提供了所有功能的统一入口点：

```go
import "amapi-pkg/pkgs/amapi"

// 创建客户端
c, err := amapi.NewClient(cfg)

// 使用导出的类型
var enterprise amapi.Enterprise
var policy amapi.Policy
var device amapi.Device
```

主要内容：
- 包级文档和快速开始示例
- 所有核心类型的导出
- 便捷的请求/响应结构体
- NewClient 函数

### 配置包 (config)

配置管理的完整文档：

```bash
go doc config
go doc config.Config
go doc config.DefaultConfig
go doc config.LoadFromFile
go doc config.AutoLoadConfig
```

主要内容：
- Config 结构体的所有字段说明
- 配置加载的各种方式
- 环境变量列表
- 配置验证规则
- 配置文件格式示例

### 客户端包 (client)

API 客户端的详细文档：

```bash
go doc client
go doc client.Client
go doc client.New
```

主要内容：
- Client 结构和方法
- 企业、策略、设备等服务的访问方法
- 错误处理
- 重试和速率限制机制

### 类型包 (types)

所有数据类型的定义：

```bash
go doc types
go doc types.Enterprise
go doc types.Policy
go doc types.Device
```

主要内容：
- 核心数据结构
- 请求和响应类型
- 错误类型
- 辅助方法

### 预设包 (presets)

策略预设模板：

```bash
go doc presets
go doc presets.GetFullyManagedPreset
```

主要内容：
- 8 种预配置策略模板
- 各预设的使用场景
- 预设自定义方法

## 文档示例

### 查看 NewClient 函数

```bash
$ go doc NewClient

package amapi // import "amapi-pkg/pkgs/amapi"

func NewClient(cfg *Config) (*Client, error)
    NewClient 创建一个新的 Android Management API 客户端。

    参数 cfg 包含客户端配置，包括项目 ID、认证凭证等。
    返回的客户端是线程安全的，可以在多个 goroutine 中共享使用。

    在使用完毕后应该调用 Close() 方法释放资源：

        c, err := NewClient(cfg)
        if err != nil {
            return err
        }
        defer c.Close()

    如果配置无效或认证失败，将返回错误。
```

### 查看 Config 结构

```bash
$ go doc config.Config

package config // import "amapi-pkg/pkgs/amapi/config"

type Config struct {
    ProjectID               string
    CredentialsFile         string
    CredentialsJSON         string
    ServiceAccountEmail     string
    Scopes                  []string
    Timeout                 time.Duration
    RetryAttempts          int
    RetryDelay             time.Duration
    EnableRetry            bool
    CallbackURL            string
    EnableCache            bool
    CacheTTL               time.Duration
    LogLevel               string
    EnableDebugLogging     bool
    RateLimit              int
    RateBurst              int
}
    Config 包含 Android Management API 客户端的所有配置选项。

    配置可以通过多种方式提供：环境变量、配置文件或程序化创建。
    使用 Validate() 方法可以验证配置的完整性和有效性。

func DefaultConfig() *Config
func LoadFromFile(path string) (*Config, error)
...
```

## 在线查看完整文档

启动 godoc 服务器后，你可以：

1. **浏览包列表** - 查看所有可用的包
2. **搜索功能** - 在页面顶部搜索类型、函数等
3. **查看源码** - 点击类型或函数名可查看源代码
4. **跳转链接** - 文档中的类型引用都是可点击的链接
5. **示例代码** - 查看内嵌的代码示例

## 文档编写规范

本项目的 godoc 注释遵循以下规范：

1. **包注释** - 每个包的 package 语句前有详细的包级文档
2. **类型注释** - 所有导出的类型都有完整说明
3. **函数注释** - 所有导出的函数都有参数、返回值和用法说明
4. **字段注释** - 结构体的导出字段都有行内注释
5. **代码示例** - 重要的功能都包含可执行的示例代码
6. **章节标题** - 使用 `#` 创建章节结构
7. **列表格式** - 使用 `-` 创建项目列表
8. **代码块** - 使用缩进创建代码示例

## 常用命令速查

```bash
# 快速查看
./docs.sh                  # 包概览
./docs.sh all              # 完整文档
./docs.sh type Client      # 特定类型

# go doc 命令
go doc                     # 当前包文档
go doc -all                # 完整文档
go doc Type                # 特定类型
go doc Function            # 特定函数
go doc package.Type        # 子包中的类型

# Web 界面
./docs.sh serve            # 启动服务器
./docs.sh serve 8080       # 指定端口

# 生成静态 HTML（高级）
godoc -url=/pkg/amapi-pkg/pkgs/amapi/ > amapi_docs.html
```

## 总结

✅ **文档已完成** - 所有主要包都有完整的 godoc 注释
✅ **多种查看方式** - 命令行、Web 界面、编辑器集成
✅ **详细说明** - 包括参数、返回值、示例代码
✅ **中文文档** - 全中文注释，易于理解
✅ **便捷工具** - 提供脚本快速查看文档

开始使用：

```bash
cd pkgs/amapi
./docs.sh serve
# 访问 http://localhost:6060/pkg/amapi-pkg/pkgs/amapi/
```

享受完整的 API 文档！ 📚

