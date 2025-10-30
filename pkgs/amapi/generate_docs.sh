#!/bin/bash
#
# generate_docs.sh - 生成 AMAPI SDK 的 godoc 文档
#
# 此脚本会：
# 1. 验证代码编译
# 2. 启动 godoc 服务器
# 3. 生成静态 HTML 文档（可选）
#

set -e

echo "========================================="
echo "  AMAPI SDK 文档生成工具"
echo "========================================="
echo ""

# 获取脚本所在目录
SCRIPT_DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )"
cd "$SCRIPT_DIR"

# 验证代码
echo "📦 验证代码编译..."
if go build ./...; then
    echo "✅ 代码编译成功"
else
    echo "❌ 代码编译失败，请修复错误后重试"
    exit 1
fi

echo ""
echo "📝 检查 go doc 是否可用..."

# 检查是否安装了 go doc
if ! command -v go &> /dev/null; then
    echo "❌ Go 未安装或不在 PATH 中"
    exit 1
fi

echo "✅ Go 工具链已就绪"
echo ""

# 显示包文档
echo "========================================="
echo "  包级别文档预览"
echo "========================================="
echo ""

go doc -all

echo ""
echo "========================================="
echo "  启动 godoc 服务器"
echo "========================================="
echo ""

# 检查是否安装了 godoc
if ! command -v godoc &> /dev/null; then
    echo "⚠️  godoc 未安装，正在安装..."
    go install golang.org/x/tools/cmd/godoc@latest

    if ! command -v godoc &> /dev/null; then
        echo "❌ godoc 安装失败"
        echo ""
        echo "请手动安装："
        echo "  go install golang.org/x/tools/cmd/godoc@latest"
        echo ""
        echo "或使用 go doc 查看文档："
        echo "  go doc -all amapi-pkg/pkgs/amapi"
        exit 1
    fi
fi

echo "✅ godoc 已就绪"
echo ""
echo "🌐 启动 godoc 服务器在 http://localhost:6060"
echo ""
echo "访问以下 URL 查看文档："
echo "  主包: http://localhost:6060/pkg/amapi-pkg/pkgs/amapi/"
echo "  配置: http://localhost:6060/pkg/amapi-pkg/pkgs/amapi/config/"
echo "  客户端: http://localhost:6060/pkg/amapi-pkg/pkgs/amapi/client/"
echo "  类型: http://localhost:6060/pkg/amapi-pkg/pkgs/amapi/types/"
echo "  预设: http://localhost:6060/pkg/amapi-pkg/pkgs/amapi/presets/"
echo ""
echo "按 Ctrl+C 停止服务器"
echo ""

# 启动 godoc 服务器
cd ../../..
godoc -http=:6060

