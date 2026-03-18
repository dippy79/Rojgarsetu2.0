# PHASE 26 — PRODUCTION ROLLOUT CHECKLIST
## Approved Plan Steps (Breakdown)

### 1. Create Base Source for Crawler Safeguards
- Create `crawler_go/internal/sources/base.go` (robots.txt, 403/429 pause/alert metric)

### 2. Update Crawler Sources (UA Rotation + Throttle)
- Edit all `crawler_go/internal/sources/*.go` (replace static UA → rotator, add domain limiter)

### 3. Enhance Retry Logic
- Edit `crawler_go/internal/retry/retry.go` (integrate 403/429 pause)

### 4. Add Prometheus Crawler Metric
- Edit `backend_go/internal/middleware/prometheus.go` (add crawler_blocks_total)

### 5. DB Deadlock Retry
- Edit `backend_go/internal/db/database.go` (add retry wrapper pgcode 40P01)

### 6. Rotate JWT_SECRET & Update Configs
- Edit `deployment/docker-compose.yml` (new 64-char secret in envs)

### 7. Enhance Monitoring Alert
- Edit `monitoring/prometheus/rules.yml` (crawler block alert)

### 8. Create Deployment Scripts
- Create `deployment/blue-green-deploy.ps1`
- Create `deployment/phase26-verify.ps1`

### 9. Add SQLi Test to README
- Edit `deployment/README.md`

### 10. Build & Test (Manual after edits)
- cd deployment && docker-compose build backend --no-cache
- docker-compose down && docker-compose up -d postgres backend prometheus grafana
- Run phase26-verify.ps1
- Grafana/Prom checks
- k6 loadtest
- Generate REPORT

**Progress: 2/10 complete** (crawler safeguards + DB retry)

