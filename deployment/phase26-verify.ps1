# PHASE 26 Verification Script
cd deployment

Write-Host "Health Checks:"
curl -k https://localhost:8443/health
curl -k https://localhost:8443/live
curl -k https://localhost:8443/ready
curl http://localhost:8083/metrics | Select-String "http_requests_total"

Write-Host "Logs:"
docker-compose logs backend | Select-String "error|tls|ready"

Write-Host "Load Test:"
k6 run loadtest/k6.js

Write-Host "Prometheus Targets: http://localhost:9090/targets"
Write-Host "Grafana: http://localhost:3002 (admin/StrongGrafana@2026)"
Write-Host "Done - check alerts!"

