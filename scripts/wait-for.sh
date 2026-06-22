#!/bin/sh
# wait-for.sh — block until TCP host:port(s) are reachable, then exec command

TIMEOUT="${WAIT_FOR_TIMEOUT:-60}"
HOSTS=""

while [ $# -gt 0 ]; do
  case "$1" in
    --) shift; break ;;
    *) HOSTS="$HOSTS $1" ;;
  esac
  shift
done

for hostport in $HOSTS; do
  host="${hostport%:*}"
  port="${hostport##*:}"
  elapsed=0
  while ! nc -z "$host" "$port" 2>/dev/null; do
    if [ "$elapsed" -ge "$TIMEOUT" ]; then
      echo "Timeout waiting for $host:$port"
      exit 1
    fi
    sleep 1
    elapsed=$((elapsed + 1))
  done
  echo "$host:$port is available"
done

exec "$@"
