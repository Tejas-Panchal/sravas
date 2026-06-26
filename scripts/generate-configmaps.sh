#!/usr/bin/env bash
# Generate and apply ConfigMaps per service from Terraform outputs
set -euo pipefail

ENV="${1:-dev}"
SCRIPT_DIR="$(cd "$(dirname "$0")" && pwd)"
TF_DIR="${SCRIPT_DIR}/../infrastructure/terraform"
NS="sravas"

cd "$TF_DIR"

echo "Reading Terraform outputs for ${ENV}..."
DB_ENDPOINT=$(terraform output -raw db_endpoint)
REDIS_ENDPOINT=$(terraform output -raw redis_endpoint)
S3_BUCKET=$(terraform output -raw upload_bucket_id)
CDN_URL=$(terraform output -raw cloudfront_domain_name)

echo "  DB_ENDPOINT  = ${DB_ENDPOINT}"
echo "  REDIS_HOST   = ${REDIS_ENDPOINT}"
echo "  S3_BUCKET    = ${S3_BUCKET}"
echo "  CDN_URL      = ${CDN_URL}"

# ── api-gateway ──────────────────────────────────────────────────
echo "Creating api-gateway-config..."
kubectl create configmap api-gateway-config -n "$NS" \
  --from-literal=PORT=8000 \
  --from-literal=AUTH_SERVICE_URL=http://auth-service:3001 \
  --from-literal=USER_SERVICE_URL=http://user-service:3002 \
  --from-literal=UPLOAD_SERVICE_URL=http://upload-service:3003 \
  --from-literal=STREAMING_SERVICE_URL=http://streaming-service:3004 \
  --from-literal=SEARCH_SERVICE_URL=http://search-service:3005 \
  --from-literal=COMMENT_SERVICE_URL=http://comment-service:3006 \
  --from-literal=ANALYTICS_SERVICE_URL=http://analytics-service:3007 \
  --from-literal=NOTIFICATION_SERVICE_URL=http://notification-service:3008 \
  --dry-run=client -o yaml | kubectl apply -f -

# ── auth-service ─────────────────────────────────────────────────
echo "Creating auth-service-config..."
kubectl create configmap auth-service-config -n "$NS" \
  --from-literal=PORT=3001 \
  --from-literal=REDIS_URL="redis://${REDIS_ENDPOINT}:6379" \
  --from-literal=JWT_EXPIRY=15m \
  --from-literal=JWT_REFRESH_EXPIRY=7d \
  --dry-run=client -o yaml | kubectl apply -f -

# ── user-service ─────────────────────────────────────────────────
echo "Creating user-service-config..."
kubectl create configmap user-service-config -n "$NS" \
  --from-literal=PORT=3002 \
  --from-literal=DB_HOST="${DB_ENDPOINT}" \
  --dry-run=client -o yaml | kubectl apply -f -

# ── upload-service ───────────────────────────────────────────────
echo "Creating upload-service-config..."
kubectl create configmap upload-service-config -n "$NS" \
  --from-literal=PORT=3003 \
  --from-literal=KAFKA_BROKERS=kafka:9092 \
  --from-literal=AWS_REGION=us-east-1 \
  --from-literal=S3_BUCKET="${S3_BUCKET}" \
  --from-literal=CDN_URL="${CDN_URL}" \
  --dry-run=client -o yaml | kubectl apply -f -

# ── streaming-service ────────────────────────────────────────────
echo "Creating streaming-service-config..."
kubectl create configmap streaming-service-config -n "$NS" \
  --from-literal=PORT=3004 \
  --from-literal=REDIS_URL="redis://${REDIS_ENDPOINT}:6379" \
  --from-literal=KAFKA_BROKERS=kafka:9092 \
  --dry-run=client -o yaml | kubectl apply -f -

# ── search-service ───────────────────────────────────────────────
echo "Creating search-service-config..."
kubectl create configmap search-service-config -n "$NS" \
  --from-literal=PORT=3005 \
  --from-literal=ELASTICSEARCH_URL=http://elasticsearch:9200 \
  --dry-run=client -o yaml | kubectl apply -f -

# ── comment-service ──────────────────────────────────────────────
echo "Creating comment-service-config..."
kubectl create configmap comment-service-config -n "$NS" \
  --from-literal=PORT=3006 \
  --from-literal=MONGODB_URL=mongodb://mongodb:27017/sravas \
  --dry-run=client -o yaml | kubectl apply -f -

# ── analytics-service ────────────────────────────────────────────
echo "Creating analytics-service-config..."
kubectl create configmap analytics-service-config -n "$NS" \
  --from-literal=PORT=3007 \
  --from-literal=KAFKA_BROKERS=kafka:9092 \
  --dry-run=client -o yaml | kubectl apply -f -

# ── notification-service ─────────────────────────────────────────
echo "Creating notification-service-config..."
kubectl create configmap notification-service-config -n "$NS" \
  --from-literal=PORT=3008 \
  --from-literal=KAFKA_BROKERS=kafka:9092 \
  --dry-run=client -o yaml | kubectl apply -f -

# ── scheduler-service ────────────────────────────────────────────
echo "Creating scheduler-service-config..."
kubectl create configmap scheduler-service-config -n "$NS" \
  --from-literal=PORT=3009 \
  --from-literal=REDIS_URL="redis://${REDIS_ENDPOINT}:6379" \
  --from-literal=DB_HOST="${DB_ENDPOINT}" \
  --dry-run=client -o yaml | kubectl apply -f -

echo ""
echo "All 11 ConfigMaps applied to namespace ${NS}."
kubectl get configmap -n "$NS" -l app.kubernetes.io/managed-by=terraform
echo ""
echo "Note: Sensitive values (JWT_SECRET, DB_PASSWORD, SENDGRID_API_KEY, AWS creds)"
echo "are NOT in ConfigMaps. Create them as Secrets in Day 2."
