#!/usr/bin/env bash
# Initialize MongoDB replica set on prod (3-member cluster)
set -euo pipefail

NS="${1:-sravas}"

echo "Waiting for all 3 MongoDB pods to be ready..."
kubectl -n "$NS" rollout status statefulset/mongodb --timeout=180s

echo "Initializing replica set on mongodb-0..."
kubectl -n "$NS" exec mongodb-0 -- mongosh --quiet --eval '
rs.initiate({
  _id: "rs0",
  members: [
    { _id: 0, host: "mongodb-0.mongodb-headless.sravas.svc.cluster.local:27017" }
  ]
})'

echo "Adding mongodb-1 to replica set..."
kubectl -n "$NS" exec mongodb-0 -- mongosh --quiet --eval \
  'rs.add("mongodb-1.mongodb-headless.sravas.svc.cluster.local:27017")'

echo "Adding mongodb-2 to replica set..."
kubectl -n "$NS" exec mongodb-0 -- mongosh --quiet --eval \
  'rs.add("mongodb-2.mongodb-headless.sravas.svc.cluster.local:27017")'

echo "Replica set status:"
kubectl -n "$NS" exec mongodb-0 -- mongosh --quiet --eval \
  'rs.status().members.forEach(m => print(m.name + " → " + m.stateStr))'
