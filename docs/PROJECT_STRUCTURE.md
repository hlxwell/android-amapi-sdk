# 项目结构说明

## 📁 目录结构

```
amapi-pkg/
├── cmd/                    # 命令行工具
│   └── amapi-cli/         # AMAPI CLI 工具
│       ├── cmd/           # CLI 命令实现
│       ├── internal/      # 内部工具
│       └── main.go        # 入口文件
│
├── pkgs/                  # 代码包
│   └── amapi/            # AMAPI SDK
│       ├── client/       # API 客户端实现
│       ├── config/       # 配置管理
│       ├── examples/     # 示例代码
│       ├── presets/      # 预设配置
│       ├── types/        # 类型定义
│       ├── utils/        # 工具函数
│       ├── amapi.go      # 主包入口
│       ├── README.md     # SDK 文档
│       └── GODOC.md      # GoDoc 指南
│
├── scripts/               # 脚本工具 ⭐
│   ├── docs.sh          # 统一文档工具（查看、验证、生成）
│   └── README.md        # 脚本使用说明
│
├── docs/                  # 项目文档 ⭐
│   ├── README.md                   # 文档索引
│   ├── QUICKSTART.md              # 快速开始
│   ├── CLI_USAGE.md               # CLI 使用手册
│   ├── USAGE_GUIDE.md             # 使用指南
│   ├── BUILD_GUIDE.md             # 构建指南
│   ├── SECURITY.md                # 安全指南
│   ├── PROJECT_SUMMARY.md         # 项目总结
│   ├── DESENSITIZATION_SUMMARY.md # 脱敏总结
│   ├── AMAPI_快速开始.md           # AMAPI 快速开始
│   ├── AMAPI_文档完成总结.md       # AMAPI 文档总结
│   └── PROJECT_STRUCTURE.md       # 项目结构（本文件）
│
├── terraform/             # Terraform 配置
│   ├── main.tf           # 主配置
│   ├── variables.tf      # 变量定义
│   ├── outputs.tf        # 输出定义
│   ├── README.md         # 完整文档
│   ├── QUICK_START.md    # 快速开始
│   ├── SUMMARY.md        # 配置总结
│   ├── Makefile          # Make 命令
│   └── terraform.tfvars.example  # 配置示例
│
├── build/                 # 构建输出
│   └── amapi-cli         # 编译后的二进制文件
│
├── config.yaml           # 配置文件
├── config.yaml.example   # 配置示例
├── sa-key.json          # Service Account Key
├── sa-key.json.example  # Key 示例
├── Makefile             # 项目构建工具
├── go.mod               # Go 模块定义
├── go.sum               # Go 依赖锁定
├── README.md            # 项目主文档
└── AGENTS.md            # 项目规范 ⭐
```

## 🎯 目录规范

根据 [AGENTS.md](../AGENTS.md) 的规定：

### ✅ 脚本文件位置
**规则**: 所有脚本文件必须放在 `/scripts` 目录下

**包含**:
- `docs.sh` - 统一文档工具（查看、验证、生成）
- 未来添加的其他脚本

**使用方法**:
```bash
# 从项目根目录运行
./scripts/docs.sh help          # 查看帮助
./scripts/docs.sh verify        # 验证编译
./scripts/docs.sh serve         # 启动服务器
./scripts/docs.sh generate      # 完整流程
```

### ✅ 文档文件位置
**规则**: 所有文档文件必须放在 `/docs` 目录下

**包含**:
- 项目级文档（快速开始、使用指南等）
- AMAPI SDK 相关文档
- Terraform 相关文档（在 `/terraform` 子目录）
- 开发和构建文档

**特殊说明**:
- `README.md` 保留在各自的模块根目录
- `pkgs/amapi/README.md` - SDK 专用文档
- `terraform/README.md` - Terraform 专用文档

### ✅ 构建验证
**规则**: 所有代码修改必须通过 `make build` 验证

**验证命令**:
```bash
make build
```

**预期输出**:
```
🔨 开始构建 amapi-cli...
✅ 构建完成: build/amapi-cli
```

## 📦 主要模块说明

