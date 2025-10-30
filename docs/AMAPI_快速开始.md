# AMAPI SDK 快速开始

欢迎使用 AMAPI SDK！本文档将帮助你快速上手。

## 📚 已完成的文档

我已经为你创建了以下完整的文档：

### 1️⃣ SDK 使用文档
- **文件**: `README.md` (40KB, 1700+ 行)
- **内容**: 完整的 SDK 使用指南，包含 50+ 个代码示例
- **涵盖**: 安装、配置、所有核心功能、API 参考、最佳实践

### 2️⃣ GoDoc 代码注释
- **文件**: `amapi.go`, `config/config.go`
- **特点**: 全中文 godoc 格式注释
- **包含**: 包文档、类型说明、函数文档、字段注释

### 3️⃣ GoDoc 使用指南
- **文件**: `GODOC.md`
- **内容**: 如何查看和使用 godoc 文档
- **包含**: 4 种查看方式、常用命令、文档特点

### 4️⃣ 文档查看工具
- **文件**: `docs.sh`, `generate_docs.sh`
- **功能**: 便捷查看和生成文档
- **特点**: 彩色输出、自动安装依赖

## 🚀 立即开始

### 查看 SDK 使用文档

```bash
# 在 VS Code 或其他编辑器中打开
cat README.md

# 或在浏览器中查看（如果有 Markdown 预览插件）
```

### 查看 GoDoc 文档（推荐）

```bash
# 方式 1: 启动 Web 界面（最佳）
cd /Users/hlxwell/projects/amapi-pkg/pkgs/amapi
./docs.sh serve

# 然后在浏览器访问:
# http://localhost:6060/pkg/amapi-pkg/pkgs/amapi/
```

### 快速查看命令

```bash
cd /Users/hlxwell/projects/amapi-pkg/pkgs/amapi

# 查看包概览
./docs.sh

# 查看完整文档
./docs.sh all

# 查看特定类型
./docs.sh type Client
./docs.sh type Config
./docs.sh type NewClient

# 查看子包
./docs.sh config
./docs.sh client

# 查看帮助
./docs.sh help
```

## 📖 文档内容概览

### README.md 包含:

1. **安装和配置**
   - 3 种配置方式（环境变量、YAML、程序化）
   - 详细的配置参数说明
   - 配置验证和最佳实践

2. **核心功能示例**
   - 企业管理（10+ 示例）
   - 策略管理（10+ 示例）
   - 设备管理（10+ 示例）
   - 注册令牌（8+ 示例）
   - 迁移令牌、Web 应用等

3. **策略预设**
   - 8 种预配置策略模板
   - 每种预设的详细说明
   - 自定义预设的方法

4. **高级功能**
   - Context 支持
   - 自动重试机制
   - 速率限制
   - 错误处理

5. **API 参考**
   - 所有客户端方法
   - 所有服务 API
   - 参数和返回值说明

6. **完整示例**
   - 端到端企业设置工作流
   - 实际应用场景

### GoDoc 文档包含:

- 包级文档和概述
- 所有导出类型的详细说明
- 所有函数的参数、返回值文档
- 代码使用示例
- 最佳实践建议

## 🎯 推荐的学习路径

### 第 1 步: 阅读快速开始
```bash
# 打开 README.md，阅读"快速开始"部分
# 了解基本的客户端创建和使用
```

### 第 2 步: 查看配置说明
```bash
# 阅读 README.md 中的"配置说明"部分
# 或查看 config 包文档:
./docs.sh config
```

### 第 3 步: 浏览核心功能
```bash
# 在 README.md 中查看感兴趣的功能
# 例如: 企业管理、策略管理、设备管理
```

### 第 4 步: 运行示例代码
```bash
# 查看 examples 目录中的示例:
cat examples/basic_usage.go
cat examples/enterprise_setup.go
```

### 第 5 步: 使用 GoDoc 深入学习
```bash
# 启动 godoc 服务器
./docs.sh serve

# 浏览详细的 API 文档
# http://localhost:6060/pkg/amapi-pkg/pkgs/amapi/
```

