# PHASE 28 — FINAL PRODUCTION SANITIZATION & GITHUB HANDOVER
## STATUS: 🟢 IN PROGRESS

====================================================================
PHASE 28 — CONFIG HARDENING & SECRET ENFORCEMENT
====================================================================

### Improvements Applied:
1. Remove Go fallbacks:
   - [ ] backend_go/config/config.go → Replace `if JWT==\"\"` with `panic(\"JWT_SECRET required\")`.
   - [ ] crawler_go/cmd/scheduler/main.go → Delete `if dbURL==\"\"` block; require DATABASE_URL.
   - [ ] crawler_go/internal/store/store.go → Same for connStr.

2. Sanitize configs:
   - [ ] api_gateway_node/src/config/index.js → Throw error if DB_PASSWORD/JWT_SECRET missing.
   - [ ] auth_service_java/src/main/resources/application.properties → Use `${DB_PASSWORD}`, `${JWT_SECRET}`.
   - [ ] deployment/docker-compose.yml → Use `${DATABASE_URL}`, `${GF_SECURITY_ADMIN_PASSWORD}`.

3. Kubernetes Secrets:
   - [ ] deployment/k8s/configmap-backend.yaml → Remove inline stringData; reference secrets only.
   - [ ] Add comment: \"# Populate via kubectl create secret\".
   - [ ] .env.example → Append 15+ new keys (Monitoring, S3, AWS creds, etc.), preserve existing.

4. New files:
   - [ ] README.md → Full rewrite with sections: Tech Stack, Security, Local Setup, Deployment, Disaster Recovery.
   - [ ] .github/workflows/production-deploy.yml → CI/CD pipeline with secrets list.

------------------------------------------------------------
### STEP 1 — File Updates
------------------------------------------------------------
- [x] Edit backend_go/config/config.go
- [x] Edit crawler_go/cmd/scheduler/main.go
- [x] Edit crawler_go/internal/store/store.go
- [x] Edit api_gateway_node/src/config/index.js
- [x] Edit auth_service_java/src/main/resources/application.properties
- [x] Edit deployment/docker-compose.yml
 - [x] Edit deployment/k8s/configmap-backend.yaml
- [ ] Edit .env.example (append)
- [ ] Rewrite README.md
- [ ] Create .github/workflows/production-deploy.yml

------------------------------------------------------------
### STEP 2 — Build & Test
------------------------------------------------------------
- [ ] docker-compose up (using .env)
- [ ] go build crawlers
- [ ] kubectl apply --dry-run=client -f deployment/k8s/

------------------------------------------------------------
### STEP 3 — Verification
------------------------------------------------------------
- [ ] Ensure no compile errors.
- [ ] Run: git grep -i secret → no hardcoded secrets.
- [ ] Verify JWT_SECRET enforced.
- [ ] Verify DATABASE_URL enforced.
- [ ] CI/CD pipeline passes with secrets injected.

------------------------------------------------------------
### STEP 4 — Git History Scrubbing & Report
------------------------------------------------------------
- [ ] Provide bfg/git filter-repo commands.
- [ ] Generate PHASE 28 REPORT.

### RULES:
- Think deeply before each step.
- Do NOT repeat previous phases.
- Do NOT skip verification.
- STOP on any error and diagnose.
- WAIT for approval before PHASE 29.

**Next Phase Preview: PHASE 29 — Continuous Ops Automation**
