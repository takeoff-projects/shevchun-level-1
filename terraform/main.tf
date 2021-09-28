terraform {
  backend "gcs" {
    bucket  = "roi-takeoff-user51-tf-state"
    prefix = "terraform-website-state"
  }
}

provider "google" {
  project = var.project_id
  region  = var.provider_region
}

resource "google_sql_database_instance" "instance" {
    name = "website-db-instance"
    region = var.provider_region
    database_version = "POSTGRES_13"
    deletion_protection = true
    settings{
        tier="db-f1-micro"
    }
}
resource "google_sql_database" "database"{
    name="website-db"
    instance=google_sql_database_instance.instance.name
}
resource "google_sql_user" "database-user" {
    instance = google_sql_database_instance.instance.name
    name = var.database_user
    password = var.database_password
}

resource "google_cloud_run_service" "default" {
  name     = "website"
  location = var.provider_region

  template {
    spec {
      containers {
        image = "gcr.io/roi-takeoff-user51/go-website:v1.8"
        ports {
            container_port = 8080 
        }
        env {
            name = "ENV"
            value = "production"
        }
        env {
            name = "GOOGLE_CLOUD_PROJECT"
            value = var.project_id
        }
        env {
            name = "DB_URL"
            value = "postgresql://${var.database_user}:${var.database_password}@/website-db?host=/cloudsql/${google_sql_database_instance.instance.connection_name}"
        }
      }
    }
    metadata {
      annotations = {
        "run.googleapis.com/cloudsql-instances"=google_sql_database_instance.instance.connection_name
      }
    }
  }

  traffic {
    percent         = 100
    latest_revision = true
  }
}

data "google_iam_policy" "noauth" {
  binding {
    role = "roles/run.invoker"
    members = [
      "allUsers",
    ]
  }
}

resource "google_cloud_run_service_iam_policy" "noauth" {
  location    = google_cloud_run_service.default.location
  project     = google_cloud_run_service.default.project
  service     = google_cloud_run_service.default.name

  policy_data = data.google_iam_policy.noauth.policy_data
}

output "public_url" {
  value = "${google_cloud_run_service.default.status[0].url}"
}