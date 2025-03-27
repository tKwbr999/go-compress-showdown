output "function_url" {
  description = "The HTTPS trigger URL for the deployed Cloud Function"
  value       = google_cloudfunctions2_function.compress_function.service_config[0].uri
}