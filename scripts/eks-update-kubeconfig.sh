#!/usr/bin/env bash
# Configure kubectl for the EKS cluster via Terraform output
set -euo pipefail

ENV="${1:-dev}"
SCRIPT_DIR="$(cd "$(dirname "$0")" && pwd)"
TF_DIR="${SCRIPT_DIR}/../infrastructure/terraform"

cd "$TF_DIR"

echo "Fetching kubeconfig command for ${ENV}..."
CMD=$(terraform output -raw kubeconfig_command -var "environment=${ENV}" 2>/dev/null || terraform output -raw kubeconfig_command)

echo "Running: $CMD"
eval "$CMD"

echo "Verifying cluster connection..."
kubectl cluster-info
kubectl get nodes -o wide
