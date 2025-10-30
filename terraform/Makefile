.PHONY: help init plan apply destroy clean fmt validate output

# 默认目标
.DEFAULT_GOAL := help

help: ## 显示此帮助信息
	@echo "Android Management API - Terraform 管理命令"
	@echo ""
	@echo "使用方法: make [target]"
	@echo ""
	@echo "可用目标:"
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "  \033[36m%-15s\033[0m %s\n", $$1, $$2}'

init: ## 初始化 Terraform
	@echo "🔧 初始化 Terraform..."
	terraform init

plan: ## 查看将要应用的更改
	@echo "📋 生成执行计划..."
	terraform plan

apply: ## 应用 Terraform 配置
	@echo "🚀 应用 Terraform 配置..."
	terraform apply

apply-auto: ## 自动应用配置(不需要确认)
	@echo "🚀 自动应用 Terraform 配置..."
	terraform apply -auto-approve

destroy: ## 销毁所有 Terraform 管理的资源
	@echo "⚠️  警告: 将销毁所有资源!"
	terraform destroy

destroy-auto: ## 自动销毁资源(不需要确认)
	@echo "⚠️  自动销毁所有资源..."
	terraform destroy -auto-approve

fmt: ## 格式化 Terraform 配置文件
	@echo "✨ 格式化 Terraform 文件..."
	terraform fmt -recursive

validate: ## 验证 Terraform 配置
	@echo "✅ 验证 Terraform 配置..."
	terraform validate

output: ## 显示所有输出
	@echo "📤 显示输出..."
	terraform output

output-json: ## 以 JSON 格式显示输出
	@echo "📤 显示 JSON 格式输出..."
	terraform output -json

show: ## 显示当前状态
	@echo "📊 显示当前状态..."
	terraform show

refresh: ## 刷新状态
	@echo "🔄 刷新状态..."
	terraform refresh

graph: ## 生成资源依赖图
	@echo "📈 生成资源依赖图..."
	terraform graph | dot -Tpng > graph.png
	@echo "✅ 依赖图已保存到 graph.png"

clean: ## 清理本地文件
	@echo "🧹 清理本地文件..."
	rm -rf .terraform
	rm -f .terraform.lock.hcl
	rm -f terraform.tfstate*
	rm -f graph.png
	@echo "✅ 清理完成"

test-cn: ## 测试 CN Topic
	@echo "🧪 测试 CN Topic..."
	@TOPIC=$$(terraform output -raw amapi_topic_cn_name); \
	echo "发布测试消息到 $$TOPIC..."; \
	gcloud pubsub topics publish $$TOPIC --message="Test from Makefile - CN"
	@SUB=$$(terraform output -raw amapi_subscription_cn_name); \
	echo "从订阅 $$SUB 拉取消息..."; \
	gcloud pubsub subscriptions pull $$SUB --auto-ack --limit=5

test-row: ## 测试 ROW Topic
	@echo "🧪 测试 ROW Topic..."
	@TOPIC=$$(terraform output -raw amapi_topic_row_name); \
	echo "发布测试消息到 $$TOPIC..."; \
	gcloud pubsub topics publish $$TOPIC --message="Test from Makefile - ROW"
	@SUB=$$(terraform output -raw amapi_subscription_row_name); \
	echo "从订阅 $$SUB 拉取消息..."; \
	gcloud pubsub subscriptions pull $$SUB --auto-ack --limit=5

test-all: test-cn test-row ## 测试所有 Topics

download-key: ## 下载 Service Account Key
	@echo "🔑 下载 Service Account Key..."
	@SA_EMAIL=$$(terraform output -raw service_account_email); \
	echo "Service Account: $$SA_EMAIL"; \
	gcloud iam service-accounts keys create sa-key.json --iam-account=$$SA_EMAIL
	@echo "✅ Key 已保存到 sa-key.json"

setup: init ## 完整设置流程
	@echo "📦 创建配置文件..."
	@if [ ! -f terraform.tfvars ]; then \
		cp terraform.tfvars.example terraform.tfvars; \
		echo "⚠️  请编辑 terraform.tfvars 设置你的 project_id"; \
	else \
		echo "✅ terraform.tfvars 已存在"; \
	fi

check: fmt validate ## 检查配置(格式化 + 验证)
	@echo "✅ 配置检查完成"

