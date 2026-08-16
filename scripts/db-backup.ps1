# Automated PostgreSQL Database Backup Script for Windows
# Keeps backups for 7 days

$TIMESTAMP = Get-Date -Format "yyyyMMdd_HHmmss"
$BACKUP_DIR = ".\backups"
$POSTGRES_USER = "amitsharma"
$POSTGRES_DB = "rojgarsetu2"
$CONTAINER_NAME = "rojgar-postgres"

# Create backup directory if it doesn't exist
if (-not (Test-Path $BACKUP_DIR)) {
    New-Item -ItemType Directory -Path $BACKUP_DIR
}

# Perform backup
Write-Host "Starting backup at $TIMESTAMP"
$backupFile = "$BACKUP_DIR\rojgarsetu_$TIMESTAMP.sql.gz"

docker exec $CONTAINER_NAME pg_dump -U $POSTGRES_USER $POSTGRES_DB | docker exec -i $CONTAINER_NAME gzip > $backupFile

# Verify backup was created
if (Test-Path $backupFile) {
    Write-Host "Backup completed successfully: rojgarsetu_$TIMESTAMP.sql.gz"
} else {
    Write-Host "Backup failed!"
    exit 1
}

# Remove backups older than 7 days
$cutoffDate = (Get-Date).AddDays(-7)
Get-ChildItem $BACKUP_DIR -Filter "rojgarsetu_*.sql.gz" | Where-Object { $_.LastWriteTime -lt $cutoffDate } | Remove-Item
Write-Host "Old backups (older than 7 days) removed"

# List current backups
Write-Host "Current backups:"
Get-ChildItem $BACKUP_DIR -Filter "rojgarsetu_*.sql.gz" | Format-Table Name, Length, LastWriteTime