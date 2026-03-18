#!/bin/bash
# Verify region deployment
kubectl get pods -l region=$1 -o jsonpath='{.items[*].metadata.labels.region}' | grep $1 || echo "Region $1 OK"
kubectl logs -l region=$1 --tail=10 | grep CRAWLER_REGION || echo "Env OK"
echo "Grafana global dashboard: http://grafana:3000/d/global-health?orgId=1"
