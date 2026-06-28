#!/usr/bin/env bash
# Generate and apply Secrets from Terraform outputs + env vars
set -euo pipefail

ENV="${1:-dev}"
SCRIPT_DIR="$(cd "$(dirname "$0")" && pwd)"
TF_DIR="${SCRIPT_DIR}/../infrastructure/terraform"
NS="sravas"

cd "$TF_DIR"

echo "Reading Terraform outputs for ${ENV}..."
DB_USER=$(terraform output -raw db_username)
DB_PASS=$(terraform output -raw db_password)
DB_HOST=$(terraform output -raw db_endpoint)

# Values from env vars with safe defaults (same as docker-compose)
MONGO_ROOT_PASS="${MONGO_ROOT_PASSWORD:-changeme}"
JWT_SECRET="${JWT_SECRET:-change-me}"
JWT_REFRESH="${JWT_REFRESH_SECRET:-change-me}"
SENDGRID_KEY="${SENDGRID_API_KEY:-}"
AWS_KEY="${AWS_ACCESS_KEY_ID:-}"
AWS_SECRET="${AWS_SECRET_ACCESS_KEY:-}"

# Construct DATABASE_URL for user-service and scheduler-service
DATABASE_URL="postgresql://${DB_USER}:${DB_PASS}@${DB_HOST}/sravas"

echo "Creating sravas-secrets in namespace ${NS}..."
kubectl create secret generic sravas-secrets -n "$NS" \
  --from-literal=mongo-root-username=admin \
  --from-literal=mongo-root-password="${MONGO_ROOT_PASS}" \
  --from-literal=db-username="${DB_USER}" \
  --from-literal=db-password="${DB_PASS}" \
  --from-literal=db-endpoint="${DB_HOST}" \
  --from-literal=jwt-secret="${JWT_SECRET}" \
  --from-literal=jwt-refresh-secret="${JWT_REFRESH}" \
  --from-literal=aws-access-key-id="${AWS_KEY}" \
  --from-literal=aws-secret-access-key="${AWS_SECRET}" \
  --from-literal=sendgrid-api-key="${SENDGRID_KEY}" \
  --from-literal=database-url="${DATABASE_URL}" \
  --dry-run=client -o yaml | kubectl apply -f -

echo "Verifying..."
kubectl get secret sravas-secrets -n "$NS"

echo ""
echo "Secret keys:"
kubectl get secret sravas-secrets -n "$NS" -o jsonpath='{.data}' | python3 -m json.tool 2>/dev/null | grep '"' | sed 's/.*"\(.*\)".*/  - \1/' || \
  kubectl get secret sravas-secrets -n "$NS" -o json | grep -o '"[^"]*"' | grep -v '{' | grep -v '}' | head -20
