variable "project_id" { type = string }
variable "location" { type = string default = "asia-southeast1-b" }
variable "cluster_name" { type = string default = "go-rabbitmq-gke-cluster" }
variable "node_count" { type = number default = 2 }
variable "machine_type" { type = string default = "e2-standard-2" }