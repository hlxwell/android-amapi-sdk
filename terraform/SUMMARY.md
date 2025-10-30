# Terraform 配置总结

## 📋 项目概述

为 Android Management API 创建的完整 Terraform 基础设施即代码(IaC)配置，采用**双区域架构**支持中国(CN)和世界其他地区(ROW)。

## 📁 文件结构

```
terraform/
├── main.tf                    # 主配置文件 - 定义所有 GCP 资源
├── variables.tf               # 变量定义
├── outputs.tf                 # 输出定义
├── terraform.tfvars.example   # 配置示例文件
├── .gitignore                 # Git 忽略规则
├── Makefile                   # 便捷命令工具
├── README.md                  # 完整文档
├── QUICK_START.md            # 快速开始指南
└── SUMMARY.md                # 本文档
```

## 🎯 核心功能

### 1. 双区域 Pub/Sub Topics

| 资源 | 名称 | 用途 |
|------|------|------|
| CN Topic | `amapi-events-cn` | 中国区域设备事件 |
| ROW Topic | `amapi-events-row` | 世界其他地区设备事件 |
| CN Dead Letter | `amapi-events-cn-deadletter` | CN 失败消息 |
| ROW Dead Letter | `amapi-events-row-deadletter` | ROW 失败消息 |

### 2. 自动配置的订阅

每个 Topic 都有对应的订阅，配置包括:
- ⏰ 20 秒确认超时
- 🔄 指数退避重试策略
- 💀 最多 5 次重试后进入 Dead Letter Queue
- 📅 7 天消息保留
- 🗑️ 31 天未使用自动过期

### 3. Service Account 和权限

创建专用 Service Account 并自动配置:
- ✅ `roles/androidmanagement.user` - AMAPI 管理权限
- ✅ `roles/pubsub.publisher` - 发布消息权限 (CN & ROW)
- ✅ `roles/pubsub.subscriber` - 订阅消息权限 (CN & ROW)
- ✅ `roles/pubsub.viewer` - 查看 Pub/Sub 资源
- ✅ 自动授权 Android Management API 服务账号

### 4. API 自动启用

- Android Management API
- Pub/Sub API
- IAM API

## 🚀 快速开始

```bash
# 1. 进入目录
cd terraform

# 2. 初始化并部署
make setup
make apply

# 3. 查看结果
make output
```

## 📊 部署资源清单

部署后会创建以下资源:

### Pub/Sub (8 个资源)
- ✅ 2 个主 Topics (CN & ROW)
- ✅ 2 个 Dead Letter Topics
- ✅ 4 个订阅

### IAM (8 个权限绑定)
- ✅ 1 个 Service Account
- ✅ 4 个 Topic IAM 绑定
- ✅ 2 个 Subscription IAM 绑定
- ✅ 1 个项目级权限

### APIs (3 个)
- ✅ Android Management API
- ✅ Pub/Sub API
- ✅ IAM API

**总计**: 约 19 个 GCP 资源

## 🔧 配置变量

| 变量 | 默认值 | 必填 | 说明 |
|------|--------|------|------|
| `project_id` | - | ✅ | GCP 项目 ID |
| `region` | `us-central1` | ❌ | GCP 区域 |
| `topic_name_prefix` | `amapi-events` | ❌ | Topic 名称前缀 |
| `service_account_id` | `amapi-service-account` | ❌ | SA ID |

## 📤 主要输出

```bash
# CN 区域
terraform output amapi_topic_cn_id
terraform output amapi_subscription_cn_name

# ROW 区域
terraform output amapi_topic_row_id
terraform output amapi_subscription_row_name

# Service Account
terraform output service_account_email

# 使用说明
terraform output setup_instructions
```

## 🎮 Makefile 命令

```bash
make help          # 显示所有可用命令
make init          # 初始化 Terraform
make plan          # 查看执行计划
make apply         # 应用配置
make output        # 显示输出
make test-cn       # 测试 CN Topic
make test-row      # 测试 ROW Topic
make test-all      # 测试所有 Topics
make download-key  # 下载 Service Account Key
make check         # 检查配置(格式化+验证)
make clean         # 清理本地文件
```

## 🔐 安全特性

1. ✅ 使用 Service Account 而非用户凭证
2. ✅ 最小权限原则
3. ✅ .gitignore 保护敏感文件
4. ✅ 支持 Workload Identity
5. ✅ 分离的 Dead Letter Queues
6. ✅ 自动重试和失败处理

## 📈 成本估算

基于标准使用:
- Pub/Sub Topics: $0.04/GB
- Pub/Sub 订阅: $0.04/GB
- Service Account: 免费
- API 调用: 按使用量计费

**预估月成本**: < $10 (低流量场景)

## 🔄 更新和维护

```bash
# 更新资源
make plan
make apply

# 查看当前状态
make show

# 格式化代码
make fmt

# 验证配置
make validate
```

## 🧪 测试流程

### 自动测试
```bash
# 测试所有区域
make test-all
```

### 手动测试
```bash
# CN 区域
gcloud pubsub topics publish amapi-events-cn --message="测试"
gcloud pubsub subscriptions pull amapi-events-cn-subscription --auto-ack

# ROW 区域
gcloud pubsub topics publish amapi-events-row --message="test"
gcloud pubsub subscriptions pull amapi-events-row-subscription --auto-ack
```

## 📚 文档索引

- **快速开始**: [QUICK_START.md](QUICK_START.md)
- **完整文档**: [README.md](README.md)
- **主项目**: [../README.md](../README.md)
- **CLI 使用**: [../docs/CLI_USAGE.md](../docs/CLI_USAGE.md)

## ⚠️ 注意事项

1. **项目 ID**: 部署前必须设置正确的 `project_id`
2. **权限要求**: 需要项目编辑者或所有者权限
3. **API 配额**: 注意 Pub/Sub API 的配额限制
4. **成本控制**: 监控 Pub/Sub 使用量避免意外费用
5. **备份**: 定期备份 Terraform 状态文件

## 🎯 下一步

1. ✅ 部署基础设施: `make apply`
2. ✅ 下载 Service Account Key: `make download-key`
3. ✅ 更新应用配置: 编辑 `../config.yaml`
4. ✅ 测试集成: `make test-all`
5. ✅ 配置监控和告警

## 📞 获取帮助

- **Terraform 问题**: [README.md#故障排查](README.md#故障排查)
- **AMAPI 问题**: [../docs/](../docs/)
- **GCP 文档**: https://cloud.google.com/docs

---

**创建时间**: 2025-10-30
**Terraform 版本要求**: >= 1.0
**Google Provider 版本**: ~> 5.0

