# PHASE 27 — POST-ROLLOUT OPS & DISASTER RECOVERY
## Plan Steps

### 1. Database Resilience
- Create scripts/db-backup-s3.ps1 (WAL-G S3 PITR).
- Create scripts/db-dry-run-restore.ps1 (temp container verify).

### 2. Data Staleness Metrics
- Edit crawler_go/... : Add `crawler_last_success_timestamp` gauge per source.
- Edit monitoring/prometheus/rules.yml (staleness alerts).

### 3. K8s Manifests
- Edit deployment/k8s/backend.yaml (HPA, ConfigMap, Ingress, PDB).
- Resource limits from k6 (CPU 500m, Mem 512Mi).

### 4. Centralized Logging
- Edit crawler_go/internal/logger/logger.go (source_domain, correlation_id).
- Create docker-compose.monitoring.yml (Loki/Promtail).

### 5. Runbooks
- Create docs/ops/RUNBOOK_CRAWLER_BLOCK.md
- Create docs/ops/RUNBOOK_DB_RESTORE.md

### 6. Grafana Ops Dashboard
- Create monitoring/grafana/provisioning/dashboards/ops-overview.json

**Progress: 0/6 complete**
