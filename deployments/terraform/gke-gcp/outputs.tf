output "project_id" { value = var.project_id }
output "location" { value = var.location }
output "cluster_name" { value = google_container_cluster.primary.name }