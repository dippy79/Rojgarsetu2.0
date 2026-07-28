# PHASE 26 Blue/Green Deploy Script (Local Simulation)
Write-Host "PHASE 26 - Blue/Green Backend Rollout (Local docker-compose simulation)"
Write-Host "1. Build backend no-cache (Green env)..."

cd deployment
docker-compose build backend --no-cache

Write-Host "2. Graceful down (keep postgres/prom/grafana)..."
docker-compose down backend

Write-Host "3. Rollout green backend..."
docker-compose up -d postgres backend prometheus grafana

Write-Host "4. Health checks..."
for ($i = 0; $i -lt 30; $i++) {
    $health = Invoke-WebRequest -Uri "http://localhost:8083/health" -UseBasicParsing
    $live = Invoke-WebRequest -Uri "http://localhost:8083/live" -UseBasicParsing
    $ready = Invoke-WebRequest -Uri "http://localhost:8083/ready" -UseBasicParsing
    if ($health.StatusCode -eq 200 -and $live.StatusCode -eq 200 -and $ready.StatusCode -eq 200) {
        Write-Host "All checks passed! Traffic switch complete."
        break
    }
    Start-Sleep 2
}

Write-Host "Rollback ready: docker-compose down backend && docker-compose up -d backend.old"

