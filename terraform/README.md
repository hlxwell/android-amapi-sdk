# Android Management API - Terraform 配置

这个 Terraform 配置用于自动化部署 Android Management API 所需的 GCP 资源。

## 功能特性

本 Terraform 配置会自动创建和配置以下资源:

### 1. 启用 API
- ✅ Android Management API (`androidmanagement.googleapis.com`)
- ✅ Pub/Sub API (`pubsub.googleapis.com`)
- ✅ IAM API (`iam.googleapis.com`)

### 2. Pub/Sub 资源
- 📨 **主 Topic**: `amapi-events` - 接收来自 Android Management API 的事件
- 💀 **Dead Letter Topic**: `amapi-events-deadletter` - 处理失败的消息
- 📬 **订阅**: 自动创建订阅,配置重试策略和 Dead Letter 队列
- ⏰ 消息保留时间: 7天
- 🔄 自动重试配置: 最多 5 次,指数退避策略

### 3. Service Account 和权限
- 👤 创建专用 Service Account
- 🔐 自动配置所需的 IAM 权限:
  - `roles/androidmanagement.user` - 管理 Android 设备
  - `roles/pubsub.publisher` - 发布消息到 Topic
  - `roles/pubsub.subscriber` - 订阅和消费消息
  - `roles/pubsub.viewer` - 查看 Pub/Sub 资源
- 🤖 自动授权 Android Management API 服务账号发布权限

## 前置要求

### 1. 安装工具
```bash
# Terraform (>= 1.0)
brew install terraform

# gcloud CLI
brew install google-cloud-sdk
```

### 2. GCP 认证
```bash
# 登录到 GCP
gcloud auth application-default login

# 设置默认项目
gcloud config set project YOUR_PROJECT_ID
```

### 3. 启用必要的 API (可选 - Terraform 会自动启用)
```bash
gcloud services enable cloudresourcemanager.googleapis.com
gcloud services enable serviceusage.googleapis.com
```

## 快速开始

### 1. 配置变量
```bash
cd terraform
cp terraform.tfvars.example terraform.tfvars
```

编辑 `terraform.tfvars` 文件,至少设置你的项目 ID:
```hcl
project_id = "enhancer-471605"  # 替换为你的项目 ID
```

### 2. 初始化 Terraform
```bash
terraform init
```

### 3. 查看计划
```bash
terraform plan
```

### 4. 应用配置
```bash
terraform apply
```

输入 `yes` 确认创建资源。

### 5. 查看输出
```bash
terraform output
```

## 配置选项

### 变量说明

| 变量名 | 类型 | 默认值 | 说明 |
|--------|------|--------|------|
| `project_id` | string | - | **必填** GCP 项目 ID |
| `region` | string | `us-central1` | GCP 区域 |
| `topic_name` | string | `amapi-events` | Pub/Sub Topic 名称 |
| `service_account_id` | string | `amapi-service-account` | Service Account ID |
| `service_account_display_name` | string | `Android Management API Service Account` | Service Account 显示名称 |
| `create_service_account_key` | bool | `false` | 是否创建 Service Account Key |
| `save_key_to_file` | bool | `false` | 是否保存 Key 到文件 |
| `service_account_key_filename` | string | `sa-key.json` | Key 文件名 |

### 示例配置

#### 基础配置 (推荐用于生产环境)
```hcl
project_id = "enhancer-471605"
region     = "us-central1"
topic_name = "amapi-events"
```

#### 开发环境配置 (包含 Service Account Key)
```hcl
project_id                 = "enhancer-471605"
region                     = "us-central1"
topic_name                 = "amapi-events"
create_service_account_key = true
save_key_to_file           = true
service_account_key_filename = "sa-key.json"
```

## 使用输出

部署完成后,你可以通过以下方式获取输出信息:

```bash
# 查看所有输出
terraform output

# 查看特定输出
terraform output service_account_email
terraform output amapi_topic_name
terraform output setup_instructions
```

### 主要输出变量

- `service_account_email` - Service Account 邮箱地址
- `amapi_topic_id` - Pub/Sub Topic 完整 ID
- `amapi_subscription_name` - Pub/Sub Subscription 名称
- `setup_instructions` - 详细的后续步骤说明

## 与项目集成

### 更新 config.yaml

部署完成后,更新项目根目录的 `config.yaml`:

```yaml
# Google Cloud 配置
project_id: "enhancer-471605"  # 使用 terraform output project_id
credentials_file: "./sa-key.json"  # 如果创建了 Key

# Pub/Sub 配置
pubsub_topic: "projects/enhancer-471605/topics/amapi-events"  # 使用 terraform output amapi_topic_id
```

### 手动下载 Service Account Key

如果没有通过 Terraform 创建 Key,可以手动创建:

