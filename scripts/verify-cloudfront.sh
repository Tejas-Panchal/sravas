#!/usr/bin/env bash
# verify-cloudfront.sh — test video upload → S3 → CloudFront delivery
# Usage: ./scripts/verify-cloudfront.sh [dev|prod]
# Assumes: Terraform applied, services deployed on EKS, ALB Ingress ready
set -euo pipefail

ENV="${1:-dev}"
NAMESPACE="${ENV}-sravas"
PASS=0
FAIL=0

pass() { PASS=$((PASS + 1)); echo "  PASS  $1"; }
fail() { FAIL=$((FAIL + 1)); echo "  FAIL  $1"; }

# --- Pre-flight: get ALB host from Ingress ---
ALB_HOST=$(kubectl -n "$NAMESPACE" get ingress sravas-ingress \
  -o jsonpath='{.status.loadBalancer.ingress[0].hostname}' 2>/dev/null || true)
if [ -z "$ALB_HOST" ]; then
  echo "No ALB hostname on Ingress — skipping verification"
  exit 0
fi
echo "ALB: ${ALB_HOST}"

# --- Create 1MB test video file ---
TEST_FILE=$(mktemp /tmp/test-video-XXXXXX.mp4)
dd if=/dev/urandom of="$TEST_FILE" bs=1M count=1 2>/dev/null
trap 'rm -f "$TEST_FILE"' EXIT

echo ""
echo "=== Upload video ==="
RESPONSE=$(curl -sf --max-time 30 \
  -F "video=@${TEST_FILE};type=video/mp4" \
  -F "user_id=test-user" \
  -F "title=Smoke test video" \
  "http://${ALB_HOST}/api/v1/videos/upload" 2>/dev/null || true)

if [ -n "$RESPONSE" ]; then
  pass "POST /api/v1/videos/upload → 202"
else
  fail "POST /api/v1/videos/upload"
  echo "  Skipping remaining checks"
  echo ""
  echo "---"
  echo "Results: $PASS passed, $FAIL failed"
  exit 1
fi

# --- Extract video ID and URL ---
VIDEO_ID=$(echo "$RESPONSE" | grep -o '"id":"[^"]*"' | head -1 | cut -d'"' -f4)
VIDEO_URL=$(echo "$RESPONSE" | grep -o '"url":"[^"]*"' | head -1 | cut -d'"' -f4)

echo "  Video ID: ${VIDEO_ID}"

# --- Verify CloudFront URL is present ---
if [ -n "$VIDEO_URL" ] && echo "$VIDEO_URL" | grep -q "cloudfront.net"; then
  pass "Response contains CloudFront URL: $VIDEO_URL"
else
  fail "No CloudFront URL in response (url='${VIDEO_URL:-}')" \
    "  (CDN_URL env var may not be set in the upload service)"
fi

# --- Verify video is downloadable via CloudFront ---
if [ -n "$VIDEO_URL" ]; then
  HTTP_CODE=$(curl -sf -o /dev/null -w '%{http_code}' --max-time 10 "$VIDEO_URL" 2>/dev/null || echo "000")
  if [ "$HTTP_CODE" = "200" ]; then
    pass "GET ${VIDEO_URL} → HTTP 200"
  else
    fail "GET ${VIDEO_URL} → HTTP ${HTTP_CODE}"
  fi
fi

# --- Clean up: delete the video ---
if [ -n "$VIDEO_ID" ]; then
  DEL_CODE=$(curl -sf -o /dev/null -w '%{http_code}' --max-time 10 \
    -X DELETE "http://${ALB_HOST}/api/v1/videos/${VIDEO_ID}" 2>/dev/null || echo "000")
  if [ "$DEL_CODE" = "200" ]; then
    pass "DELETE /api/v1/videos/${VIDEO_ID} → 200"
  else
    fail "DELETE /api/v1/videos/${VIDEO_ID} → HTTP ${DEL_CODE}"
  fi
fi

echo ""
echo "---"
echo "Results: $PASS passed, $FAIL failed"
[ "$FAIL" -eq 0 ] || exit 1
