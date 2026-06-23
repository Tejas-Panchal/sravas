#!/bin/sh
# verify-network.sh — cross-service connectivity tests on sravas-net
# Assumes `docker compose up -d` is running and all services are healthy.
# Notes: frontend (nginx:alpine) has no curl so it's source-skipped;
#        infra services (postgres, redis, etc.) have HTTP health endpoints
#        tested as targets but lack curl for outbound tests.

SERVICES="api-gateway:8000 auth-service:3001 user-service:3002 upload-service:3003 streaming-service:3004 search-service:3005 comment-service:3006 analytics-service:3007 notification-service:3008 scheduler-service:3009 frontend:80"

PASS=0
FAIL=0

echo "Verifying cross-service connectivity on sravas-net..."
echo ""

for src in $SERVICES; do
  src_name="${src%:*}"
  # frontend has wget but not curl — skip as source
  [ "$src_name" = "frontend" ] && continue

  for tgt in $SERVICES; do
    tgt_name="${tgt%:*}"
    tgt_port="${tgt#*:}"
    [ "$src_name" = "$tgt_name" ] && continue

    url="http://$tgt_name:$tgt_port/health"
    [ "$tgt_name" = "frontend" ] && url="http://$tgt_name:$tgt_port/"

    if docker compose exec -T "$src_name" curl -f -s -o /dev/null "$url" 2>/dev/null; then
      echo "  PASS  $src_name → $tgt_name:$tgt_port"
      PASS=$((PASS + 1))
    else
      echo "  FAIL  $src_name → $tgt_name:$tgt_port"
      FAIL=$((FAIL + 1))
    fi
  done
done

echo ""
echo "---"
echo "Results: $PASS passed, $FAIL failed"
[ "$FAIL" -eq 0 ] || exit 1
