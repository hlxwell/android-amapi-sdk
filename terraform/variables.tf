variable "project_id" {
  description = "GCP 项目 ID"
  type        = string
}

variable "region" {
  description = "GCP 区域"
  type        = string
  default     = "us-central1"
}

variable "topic_name_prefix" {
  description = "Pub/Sub Topic 名称前缀 (将会创建 {prefix}-cn 和 {prefix}-row)"
  type        = string
  default     = "amapi-events"
}

variable "service_account_id" {
  description = "Service Account ID"
  type        = string
  default     = "amapi-service-account"
}

variable "service_account_display_name" {
  description = "Service Account 显示名称"
  type        = string
  default     = "Android Management API Service Account"
}

variable "create_service_account_key" {
  description = "是否创建 Service Account Key"
  type        = bool
  default     = false
}

variable "save_key_to_file" {
  description = "是否将 Service Account Key 保存到文件"
  type        = bool
  default     = false
}

variable "service_account_key_filename" {
  description = "Service Account Key 文件名"
  type        = string
  default     = "sa-key.json"
}

