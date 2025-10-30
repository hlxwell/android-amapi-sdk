terraform {
  required_version = ">= 1.0"

  required_providers {
    google = {
      source  = "hashicorp/google"
      version = "~> 5.0"
    }
  }
}

provider "google" {
  project = var.project_id
  region  = var.region
}

# 1. 启用 Android Management API
resource "google_project_service" "androidmanagement" {
  project = var.project_id
  service = "androidmanagement.googleapis.com"

  disable_on_destroy = false
}

# 启用 Pub/Sub API (Topic 需要)
resource "google_project_service" "pubsub" {
  project = var.project_id
  service = "pubsub.googleapis.com"

  disable_on_destroy = false
}

# 启用 IAM API
resource "google_project_service" "iam" {
  project = var.project_id
  service = "iam.googleapis.com"

  disable_on_destroy = false
}

# 2. 创建 AMAPI Pub/Sub Topics (CN 和 ROW 区域)
# CN Topic - 中国区域
resource "google_pubsub_topic" "amapi_events_cn" {
  name    = "${var.topic_name_prefix}-cn"
  project = var.project_id

  message_retention_duration = "604800s" # 7 days

  labels = {
    purpose = "amapi-events"
    region  = "cn"
    managed = "terraform"
  }

  depends_on = [google_project_service.pubsub]
}

# ROW Topic - 世界其他地区
resource "google_pubsub_topic" "amapi_events_row" {
  name    = "${var.topic_name_prefix}-row"
  project = var.project_id

  message_retention_duration = "604800s" # 7 days

  labels = {
    purpose = "amapi-events"
    region  = "row"
    managed = "terraform"
  }

  depends_on = [google_project_service.pubsub]
}

# 创建 Dead Letter Topics
resource "google_pubsub_topic" "amapi_events_cn_deadletter" {
  name    = "${var.topic_name_prefix}-cn-deadletter"
  project = var.project_id

  message_retention_duration = "604800s"

  labels = {
    purpose = "amapi-events-deadletter"
    region  = "cn"
    managed = "terraform"
  }

  depends_on = [google_project_service.pubsub]
}

resource "google_pubsub_topic" "amapi_events_row_deadletter" {
  name    = "${var.topic_name_prefix}-row-deadletter"
  project = var.project_id

  message_retention_duration = "604800s"

  labels = {
    purpose = "amapi-events-deadletter"
    region  = "row"
    managed = "terraform"
  }

  depends_on = [google_project_service.pubsub]
}

# 创建订阅 - CN
resource "google_pubsub_subscription" "amapi_events_cn_sub" {
  name    = "${var.topic_name_prefix}-cn-subscription"
  topic   = google_pubsub_topic.amapi_events_cn.name
  project = var.project_id

  # 消息确认截止时间
  ack_deadline_seconds = 20

  # 消息保留时间
  message_retention_duration = "604800s"

  # 重试策略
  retry_policy {
    minimum_backoff = "10s"
    maximum_backoff = "600s"
  }

  # Dead Letter 配置
  dead_letter_policy {
    dead_letter_topic     = google_pubsub_topic.amapi_events_cn_deadletter.id
    max_delivery_attempts = 5
  }

  # 过期策略 (31天未使用则删除)
  expiration_policy {
    ttl = "2678400s" # 31 days
  }

  labels = {
    purpose = "amapi-events"
    region  = "cn"
    managed = "terraform"
  }
}

# 创建订阅 - ROW
resource "google_pubsub_subscription" "amapi_events_row_sub" {
  name    = "${var.topic_name_prefix}-row-subscription"
  topic   = google_pubsub_topic.amapi_events_row.name
  project = var.project_id

  # 消息确认截止时间
  ack_deadline_seconds = 20

  # 消息保留时间
  message_retention_duration = "604800s"

  # 重试策略
  retry_policy {
    minimum_backoff = "10s"
    maximum_backoff = "600s"
  }

  # Dead Letter 配置
  dead_letter_policy {
    dead_letter_topic     = google_pubsub_topic.amapi_events_row_deadletter.id
    max_delivery_attempts = 5
  }

  # 过期策略 (31天未使用则删除)
  expiration_policy {
    ttl = "2678400s" # 31 days
  }

  labels = {
    purpose = "amapi-events"
    region  = "row"
    managed = "terraform"
  }
}

# 为 Dead Letter Topics 创建订阅
resource "google_pubsub_subscription" "amapi_events_cn_deadletter_sub" {
  name    = "${var.topic_name_prefix}-cn-deadletter-subscription"
  topic   = google_pubsub_topic.amapi_events_cn_deadletter.name
  project = var.project_id

  ack_deadline_seconds       = 20
  message_retention_duration = "604800s"

  labels = {
    purpose = "amapi-events-deadletter"
    region  = "cn"
    managed = "terraform"
  }
}

