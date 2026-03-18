# PHASE 24 — PRODUCTION DEPLOYMENT CHECKLIST TODO
Status: [IN PROGRESS]

## Steps:
- [ ] 1. Read config.go and middleware/logging.go for TLS/logging
- [ ] 2. Edit backend_go/cmd/server/main.go: Add HTTPS :8443, /live /ready /health enhanced, request ID middleware
- [ ] 3. Edit deployment/Dockerfile.backend: Fix HEALTHCHECK to HTTP /health
- [ ] 4. cd backend_go && go mod tidy && go build ./cmd/server/main.go
- [ ] 5. cd deployment && docker-compose build backend --no-cache
- [ ] 6. docker-compose down && docker-compose up -d postgres backend prometheus
- [ ] 7. Verify: curl -k https://localhost:8443/health; curl http://localhost:8083/{health,live,ready,metrics}; logs
- [ ] 8. Test HTTP errors (invalid JWT etc.)
- [ ] 9. Check DB: goose/psql status
- [ ] 10. Produce PHASE 24 REPORT

