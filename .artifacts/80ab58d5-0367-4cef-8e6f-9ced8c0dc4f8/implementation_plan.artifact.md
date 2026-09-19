# Enterprise System Wipe & VPS Production Readiness Plan

This plan consolidates the **Security Audit remediation** with a comprehensive **System Wipe and VPS Deployment preparation** strategy. It ensures the environment is clean, secure, and ready for professional hosting.

## User Review Required

> [!CAUTION]
> **Data Loss Warning:** Phase 1 includes a full Docker system prune. We will backup the database first, but ensure all other transient data is secured.
> **Production Secrets:** We will no longer use mock secrets. You must provide actual values for `.env` before the final deployment build.

## Proposed Changes

### Phase 1: Total Docker & Residual Deep Clean
- **Step 1.0: Database Preservation**
  - Execute `pg_dump` within the `rojgar-postgres` container to save a full backup (schema + data) to `deployment/backups/local_dump.sql`.
- **Step 1.1: Docker Prune**
  - Execute `docker system prune -af --volumes` to ensure no stale layers or cached configurations interfere with the production-ready build.
- **Step 1.2: Host Residual Clean**
  - Verify and remove any accidental `.gradle`, `.android`, or project caches on `C:` drive, maintaining strict `F:` drive residency.

### Phase 2: Security & Performance Hardening (P0/P1 Audit Items)
- **Step 2.1: Cookie & CORS Security**
  - [auth_service.go](file:///F:/Rojgarsetu2.0/rojgarsetu2/backend_go/internal/services/auth_service.go): Enforce `Secure=true` and `SameSite=Strict`.
  - [index.js](file:///F:/Rojgarsetu2.0/rojgarsetu2/services/api-gateway-node/src/index.js): Implement "Fail-Closed" CORS.
- **Step 2.2: Frontend Resilience**
  - [ErrorBoundary.jsx](file:///F:/Rojgarsetu2.0/rojgarsetu2/frontend/src/components/ErrorBoundary.jsx): Implement global React Error Boundaries.
- **Step 2.3: Database Performance**
  - [000029_performance_indices.up.sql](file:///F:/Rojgarsetu2.0/rojgarsetu2/backend_go/migrations/000029_performance_indices.up.sql): Add GIN trigram indexes for Govt job searches.
- **Step 2.4: Crawler Rate Limiting**
  - Implement request delays (jitter) and user agent rotation in `services/crawler-go/internal/scheduler/run.go` and `browser.go`.
- **Step 2.5: PII Log Sanitization**
  - Implement middleware/logic in the Go backend (`internal/logger/logger.go`) to mask email addresses in logs.

### Phase 3: Zero-Cache Build & Migration
- **Step 3.1: Rebuild Infrastructure**
  - Execute `docker compose build --no-cache` for all services.
- **Step 3.2: Database Initialization**
  - Run all migrations up to Version 29.
- **Step 3.3: Database Seeding (Conditional)**
  - Execute database seeder script to populate mandatory initial roles/configs **only if** the backup is not restored.
- **Step 3.4: Database Restoration**
  - Restore `deployment/backups/local_dump.sql` into the fresh database container to preserve testing state.

### Phase 4: Autonomous Testing & Self-Healing
- **Step 4.1: Environment Specification**
  - Create a strictly defined `.env.example` with all required production keys.
  - Implement a startup validator that fails gracefully if required secrets are missing. **No mock secrets will be used.**

### Phase 5: VPS Production Asset Generation
- **Step 5.1: Certbot & SSL Integration**
  - Update `nginx.conf` with `.well-known/acme-challenge/` paths.
  - Add `certbot` service to `docker-compose.prod.yml` for automated SSL.
- **Step 5.2: Persistent Logging**
  - Configure external volume mounts in `docker-compose.prod.yml` for the Go backend and API Gateway.

## Verification Plan

### Automated Tests
- `npm run build` (Frontend) to verify Error Boundaries.
- `go test ./...` (Backend) to verify security middleware and log sanitization.
- Verify crawler includes `time.Sleep` jitter logic.

### Manual Verification
- Verify `deployment/backups/local_dump.sql` exists and is successfully restored.
- Check Go backend logs to ensure emails appear as `a***@example.com`.
- Verify cookies are blocked over unencrypted connections.
