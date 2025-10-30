#!/bin/bash
#
# docs.sh - 快速查看 AMAPI SDK 文档
#

# 设置颜色
GREEN='\033[0;32m'
BLUE='\033[0;34m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

echo ""
echo -e "${BLUE}=========================================${NC}"
echo -e "${BLUE}   AMAPI SDK 文档${NC}"
echo -e "${BLUE}=========================================${NC}"
echo ""

# 获取脚本所在目录和项目根目录
SCRIPT_DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )"
PROJECT_ROOT="$( cd "$SCRIPT_DIR/.." && pwd )"

# 进入 AMAPI 包目录
cd "$PROJECT_ROOT/pkgs/amapi"

# 选项1: 显示包文档
if [ "$1" == "show" ] || [ "$1" == "" ]; then
    echo -e "${GREEN}📚 主包文档:${NC}"
    echo ""
    go doc
    echo ""
    echo "---"
    echo ""

    echo -e "${GREEN}📋 导出的类型和函数:${NC}"
    echo ""
    go doc -short
    echo ""
fi

# 选项2: 显示特定类型或函数的文档
if [ "$1" == "type" ] && [ "$2" != "" ]; then
    echo -e "${GREEN}📖 类型文档: $2${NC}"
    echo ""
    go doc "$2"
    echo ""
fi

# 选项3: 显示所有文档
if [ "$1" == "all" ]; then
    echo -e "${GREEN}📚 完整文档:${NC}"
    echo ""
    go doc -all
    echo ""
fi

# 选项4: 启动 godoc 服务器
if [ "$1" == "serve" ]; then
    PORT="${2:-6060}"

    # 检查 godoc 是否安装
    if ! command -v godoc &> /dev/null; then
        echo -e "${YELLOW}⚠️  godoc 未安装，正在安装...${NC}"
        go install golang.org/x/tools/cmd/godoc@latest
    fi

    if command -v godoc &> /dev/null; then
        echo -e "${GREEN}🌐 启动 godoc 服务器在 http://localhost:$PORT${NC}"
        echo ""
        echo "访问以下 URL 查看文档："
        echo -e "  ${BLUE}主包:${NC} http://localhost:$PORT/pkg/amapi-pkg/pkgs/amapi/"
        echo -e "  ${BLUE}配置:${NC} http://localhost:$PORT/pkg/amapi-pkg/pkgs/amapi/config/"
        echo -e "  ${BLUE}客户端:${NC} http://localhost:$PORT/pkg/amapi-pkg/pkgs/amapi/client/"
        echo -e "  ${BLUE}类型:${NC} http://localhost:$PORT/pkg/amapi-pkg/pkgs/amapi/types/"
        echo -e "  ${BLUE}预设:${NC} http://localhost:$PORT/pkg/amapi-pkg/pkgs/amapi/presets/"
        echo ""
        echo -e "${YELLOW}按 Ctrl+C 停止服务器${NC}"
        echo ""

        # 切换到项目根目录
        cd "$PROJECT_ROOT"
        godoc -http=:$PORT
    else
        echo -e "${YELLOW}❌ godoc 安装失败${NC}"
        echo ""
        echo "请手动安装："
        echo "  go install golang.org/x/tools/cmd/godoc@latest"
        exit 1
    fi
fi

# 选项5: 显示配置包文档
if [ "$1" == "config" ]; then
    echo -e "${GREEN}📋 Config 包文档:${NC}"
    echo ""
    go doc config
    echo ""
fi

# 选项6: 显示客户端包文档
if [ "$1" == "client" ]; then
    echo -e "${GREEN}📋 Client 包文档:${NC}"
    echo ""
    go doc client
    echo ""
fi

# 帮助信息
if [ "$1" == "help" ] || [ "$1" == "-h" ] || [ "$1" == "--help" ]; then
    echo "用法: ./docs.sh [命令] [参数]"
    echo ""
    echo "命令:"
    echo "  show, (默认)  - 显示包概览和导出的类型"
    echo "  all           - 显示完整文档（包括所有详细信息）"
    echo "  type <名称>   - 显示特定类型的文档"
    echo "  config        - 显示 config 包文档"
    echo "  client        - 显示 client 包文档"
    echo "  serve [端口]  - 启动 godoc Web 服务器（默认端口 6060）"
    echo "  help          - 显示此帮助信息"
    echo ""
    echo "示例:"
    echo "  ./docs.sh                    # 显示包概览"
    echo "  ./docs.sh all                # 显示完整文档"
    echo "  ./docs.sh type Client        # 显示 Client 类型文档"
    echo "  ./docs.sh type NewClient     # 显示 NewClient 函数文档"
    echo "  ./docs.sh config             # 显示 config 包文档"
    echo "  ./docs.sh serve              # 启动 godoc 服务器在 6060 端口"
    echo "  ./docs.sh serve 8080         # 启动 godoc 服务器在 8080 端口"
    echo ""
fi

