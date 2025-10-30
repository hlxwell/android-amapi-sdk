#!/bin/bash
#
# docs.sh - AMAPI SDK 文档工具
#
# 功能：
# - 查看文档（show, all, type, config, client）
# - 验证代码编译（verify）
# - 启动 godoc 服务器（serve）
# - 生成完整文档（generate）
#

# 设置颜色
GREEN='\033[0;32m'
BLUE='\033[0;34m'
YELLOW='\033[1;33m'
RED='\033[0;31m'
NC='\033[0m' # No Color

# 获取脚本所在目录和项目根目录
SCRIPT_DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )"
PROJECT_ROOT="$( cd "$SCRIPT_DIR/.." && pwd )"

# 显示头部
show_header() {
    echo ""
    echo -e "${BLUE}=========================================${NC}"
    echo -e "${BLUE}   AMAPI SDK 文档工具${NC}"
    echo -e "${BLUE}=========================================${NC}"
    echo ""
}

# 验证代码编译
verify_build() {
    echo -e "${BLUE}📦 验证代码编译...${NC}"
    cd "$PROJECT_ROOT/pkgs/amapi"

    if go build ./...; then
        echo -e "${GREEN}✅ 代码编译成功${NC}"
        return 0
    else
        echo -e "${RED}❌ 代码编译失败，请修复错误后重试${NC}"
        return 1
    fi
}

# 检查并安装 godoc
check_godoc() {
    if ! command -v godoc &> /dev/null; then
        echo -e "${YELLOW}⚠️  godoc 未安装，正在安装...${NC}"
        go install golang.org/x/tools/cmd/godoc@latest

        if ! command -v godoc &> /dev/null; then
            echo -e "${RED}❌ godoc 安装失败${NC}"
            echo ""
            echo "请手动安装："
            echo "  go install golang.org/x/tools/cmd/godoc@latest"
            return 1
        fi
    fi
    return 0
}

# 启动 godoc 服务器
start_server() {
    local PORT="${1:-6060}"

    if ! check_godoc; then
        exit 1
    fi

    echo -e "${GREEN}🌐 启动 godoc 服务器在 http://localhost:$PORT${NC}"
    echo ""
    echo "访问以下 URL 查看文档："
    echo -e "  ${BLUE}主包:${NC} http://localhost:$PORT/pkg/amapi-pkg/pkgs/amapi/"
    echo -e "  ${BLUE}配置:${NC} http://localhost:$PORT/pkg/amapi-pkg/pkgs/amapi/config/"
    echo -e "  ${BLUE}客户端:${NC} http://localhost:$PORT/pkg/amapi-pkg/pkgs/amapi/client/"
    echo -e "  ${BLUE}类型:${NC} http://localhost:$PORT/pkg/amapi-pkg/pkgs/amapi/types/"
    echo -e "  ${BLUE}预设:${NC} http://localhost:$PORT/pkg/amapi-pkg/pkgs/amapi/presets/"
    echo -e "  ${BLUE}工具:${NC} http://localhost:$PORT/pkg/amapi-pkg/pkgs/amapi/utils/"
    echo ""
    echo -e "${YELLOW}按 Ctrl+C 停止服务器${NC}"
    echo ""

    cd "$PROJECT_ROOT"
    godoc -http=:$PORT
}

# 显示帮助信息
show_help() {
    echo "用法: ./scripts/docs.sh [命令] [参数]"
    echo ""
    echo "命令:"
    echo -e "  ${GREEN}show${NC}, (默认)    - 显示包概览和导出的类型"
    echo -e "  ${GREEN}all${NC}             - 显示完整文档（包括所有详细信息）"
    echo -e "  ${GREEN}type${NC} <名称>     - 显示特定类型的文档"
    echo -e "  ${GREEN}config${NC}          - 显示 config 包文档"
    echo -e "  ${GREEN}client${NC}          - 显示 client 包文档"
    echo -e "  ${GREEN}types${NC}           - 显示 types 包文档"
    echo -e "  ${GREEN}verify${NC}          - 验证代码编译"
    echo -e "  ${GREEN}serve${NC} [端口]    - 启动 godoc Web 服务器（默认端口 6060）"
    echo -e "  ${GREEN}generate${NC}        - 验证代码、显示文档预览并启动服务器"
    echo -e "  ${GREEN}help${NC}            - 显示此帮助信息"
    echo ""
    echo "示例:"
    echo "  ./scripts/docs.sh                    # 显示包概览"
    echo "  ./scripts/docs.sh all                # 显示完整文档"
    echo "  ./scripts/docs.sh type Client        # 显示 Client 类型文档"
    echo "  ./scripts/docs.sh config             # 显示 config 包文档"
    echo "  ./scripts/docs.sh verify             # 验证代码编译"
    echo "  ./scripts/docs.sh serve              # 启动 godoc 服务器在 6060 端口"
    echo "  ./scripts/docs.sh serve 8080         # 启动 godoc 服务器在 8080 端口"
    echo "  ./scripts/docs.sh generate           # 完整的文档生成流程"
    echo ""
}