## 💡 常用操作

### 查看特定功能的文档

```bash
# 查看企业管理相关
# 在 README.md 中搜索 "### 企业管理"

# 查看策略管理相关
# 在 README.md 中搜索 "### 策略管理"

# 查看设备管理相关
# 在 README.md 中搜索 "### 设备管理"
```

### 查看 API 方法

```bash
# 使用 go doc
go doc client.Client
go doc config.Config

# 或使用脚本
./docs.sh type Client
./docs.sh type Config
```

### 查看配置选项

```bash
# 查看所有配置参数
./docs.sh type Config

# 或在 README.md 中查看"配置参数详解"表格
```

## 📝 文档文件说明

| 文件 | 大小 | 说明 |
|------|------|------|
| `README.md` | 40KB | SDK 完整使用文档 |
| `GODOC.md` | 7.3KB | GoDoc 查看指南 |
| `文档完成总结.md` | 8.2KB | 文档创建总结 |
| `快速开始.md` | 本文件 | 快速入门指南 |
| `docs.sh` | 3.8KB | 文档查看脚本 |
| `generate_docs.sh` | 2.2KB | 文档生成脚本 |

## 🔧 工具脚本说明

### docs.sh

**用途**: 快速查看文档

**使用方法**:
```bash
./docs.sh              # 查看包概览
./docs.sh all          # 查看完整文档
./docs.sh type <名称>  # 查看特定类型
./docs.sh serve        # 启动 Web 服务器
./docs.sh help         # 查看帮助
```

### generate_docs.sh

**用途**: 生成和验证文档

**使用方法**:
```bash
# 从项目根目录运行
./scripts/generate_docs.sh     # 验证代码并启动 godoc
```

## 🌟 推荐的使用方式

### 对于初学者

1. 先阅读 `README.md` 的"快速开始"部分
2. 运行几个基本示例
3. 查看"核心功能"部分学习具体用法
4. 使用 `./docs.sh serve` 查看详细文档

### 对于有经验的开发者

1. 快速浏览 `README.md` 了解 SDK 结构
2. 直接查看"API 参考"部分
3. 使用 `go doc` 或 `./docs.sh` 快速查询
4. 根据需要查看具体功能的示例

### 对于集成开发

1. 查看"配置说明"设置开发环境
2. 参考"完整示例"部分的工作流
3. 使用"最佳实践"部分的建议
4. 查看"错误处理"部分处理异常

## ❓ 常见问题

### 如何开始使用 SDK?

```bash
# 1. 查看快速开始
cat README.md | grep -A 50 "## 快速开始"

# 2. 或启动 godoc
./docs.sh serve
```

### 如何配置客户端?

```bash
# 查看配置文档
./docs.sh type Config

# 或查看 README.md 的配置章节
```

### 如何查看某个功能的用法?

```bash
# 在 README.md 中搜索相关章节
# 或使用 go doc:
go doc <类型或函数名>

# 或使用脚本:
./docs.sh type <名称>
```

### 在哪里可以找到代码示例?

```bash
# README.md 包含 50+ 个示例
# examples/ 目录包含完整的示例程序
cat examples/basic_usage.go
```

## 🎉 开始使用

现在你已经了解了所有可用的文档，可以开始使用 AMAPI SDK 了！

**推荐第一步**:

```bash
# 启动 godoc Web 界面
cd /Users/hlxwell/projects/amapi-pkg/pkgs/amapi
./docs.sh serve

# 然后在浏览器中打开:
# http://localhost:6060/pkg/amapi-pkg/pkgs/amapi/
```

**或者**:

```bash
# 在你喜欢的编辑器中打开 README.md
# 从"快速开始"章节开始阅读
```

祝你使用愉快！ 🚀

---

如有任何问题，请查看：
- `README.md` - 完整使用文档
- `GODOC.md` - GoDoc 使用指南
- `文档完成总结.md` - 文档创建详情
- 运行 `./docs.sh serve` - 查看在线文档

