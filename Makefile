.PHONY: build build-all clean install test help

# 变量定义
PROJECT_NAME=amapi-cli
BUILD_DIR=build
VERSION=$(shell git describe --tags --always --dirty 2>/dev/null || echo "dev")
COMMIT=$(shell git rev-parse --short HEAD 2>/dev/null || echo "unknown")
BUILD_DATE=$(shell date -u '+%Y-%m-%dT%H:%M:%SZ')

# 构建标志
LDFLAGS=-s -w -X main.version=${VERSION} -X main.commit=${COMMIT} -X main.buildDate=${BUILD_DATE}

# 默认目标
.DEFAULT_GOAL := help

## build: 构建 CLI 工具（当前平台）
build:
	@echo "🔨 开始构建 ${PROJECT_NAME}..."
	@mkdir -p ${BUILD_DIR}
	@go build -ldflags="${LDFLAGS}" -o ${BUILD_DIR}/${PROJECT_NAME} ./cmd/amapi-cli
	@echo "✅ 构建完成: ${BUILD_DIR}/${PROJECT_NAME}"
	@ls -lh ${BUILD_DIR}/${PROJECT_NAME}

## build-all: 构建所有平台的二进制文件
build-all:
	@echo "🔨 开始跨平台构建..."
	@mkdir -p ${BUILD_DIR}
	@echo "Building for linux/amd64..."
	@GOOS=linux GOARCH=amd64 go build -ldflags="${LDFLAGS}" -o ${BUILD_DIR}/${PROJECT_NAME}-linux-amd64 ./cmd/amapi-cli
	@echo "Building for linux/arm64..."
	@GOOS=linux GOARCH=arm64 go build -ldflags="${LDFLAGS}" -o ${BUILD_DIR}/${PROJECT_NAME}-linux-arm64 ./cmd/amapi-cli
	@echo "Building for darwin/amd64..."
	@GOOS=darwin GOARCH=amd64 go build -ldflags="${LDFLAGS}" -o ${BUILD_DIR}/${PROJECT_NAME}-darwin-amd64 ./cmd/amapi-cli
	@echo "Building for darwin/arm64..."
	@GOOS=darwin GOARCH=arm64 go build -ldflags="${LDFLAGS}" -o ${BUILD_DIR}/${PROJECT_NAME}-darwin-arm64 ./cmd/amapi-cli
	@echo "Building for windows/amd64..."
	@GOOS=windows GOARCH=amd64 go build -ldflags="${LDFLAGS}" -o ${BUILD_DIR}/${PROJECT_NAME}-windows-amd64.exe ./cmd/amapi-cli
	@echo "✅ 跨平台构建完成"
	@ls -lh ${BUILD_DIR}/

## clean: 清理构建输出
clean:
	@echo "🧹 清理构建文件..."
	@rm -rf ${BUILD_DIR}
	@rm -f amapi-cli amapi-demo amapi-simple
	@echo "✅ 清理完成"

## install: 安装到系统路径
install: build
	@echo "📦 安装到系统..."
	@sudo cp ${BUILD_DIR}/${PROJECT_NAME} /usr/local/bin/
	@echo "✅ 已安装到 /usr/local/bin/${PROJECT_NAME}"

## test: 运行测试
test:
	@echo "🧪 运行测试..."
	@go test -v ./...

## test-coverage: 运行测试并显示覆盖率
test-coverage:
	@echo "🧪 运行测试（含覆盖率）..."
	@go test -cover ./...

## fmt: 格式化代码
fmt:
	@echo "💅 格式化代码..."
	@go fmt ./...
	@echo "✅ 格式化完成"

## lint: 代码检查
lint:
	@echo "🔍 代码检查..."
	@golangci-lint run || echo "提示: 安装 golangci-lint: brew install golangci-lint"

## deps: 下载依赖
deps:
	@echo "📦 下载依赖..."
	@go mod download
	@go mod tidy
	@echo "✅ 依赖已更新"

## run: 运行 CLI 工具（显示帮助）
run: build
	@${BUILD_DIR}/${PROJECT_NAME} --help

## dev: 开发模式构建（不优化）
dev:
	@echo "🔧 开发模式构建..."
	@mkdir -p ${BUILD_DIR}
	@go build -gcflags="all=-N -l" -o ${BUILD_DIR}/${PROJECT_NAME} ./cmd/amapi-cli
	@echo "✅ 开发版本构建完成"

## release: 创建发布包
release: build-all
	@echo "📦 创建发布包..."
	@mkdir -p ${BUILD_DIR}/release
	@cd ${BUILD_DIR} && tar -czf release/${PROJECT_NAME}-${VERSION}-linux-amd64.tar.gz ${PROJECT_NAME}-linux-amd64
	@cd ${BUILD_DIR} && tar -czf release/${PROJECT_NAME}-${VERSION}-linux-arm64.tar.gz ${PROJECT_NAME}-linux-arm64
	@cd ${BUILD_DIR} && tar -czf release/${PROJECT_NAME}-${VERSION}-darwin-amd64.tar.gz ${PROJECT_NAME}-darwin-amd64
	@cd ${BUILD_DIR} && tar -czf release/${PROJECT_NAME}-${VERSION}-darwin-arm64.tar.gz ${PROJECT_NAME}-darwin-arm64
	@cd ${BUILD_DIR} && zip -q release/${PROJECT_NAME}-${VERSION}-windows-amd64.zip ${PROJECT_NAME}-windows-amd64.exe
	@echo "✅ 发布包已创建:"
	@ls -lh ${BUILD_DIR}/release/

## version: 显示版本信息
version:
	@echo "项目: ${PROJECT_NAME}"
	@echo "版本: ${VERSION}"
	@echo "提交: ${COMMIT}"
	@echo "日期: ${BUILD_DATE}"

## help: 显示帮助信息
help:
	@echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
	@echo "📦 ${PROJECT_NAME} Makefile 帮助"
	@echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
	@echo ""
	@echo "可用命令："
	@echo ""
	@grep -E '^## ' Makefile | sed 's/## /  make /' | column -t -s ':'
	@echo ""
	@echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
	@echo "示例："
	@echo "  make build        # 构建当前平台"
	@echo "  make build-all    # 构建所有平台"
	@echo "  make clean        # 清理构建文件"
	@echo "  make test         # 运行测试"
	@echo "  make install      # 安装到系统"
	@echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"

