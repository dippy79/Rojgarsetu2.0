# PHASE 29 — CONTINUOUS OPS AUTOMATION
## Approved Plan Implementation Steps

### STEP 1 — Create New Files
- [x] deployment/k8s/hpa-crawler.yaml
- [x] deployment/k8s/pdb-crawler.yaml  
- [x] deployment/k8s/vpa-crawler.yaml
- [x] docs/ops/RUNBOOK_INCIDENT.md
- [x] docs/ops/DR_DRILL.md
- [x] monitoring/grafana/provisioning/dashboards/backup-verify.json
- [x] monitoring/grafana/provisioning/dashboards/incident-timeline.json

### STEP 2 — Edit Existing Files
- [x] .github/workflows/prod-ci-cd.yml (add weekly backup verify job)
- [x] crawler_go/internal/sources/base.go (UA rotation, adaptive throttle, new metric)
- [x] monitoring/prometheus/rules.yml (Slack webhook for alerts)
- [x] backend_go/internal/logger/logger.go (add traceID correlation)
- [ ] crawler_go/internal/logger/logger.go (add traceID if exists, else backend middleware)
- [ ] docker-compose.monitoring.yml (Promtail app log scraping)
- [x] deployment/loadtest/k6.js (2000 VUs)

### STEP 3 — Build & Dry-Run Tests
- [ ] Rebuild crawler_go: cd crawler_go && go mod tidy && go build
- [ ] kubectl apply --dry-run=client -f deployment/k8s/
- [ ] docker-compose -f docker-compose.monitoring.yml up -d

### STEP 4 — Verification
- [ ] k6 run --vus 2000 --duration 30s deployment/loadtest/k6.js
- [ ] Check Grafana: new dashboards visible
- [ ] Prometheus: new alerts fire (simulate)
- [ ] Simulate DB outage: docker-compose stop postgres → follow RUNBOOK_INCIDENT.md
- [ ] Simulate TLS: fake cert expiry alert → verify runbook
- [ ] Crawler surge: manual block trigger → adaptive throttle metric

### STEP 5 — CI Test & Report
- [ ] Push to GitHub → verify weekly job runs
- [ ] Generate PHASE 29 REPORT.md
- [ ] Mark complete

**Progress: 0/XX**  
**Status: Implementation started**
