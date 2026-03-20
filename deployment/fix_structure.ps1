# RojgarSetu 2.0 - Fix Directory Structure for Docker Compose
# Run from deployment/ directory: .\fix_structure.ps1

Write-Host "=== RojgarSetu 2.0 Docker Structure Fix ===" -ForegroundColor Green
Write-Host "Current dir: $(Get-Location)" -ForegroundColor Yellow

# Remove existing services/ if any
if (Test-Path .\services) {
    Write-Host "Removing existing .\services\" -ForegroundColor Yellow
    Remove-Item .\services -Recurse -Force
}

# Move services from root to deployment/services/
Write-Host "Moving ..\services\ -> .\services\" -ForegroundColor Green
robocopy ..\services .\services /E /MOV /R:0 /W:5 /NFL /NDL | Out-Null
if ($LASTEXITCODE -gt 7) {
    Write-Error "Robocopy failed with code $LASTEXITCODE"
    exit 1
}

# Copy .env from root to deployment/.env
Write-Host "Copying ..\..\.env -> .\.env" -ForegroundColor Green
Copy-Item ..\..\.env .\.env -Force -ErrorAction SilentlyContinue

# Append/ensure critical vars to .env (backup first)
if (Test-Path .\.env.bak) { Remove-Item .\.env.bak }
Copy-Item .\.env .\.env.bak -Force

Add-Content -Path .\.env -Value "`n# Critical Docker vars for connection fix"
Add-Content -Path .\.env -Value "DATABASE_URL=postgres://postgres:postgres@postgres:5432/rojgarsetu?sslmode=disable"
Add-Content -Path .\.env -Value "REDIS_URL=redis://redis:6379"
Add-Content -Path .\.env -Value "POSTGRES_PASSWORD=postgres"
Add-Content -Path .\.env -Value "JWT_SECRET=your-super-secret-jwt-key-change-me"
Add-Content -Path .\.env -Value "YOUTUBE_API_KEY=your-key-optional"
Add-Content -Path .\.env -Value "GF_SECURITY_ADMIN_PASSWORD=admin"

Write-Host "=== Structure fixed successfully! ===" -ForegroundColor Green
Write-Host "Next: Update TODO.md [1] to [x], run docker compose config to validate." -ForegroundColor Cyan
Write-Host "Verify: ls services/, cat .env | head" -ForegroundColor Cyan
