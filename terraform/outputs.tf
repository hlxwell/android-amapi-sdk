output "project_id" {
  description = "GCP 项目 ID"
  value       = var.project_id
}

# CN Topic 输出
output "amapi_topic_cn_name" {
  description = "AMAPI Pub/Sub Topic 名称 - CN"
  value       = google_pubsub_topic.amapi_events_cn.name
}

output "amapi_topic_cn_id" {
  description = "AMAPI Pub/Sub Topic 完整 ID - CN"
  value       = google_pubsub_topic.amapi_events_cn.id
}

output "amapi_subscription_cn_name" {
  description = "AMAPI Pub/Sub Subscription 名称 - CN"
  value       = google_pubsub_subscription.amapi_events_cn_sub.name
}

output "amapi_deadletter_topic_cn_name" {
  description = "AMAPI Dead Letter Topic 名称 - CN"
  value       = google_pubsub_topic.amapi_events_cn_deadletter.name
}

# ROW Topic 输出
output "amapi_topic_row_name" {
  description = "AMAPI Pub/Sub Topic 名称 - ROW"
  value       = google_pubsub_topic.amapi_events_row.name
}

output "amapi_topic_row_id" {
  description = "AMAPI Pub/Sub Topic 完整 ID - ROW"
  value       = google_pubsub_topic.amapi_events_row.id
}

output "amapi_subscription_row_name" {
  description = "AMAPI Pub/Sub Subscription 名称 - ROW"
  value       = google_pubsub_subscription.amapi_events_row_sub.name
}

output "amapi_deadletter_topic_row_name" {
  description = "AMAPI Dead Letter Topic 名称 - ROW"
  value       = google_pubsub_topic.amapi_events_row_deadletter.name
}

output "service_account_email" {
  description = "Service Account 邮箱地址"
  value       = google_service_account.amapi_sa.email
}

output "service_account_name" {
  description = "Service Account 名称"
  value       = google_service_account.amapi_sa.name
}

output "service_account_unique_id" {
  description = "Service Account 唯一 ID"
  value       = google_service_account.amapi_sa.unique_id
}

output "service_account_key_path" {
  description = "Service Account Key 文件路径 (如果创建)"
  value       = var.create_service_account_key && var.save_key_to_file ? local_file.amapi_sa_key_file[0].filename : "未创建"
}

output "androidmanagement_service_account" {
  description = "Android Management API 的 Google 管理服务账号"
  value       = "service-${data.google_project.project.number}@gcp-sa-androidmanagement.iam.gserviceaccount.com"
}

output "setup_instructions" {
  description = "设置说明"
  value = <<-EOT
    ========================================
    Android Management API Terraform 部署完成
    ========================================

    📍 CN 区域 (中国):
    - Topic: ${google_pubsub_topic.amapi_events_cn.id}
    - Subscription: ${google_pubsub_subscription.amapi_events_cn_sub.id}

    🌍 ROW 区域 (世界其他地区):
    - Topic: ${google_pubsub_topic.amapi_events_row.id}
    - Subscription: ${google_pubsub_subscription.amapi_events_row_sub.id}

    👤 Service Account: ${google_service_account.amapi_sa.email}

    下一步操作:
    -----------
    1. 如果需要下载 Service Account Key:
       gcloud iam service-accounts keys create sa-key.json \
         --iam-account=${google_service_account.amapi_sa.email}

    2. 在您的应用配置中使用以下设置:
       - PROJECT_ID: ${var.project_id}
       - TOPIC_CN: ${google_pubsub_topic.amapi_events_cn.name}
       - TOPIC_ROW: ${google_pubsub_topic.amapi_events_row.name}
       - SERVICE_ACCOUNT: ${google_service_account.amapi_sa.email}

    3. 测试 Pub/Sub (CN):
       gcloud pubsub topics publish ${google_pubsub_topic.amapi_events_cn.name} \
         --message="Test message for CN"

       gcloud pubsub subscriptions pull ${google_pubsub_subscription.amapi_events_cn_sub.name} \
         --auto-ack --limit=10

    4. 测试 Pub/Sub (ROW):
       gcloud pubsub topics publish ${google_pubsub_topic.amapi_events_row.name} \
         --message="Test message for ROW"

       gcloud pubsub subscriptions pull ${google_pubsub_subscription.amapi_events_row_sub.name} \
         --auto-ack --limit=10
  EOT
}

