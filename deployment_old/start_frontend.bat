@echo off
REM RojgarSetu 2.0 - Quick Frontend + API Start (Fixes Connection Refused)
REM Run after starting Docker Desktop

echo [1/6] Checking Docker daemon...
docker ps >nul 2>&1
if errorlevel 1 (
  echo ERROR: Docker daemon not running. Open Docker Desktop first.
  pause
  exit /b 1
)

echo [2/6] Setup DB if needed...
powershell -ExecutionPolicy Bypass -File setup_local_db.ps1

echo [3/6] Starting minimal services: postgres redis api-gateway frontend...
docker compose up -d postgres redis api-gateway frontend

echo [4/6] Building frontend with fixed nginx.conf...
docker compose build --no-cache frontend

echo [5/6] Restarting frontend...
docker compose restart frontend

echo [6/6] Checking ports...
netstat -ano | findstr :3001
netstat -ano | findstr :3000
echo.
echo SUCCESS: Frontend http://localhost:3001 ^| API http://localhost:3000
echo Test in browser console for no iframe errors.
echo Logs: docker compose logs -f frontend api-gateway
echo Stop: docker compose down
pause
