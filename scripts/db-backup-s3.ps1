# PHASE 27 DB Backup to S3 (WAL-G)
# Assumes rclone setup with 's3-backup' remote (no hardcoded bucket)

param (
  [string]$BackupDir = "f:/Rojgarsetu2.0/rojgarsetu2/backups",
  [string]$RemoteName = "s3-backup",
  [string]$Region = "us"
)

$BucketSuffix = switch ($Region) {
  "us" { "us-east-1" }
  "eu" { "eu-west-1" }
  "apac" { "ap-south-1" }
  default { "default" }
}


cd deployment

# Full logical backup
docker-compose exec -T postgres pg_dump -U rojgarsetu rojgarsetu > "$BackupDir/full-$(Get-Date -Format 'yyyyMMdd-HHmmss').sql"

# WAL archive (continuous PITR)
docker-compose exec postgres pgbackrest --stanza=main backup --type=incr

# Copy to S3 via rclone
rclone sync "$BackupDir" "$RemoteName":rojgarsetu-backups/ --progress

Write-Host "Backup complete. Dry-run restore: .\scripts\db-dry-run-restore.ps1"

