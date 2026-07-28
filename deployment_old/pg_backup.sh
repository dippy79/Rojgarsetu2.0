#!/bin/bash
# Daily Postgres Backup Script - RojgarSetu
# Cron: 2 0 * * * /opt/rojgarsetu/deployment/pg_backup.sh

set -e

BACKUP_DIR="/backups/rojgarsetu"
DATE=$(date +%Y%m%d_%H%M%S)
DB_NAME="rojgardb"

cd "$(dirname "$0")"/..

# Load .env.production
export $(grep -v &#x27;^#&#x27; .env.production | xargs)

mkdir -p ${BACKUP_DIR}

echo "🗄️  Backup ${DB_NAME} to ${BACKUP_DIR}/${DATE}.sql.gz"

# pg_dump (docker exec or local psql)
docker compose exec -T postgres pg_dump -U rojgar ${DB_NAME} | gzip > ${BACKUP_DIR}/${DATE}.sql.gz

# Verify size
if [ ! -s ${BACKUP_DIR}/${DATE}.sql.gz ]; then
  echo "❌ Backup empty!"
  exit 1
fi

echo "✅ Backup complete: $(du -h ${BACKUP_DIR}/${DATE}.sql.gz)"

# Keep last 7 days
find ${BACKUP_DIR} -name "*.sql.gz" -mtime +7 -delete

echo "🧹 Old backups cleaned"

ls -la ${BACKUP_DIR} | head -5

# Optional: Upload to S3
# aws s3 cp ${BACKUP_DIR}/${DATE}.sql.gz s3://your-backup-bucket/ --storage-class STANDARD_IA