resource "google_pubsub_subscription" "amapi_events_row_deadletter_sub" {
  name    = "${var.topic_name_prefix}-row-deadletter-subscription"
  topic   = google_pubsub_topic.amapi_events_row_deadletter.name
  project = var.project_id

  ack_deadline_seconds       = 20
  message_retention_duration = "604800s"

  labels = {
    purpose = "amapi-events-deadletter"
    region  = "row"
    managed = "terraform"
  }
}

# 3. 创建 Service Account
resource "google_service_account" "amapi_sa" {
  account_id   = var.service_account_id
  display_name = var.service_account_display_name
  description  = "Service account for Android Management API and Pub/Sub operations"
  project      = var.project_id

  depends_on = [google_project_service.iam]
}

# 授予 Service Account Android Management API 权限
resource "google_project_iam_member" "amapi_admin" {
  project = var.project_id
  role    = "roles/androidmanagement.user"
  member  = "serviceAccount:${google_service_account.amapi_sa.email}"

  depends_on = [google_project_service.androidmanagement]
}

# 授予 Service Account Pub/Sub Publisher 权限 - CN Topic
resource "google_pubsub_topic_iam_member" "amapi_publisher_cn" {
  project = var.project_id
  topic   = google_pubsub_topic.amapi_events_cn.name
  role    = "roles/pubsub.publisher"
  member  = "serviceAccount:${google_service_account.amapi_sa.email}"
}

# 授予 Service Account Pub/Sub Publisher 权限 - ROW Topic
resource "google_pubsub_topic_iam_member" "amapi_publisher_row" {
  project = var.project_id
  topic   = google_pubsub_topic.amapi_events_row.name
  role    = "roles/pubsub.publisher"
  member  = "serviceAccount:${google_service_account.amapi_sa.email}"
}

# 授予 Service Account Pub/Sub Subscriber 权限 - CN Subscription
resource "google_pubsub_subscription_iam_member" "amapi_subscriber_cn" {
  project      = var.project_id
  subscription = google_pubsub_subscription.amapi_events_cn_sub.name
  role         = "roles/pubsub.subscriber"
  member       = "serviceAccount:${google_service_account.amapi_sa.email}"
}

# 授予 Service Account Pub/Sub Subscriber 权限 - ROW Subscription
resource "google_pubsub_subscription_iam_member" "amapi_subscriber_row" {
  project      = var.project_id
  subscription = google_pubsub_subscription.amapi_events_row_sub.name
  role         = "roles/pubsub.subscriber"
  member       = "serviceAccount:${google_service_account.amapi_sa.email}"
}

# 授予 Service Account 查看者权限 (用于订阅管理)
resource "google_project_iam_member" "amapi_viewer" {
  project = var.project_id
  role    = "roles/pubsub.viewer"
  member  = "serviceAccount:${google_service_account.amapi_sa.email}"
}

# 为 Android Management API 授权发布到 Topics
# Android Management API 使用 Google-managed service account
data "google_project" "project" {
  project_id = var.project_id
}

# 授予 Android Management API service account 发布权限 - CN Topic
resource "google_pubsub_topic_iam_member" "amapi_service_publisher_cn" {
  project = var.project_id
  topic   = google_pubsub_topic.amapi_events_cn.name
  role    = "roles/pubsub.publisher"
  member  = "serviceAccount:service-${data.google_project.project.number}@gcp-sa-androidmanagement.iam.gserviceaccount.com"

  depends_on = [google_project_service.androidmanagement]
}

# 授予 Android Management API service account 发布权限 - ROW Topic
resource "google_pubsub_topic_iam_member" "amapi_service_publisher_row" {
  project = var.project_id
  topic   = google_pubsub_topic.amapi_events_row.name
  role    = "roles/pubsub.publisher"
  member  = "serviceAccount:service-${data.google_project.project.number}@gcp-sa-androidmanagement.iam.gserviceaccount.com"

  depends_on = [google_project_service.androidmanagement]
}

# 创建 Service Account Key (可选,用于本地开发)
resource "google_service_account_key" "amapi_sa_key" {
  count = var.create_service_account_key ? 1 : 0

  service_account_id = google_service_account.amapi_sa.name
}

# 将 key 保存到文件 (可选)
resource "local_file" "amapi_sa_key_file" {
  count = var.create_service_account_key && var.save_key_to_file ? 1 : 0

  content  = base64decode(google_service_account_key.amapi_sa_key[0].private_key)
  filename = "${path.module}/${var.service_account_key_filename}"

  file_permission = "0600"
}

