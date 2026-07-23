# Non-interactive backend_go cleanup + rebuild + deploy script
# Run from repo root: powershell -ExecutionPolicy Bypass -File run_deploy_backend.ps1
# Or paste directly into PowerShell

$root='F:\Rojgarsetu2.0\rojgarsetu2\services\backend-go'
Set-Location $root
$ts=Get-Date -Format 'yyyyMMdd_HHmmss'
$backup=Join-Path $root "deployment\backup_temp_stubs_$ts"
New-Item -Path $backup -ItemType Directory -Force | Out-Null
Get-ChildItem -Path (Join-Path $root 'internal\db') -Filter '*temp*.go','*extra*.go','*stubs*.go' -ErrorAction SilentlyContinue |
  ForEach-Object { Move-Item -Path $_.FullName -Destination $backup -Force -ErrorAction SilentlyContinue }
Write-Output "Moved temp stubs (if any) to: $backup"

# Ensure logs dir exists
New-Item -Path (Join-Path $root 'deployment\logs') -ItemType Directory -Force | Out-Null

# Clean Go caches
Write-Output "Cleaning Go caches..."
go clean -modcache 2>&1 | Out-File -FilePath (Join-Path $root 'deployment\logs\go_clean_modcache.txt') -Encoding UTF8
go clean -cache 2>&1 | Out-File -FilePath (Join-Path $root 'deployment\logs\go_clean_cache.txt') -Encoding UTF8

# Docker prune (noninteractive)
Write-Output "Pruning Docker..."
docker system prune --all --volumes --force 2>&1 | Out-File -FilePath (Join-Path $root 'deployment\logs\docker_prune.txt') -Encoding UTF8

# Go tidy and build
Write-Output "Running go mod tidy and build..."
go mod tidy 2>&1 | Out-File -FilePath (Join-Path $root 'deployment\logs\go_mod_tidy.txt') -Encoding UTF8
go build ./cmd/server 2>&1 | Tee-Object (Join-Path $root 'deployment\logs\go_build_after_cleanup.txt')
if ($LASTEXITCODE -ne 0) { 
  Write-Output 'go build failed; see deployment\logs\go_build_after_cleanup.txt'; 
  exit 1 
}

# Docker build
Write-Output "Building Docker image..."
Set-Location (Split-Path $root -Parent)
docker build --progress=plain -f .\deployment\Dockerfile.backend -t backend-test ..\services\backend-go > .\deployment\logs\docker_build_plain.txt 2>&1
if ($LASTEXITCODE -ne 0) { 
  Write-Output 'docker build failed; see deployment\logs\docker_build_plain.txt'; 
  exit 2 
}

# Run container (replace existing)
Write-Output "Starting container..."
docker rm -f backend_test_run -ErrorAction SilentlyContinue | Out-Null
docker run -d --name backend_test_run -p 8083:8083 backend-test | Out-Null
Start-Sleep -Seconds 3
docker logs backend_test_run --tail 200 | Out-File -FilePath .\deployment\logs\docker_run_logs.txt -Encoding UTF8

# Health check
Write-Output "Health check..."
try {
  Invoke-WebRequest -Uri http://localhost:8083/health -UseBasicParsing -TimeoutSec 5 2>&1 | Out-File -FilePath .\deployment\logs\health_check.txt -Encoding UTF8
  Write-Output "Health check passed"
} catch {
  'health check failed or timed out' | Out-File -FilePath .\deployment\logs\health_check.txt -Encoding UTF8
  Write-Output "Health check failed - check logs"
}

Write-Output 'Completed noninteractive cleanup+rebuild. Check logs:'
Write-Output '  deployment\logs\go_clean_modcache.txt'
Write-Output '  deployment\logs\go_clean_cache.txt'
Write-Output '  deployment\logs\docker_prune.txt'
Write-Output '  deployment\logs\go_mod_tidy.txt'
Write-Output '  deployment\logs\go_build_after_cleanup.txt'
Write-Output '  deployment\logs\docker_build_plain.txt'
Write-Output '  deployment\logs\docker_run_logs.txt'
Write-Output '  deployment\logs\health_check.txt'

Write-Output "Backend deployed on port 8083. Test: http://localhost:8083/health"
