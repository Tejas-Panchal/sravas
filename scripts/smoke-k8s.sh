#!/usr/bin/env bash
# smoke-k8s.sh — cluster smoke tests: pod readiness + /health via Ingress
# Usage: ./scripts/smoke-k8s.sh [dev|prod]
set -euo pipefail

ENV="${1:-dev}"
NAMESPACE="${ENV}-sravas"
PASS=0
FAIL=0

pass() { PASS=$((PASS + 1)); echo "  PASS  $1"; }
fail() { FAIL=$((FAIL + 1)); echo "  FAIL  $1"; }

echo "=== Phase 1: Deployment rollout status ==="
for deploy in api-gateway auth-service user-service upload-service streaming-service \
              search-service comment-service analytics-service notification-service \
              scheduler-service frontend; do
  if kubectl -n "$NAMESPACE" rollout status deployment/"${deploy}" --timeout=30s >/dev/null 2>&1; then
    pass "$deploy ready"
  else
    fail "$deploy not ready"
  fi
done

echo ""
echo "=== Phase 2: HPA resources ==="
for hpa in api-gateway auth-service user-service upload-service streaming-service \
           search-service comment-service analytics-service notification-service \
           scheduler-service; do
  if kubectl -n "$NAMESPACE" get hpa "$hpa" >/dev/null 2>&1; then
    pass "HPA $hpa exists"
  else
    fail "HPA $hpa missing"
  fi
done

echo ""
echo "=== Phase 3: Ingress ALB address ==="
ALB_HOST=$(kubectl -n "$NAMESPACE" get ingress sravas-ingress \
  -o jsonpath='{.status.loadBalancer.ingress[0].hostname}' 2>/dev/null || true)
if [ -n "$ALB_HOST" ]; then
  pass "ALB host: $ALB_HOST"
else
  fail "No ALB hostname on Ingress"
  echo "  Skipping Phase 4 (no ALB to test)"
  echo ""
  echo "---"
  echo "Results: $PASS passed, $FAIL failed"
  [ "$FAIL" -eq 0 ] || exit 1
  exit 0
fi

echo ""
echo "=== Phase 4: Health via Ingress ==="
# Frontend
if curl -sf -o /dev/null --max-time 10 "http://${ALB_HOST}/" 2>/dev/null; then
  pass "GET / → frontend (200)"
else
  fail "GET / → frontend"
fi

# API Gateway health
if curl -sf -o /dev/null --max-time 10 "http://${ALB_HOST}/api/health" 2>/dev/null; then
  pass "GET /api/health → api-gateway (200)"
else
  fail "GET /api/health → api-gateway"
fi

# API Gateway ready
if curl -sf -o /dev/null --max-time 10 "http://${ALB_HOST}/api/ready" 2>/dev/null; then
  pass "GET /api/ready → api-gateway (200)"
else
  fail "GET /api/ready → api-gateway"
fi

echo ""
echo "---"
echo "Results: $PASS passed, $FAIL failed"
[ "$FAIL" -eq 0 ] || exit 1
