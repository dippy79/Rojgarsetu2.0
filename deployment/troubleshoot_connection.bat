@echo off
REM Troubleshoot ERR_CONNECTION_REFUSED localhost:3001

echo [1] Check Docker:
docker ps

echo [2] Check ports 3001/3000:
netstat -ano | findstr :3001
netstat -ano | findstr :3000

echo [3] Check Windows Firewall:
netsh advfirewall firewall show rule name=all | findstr Docker

echo [4] Check proxy:
powershell -c "netsh winhttp show proxy"

echo [5] Kill conflicting processes if any:
REM Manual: taskkill /PID ^<pid^> /F for port processes

echo [6] Docker Desktop status:
tasklist | findstr Docker

pause
echo Run: Start Docker Desktop ^> cd deployment ^> start_frontend.bat
