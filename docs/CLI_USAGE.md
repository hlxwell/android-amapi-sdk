# AMAPI 命令行工具使用手册

本文档介绍如何使用 amapi-cli 命令行工具来管理 Android Management API。

## 目录

- [安装和配置](#安装和配置)
- [基本用法](#基本用法)
- [企业管理](#企业管理)
- [策略管理](#策略管理)
- [设备管理](#设备管理)
- [注册令牌管理](#注册令牌管理)
- [配置管理](#配置管理)
- [健康检查](#健康检查)
- [故障排除](#故障排除)

## 安装和配置

### 1. 构建工具

```bash
# 在项目根目录下
go build -o amapi-cli ./cmd/amapi-cli
```

### 2. 配置认证

#### 方法一：使用配置文件

创建 `amapi.yaml` 配置文件：

```yaml
project_id: "your-project-id"
credentials_file: "/path/to/service-account.json"
timeout: 30s
retry_attempts: 3
log_level: "info"
```

#### 方法二：使用环境变量

```bash
export GOOGLE_CLOUD_PROJECT="your-project-id"
export GOOGLE_APPLICATION_CREDENTIALS="/path/to/service-account.json"
```

### 3. 验证配置

```bash
./amapi-cli config show
./amapi-cli health check
```

## 基本用法

### 查看帮助

```bash
# 查看主帮助
./amapi-cli --help

# 查看子命令帮助
./amapi-cli enterprise --help
./amapi-cli policy --help
```

### 全局选项

```bash
# 启用调试模式
./amapi-cli --debug <command>

# 指定输出格式
./amapi-cli <command> --output json    # JSON 格式（默认）
./amapi-cli <command> --output yaml    # YAML 格式
./amapi-cli <command> --output table   # 表格格式
```

## 企业管理

### 创建企业

```bash
# 基本创建
./amapi-cli enterprise create --display-name "我的公司" --project-id my-project

# 指定回调URL
./amapi-cli enterprise create \
  --display-name "测试企业" \
  --project-id test-project \
  --callback https://example.com/callback
```

### 查看企业

```bash
# 获取特定企业
./amapi-cli enterprise get enterprises/LC12345678

# 列出所有企业
./amapi-cli enterprise list my-project

# 以表格格式显示
./amapi-cli enterprise list my-project --output table
```

### 更新企业

```bash
# 更新显示名称
./amapi-cli enterprise update enterprises/LC12345678 --display-name "新公司名"

# 更新主颜色
./amapi-cli enterprise update enterprises/LC12345678 --primary-color 0xFF0000
```

### 删除企业

```bash
# 删除企业（需要确认）
./amapi-cli enterprise delete enterprises/LC12345678

# 强制删除（跳过确认）
./amapi-cli enterprise delete enterprises/LC12345678 --force
```

### 生成注册URL

```bash
# 生成企业注册URL
./amapi-cli enterprise signup-url --project-id my-project

# 指定回调URL
./amapi-cli enterprise signup-url \
  --project-id my-project \
  --callback https://example.com/callback
```

### 通知管理

```bash
# 启用通知
./amapi-cli enterprise notifications enable enterprises/LC12345678

# 禁用通知
./amapi-cli enterprise notifications disable enterprises/LC12345678
```

### 查看企业应用

```bash
./amapi-cli enterprise applications enterprises/LC12345678
```

## 策略管理

### 创建策略

```bash
# 创建基本策略
./amapi-cli policy create \
  --enterprise LC12345678 \
  --policy-id basic-policy \
  --name "基本策略"

# 从预设创建策略
./amapi-cli policy create \
  --enterprise LC12345678 \
  --policy-id work-profile \
  --from-preset work_profile
```

### 查看策略

```bash
# 获取特定策略
./amapi-cli policy get enterprises/LC12345678/policies/basic-policy

# 列出所有策略
./amapi-cli policy list --enterprise LC12345678

# 以表格格式显示
./amapi-cli policy list --enterprise LC12345678 --output table
```

### 更新策略

```bash
# 禁用摄像头
./amapi-cli policy update enterprises/LC12345678/policies/basic-policy --camera-disabled

# 启用Kiosk模式
./amapi-cli policy update enterprises/LC12345678/policies/basic-policy --kiosk-mode

# 禁用蓝牙
./amapi-cli policy update enterprises/LC12345678/policies/basic-policy --bluetooth-disabled
```

### 策略预设

```bash
# 查看所有可用预设
./amapi-cli policy presets

# 以表格格式显示预设
./amapi-cli policy presets --output table

# 应用预设到策略
./amapi-cli policy apply-preset \
  --enterprise LC12345678 \
  --policy-id work-policy \
  --preset work_profile
```

### 删除策略

```bash
# 删除策略（需要确认）
./amapi-cli policy delete enterprises/LC12345678/policies/basic-policy

# 强制删除
./amapi-cli policy delete enterprises/LC12345678/policies/basic-policy --force
```

## 设备管理

### 查看设备

```bash
# 列出所有设备
./amapi-cli device list --enterprise LC12345678

# 限制结果数量
./amapi-cli device list --enterprise LC12345678 --page-size 10

# 按状态过滤
./amapi-cli device list --enterprise LC12345678 --filter "state=ACTIVE"

# 获取特定设备信息
./amapi-cli device get enterprises/LC12345678/devices/device123
```

### 设备控制

```bash
# 锁定设备
./amapi-cli device lock enterprises/LC12345678/devices/device123

# 重启设备
./amapi-cli device reboot enterprises/LC12345678/devices/device123

# 恢复出厂设置（危险操作）
./amapi-cli device reset enterprises/LC12345678/devices/device123

# 所有控制命令都支持 --force 跳过确认
./amapi-cli device lock enterprises/LC12345678/devices/device123 --force
```

### 丢失模式

```bash
# 启用丢失模式
./amapi-cli device lost-mode start enterprises/LC12345678/devices/device123

# 停用丢失模式
./amapi-cli device lost-mode stop enterprises/LC12345678/devices/device123
```

### 清除应用数据

```bash
# 清除特定应用数据
./amapi-cli device clear-data enterprises/LC12345678/devices/device123 \
  --package com.example.app

# 强制清除（跳过确认）
./amapi-cli device clear-data enterprises/LC12345678/devices/device123 \
  --package com.example.app --force
```

### 设备筛选

```bash
# 活跃设备
./amapi-cli device filter active --enterprise LC12345678

# 合规设备
./amapi-cli device filter compliant --enterprise LC12345678

# 非合规设备
./amapi-cli device filter non-compliant --enterprise LC12345678

# 按用户筛选
./amapi-cli device filter by-user username --enterprise LC12345678
```

## 注册令牌管理

### 创建注册令牌

```bash
# 创建基本令牌
./amapi-cli enrollment create \
  --enterprise LC12345678 \
  --policy basic-policy

# 创建24小时有效令牌
./amapi-cli enrollment create \
  --enterprise LC12345678 \
  --policy basic-policy \
  --duration 24h

# 创建工作配置文件令牌
./amapi-cli enrollment create \
  --enterprise LC12345678 \
  --policy work-policy \
  --work-profile

# 创建一次性令牌
./amapi-cli enrollment create \
  --enterprise LC12345678 \
  --policy secure-policy \
  --one-time

# 快速创建（24小时有效期）
./amapi-cli enrollment quick \
  --enterprise LC12345678 \
  --policy basic-policy
```

### 查看令牌

```bash
# 获取特定令牌
./amapi-cli enrollment get enterprises/LC12345678/enrollmentTokens/token123

# 列出所有令牌
./amapi-cli enrollment list --enterprise LC12345678

# 只显示活跃令牌
./amapi-cli enrollment list --enterprise LC12345678 --active-only

# 以表格格式显示
./amapi-cli enrollment list --enterprise LC12345678 --output table
```

### 生成QR码

```bash
# 生成基本QR码
./amapi-cli enrollment qr-code enterprises/LC12345678/enrollmentTokens/token123

# 包含WiFi配置的QR码
./amapi-cli enrollment qr-code enterprises/LC12345678/enrollmentTokens/token123 \
  --wifi-ssid MyNetwork \
  --wifi-password mypassword \
  --wifi-security WPA

# 跳过设置向导
./amapi-cli enrollment qr-code enterprises/LC12345678/enrollmentTokens/token123 \
  --skip-setup

# 设置语言
./amapi-cli enrollment qr-code enterprises/LC12345678/enrollmentTokens/token123 \
  --locale zh-CN

# 保存QR码到文件
./amapi-cli enrollment qr-code enterprises/LC12345678/enrollmentTokens/token123 \
  --save qr-code.png
```

### 撤销令牌

```bash
# 撤销令牌（需要确认）
./amapi-cli enrollment revoke enterprises/LC12345678/enrollmentTokens/token123

# 强制撤销
./amapi-cli enrollment revoke enterprises/LC12345678/enrollmentTokens/token123 --force
```

### 令牌统计

```bash
# 查看令牌使用统计
./amapi-cli enrollment stats --enterprise LC12345678
```

## 配置管理

### 查看配置

```bash
# 显示当前配置
./amapi-cli config show

# 显示包含敏感信息的配置
./amapi-cli config show --show-sensitive

# 以表格格式显示
./amapi-cli config show --output table
```

### 设置配置

```bash
# 设置超时时间
./amapi-cli config set timeout 60s

# 设置重试次数
./amapi-cli config set retry-attempts 5

# 设置日志级别
./amapi-cli config set log-level debug

# 设置项目ID
./amapi-cli config set project-id my-project
```

### 验证配置

```bash
# 验证配置有效性
./amapi-cli config validate
```

### 初始化配置

```bash
# 交互式初始化
./amapi-cli config init --interactive

# 指定参数初始化
./amapi-cli config init \
  --project-id my-project \
  --credentials-file /path/to/creds.json
```

### 查看环境变量

```bash
# 显示所有支持的环境变量
./amapi-cli config environment

# 以表格格式显示
./amapi-cli config environment --output table
```

## 健康检查

### 基本检查

```bash
# 完整健康检查
./amapi-cli health check

# 详细健康检查
./amapi-cli health check --detailed

# 快速检查
./amapi-cli health quick
```

### 特定检查

```bash
# 检查API连接
./amapi-cli health connection

# 检查配置状态
./amapi-cli health config
```

## 故障排除

### 常见问题

#### 1. 认证失败

```bash
# 检查凭证文件
ls -la $GOOGLE_APPLICATION_CREDENTIALS

# 验证项目ID
./amapi-cli config show | grep project_id

# 测试连接
./amapi-cli health connection
```

#### 2. 权限不足

确保服务账号具有以下权限：
- Android Management API Service Agent
- 或者自定义角色包含必要的权限

#### 3. API未启用

在 Google Cloud Console 中启用 Android Management API：
```
https://console.cloud.google.com/apis/library/androidmanagement.googleapis.com
```

#### 4. 网络连接问题

```bash
# 检查网络连接
curl -I https://androidmanagement.googleapis.com/

# 使用调试模式查看详细错误
./amapi-cli --debug health connection
```

### 调试技巧

#### 启用详细日志

```bash
# 设置调试级别
./amapi-cli config set log-level debug

# 使用调试模式运行命令
./amapi-cli --debug <command>
```

#### 检查配置

```bash
# 显示完整配置（包含敏感信息）
./amapi-cli config show --show-sensitive

# 验证配置有效性
./amapi-cli config validate
```

#### 测试连接

```bash
# 完整健康检查
./amapi-cli health check --detailed

# 单独测试API连接
./amapi-cli health connection
```

## 输出格式

所有命令都支持多种输出格式：

```bash
# JSON 格式（默认）
./amapi-cli enterprise list my-project

# YAML 格式
./amapi-cli enterprise list my-project --output yaml

# 表格格式（便于阅读）
./amapi-cli enterprise list my-project --output table
```

## 批量操作

### 使用脚本自动化

```bash
#!/bin/bash

# 批量创建策略
policies=("basic" "kiosk" "work")
for policy in "${policies[@]}"; do
  ./amapi-cli policy create \
    --enterprise LC12345678 \
    --policy-id "${policy}-policy" \
    --from-preset "${policy}"
done

# 批量生成注册令牌
for policy in "${policies[@]}"; do
  ./amapi-cli enrollment create \
    --enterprise LC12345678 \
    --policy "${policy}-policy" \
    --duration 7d
done
```

### 使用配置文件

可以通过环境变量或配置文件预设常用参数：

```yaml
# ~/.config/amapi/config.yaml
project_id: "my-default-project"
enterprise_id: "LC12345678"
default_policy: "basic-policy"
```

## 高级用法

### 管道和过滤

```bash
# 使用 jq 过滤 JSON 输出
./amapi-cli device list --enterprise LC12345678 --output json | \
  jq '.items[] | select(.state == "ACTIVE") | .name'

# 提取特定字段
./amapi-cli enterprise list my-project --output json | \
  jq -r '.items[] | .displayName'
```

### 定时任务

```bash
# 定期检查设备合规性
*/30 * * * * /path/to/amapi-cli device filter non-compliant \
  --enterprise LC12345678 --output table > /var/log/non-compliant-devices.log
```

这个手册涵盖了 amapi-cli 工具的所有主要功能。如需更多帮助，请使用 `--help` 选项查看详细的命令说明。