### 1. CLI 工具 (`cmd/amapi-cli/`)
命令行工具，提供企业管理、设备控制、策略配置等功能。

**主要命令**:
- `enterprise` - 企业管理
- `device` - 设备管理
- `policy` - 策略管理
- `enrollment` - 注册管理
- `health` - 健康检查

### 2. AMAPI SDK (`pkgs/amapi/`)
Go 客户端库，提供完整的 Android Management API 功能。

**主要包**:
- `client` - API 客户端
- `config` - 配置管理
- `types` - 类型定义
- `presets` - 预设策略
- `utils` - 工具函数

### 3. Terraform 配置 (`terraform/`)
基础设施即代码配置，用于自动化部署 GCP 资源。

**创建的资源**:
- Pub/Sub Topics (CN & ROW)
- Service Account
- IAM 权限
- API 启用

### 4. 脚本工具 (`scripts/`)
开发和文档相关的脚本工具。

**可用脚本**:
- `docs.sh` - 统一文档工具（包含所有文档功能）

### 5. 文档 (`docs/`)
所有项目文档的集中位置。

**文档分类**:
- 快速开始和使用指南
- 开发和构建文档
- 安全和最佳实践
- AMAPI SDK 文档
- 项目总结和报告

## 🔄 工作流程

### 开发流程
```bash
# 1. 修改代码
vim pkgs/amapi/client/devices.go

# 2. 验证构建
make build

# 3. 测试
./build/amapi-cli health

# 4. 生成文档（可选）
./scripts/docs.sh generate
```

### 文档更新流程
```bash
# 1. 更新文档
vim docs/USAGE_GUIDE.md

# 2. 更新脚本
vim scripts/docs.sh

# 3. 验证构建（确保没有破坏代码）
make build
```

### Terraform 部署流程
```bash
# 1. 进入 terraform 目录
cd terraform

# 2. 初始化
terraform init

# 3. 部署
terraform apply
```

## 📝 文件命名规范

### Markdown 文件
- **大写**: 英文标题（如 `README.md`, `QUICKSTART.md`）
- **中文**: 使用下划线分隔（如 `AMAPI_快速开始.md`）

### Shell 脚本
- **小写**: 使用下划线或连字符分隔（如 `docs.sh`）
- **权限**: 确保可执行 `chmod +x`

### Go 文件
- **小写**: 使用下划线分隔（如 `rate_limiter.go`）
- **测试**: `_test.go` 后缀（如 `config_test.go`）

## 🎨 最佳实践

### 1. 添加新功能
```
1. 在 pkgs/amapi/ 中实现功能
2. 在 cmd/amapi-cli/ 中添加 CLI 命令
3. 在 docs/ 中更新文档
4. 运行 make build 验证
```

### 2. 添加新脚本
```
1. 创建脚本文件在 scripts/
2. 添加执行权限
3. 在 scripts/README.md 中记录
4. 在相关文档中引用
```

### 3. 添加新文档
```
1. 创建文档文件在 docs/
2. 在 docs/README.md 中添加链接
3. 更新相关的交叉引用
```

## 🔍 快速查找

### 我想...

**查看使用文档** → `docs/README.md`

**快速开始** → `docs/QUICKSTART.md`

**使用 CLI** → `docs/CLI_USAGE.md`

**查看 SDK API** → `pkgs/amapi/README.md`

**部署到 GCP** → `terraform/README.md`

**运行脚本** → `scripts/README.md`

**了解安全** → `docs/SECURITY.md`

**构建项目** → `docs/BUILD_GUIDE.md`

## 📊 统计信息

```
文件总数: ~80 个文件
文档文件: ~15 个 Markdown 文件
Go 代码: ~30 个 Go 文件
脚本文件: 2 个 Shell 脚本
Terraform: 8 个配置文件
```

## 🔗 相关链接

- [项目主页](../README.md)
- [文档索引](README.md)
- [脚本说明](../scripts/README.md)
- [Terraform 配置](../terraform/README.md)
- [项目规范](../AGENTS.md)

---

**最后更新**: 2025-10-30
**维护**: 按照 AGENTS.md 规范维护