```bash
# 获取 Service Account 邮箱
SA_EMAIL=$(terraform output -raw service_account_email)

# 创建 Key
gcloud iam service-accounts keys create sa-key.json \
  --iam-account=$SA_EMAIL

# 复制到项目根目录
cp sa-key.json ../sa-key.json
```

## 测试部署

### 1. 测试 Pub/Sub Topic
```bash
# 获取 Topic 名称
TOPIC_NAME=$(terraform output -raw amapi_topic_name)

# 发布测试消息
gcloud pubsub topics publish $TOPIC_NAME --message="Test message from Terraform"
```

### 2. 测试订阅
```bash
# 获取订阅名称
SUB_NAME=$(terraform output -raw amapi_subscription_name)

# 拉取消息
gcloud pubsub subscriptions pull $SUB_NAME --auto-ack --limit=10
```

### 3. 测试 Service Account 权限
```bash
# 使用 Service Account 认证
export GOOGLE_APPLICATION_CREDENTIALS="./sa-key.json"

# 运行 AMAPI CLI 命令
cd ..
./build/amapi-cli health

# 或者使用 go run
go run cmd/amapi-cli/main.go health
```

## 更新和维护

### 更新资源
```bash
# 修改 terraform.tfvars 或 *.tf 文件后
terraform plan
terraform apply
```

### 查看当前状态
```bash
terraform show
```

### 格式化代码
```bash
terraform fmt
```

### 验证配置
```bash
terraform validate
```

## 清理资源

⚠️ **警告**: 这将删除所有通过 Terraform 创建的资源!

```bash
terraform destroy
```

### 选择性删除

如果你只想删除某些资源,可以使用:

```bash
# 删除特定资源
terraform destroy -target=google_pubsub_topic.amapi_events_deadletter

# 从状态中移除但不删除资源
terraform state rm google_service_account.amapi_sa
```

## 高级用法

### 使用不同的后端

默认情况下,Terraform 使用本地后端。对于团队协作,建议使用远程后端:

#### GCS Backend 示例

在 `main.tf` 中添加:

```hcl
terraform {
  backend "gcs" {
    bucket = "your-terraform-state-bucket"
    prefix = "terraform/amapi/state"
  }
}
```

### 使用 Workspaces

```bash
# 创建新的 workspace
terraform workspace new development
terraform workspace new production

# 切换 workspace
terraform workspace select development

# 列出 workspaces
terraform workspace list
```

### 导入现有资源

如果你已经手动创建了一些资源,可以导入到 Terraform:

```bash
# 导入 Service Account
terraform import google_service_account.amapi_sa \
  projects/enhancer-471605/serviceAccounts/amapi-demo-sa@enhancer-471605.iam.gserviceaccount.com

# 导入 Topic
terraform import google_pubsub_topic.amapi_events \
  projects/enhancer-471605/topics/amapi-events
```

## 故障排查

### 权限错误

如果遇到权限错误:

```bash
# 确认你有足够的权限
gcloud projects get-iam-policy YOUR_PROJECT_ID

# 授予必要的角色
gcloud projects add-iam-policy-binding YOUR_PROJECT_ID \
  --member="user:YOUR_EMAIL" \
  --role="roles/editor"
```

### API 未启用

```bash
# 手动启用必要的 API
gcloud services enable cloudresourcemanager.googleapis.com
gcloud services enable serviceusage.googleapis.com
gcloud services enable iam.googleapis.com
```

### 状态锁定

如果 Terraform 状态被锁定:

```bash
# 强制解锁 (谨慎使用!)
terraform force-unlock LOCK_ID
```

## 安全最佳实践

### 1. Service Account Key 管理

- ✅ **推荐**: 在生产环境使用 Workload Identity 或 GKE Workload Identity
- ✅ **推荐**: 使用 Secret Manager 存储 keys
- ⚠️ **避免**: 将 keys 提交到版本控制系统
- ⚠️ **避免**: 在不安全的地方存储 keys

### 2. Terraform 状态文件

- 状态文件包含敏感信息,不要提交到版本控制
- 使用远程后端 (如 GCS) 并启用加密
- 定期备份状态文件

### 3. 最小权限原则

- 只授予必要的权限
- 定期审查和更新 IAM 策略
- 使用不同的 Service Accounts 用于不同的环境

## 相关文档

- [Terraform Google Provider](https://registry.terraform.io/providers/hashicorp/google/latest/docs)
- [Android Management API](https://developers.google.com/android/management)
- [Google Cloud Pub/Sub](https://cloud.google.com/pubsub/docs)
- [Service Account 最佳实践](https://cloud.google.com/iam/docs/best-practices-for-securing-service-accounts)

## 支持

如有问题,请参考:
- 项目主 README: `../README.md`
- CLI 使用文档: `../docs/CLI_USAGE.md`
- 快速开始: `../docs/QUICKSTART.md`

## License

同项目主 License

