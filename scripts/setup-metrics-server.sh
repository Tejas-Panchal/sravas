#!/usr/bin/env bash
# Install metrics-server on EKS (required for HPA)
set -euo pipefail

ENV="${1:-dev}"
SCRIPT_DIR="$(cd "$(dirname "$0")" && pwd)"
K8S_DIR="${SCRIPT_DIR}/../infrastructure/kubernetes"

echo "Installing metrics-server for ${ENV}..."
kubectl apply -k "${K8S_DIR}/addons/metrics-server"

echo "Waiting for metrics-server to become ready..."
kubectl -n kube-system rollout status deployment/metrics-server --timeout=120s

echo "Verifying metrics..."
kubectl top nodes 2>/dev/null || echo "metrics not yet available (may take a few seconds)"
echo "metrics-server installation complete."
