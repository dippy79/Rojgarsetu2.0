# 👑 ROJGARSETU 2.0 - MASTER TODO (CONSOLIDATED)

## 🎨 FRONTEND (React / Next.js)
- [ ] Implement request retry logic for failed API calls.
- [ ] Add request cancellation for component unmounting.
- [ ] Implement proper code splitting and lazy loading for heavy routes.
- [ ] Fix image optimization (currently set to `unoptimized: true`).
- [ ] Add Service Worker for offline support.
- [ ] Standardize all file extensions to `.tsx`.
- [ ] Revamp CSS with per-screen gradients and glassmorphism (Progressed).
- [ ] Implement Mobile App screens (Candidate Profile, Company Profile, My Applications, Job Applications, Gov Jobs, Courses, Videos).

## ⚙️ BACKEND (Go / Java / Node.js)
- [ ] Consolidate duplicate routes (e.g., `/private-jobs` vs `/priv-jobs`).
- [ ] Optimize database indexes based on query analysis.
- [ ] Fix N+1 query issues in complex joins.
- [ ] Implement token revocation mechanism beyond basic expiry.
- [ ] Add Multi-Factor Authentication (MFA).
- [ ] Move rate limiting to Redis (currently in-memory).
- [ ] Integrate centralized error tracking/alerting (e.g., Sentry).

## 🤖 AI ENGINE (Python)
- [ ] Implement fallback logic for when Gemini API is unavailable.
- [ ] Optimize scoring algorithm (currently simple Jaccard similarity).
- [ ] Implement recommendation caching.
- [ ] Add input validation for all recommendation parameters.
- [ ] Implement cost controls and usage monitoring for AI APIs.

## 🕷️ CRAWLER (Go)
- [ ] Implement adaptive rate limiting based on target source response.
- [ ] Add circuit breaker pattern for failing sources.
- [ ] Implement automatic source health monitoring.
- [ ] Verify scraped data quality before saving to DB.
- [ ] Restore local codebase parity with running container (Ongoing).
- [ ] Complete Task D Batch Restructuring (Move to `shared/`, `jobs/`, `courses/`, `videos/`).

## 🛠️ INFRASTRUCTURE & DEVOPS
- [ ] Implement resource limits (CPU/Memory) in `docker-compose.yml`.
- [ ] Harden logging configuration (max-size, max-file).
- [ ] Complete blue-green deployment scripts.
- [ ] Finalize K8s manifests (HPA, VPA, PDB).
- [ ] Set up weekly database backup verification job.
- [ ] Implement secret rotation mechanism.

## 🔐 SECURITY
- [ ] Replace `localStorage` with `HttpOnly` cookies for token storage.
- [ ] Implement CSRF protection.
- [ ] Add Content Security Policy (CSP) headers.
- [ ] Implement request signing for inter-service communication.
- [ ] Add brute force protection.

## 📊 MONITORING & ANALYTICS
- [ ] Finalize Grafana dashboards (Throughput, DB Pool, TLS Health, Containers).
- [ ] Configure Slack webhooks for critical Prometheus alerts.
- [ ] Add staleness metrics for crawler sources.

---
*Note: This file was consolidated from 13+ scattered TODO files on 2026-08-29.*
