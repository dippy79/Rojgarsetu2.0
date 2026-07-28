# Local services launcher for RojgarSetu 2.0 - No Docker
# Usage: powershell deployment/run_local_services.ps1
# Starts backend Go, with instructions for others.
# Logs to logs/<service>_run.txt

$logDir = 'logs'
New-Item -Path $logDir -ItemType Directory -Force | Out-Null

Write-Output \"Starting local backend...\"

# Backend Go
Set-Location backend_go
go clean -cache -modcache | Tee-Object -FilePath \"../deployment/$logDir/backend_prep.txt\"
go mod tidy | Tee-Object -FilePath \"../deployment/$logDir/backend_prep.txt\" -Append
go build -o backend_server ./cmd/server | Tee-Object -FilePath \"../deployment/$logDir/backend_build.txt\"
if ($LASTEXITCODE -eq 0) {
  Start-Process powershell -ArgumentList '-NoExit', '-Command', \"& './backend_server'\" -WorkingDirectory (Get-Location)
  Write-Output \"Backend running on :8083. Logs: deployment/logs/backend_*.txt\"
} else {
  Write-Error \"Backend build failed. See logs.\"
}

# Other services - manual start instructions
Write-Output \"\`
Manual start for other services:

1. Auth Java (port 8081):
 cd auth_service_java
 mvn clean spring-boot:run -Dspring.datasource.url='jdbc:postgresql://localhost:5432/rojgarsetu' -Dspring.datasource.username=rojgarsetu -Dspring.datasource.password=rojgarsetu_secret | tee ../deployment/logs/auth_run.txt

2. API Gateway Node (port 3000):
 cd api_gateway_node
 npm ci
 npm start | tee ../deployment/logs/gateway_run.txt

3. Frontend (port 3000 or 3001):
 cd frontend
 npm ci
 npm start | tee ../deployment/logs/frontend_run.txt

4. Test connectivity:
 curl http://localhost:8083/health
 open http://localhost:3001 or npm run dev

Redis optional: redis-server (port 6379)
\`\"

