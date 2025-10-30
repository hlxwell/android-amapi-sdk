# 项目脚本

本目录包含项目相关的脚本工具。

## 📜 可用脚本

### generate_docs.sh

生成 AMAPI SDK 的 godoc 文档。

**用途**:
- 验证代码编译
- 显示包文档预览
- 启动 godoc 服务器

**使用方法**:
```bash
# 从项目根目录运行
./scripts/generate_docs.sh

# 或从任何位置运行（使用绝对路径）
/path/to/amapi-pkg/scripts/generate_docs.sh
```

**功能**:
1. ✅ 验证代码编译成功
2. 📝 显示包级别文档预览
3. 🌐 启动 godoc 服务器在 http://localhost:6060

**访问文档**:
- 主包: http://localhost:6060/pkg/amapi-pkg/pkgs/amapi/
- 配置: http://localhost:6060/pkg/amapi-pkg/pkgs/amapi/config/
- 客户端: http://localhost:6060/pkg/amapi-pkg/pkgs/amapi/client/
- 类型: http://localhost:6060/pkg/amapi-pkg/pkgs/amapi/types/
- 预设: http://localhost:6060/pkg/amapi-pkg/pkgs/amapi/presets/

按 `Ctrl+C` 停止服务器。

### docs.sh

简化版文档工具（如果存在）。

## 🔧 脚本开发规范

1. **Shebang**: 所有脚本应以 `#!/bin/bash` 开头
2. **错误处理**: 使用 `set -e` 在错误时退出
3. **注释**: 在脚本开头添加功能说明
4. **可执行权限**: 确保脚本有执行权限 `chmod +x`

## 📝 添加新脚本

添加新脚本时:

1. 将脚本文件放在此目录
2. 添加执行权限: `chmod +x scripts/your_script.sh`
3. 更新此 README，说明脚本功能和用法
4. 确保脚本通过 `make build` 验证

## 🗂️ 项目结构

```
/
├── scripts/          # 所有脚本文件
├── docs/            # 所有文档文件
├── pkgs/            # 代码包
├── cmd/             # 命令行工具
└── terraform/       # Terraform 配置
```

遵循 [AGENTS.md](../AGENTS.md) 中的项目规范。

