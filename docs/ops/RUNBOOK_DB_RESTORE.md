# RUNBOOK: DB Restore from Corruption

## Symptoms
- pg_up=0 Prometheus alert
- Query timeouts or empty gov_jobs

## Steps
1. **Pause Services**
   ```
   docker-compose scale backend=0 api-gateway=0 crawler-service=0
   ```

2. **Identify Backup**
   rclone ls s3-backup:rojgarsetu-backups/

3. **Stop Primary**
   docker-compose stop postgres

4. **Restore**
   .\scripts\db-dry-run-restore.ps1 -Backup latest.sql  # Verify first

5. **Full Restore**
   docker exec -it rojgar-postgres pg_restore -U rojgarsetu -d rojgarsetu /backup.sql --clean

6. **Verify Data**
   SELECT COUNT(*) FROM gov_jobs;  # Matches expected

7. **Resume**
   docker-compose up -d postgres backend api-gateway crawler-service

**PITR**: pgbackrest restore --type=time --target="2026-03-13 12:00"

**Escalation**: If WAL >24h lost, declare data loss, trigger insurance.

