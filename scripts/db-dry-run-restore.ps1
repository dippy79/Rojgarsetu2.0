# PHASE 27 Dry-Run Restore Verification (Idempotent)
# Restores latest backup to temp container, queries row count, destroys.

param (
  [string]$BackupDir = "f:/Rojgarsetu2.0/rojgarsetu2/backups"
)

$latestBackup = Get-ChildItem "$BackupDir" -Filter "*.sql" | Sort-Object LastWriteTime -Descending | Select-Object -First 1
if (-not $latestBackup) { Write-Error "No backup found"; exit 1 }

$tempDbName = "restore-test-$(Get-Date -Format 'yyyyMMddHHmmss')"

docker run --name temp-restore -e POSTGRES_PASSWORD=temp123 -p 5433:5432 -d postgres:15-alpine postgres

Start-Sleep 10

docker cp $latestBackup.FullName temp-restore:/backup.sql

docker exec temp-restore psql -U postgres -f /backup.sql

$govCount = docker exec temp-restore psql -U postgres -t -c "SELECT COUNT(*) FROM gov_jobs;"
$privCount = docker exec temp-restore psql -U postgres -t -c "SELECT COUNT(*) FROM priv_jobs;"

docker stop temp-restore && docker rm temp-restore

if ([int]$govCount -gt 0 -and [int]$privCount -gt 0) {
  Write-Host "Restore verified: GovJobs=$govCount PrivJobs=$privCount"
} else {
  Write-Error "Restore failed: Low row counts"
}

Write-Host "Dry-run complete (<30min)"

