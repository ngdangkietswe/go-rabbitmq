provider "google" {
  project = var.project_id
}

resource "google_container_cluster" "primary" {
  name                     = var.cluster_name
  location                 = var.location
  remove_default_node_pool = true
  initial_node_count       = 1

  release_channel { channel = "REGULAR" }
  networking_mode = "VPC_NATIVE"
  ip_allocation_policy {}
}

resource "google_container_node_pool" "primary_nodes" {
  name       = "primary-nodes"
  location   = var.location
  cluster    = google_container_cluster.primary.name
  node_count = var.node_count

  node_config {
    machine_type = var.machine_type
    oauth_scopes = [
      "https://www.googleapis.com/auth/cloud-platform"
    ]
    # For private images from Artifact Registry (not used here), add Workload Identity / GCR access
  }
}

# Kubeconfig output (works with kubectl and CI kubeconfig secret)
output "kubeconfig" {
  value     = <<EOT
apiVersion: v1
kind: Config
clusters:
- cluster:
    certificate-authority-data: ${google_container_cluster.primary.master_auth.0.cluster_ca_certificate}
    server: https://${google_container_cluster.primary.endpoint}
  name: gke_${var.project_id}_${var.location}_${var.cluster_name}
contexts:
- context:
  cluster: gke_${var.project_id}_${var.location}_${var.cluster_name}
  user: gke_${var.project_id}_${var.location}_${var.cluster_name}
  name: gke_${var.project_id}_${var.location}_${var_cluster_name}
current-context: gke_${var.project_id}_${var.location}_${var.cluster_name}
users:
- name: gke_${var.project_id}_${var.location}_${var.cluster_name}
  user:
    auth-provider:
      name: gcp
EOT
  sensitive = true
}