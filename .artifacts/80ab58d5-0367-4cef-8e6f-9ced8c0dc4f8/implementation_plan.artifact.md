# System Wipe, Restructuring & VPS Readiness Plan

This plan outlines the steps to perform a complete system cleanup, enterprise restructuring, and production asset generation for the RojgarSetu 2.0 project.

## User Review Required

> [!CAUTION]
> **Data Loss Warning**: Phase 1 involves `docker system prune -a --volumes -f` and aggressive deletion of local build artifacts (`node_modules`, `target`, etc.). Ensure no unsaved build-related data is needed locally.

## Proposed Changes

### Phase 1: Total Docker & Residual Deep Clean
- Run Docker cleanup commands.
- Delete `node_modules`, `target`, `.next`, `__pycache__`, `.venv`, and `.log` files across the repository.

### Phase 2: Enterprise Structure & Orphan Audit
- Move loose root scripts to `[scripts/](file:///F:/Rojgarsetu2.0/rojgarsetu2/scripts)`.
- Audit [`.gitignore`](file:///F:/Rojgarsetu2.0/rojgarsetu2/.gitignore) and [`.dockerignore`](file:///F:/Rojgarsetu2.0/rojgarsetu2/.dockerignore).
- Remove orphan files (deprecated reports/logs in root).

### Phase 3: Zero-Cache Build & Migration
- Execute `docker compose build --no-cache`.
- Boot services with `docker compose up -d`.
- Run database migrations after 10s wait.

### Phase 4: Autonomous Testing & Self-Healing
- Verify all containers are running.
- Run `.\scripts\verify-all.ps1`.
- Inject mock secrets (64+ chars) into `.env` if validations fail.

### Phase 5: VPS Production Asset Generation
- **[NEW]** `deployment/vps/docker-compose.prod.yml`: Production-hardened compose file.
- **[NEW]** `deployment/vps/nginx.conf`: Reverse proxy with SSL termination and security headers.

## Verification Plan

### Automated Tests
- Run `.\scripts\verify-all.ps1` from the host PowerShell.
- Verify 100% pass rate for all microservices.

### Manual Verification
- Check `docker ps` for healthy status.
- Inspect `deployment/vps/` for generated assets.
