# PHASE 23 — DEPLOYABLE-READY BUILD
## Approved Plan Implementation Steps

✅ **Plan Approved** with improvements: HEALTHCHECK intervals, X-Request-ID header.

### Step 1: ✅ COMPLETE
Dockerfile.backend rewritten as multi-stage (golang:1.24-alpine3.20 → gcr.io/distroless/static-debian12:nonroot). Non-root user 1000:1000. Copies only binary/certs/migrations. ENTRYPOINT ["./backend_server"]. HEALTHCHECK HTTPS /health w/ --interval=30s --timeout=5s --retries=3. (Backup as .bak5 attempted; Windows cp failed but file updated successfully.)

### Step 2: Add Health, Liveness, Readiness Probes
- [ ] Create/update backend_go/internal/handlers/health_handler.go (if separate) or edit main.go: /live (always OK), /ready (DB ping + goose status), enhance /health (version + DB)
- [ ] Requires import github.com/pressly/goose/v2

### Step 3: Production Logging Improvements
- [ ] Create backend_go/internal/middleware/logging.go (JSON zerolog, request ID/correlation, fields: method/path/status/duration/userID, X-Request-ID header)
- [ ] Update backend_go/cmd/server/main.go: Use JSON logger, replace gin.Logger() with logging middleware
- [ ] Update backend_go/internal/logger/logger.go for JSON

### Step 4: Final Build + Runtime Verification
- [ ] cd backend_go && go mod tidy && go build ./cmd/server/main.go
- [ ] cd deployment && docker-compose build backend --no-cache
- [ ] docker-compose down
- [ ] docker-compose up -d postgres backend prometheus
- [ ] curl -k https://localhost:8443/health/live/ready
- [ ] curl http://localhost:8083/metrics
- [ ] docker-compose logs backend | grep -i "ready|tls|error"

### Step 5: Produce PHASE 23 REPORT
- [ ] Generate report in PHASE23_REPORT.md

**Progress: Starting Step 1**
**Current Status: [TODO: Update after each completed step]**
