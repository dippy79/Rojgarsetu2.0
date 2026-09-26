$TIMESTAMP = Get-Date -Format "yyyyMMdd_HHmmss"
$BACKUP_DIR = "./backups"
$POSTGRES_USER = if ($env:POSTGRES_USER) { $env:POSTGRES_USER } else { "postgres" }
$POSTGRES_DB = if ($env:POSTGRES_DB) { $env:POSTGRES_DB } else { "rojgarsetu2" }
$CONTAINER_NAME = "rojgar-postgres"

if (-not (Test-Path $BACKUP_DIR)) {
    New-Item -ItemType Directory -Path $BACKUP_DIR
}

Write-Output "Starting backup at $TIMESTAMP"
docker exec $CONTAINER_NAME pg_dump -U $POSTGRES_USER $POSTGRES_DB | Out-File -FilePath "$BACKUP_DIR/rojgarsetu_$TIMESTAMP.sql"

if (Test-Path "$BACKUP_DIR/rojgarsetu_$TIMESTAMP.sql") {
    Write-Output "Backup completed successfully: rojgarsetu_$TIMESTAMP.sql"
} else {
    Write-Error "Backup failed!"
}
