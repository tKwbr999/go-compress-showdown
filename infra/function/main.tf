terraform {
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

# Cloud Functions のソースコードをアップロードするための GCS バケット
resource "google_storage_bucket" "source_bucket" {
  name          = "${var.project_id}-function-source"
  location      = var.region
  force_destroy = true # 本番環境では false を推奨
}

# Cloud Functions (第2世代) の定義
resource "google_cloudfunctions2_function" "compress_function" {
  name     = "compress-showdown-function"
  location = var.region
  project  = var.project_id

  build_config {
    runtime     = "go122" # go.mod で確認したバージョンに対応
    entry_point = ""      # HTTP 関数のため省略 (funcframework が処理)
    source {
      storage_source {
        bucket = google_storage_bucket.source_bucket.name
        object = var.source_archive_object # アップロードされる zip ファイル名
      }
    }
  }

  service_config {
    max_instance_count = 10 # 必要に応じて調整
    min_instance_count = 0  # コスト削減のためアイドル時は 0 に
    available_memory   = "256Mi" # 必要に応じて調整
    timeout_seconds    = 60      # 必要に応じて調整
    ingress_settings   = "ALLOW_ALL" # HTTP トリガーのため全許可
    # environment_variables = {
    #   KEY = "VALUE"
    # }
  }

  # IAM 設定: Cloud Run Invoker ロールを付与して公開アクセスを許可
  # より厳密な制御が必要な場合は適宜変更してください
  lifecycle {
    prevent_destroy = false # 本番環境では true を推奨
  }
}

# Cloud Functions を公開アクセス可能にするための IAM バインディング
resource "google_cloudfunctions2_function_iam_member" "invoker" {
  project        = google_cloudfunctions2_function.compress_function.project
  location       = google_cloudfunctions2_function.compress_function.location
  cloud_function = google_cloudfunctions2_function.compress_function.name
  role           = "roles/run.invoker"
  member         = "allUsers"
}
