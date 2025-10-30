# 快速开始指南

## 5 分钟部署 AMAPI 基础设施

### 步骤 1: 准备配置文件

```bash
cd terraform
cp terraform.tfvars.example terraform.tfvars
```

编辑 `terraform.tfvars`:
```hcl
project_id = "enhancer-471605"  # 替换为你的项目 ID
```

### 步骤 2: 初始化并部署

```bash
# 初始化 Terraform
terraform init

# 查看将要创建的资源
terraform plan

# 执行部署
terraform apply
# 输入 'yes' 确认
```

### 步骤 3: 查看部署结果

```bash
# 查看所有输出
terraform output

# 查看设置说明
terraform output setup_instructions
```

## 创建的资源

部署后会自动创建:

### 📨 Pub/Sub Topics
- ✅ `amapi-events-cn` - 中国区域事件
- ✅ `amapi-events-row` - 世界其他地区事件
- ✅ 对应的 Dead Letter Topics
- ✅ 自动配置的订阅

### 👤 Service Account
- ✅ 专用 Service Account
- ✅ 完整的 AMAPI 和 Pub/Sub 权限
- ✅ 自动授权给 Android Management API

### ⚙️ API 启用
- ✅ Android Management API
- ✅ Pub/Sub API
- ✅ IAM API

## 快速测试

### 测试 CN Topic
```bash
# 发布测试消息
gcloud pubsub topics publish amapi-events-cn --message="测试消息"

# 查看消息
gcloud pubsub subscriptions pull amapi-events-cn-subscription --auto-ack
```

### 测试 ROW Topic
```bash
# 发布测试消息
gcloud pubsub topics publish amapi-events-row --message="Test message"

# 查看消息
gcloud pubsub subscriptions pull amapi-events-row-subscription --auto-ack
```

## 获取 Service Account Key

```bash
# 方法 1: 通过 Terraform (重新部署并启用 key 创建)
# 编辑 terraform.tfvars，添加:
create_service_account_key = true
save_key_to_file = true

# 重新应用
terraform apply

# 方法 2: 手动创建
SA_EMAIL=$(terraform output -raw service_account_email)
gcloud iam service-accounts keys create sa-key.json --iam-account=$SA_EMAIL
```

## 集成到应用

### 更新 config.yaml

```yaml
project_id: "enhancer-471605"
credentials_file: "./sa-key.json"

# 根据设备区域使用不同的 Topic
pubsub_topic_cn: "projects/enhancer-471605/topics/amapi-events-cn"
pubsub_topic_row: "projects/enhancer-471605/topics/amapi-events-row"
```

### 在代码中使用

```go
// 根据设备区域选择 Topic
func getTopicForDevice(deviceRegion string) string {
    if deviceRegion == "CN" {
        return "projects/enhancer-471605/topics/amapi-events-cn"
    }
    return "projects/enhancer-471605/topics/amapi-events-row"
}
```

## 常用命令

```bash
# 查看当前状态
terraform show

# 查看特定输出
terraform output service_account_email
terraform output amapi_topic_cn_id
terraform output amapi_topic_row_id

# 格式化配置文件
terraform fmt

# 验证配置
terraform validate

# 更新资源
terraform apply

# 销毁资源 (谨慎使用!)
terraform destroy
```

## 故障排查

### 权限不足
```bash
# 确认当前认证
gcloud auth list

# 重新认证
gcloud auth application-default login
```

### API 未启用
```bash
# 手动启用 API
gcloud services enable cloudresourcemanager.googleapis.com
gcloud services enable serviceusage.googleapis.com
```

### 查看详细日志
```bash
# 启用详细日志
export TF_LOG=DEBUG
terraform plan
```

## 下一步

1. ✅ 部署完成后，查看 [README.md](README.md) 了解更多配置选项
2. ✅ 参考 [主项目文档](../README.md) 了解如何使用 AMAPI CLI
3. ✅ 查看 [安全最佳实践](README.md#安全最佳实践)

## 需要帮助?

- 📖 完整文档: [README.md](README.md)
- 🔧 CLI 使用: [../docs/CLI_USAGE.md](../docs/CLI_USAGE.md)
- 🚀 项目快速开始: [../docs/QUICKSTART.md](../docs/QUICKSTART.md)

