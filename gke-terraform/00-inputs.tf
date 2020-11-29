variable "project_id" {
  description = "project id"
  default     = "[PROJECT]"
}

data "google_project" "project" {
  project_id = var.project_id
}
