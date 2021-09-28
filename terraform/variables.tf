variable "project_id" {
  type = string
}

variable "provider_region" {
  type    = string
  default = "us-central1"
}

variable "database_user" {
  type    = string
  default = "root"
}
variable "database_password" {
  type    = string
}
