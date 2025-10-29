# 快速开始指南

本指南帮助你快速配置和运行 AMAPI CLI 工具。

## 📋 前提条件

- Go 1.19 或更高版本
- Google Cloud Platform 账号
- Android Management API 已启用

## 🚀 5 分钟快速设置

### 第 1 步：克隆项目

```bash
git clone https://github.com/hlxwell/android-api-demo.git
cd android-api-demo/amapi-pkg
```

### 第 2 步：安装依赖

```bash
go mod tidy
```

### 第 3 步：配置 Google Cloud

#### 创建服务账号并下载密钥

```bash
# 设置项目 ID（替换为你的实际项目 ID）
export PROJECT_ID="your-project-id"

# 创建服务账号
gcloud iam service-accounts create amapi-service-account \
  --project=$PROJECT_ID \
  --display-name="Android Management API Service Account"

# 分配权限
gcloud projects add-iam-policy-binding $PROJECT_ID \
  --member="serviceAccount:amapi-service-account@${PROJECT_ID}.iam.gserviceaccount.com" \
  --role="roles/androidmanagement.user"

# 下载密钥文件
gcloud iam service-accounts keys create sa-key.json \
  --project=$PROJECT_ID \
  --iam-account=amapi-service-account@${PROJECT_ID}.iam.gserviceaccount.com

echo "✅ 服务账号密钥已创建: sa-key.json"
```

### 第 4 步：创建配置文件

#### 选项 A：使用 YAML 配置（推荐）

```bash
# 复制示例配置文件
cp config.yaml.example config.yaml

# 编辑配置文件，填写你的项目 ID
cat > config.yaml << EOF
project_id: "${PROJECT_ID}"
credentials_file: "./sa-key.json"
callback_url: "https://your-app.com/callback"
timeout: 30s
retry_attempts: 3
enable_retry: true
log_level: "info"
EOF

echo "✅ 配置文件已创建: config.yaml"
```

#### 选项 B：使用环境变量

```bash
# 设置环境变量
export GOOGLE_CLOUD_PROJECT="${PROJECT_ID}"
export GOOGLE_APPLICATION_CREDENTIALS="./sa-key.json"
export AMAPI_LOG_LEVEL="info"

echo "✅ 环境变量已设置"
```

### 第 5 步：构建 CLI 工具

```bash
# 使用 Makefile 构建（推荐）
make build

# 或手动构建
# go build -o build/amapi-cli ./cmd/amapi-cli

# 验证构建
./build/amapi-cli --help

echo "✅ CLI 工具构建成功"
```

### 第 6 步：测试连接

```bash
# 健康检查
./build/amapi-cli health check

# 如果成功，你会看到：
# ✓ 配置检查通过
# ✓ API连接正常
```

## 🎯 常用命令

### 企业管理

```bash
# 生成企业注册 URL
./build/amapi-cli enterprise signup-url --project-id $PROJECT_ID

# 列出企业
./build/amapi-cli enterprise list $PROJECT_ID

# 获取企业详情
./build/amapi-cli enterprise get enterprises/LC12345678
```

### 策略管理

```bash
# 查看可用的策略预设
./build/amapi-cli policy presets -o table

# 创建策略（从预设）
./build/amapi-cli policy create \
  --enterprise LC12345678 \
  --policy-id my-policy \
  --from-preset fully_managed

# 列出策略
./build/amapi-cli policy list --enterprise LC12345678
```

### 设备管理

```bash
# 列出设备
./build/amapi-cli device list --enterprise LC12345678

# 获取设备详情
./build/amapi-cli device get enterprises/LC12345678/devices/DEVICE_ID

# 锁定设备（10分钟）
./build/amapi-cli device lock \
  --enterprise LC12345678 \
  --device DEVICE_ID \
  --duration PT10M
```

### 注册令牌管理

```bash
# 创建注册令牌
./build/amapi-cli enrollment create \
  --enterprise LC12345678 \
  --policy my-policy \
  --duration 24h

# 生成 QR 码
./build/amapi-cli enrollment qrcode \
  --enterprise LC12345678 \
  --token TOKEN_ID \
  --wifi-ssid "CompanyWiFi" \
  --wifi-password "password123"
```

## ⚠️ 安全提醒

1. **保护密钥文件**
   - `sa-key.json` 包含敏感信息
   - 已在 `.gitignore` 中，不会被提交到 Git
   - 不要在公开场合分享此文件

2. **检查 Git 状态**
   ```bash
   # 确认敏感文件不在 Git 跟踪中
   git status

   # 应该看不到 sa-key.json
   ```

3. **如果意外提交了密钥**
   ```bash
   # 立即撤销密钥
   gcloud iam service-accounts keys delete KEY_ID \
     --iam-account=SERVICE_ACCOUNT_EMAIL

   # 创建新密钥
   gcloud iam service-accounts keys create sa-key.json \
     --iam-account=amapi-service-account@${PROJECT_ID}.iam.gserviceaccount.com
   ```

## 🐛 故障排除

### 问题 1：找不到配置文件

```
错误: 加载配置失败: configuration file not found
```

**解决方法**：
```bash
# 检查配置文件是否存在
ls -la amapi.yaml

# 或使用环境变量
export GOOGLE_CLOUD_PROJECT="your-project-id"
export GOOGLE_APPLICATION_CREDENTIALS="./sa-key.json"
```

### 问题 2：认证失败

```
错误: failed to load credentials
```

**解决方法**：
```bash
# 检查密钥文件
ls -la sa-key.json

# 测试认证
gcloud auth activate-service-account \
  --key-file=sa-key.json
```

### 问题 3：权限不足

```
错误: permission denied
```

**解决方法**：
```bash
# 检查服务账号权限
gcloud projects get-iam-policy $PROJECT_ID \
  --flatten="bindings[].members" \
  --filter="bindings.members:serviceAccount:amapi-service-account@${PROJECT_ID}.iam.gserviceaccount.com"
```

## 📚 更多资源

- [完整文档](../README.md)
- [文档索引](README.md)
- [构建指南](BUILD_GUIDE.md)
- [安全指南](SECURITY.md)
- [使用手册](USAGE_GUIDE.md)
- [CLI 使用手册](CLI_USAGE.md)
- [API 文档](https://developers.google.com/android/management)

## 💡 提示

1. **使用 Tab 补全**（如果安装了）
   ```bash
   ./build/amapi-cli completion bash > /etc/bash_completion.d/amapi-cli
   ```

2. **设置别名**
   ```bash
   alias amapi='./build/amapi-cli'
   echo "alias amapi='$(pwd)/build/amapi-cli'" >> ~/.bashrc
   ```

3. **JSON 输出**
   ```bash
   # 大多数命令支持 JSON 输出
   ./build/amapi-cli enterprise list $PROJECT_ID -o json | jq .
   ```

---

🎉 **恭喜！你已经完成了快速设置。现在可以开始使用 AMAPI CLI 工具了！**

