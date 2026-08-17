#!/bin/sh
# Automated PostgreSQL Database Backup Script
# Keeps backups for 7 days

TIMESTAMP=$(date +%Y%m%d_%H%M%S)
BACKUP_DIR="/backups"
POSTGRES_USER="amitsharma"
POSTGRES_DB="rojgarsetu2"
CONTAINER_NAME="rojgar-postgres"

# Create backup directory if it doesn't exist
mkdir -p $BACKUP_DIR

# Perform backup
echo "Starting backup at $TIMESTAMP"
docker exec $CONTAINER_NAME pg_dump -U $POSTGRES_USER $POSTGRES_DB | gzip > $BACKUP_DIR/rojgarsetu_$TIMESTAMP.sql.gz

# Verify backup was created
if [ -f "$BACKUP_DIR/rojgarsetu_$TIMESTAMP.sql.gz" ]; then
    echo "Backup completed successfully: rojgarsetu_$TIMESTAMP.sql.gz"
else
    echo "Backup failed!"
    exit 1
fi

# Remove backups older than 7 days
find $BACKUP_DIR -name "rojgarsetu_*.sql.gz" -mtime +7 -delete
echo "Old backups (older than 7 days) removed"

# List current backups
echo "Current backups:"
ls -lh $BACKUP_DIR/rojgarsetu_*.sql.gz