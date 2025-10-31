# AMAPI 命令行工具构建指南

本文档介绍如何构建和部署 amapi-cli 命令行工具。

## 目录

- [前提条件](#前提条件)
- [项目结构](#项目结构)
- [构建说明](#构建说明)
- [部署选项](#部署选项)
- [故障排除](#故障排除)
- [开发指南](#开发指南)

## 前提条件

### 系统要求

- Go 1.24 或更高版本
- Git（用于版本控制）
- 网络连接（用于下载依赖）

### Google Cloud 配置

1. **创建 Google Cloud 项目**
   ```bash
   gcloud projects create your-project-id
   gcloud config set project your-project-id
   ```

2. **启用 Android Management API**
   ```bash
   gcloud services enable androidmanagement.googleapis.com
   ```

3. **创建服务账号**
   ```bash
   gcloud iam service-accounts create amapi-service-account \
     --display-name="Android Management API Service Account"
   ```

4. **分配权限**
   ```bash
   gcloud projects add-iam-policy-binding your-project-id \
     --member="serviceAccount:amapi-service-account@your-project-id.iam.gserviceaccount.com" \
     --role="roles/androidmanagement.user"
   ```

5. **生成密钥文件**
   ```bash
   gcloud iam service-accounts keys create sa-key.json \
     --iam-account=amapi-service-account@your-project-id.iam.gserviceaccount.com
   ```

   **⚠️ 重要**：`sa-key.json` 包含敏感信息，已在 `.gitignore` 中，不要提交到版本控制！

## 项目结构

```
amapi-pkg/
├── cmd/amapi-cli/          # 命令行工具源码
│   ├── main.go            # 主程序入口
│   ├── cmd/               # 子命令实现
│   │   ├── enterprise.go  # 企业管理命令
│   │   ├── policy.go      # 策略管理命令
│   │   ├── device.go      # 设备管理命令
│   │   ├── enrollment.go  # 注册令牌命令
│   │   ├── config.go      # 配置管理命令
│   │   ├── health.go      # 健康检查命令
│   │   └── version.go     # 版本信息命令
│   └── internal/          # 内部工具函数
│       └── utils.go       # 通用工具函数
├── pkgs/amapi/            # AMAPI SDK
│   ├── client/            # API 客户端
│   ├── config/            # 配置管理
│   ├── types/             # 类型定义
│   ├── utils/             # 工具函数
│   └── presets/           # 策略预设
├── config.yaml.example    # YAML 配置模板
├── .env.example           # 环境变量模板
├── sa-key.json.example    # 服务账号密钥模板
├── go.mod                 # Go 模块定义
├── go.sum                 # 依赖锁定文件
├── README.md              # 项目说明
├── CLI_USAGE.md           # 命令行工具使用手册
└── BUILD_GUIDE.md         # 构建指南（本文件）
```

## 构建说明

### 1. 获取源码

```bash
# 克隆项目（如果从 Git 仓库）
git clone https://github.com/hlxwell/android-api-demo.git
cd android-api-demo/pkgs/amapi

# 或者如果是本地项目
cd /path/to/amapi-pkg
```

### 2. 安装依赖

```bash
# 下载和整理依赖
go mod tidy

# 验证依赖
go mod verify
```

### 3. 构建选项

#### 使用 Makefile（推荐）

```bash
# 查看所有可用命令
make help

# 构建当前平台
make build

# 开发模式构建（包含调试信息）
make dev

# 跨平台构建
make build-all

# 清理构建文件
make clean

# 运行测试
make test
```

构建后的文件位于 `build/` 目录。

#### 手动构建

```bash
# 基本构建
go build -o build/amapi-cli ./cmd/amapi-cli

# 带调试信息的构建
go build -gcflags="all=-N -l" -o build/amapi-cli ./cmd/amapi-cli
```

#### 生产构建

使用 Makefile 会自动包含版本信息和优化：

```bash
# 生产构建（已优化）
make build

# 查看版本信息
make version
```

手动构建：

```bash
# 优化构建（减小二进制大小）
go build -ldflags="-s -w" -o build/amapi-cli ./cmd/amapi-cli

# 带版本信息的构建
VERSION=$(git describe --tags --always --dirty)
COMMIT=$(git rev-parse --short HEAD)
BUILD_DATE=$(date -u '+%Y-%m-%dT%H:%M:%SZ')

go build \
  -ldflags="-s -w -X main.version=${VERSION} -X main.commit=${COMMIT} -X main.buildDate=${BUILD_DATE}" \
  -o build/amapi-cli ./cmd/amapi-cli
```

#### 跨平台构建

使用 Makefile（推荐）：

```bash
# 构建所有平台
make build-all

# 所有平台的二进制文件会在 build/ 目录中
ls -lh build/
```

手动跨平台构建：

```bash
# Linux (amd64)
GOOS=linux GOARCH=amd64 go build -o build/amapi-cli-linux-amd64 ./cmd/amapi-cli

# Linux (arm64)
GOOS=linux GOARCH=arm64 go build -o build/amapi-cli-linux-arm64 ./cmd/amapi-cli

# macOS (amd64)
GOOS=darwin GOARCH=amd64 go build -o build/amapi-cli-darwin-amd64 ./cmd/amapi-cli

# macOS (arm64/M1)
GOOS=darwin GOARCH=arm64 go build -o build/amapi-cli-darwin-arm64 ./cmd/amapi-cli

# Windows (amd64)
GOOS=windows GOARCH=amd64 go build -o build/amapi-cli-windows-amd64.exe ./cmd/amapi-cli
```

### 4. 构建脚本

创建 `build.sh` 脚本来自动化构建过程：

```bash
#!/bin/bash

set -e

# 配置
PROJECT_NAME="amapi-cli"
VERSION=$(git describe --tags --always --dirty 2>/dev/null || echo "dev")
COMMIT=$(git rev-parse --short HEAD 2>/dev/null || echo "unknown")
BUILD_DATE=$(date -u '+%Y-%m-%dT%H:%M:%SZ')

# 构建标志
LDFLAGS="-s -w -X main.Version=${VERSION} -X main.GitCommit=${COMMIT} -X main.BuildDate=${BUILD_DATE}"

# 清理旧构建
rm -f ${PROJECT_NAME}*

echo "Building ${PROJECT_NAME} ${VERSION}..."

# 构建目标平台
platforms=(
    "linux/amd64"
    "linux/arm64"
    "darwin/amd64"
    "darwin/arm64"
    "windows/amd64"
)

for platform in "${platforms[@]}"; do
    IFS='/' read -r -a array <<< "$platform"
    GOOS="${array[0]}"
    GOARCH="${array[1]}"

    output_name="${PROJECT_NAME}-${GOOS}-${GOARCH}"
    if [ "$GOOS" = "windows" ]; then
        output_name+=".exe"
    fi

    echo "Building for ${GOOS}/${GOARCH}..."
    GOOS=$GOOS GOARCH=$GOARCH go build -ldflags="$LDFLAGS" -o "$output_name" ./cmd/amapi-cli

    if [ $? -ne 0 ]; then
        echo "Failed to build for ${GOOS}/${GOARCH}"
        exit 1
    fi
done

echo "Build completed successfully!"
ls -la ${PROJECT_NAME}*
```

使用构建脚本：

```bash
chmod +x build.sh
./build.sh
```

## 部署选项

### 1. 本地安装

```bash
# 构建并安装到 GOPATH/bin
go install ./cmd/amapi-cli

# 或者复制到系统路径
sudo cp amapi-cli /usr/local/bin/
```

### 2. Docker 部署

创建 `Dockerfile`：

```dockerfile
# 多阶段构建
FROM golang:1.24-alpine AS builder

WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o amapi-cli ./cmd/amapi-cli

# 最终镜像
FROM alpine:latest
RUN apk --no-cache add ca-certificates
WORKDIR /root/

COPY --from=builder /app/amapi-cli .

ENTRYPOINT ["./amapi-cli"]
```

构建和运行 Docker 镜像：

```bash
# 构建镜像
docker build -t amapi-cli .

# 运行容器
docker run --rm -v $(pwd)/config:/config amapi-cli --help

# 使用配置文件运行
docker run --rm \
  -v $(pwd)/service-account-key.json:/creds.json \
  -e GOOGLE_APPLICATION_CREDENTIALS=/creds.json \
  -e GOOGLE_CLOUD_PROJECT=your-project-id \
  amapi-cli health check
```

### 3. 打包发布

#### 创建发布包

```bash
#!/bin/bash

VERSION="v1.0.0"
RELEASE_DIR="release/${VERSION}"

# 创建发布目录
mkdir -p "$RELEASE_DIR"

# 复制二进制文件
cp amapi-cli-* "$RELEASE_DIR/"

# 复制文档
cp README.md CLI_USAGE.md BUILD_GUIDE.md "$RELEASE_DIR/"

# 创建压缩包
cd release
tar -czf "amapi-cli-${VERSION}.tar.gz" "${VERSION}/"
zip -r "amapi-cli-${VERSION}.zip" "${VERSION}/"

echo "Release packages created:"
ls -la amapi-cli-${VERSION}.*
```

#### 生成校验和

```bash
# 生成 SHA256 校验和
cd release
sha256sum amapi-cli-${VERSION}.* > SHA256SUMS

# 验证校验和
sha256sum -c SHA256SUMS
```

## 故障排除

### 1. 构建问题

#### 依赖下载失败

```bash
# 清理模块缓存
go clean -modcache

# 重新下载依赖
go mod download

# 使用代理
GOPROXY=https://proxy.golang.org,direct go mod download
```

#### 版本冲突

```bash
# 更新到最新版本
go get -u ./...

# 整理依赖
go mod tidy
```

### 2. 运行时问题

#### 权限错误

```bash
# 检查文件权限
ls -la amapi-cli
chmod +x amapi-cli

# 检查 SELinux（如果适用）
sestatus
setsebool -P httpd_can_network_connect 1
```

#### 依赖库问题

```bash
# 检查动态库依赖
ldd amapi-cli

# 静态编译（避免依赖问题）
CGO_ENABLED=0 go build -a -ldflags '-extldflags "-static"' ./cmd/amapi-cli
```

### 3. 配置问题

#### 认证失败

```bash
# 验证服务账号密钥
gcloud auth activate-service-account --key-file=service-account-key.json

# 测试权限
gcloud projects get-iam-policy your-project-id

# 检查 API 启用状态
gcloud services list --enabled | grep androidmanagement
```

## 开发指南

### 1. 开发环境设置

```bash
# 安装开发工具
go install golang.org/x/tools/cmd/goimports@latest
go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest

# 设置 Git hooks
cat > .git/hooks/pre-commit << 'EOF'
#!/bin/bash
goimports -w .
golangci-lint run
go test ./...
EOF

chmod +x .git/hooks/pre-commit
```

### 2. 代码质量检查

```bash
# 格式化代码
gofmt -w .
goimports -w .

# 静态分析
golangci-lint run

# 运行测试
go test ./...
go test -race ./...
go test -cover ./...
```

### 3. 添加新命令

要添加新的命令，请遵循以下步骤：

1. 在 `cmd/amapi-cli/cmd/` 目录下创建新文件
2. 实现命令逻辑
3. 在 `main.go` 中注册命令
4. 更新文档

示例：

```go
// cmd/amapi-cli/cmd/newcommand.go
package cmd

import "github.com/spf13/cobra"

func NewMyCommand() *cobra.Command {
    return &cobra.Command{
        Use:   "mycommand",
        Short: "My new command",
        RunE: func(cmd *cobra.Command, args []string) error {
            // 实现命令逻辑
            return nil
        },
    }
}
```

```go
// cmd/amapi-cli/main.go
func main() {
    rootCmd.AddCommand(
        // ... 其他命令
        cmd.NewMyCommand(),
    )
}
```

### 4. 测试

```bash
# 单元测试
go test ./cmd/amapi-cli/...

# 集成测试
go test -tags=integration ./...

# 基准测试
go test -bench=. ./...
```

### 5. 文档更新

更新代码后，请确保同步更新：

- `README.md` - 项目概述
- `CLI_USAGE.md` - 使用手册
- `BUILD_GUIDE.md` - 构建指南
- 内联文档和注释

这个构建指南涵盖了从开发到部署的完整流程。如果遇到特定问题，请查看相关的故障排除部分或查阅 Go 官方文档。
