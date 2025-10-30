output "project_id" {
  description = "GCP 项目 ID"
  value       = var.project_id
}

output "amapi_topic_name" {
  description = "AMAPI Pub/Sub Topic 名称"
  value       = google_pubsub_topic.amapi_events.name
}

output "amapi_topic_id" {
  description = "AMAPI Pub/Sub Topic 完整 ID"
  value       = google_pubsub_topic.amapi_events.id
}

output "amapi_deadletter_topic_name" {
  description = "AMAPI Dead Letter Topic 名称"
  value       = google_pubsub_topic.amapi_events_deadletter.name
}

output "amapi_subscription_name" {
  description = "AMAPI Pub/Sub Subscription 名称"
  value       = google_pubsub_subscription.amapi_events_sub.name
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

    1. Pub/Sub Topic: ${google_pubsub_topic.amapi_events.id}
    2. Subscription: ${google_pubsub_subscription.amapi_events_sub.id}
    3. Service Account: ${google_service_account.amapi_sa.email}

    下一步操作:
    -----------
    1. 如果需要下载 Service Account Key:
       gcloud iam service-accounts keys create sa-key.json \
         --iam-account=${google_service_account.amapi_sa.email}

    2. 在您的应用配置中使用以下设置:
       - PROJECT_ID: ${var.project_id}
       - TOPIC_NAME: ${google_pubsub_topic.amapi_events.name}
       - SERVICE_ACCOUNT: ${google_service_account.amapi_sa.email}

    3. 测试 Pub/Sub:
       gcloud pubsub topics publish ${google_pubsub_topic.amapi_events.name} \
         --message="Test message"

    4. 查看订阅消息:
       gcloud pubsub subscriptions pull ${google_pubsub_subscription.amapi_events_sub.name} \
         --auto-ack --limit=10
  EOT
}

