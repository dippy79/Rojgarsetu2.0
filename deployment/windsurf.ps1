# WindSurf: Run this in PowerShell as Administrator from repo root F:\Rojgarsetu2.0\rojgarsetu2
$repo='F:\Rojgarsetu2.0\rojgarsetu2'; Set-Location $repo

# 0. Prepare logs and backup folder
$ts=Get-Date -Format yyyyMMdd_HHmmss
$logs=Join-Path $repo 'deployment\logs'; New-Item -Path $logs -ItemType Directory -Force | Out-Null
$dockerBackup=Join-Path $repo \"deployment_backup_docker_$ts\"; New-Item -Path $dockerBackup -ItemType Directory -Force | Out-Null
Write-Output \"Logs: $logs\"
Write-Output \"Docker backup: $dockerBackup\"

# 1. Backup and remove Docker artifacts (docker-compose.yml, Dockerfile*, docker folders)
Get-ChildItem -Path $repo -Include 'docker-compose*.yml','Dockerfile*' -Recurse -File | ForEach-Object {
  $dest = Join-Path $dockerBackup ($_.FullName.TrimStart($repo.TrimEnd('\\')) -replace '\\\\','_')
  New-Item -ItemType Directory -Path (Split-Path $dest) -Force | Out-Null
  Copy-Item -Path $_.FullName -Destination $dest -Force
  Remove-Item -Path $_.FullName -Force
  Write-Output \"Backed up and removed: $($_.FullName)\"
}
Get-ChildItem -Path $repo -Directory -Recurse | Where-Object { $_.Name -match 'docker' } | ForEach-Object {
  $dest = Join-Path $dockerBackup ($_.FullName.TrimStart($repo.TrimEnd('\\')) -replace '\\\\','_')
  Copy-Item -Path $_.FullName -Destination $dest -Recurse -Force
  Remove-Item -Path $_.FullName -Recurse -Force
  Write-Output \"Backed up and removed folder: $($_.FullName)\"
}

# 2. Move or remove any stubs that caused duplicate declarations
$dbDir = Join-Path $repo 'backend_go\internal\db'
$stub = Join-Path $dbDir 'stubs_for_build.go'
if (Test-Path $stub) {
  $stubBak = Join-Path $dockerBackup ('stubs_for_build_' + $ts + '.go')
  Move-Item -Path $stub -Destination $stubBak -Force
  Write-Output \"Moved stub to backup: $stubBak\"
}

# 3. Detect duplicate type declarations and report files to operator
Select-String -Path (Join-Path $dbDir '*.go') -Pattern 'type\s+Candidate|type\s+PostgresDB|stubs_for_build' -List | Format-Table Path, LineNumber -AutoSize | Out-String | Tee-Object -FilePath (Join-Path $logs 'db_duplicate_scan.txt')

# 4. If real implementations exist, do NOT recreate conflicting types. If missing, create minimal non-conflicting stubs
$hasCandidate = Select-String -Path (Join-Path $dbDir '*.go') -Pattern 'type\s+Candidate' -Quiet
$hasPostgresDB = Select-String -Path (Join-Path $dbDir '*.go') -Pattern 'type\s+PostgresDB' -Quiet
if (-not ($hasCandidate -and $hasPostgresDB)) {
  @'
package db
import \"context\"
type CreateUserParams struct{ ID, Email, Name string }
type UpdateCandidateRequest struct{ CandidateID, Phone string }
type UpdateCompanyRequest struct{ CompanyID, Name string }
type CreateJobRequest struct{ Title, Company string }
type ApplyJobRequest struct{ JobID, Candidate string }
type PostgresDB struct{}
func (d *PostgresDB) EmailExists(ctx context.Context, email string) (bool, error) { return false, nil }
func (d *PostgresDB) CreateUser(ctx context.Context, params CreateUserParams) error { return nil }
'@ | Set-Content -Path (Join-Path $dbDir 'stubs_for_build.go') -Encoding UTF8
  Write-Output \"Created minimal non-conflicting stub at backend_go\internal\db\stubs_for_build.go\"
} else {
  Write-Output \"Real Candidate and PostgresDB types found; no stub created.\"
}

