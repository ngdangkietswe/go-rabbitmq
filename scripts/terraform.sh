#!/usr/bin/env sh

set -e  # stop script on error

# Move to terraform directory
cd ../deployments/terraform/gke-gcp || exit 1

# Set your GCP project ID
export TF_VAR_project_id="your-gcp-project"

# Initialize terraform
terraform init

# Apply terraform changes
terraform apply -auto-approve

# Export kubeconfig
terraform output -raw kubeconfig > kubeconfig_gke

# Optional: encode kubeconfig to base64 for GitHub secrets
if command -v base64 >/dev/null 2>&1; then
  # macOS uses `-b`, Linux uses `-w`
  if base64 --help 2>&1 | grep -q -- "-w"; then
    base64 -w0 kubeconfig_gke > kubeconfig_gke.b64
  else
    base64 -b0 kubeconfig_gke > kubeconfig_gke.b64
  fi
fi
