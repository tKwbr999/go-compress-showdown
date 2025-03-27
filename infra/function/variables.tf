variable "project_id" {
  description = "Google Cloud Project ID"
  type        = string
}

variable "region" {
  description = "Google Cloud Region for the function"
  type        = string
  default     = "asia-northeast1" # 必要に応じて変更
}

variable "source_archive_object" {
  description = "Name of the zipped source code object in the GCS bucket"
  type        = string
  default     = "function-source.zip"
}