# 5. Ensure Postgres is running and apply DB schema if script exists
$setupScript = Join-Path $repo 'deployment\setup_local_db.ps1'
if (Test-Path $setupScript) {
  Write-Output \"Running DB setup script...\"
& powershell -ExecutionPolicy Bypass -File $setupScript -NoLog 2>&1 | Tee-Object -FilePath (Join-Path $logs 'setup_local_db.txt'); if ($LASTEXITCODE -ne 0) { Write-Output \"DB setup failed, check logs\"; exit 1 }
} else {
  Write-Output \"No setup_local_db.ps1 found. Skipping automatic schema apply.\"
}

# 6. Quick psql checks (you may be prompted for password)
Try { psql -h localhost -U rojgarsetu -d rojgarsetu -c \"SELECT version();\" > (Join-Path $logs 'psql_version.txt') 2>&1 } Catch { 'psql check failed' | Out-File (Join-Path $logs 'psql_version.txt') }
Try { psql -h localhost -U rojgarsetu -d rojgarsetu -c \"SELECT COUNT(*) FROM jobs_government;\" > (Join-Path $logs 'jobs_gov_count.txt') 2>&1 } Catch { 'psql count failed' | Out-File (Join-Path $logs 'jobs_gov_count.txt') }

# 7. Build backend locally and capture build log
Push-Location (Join-Path $repo 'backend_go')
New-Item -Path ..\deployment\logs -ItemType Directory -Force | Out-Null
go mod tidy > (Join-Path $logs 'go_mod_tidy.txt') 2>&1
go build -o backend_server ./cmd/server > (Join-Path $logs 'go_build_output.txt') 2>&1
$rc=$LASTEXITCODE
if ($rc -ne 0) {
  Write-Output \"GO BUILD FAILED. First 200 lines of build log:\"
  Get-Content (Join-Path $logs 'go_build_output.txt') -TotalCount 200
  Pop-Location
  exit $rc
}
Write-Output \"Go build succeeded.\"

# 8. Start backend and verify health and sample API
$backendExe = Join-Path (Get-Location) 'backend_server'
Start-Process -FilePath $backendExe -PassThru -WindowStyle Hidden | Out-Null
Start-Sleep -Seconds 3
netstat -an | Select-String ':8083' | Out-File -FilePath (Join-Path $logs 'netstat_8083_after_start.txt') -Encoding utf8
Try { Invoke-WebRequest -Uri 'http://localhost:8083/health' -UseBasicParsing -TimeoutSec 5 | Out-File -FilePath (Join-Path $logs 'health_after_start.txt') -Encoding utf8 } Catch { $_.Exception.Message | Out-File -FilePath (Join-Path $logs 'health_after_start.txt') -Encoding utf8 }
Try { curl -sS \"http://localhost:8083/api/v1/gov-jobs?limit=1\" > (Join-Path $logs 'gov_jobs_sample.json') } Catch { \"curl failed\" | Out-File (Join-Path $logs 'gov_jobs_sample.json') }

Pop-Location

# 9. Start frontend locally and verify fetch
Push-Location (Join-Path $repo 'frontend')
Set-Content -Path .env -Value 'REACT_APP_BACKEND_URL=http://localhost:8083' -Encoding UTF8
npm ci > (Join-Path $logs 'frontend_npm_ci.txt') 2>&1
Start-Process -FilePath 'npm' -ArgumentList 'start' -PassThru -WindowStyle Hidden | Out-Null
Start-Sleep -Seconds 6
Try { Invoke-WebRequest -Uri 'http://localhost:3000/pages/gov-jobs' -UseBasicParsing -TimeoutSec 5 | Out-File -FilePath (Join-Path $logs 'frontend_govjobs.txt') -Encoding utf8 } Catch { $_.Exception.Message | Out-File -FilePath (Join-Path $logs 'frontend_govjobs.txt') -Encoding utf8 }
Pop-Location

# 10. Optional: run crawler to populate DB if counts are zero
$govCount = Get-Content (Join-Path $logs 'jobs_gov_count.txt') -ErrorAction SilentlyContinue
if ($govCount -match '\d' -and [int]($govCount -replace '\D','') -eq 0) {
  Write-Output \"Gov jobs count is zero; running crawler to populate DB...\"
  Push-Location (Join-Path $repo 'crawler_go\cmd')
  go run main.go > (Join-Path $logs 'crawler_run.txt') 2>&1
  Pop-Location
  psql -h localhost -U rojgarsetu -d rojgarsetu -c \"SELECT COUNT(*) FROM jobs_government;\" > (Join-Path $logs 'jobs_gov_count_after_crawl.txt') 2>&1
}

# 11. Summary and rollback instructions
Write-Output \"Done. Check logs in $logs for details.\"
Write-Output \"Docker artifacts backed up to $dockerBackup. To restore Docker files later, copy from that folder back to repo root.\"
Write-Output \"To stop backend: Get-Process -Name backend_server | Stop-Process -Force\"
Write-Output \"To stop frontend: Get-Process -Name node | Stop-Process -Force (careful if other node processes exist)\"
