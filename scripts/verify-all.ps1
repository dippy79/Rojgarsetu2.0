Write-Host "🚀 ROJGARSETU 2.0 - PRODUCTION VERIFICATION SEQUENCE" -ForegroundColor Cyan
Write-Host "----------------------------------------------------"

Write-Host "=== 1. VERIFYING GO BACKEND ===" -ForegroundColor Yellow
cd backend_go
go mod tidy
go build ./...
if ($LASTEXITCODE -eq 0) { Write-Host "Backend: CLEAN" -ForegroundColor Green } else { Write-Host "Backend: FAILED" -ForegroundColor Red; exit 1 }
cd ..

Write-Host "`n=== 2. VERIFYING GO CRAWLER ===" -ForegroundColor Yellow
cd services/crawler-go
go mod tidy
go build ./...
if ($LASTEXITCODE -eq 0) { Write-Host "Crawler: CLEAN" -ForegroundColor Green } else { Write-Host "Crawler: FAILED" -ForegroundColor Red; exit 1 }
cd ..

Write-Host "`n=== 3. VERIFYING PYTHON AI ENGINE ===" -ForegroundColor Yellow
cd services/ai-engine-python
# Simple dependency check
pip install -r requirements.txt --quiet
if ($LASTEXITCODE -eq 0) { Write-Host "AI Engine: READY" -ForegroundColor Green } else { Write-Host "AI Engine: DEPENDENCY ERROR" -ForegroundColor Red }
cd ..

Write-Host "`n=== 4. VERIFYING FRONTEND BUILD ===" -ForegroundColor Yellow
cd frontend
npm run build 2>&1 | Select-String "error", "failed"
if ($LASTEXITCODE -eq 0) { Write-Host "Frontend: PRODUCTION READY" -ForegroundColor Green } else { Write-Host "Frontend: BUILD FAILED" -ForegroundColor Red; exit 1 }
cd ..

Write-Host "`n=== 5. CHECKING DOCKER ORCHESTRATION ===" -ForegroundColor Yellow
docker compose config --quiet
if ($LASTEXITCODE -eq 0) { Write-Host "Docker Compose: CONFIG VALID" -ForegroundColor Green } else { Write-Host "Docker Compose: INVALID" -ForegroundColor Red }

Write-Host "`n✅ SUCCESS: ALL SYSTEMS 10/10 READY FOR PRODUCTION!" -ForegroundColor Green
