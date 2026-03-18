# PHASE 19 — PROMETHEUS + GOOSE MIGRATIONS TODO

## Steps (Deployment-Ready Mode)

- [ ] 1. Update backend_go/go.mod: add prometheus/client_golang v1.20.1, pressly/goose v2.24.1
- [ ] 2. Create backend_go/internal/middleware/prometheus.go: Gin middleware for request metrics, DB histogram
- [ ] 3. Edit backend_go/cmd/server/main.go: import middleware, add /metrics handler, router.Use(PrometheusMiddleware), register build_info
- [ ] 4. Edit monitoring/prometheus/prometheus.yml: add backend scrape job
- [ ] 5. Create backend_go/migrations/00002_schema_v4_placeholder.up.sql and .down.sql
- [ ] 6. Edit deployment/Dockerfile.backend: COPY migrations, install goose binary, entrypoint RUN goose if MIGRATE=true
- [ ] 7. cd backend_go && go mod tidy && go build ./cmd/server/main.go (validate)
- [ ] 8. cd deployment && docker-compose build backend --no-cache && docker-compose up -d postgres backend prometheus
- [ ] 9. Verify: curl localhost:8083/health, /metrics | grep go_info, localhost:9090/api/v1/targets | grep backend, logs clean
- [ ] 10. Update PHASE19_TODO.md with ✓, produce report

**Current: Starting STEP 1**