# 主程序
main() {
    local COMMAND="${1:-show}"

    # 进入 AMAPI 包目录（除了某些特殊命令）
    if [[ "$COMMAND" != "serve" && "$COMMAND" != "generate" && "$COMMAND" != "verify" ]]; then
        cd "$PROJECT_ROOT/pkgs/amapi"
    fi

    case "$COMMAND" in
        show|"")
            show_header
            echo -e "${GREEN}📚 主包文档:${NC}"
            echo ""
            cd "$PROJECT_ROOT/pkgs/amapi"
            go doc
            echo ""
            echo "---"
            echo ""
            echo -e "${GREEN}📋 导出的类型和函数:${NC}"
            echo ""
            go doc -short
            echo ""
            ;;

        all)
            show_header
            echo -e "${GREEN}📚 完整文档:${NC}"
            echo ""
            cd "$PROJECT_ROOT/pkgs/amapi"
            go doc -all
            echo ""
            ;;

        type)
            if [ "$2" == "" ]; then
                echo -e "${RED}错误: 请指定类型名称${NC}"
                echo "用法: ./scripts/docs.sh type <名称>"
                exit 1
            fi
            show_header
            echo -e "${GREEN}📖 类型文档: $2${NC}"
            echo ""
            cd "$PROJECT_ROOT/pkgs/amapi"
            go doc "$2"
            echo ""
            ;;

        config)
            show_header
            echo -e "${GREEN}📋 Config 包文档:${NC}"
            echo ""
            cd "$PROJECT_ROOT/pkgs/amapi"
            go doc config
            echo ""
            ;;

        client)
            show_header
            echo -e "${GREEN}📋 Client 包文档:${NC}"
            echo ""
            cd "$PROJECT_ROOT/pkgs/amapi"
            go doc client
            echo ""
            ;;

        types)
            show_header
            echo -e "${GREEN}📋 Types 包文档:${NC}"
            echo ""
            cd "$PROJECT_ROOT/pkgs/amapi"
            go doc types
            echo ""
            ;;

        verify)
            show_header
            if verify_build; then
                echo ""
                echo -e "${GREEN}✅ 验证完成${NC}"
            else
                exit 1
            fi
            ;;

        serve)
            show_header
            start_server "${2:-6060}"
            ;;

        generate)
            show_header
            echo -e "${BLUE}🔨 完整文档生成流程${NC}"
            echo ""

            # 步骤 1: 验证编译
            if ! verify_build; then
                exit 1
            fi
            echo ""

            # 步骤 2: 显示文档预览
            echo -e "${BLUE}=========================================${NC}"
            echo -e "${BLUE}   包级别文档预览${NC}"
            echo -e "${BLUE}=========================================${NC}"
            echo ""
            cd "$PROJECT_ROOT/pkgs/amapi"
            go doc -all
            echo ""

            # 步骤 3: 启动服务器
            echo -e "${BLUE}=========================================${NC}"
            echo -e "${BLUE}   启动 godoc 服务器${NC}"
            echo -e "${BLUE}=========================================${NC}"
            echo ""
            start_server 6060
            ;;

        help|-h|--help)
            show_header
            show_help
            ;;

        *)
            show_header
            echo -e "${RED}错误: 未知命令 '$COMMAND'${NC}"
            echo ""
            show_help
            exit 1
            ;;
    esac
}

# 运行主程序
main "$@"
