# RojgarSetu 2.0 - Docker Compose Verification Script
# Run from deployment/: .\verify_compose.ps1

Write-Host "=== Docker Compose Verification ===" -ForegroundColor Green

# Stop and clean
Write-Host "`n1. Cleaning old containers/volumes..." -ForegroundColor Cyan
docker compose down --volumes --remove-orphans --timeout 30 2>&1 | Tee-Object -FilePath logs/verify_cleanup.log

# Rebuild
Write-Host "`n2. Building (no-cache)..." -ForegroundColor Cyan
docker compose build --no-cache --progress=plain 2>&1 | Tee-Object -FilePath logs/verify_build.log

# Start
Write-Host "`n3. Starting services..." -ForegroundColor Cyan
docker compose up -d 2>&1 | Tee-Object -FilePath logs/verify_up.log

# Wait and status
Start-Sleep 10
Write-Host "`n4. Container status:" -ForegroundColor Cyan
docker compose ps | Out-Host

# Logs
Write-Host "`n5. Backend logs (tail 50):" -ForegroundColor Cyan
docker compose logs --tail=50 backend 2>&1 | Out-Host

# Health checks
Write-Host "`n6. Health checks:" -ForegroundColor Cyan
$healthEndpoints = @(
    "http://localhost:8083/health",  # backend
    "http://localhost:3000/health",  # api-gateway
    "http://localhost:8081/actuator/health",  # auth
    "http://localhost:8082/health"   # crawler
)

foreach ($url in $healthEndpoints) {
    try {
        $response = Invoke-WebRequest -Uri $url -Method Get -UseBasicParsing -TimeoutSec 5
        Write-Host "$url : OK ($($response.StatusCode))" -ForegroundColor Green
    } catch {
        Write-Host "$url : FAILED ($($_.Exception.Message.Split(':')[0]))" -ForegroundColor Red
    }
}

Write-Host "`n=== Verification complete! Check logs/verify_*.log ===" -ForegroundColor Green
Write-Host "Frontend: http://localhost:3001" -ForegroundColor Cyan
Write-Host "Grafana: http://localhost:3002 (admin/admin)" -ForegroundColor Cyan
Write-Host "Prometheus: http://localhost:9090" -ForegroundColor Cyan
