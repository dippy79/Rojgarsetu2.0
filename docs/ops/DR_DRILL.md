# Disaster Recovery Drills

## Quarterly Schedule
- Q1: 2026-03-31 10AM
- Q2: 2026-06-30 10AM
- Q3: 2026-09-30 10AM
- Q4: 2026-12-31 10AM

## Drill Procedure (Staging)
1. **Preparation:** Tag staging DB, note rowcounts/schemas
2. **Simulate Failure:** `docker-compose down postgres`
3. **Execute Restore:** `.\scripts\db-dry-run-restore.ps1 -Backup latest-staging.sql`
4. **Measure:** Time from failure to query success (<30min target)
5. **Verify:** Rowcounts match ±5%, schema identical
6. **Document:** Update table below
7. **Cleanup:** Destroy temp containers

## Results Log

| Date | Backup Size | Restore Time | GovJobs Count | PrivJobs Count | Status | Notes |
|------|-------------|--------------|---------------|----------------|--------|-------|
| YYYY-MM-DD | 150MB | 12min | 125,430 | 89,120 | PASS | WAL intact |
|      |     |      |               |                |        |       |

## Escalation
- RTO >30min → Review WAL-G retention
- Data loss >1% → Insurance claim process
