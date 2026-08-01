# Phase-by-Phase Implementation TODO

## Phase 1 — Crawler: Chrome path
- [x] Edit `services/crawler-go/Dockerfile` — Add chromium + CHROME_BIN env
- [x] Edit `docker-compose.yml` — Add CHROME_BIN env to crawler service
- [ ] Verify: `docker compose build crawler --no-cache && docker compose up -d crawler && docker compose ps crawler`

## Phase 2 — Auth: BCrypt + JWT deps
- [ ] Edit `services/auth-java/pom.xml` — Add spring-boot-starter-security, jjwt-impl, jjwt-jackson
- [ ] Edit `services/auth-java/src/main/java/.../AuthController.java` — Add BCrypt hashing/verify
- [ ] Verify: register user, restart container, login still works, `docker compose ps auth-service` shows healthy

## Phase 3 — AI: DB-backed recommendations
- [ ] Edit `services/ai-engine-python/requirements.txt` — Add psycopg2-binary
- [ ] Edit `services/ai-engine-python/recommender/service.py` — Add DB query, real recommendations
- [ ] Verify: POST /recommend/jobs returns non-empty ranked list

## Phase 4 — Gateway → AI
- [ ] Edit `services/api-gateway-node/src/index.js` — Add /api/jobs/recommendations/me route
- [ ] Verify: curl through gateway reaches AI engine

## Phase 5 — Backend search
- [ ] Create `backend_go/migrations/000009_fulltext_search.up.sql` + .down.sql
- [ ] Create `backend_go/internal/services/search_service.go`
- [ ] Create `backend_go/internal/handlers/search_handler.go`
- [ ] Edit `backend_go/cmd/server/main.go` — Wire search routes
- [ ] Verify: POST /search returns matching jobs

## Phase 6 — Company dedup + transactions
- [ ] Create `backend_go/migrations/000010_company_case_insensitive.up.sql` + .down.sql (with dedup cleanup)
- [ ] Edit `services/crawler-go/internal/store/store.go` — Wrap in sql.Tx transaction
- [ ] Verify: no duplicate company for mixed-case names, rollback on failure

## Phase 7 — Hardening
- [ ] Edit `services/crawler-go/internal/sources/naukri.go` — Sleep → chromedp.WaitVisible
- [ ] Edit `services/crawler-go/internal/sources/base.go` — Add sanitizeString()
- [ ] Edit all source files — Apply sanitization to scraped text
- [ ] Edit `backend_go/config/config.go` — JWT secret >= 32 bytes validation
- [ ] Create `monitoring/promtail-config.yml` — Minimal config
- [ ] Verify: all checks pass

## FINAL VERIFICATION
- [ ] `docker compose down -v`
- [ ] `docker compose build --no-cache`
- [ ] `docker compose up -d`
- [ ] `docker compose ps` — all services healthy
- [ ] Register + login + restart + login again
- [ ] Hit /recommend/jobs through gateway
- [ ] Hit /search with real query
- [ ] Run crawler, confirm no duplicate companies
