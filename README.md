Markdown# 🚀 RojgarSetu 2.0 Engine
### Enterprise-Grade Civic Tech Job & Course Aggregation Platform — Production Ready ✅

[![CI Build](https://github.com/dippy79/Rojgarsetu2.0/actions/workflows/ci.yml/badge.svg)](https://github.com/dippy79/Rojgarsetu2.0/actions)
[![Go Version](https://img.shields.io/badge/Go-1.25-00ADD8.svg?style=flat&logo=go)](https://golang.org)
[![Docker Image Size](https://img.shields.io/badge/Docker%20Image-29.9%20MB-blue.svg?style=flat&logo=docker)](https://docker.com)
[![Flutter](https://img.shields.io/badge/Flutter-3.41.1-02569B.svg?style=flat&logo=flutter)](https://flutter.dev)
[![PostgreSQL](https://img.shields.io/badge/PostgreSQL-16-green.svg?style=flat&logo=postgresql)](https://postgresql.org)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](LICENSE)

**RojgarSetu 2.0** is an enterprise-grade, microservices-driven job aggregation platform engineered for high-performance job scraping, anti-fake verification, real-time notifications, and AI-driven candidate recommendation across India.

---

## 📊 Production Health Audit (7/7 Containers Healthy)

| Service Name | Image / Component | Port Mapping | Health Status | Highlights |
| :--- | :--- | :--- | :--- | :--- |
| **`rojgar-backend`** | `rojgarsetu2-backend` | `8083:8083` | 🟢 Healthy | Go 1.25 multi-stage Alpine binary (`./rojgar_api`) — **29.9 MB** |
| **`rojgar-frontend`** | `rojgarsetu2-frontend` | `8080:80` | 🟢 Healthy | React + Vite + Tailwind CSS SPA |
| **`rojgar-api-gateway`** | `rojgarsetu2-api-gateway` | `3001:3000` | 🟢 Healthy | Node.js rate-limited API gateway |
| **`rojgar-ai-engine`** | `rojgarsetu2-ai-engine` | `8000:8000` | 🟢 Healthy | Python FastAPI recommendation engine |
| **`rojgar-auth-service`** | `rojgar-auth-service` | `8081:8081` | 🟢 Healthy | Microservice JWT authentication service |
| **`rojgar-postgres`** | `postgres:16-alpine` | `5432:5432` | 🟢 Healthy | Full-text search with `tsvector` & PostGIS |
| **`rojgar-redis`** | `redis:7-alpine` | `6379:6379` | 🟢 Healthy | Caching & rate-limiting layer |

---

## 🌟 Key Architecture & Engine Enhancements

### 🛡️ Anti-Fake & Verification Engine
* **Canonical Domain Strictness:** Restricts government job ingestion exclusively to official government domains (`.gov.in`, `.nic.in`).
* **SHA-256 Deduplication:** Generates unique hash signatures `SHA256(Title + Company + URL)` to eliminate duplicate job listings across multi-source scrapers.
* **Scam Keyword Filtering:** Real-time NLP-assisted scanning to reject fraudulent private job postings before writing to PostgreSQL.

### 🕷️ High-Performance Scraper Suite & Diagnostics
* **Multi-Source Crawlers:** Ingests public sector jobs (UPSC, SSC, RRB) and private job feeds (Apna, Internshala, Shine, Adzuna, Jooble, NCS).
* **Browser Pool Memory Management:** Bounded `chromedp` browser pool prevents RAM spikes during deep headless crawls.
* **Diagnostic CLI Tool (`cmd/diag`):** Standalone CLI tool to execute dry-run crawler diagnostics in isolated memory mode without touching DB records.

### 🔐 Zero-Leak Vaulting & Security
* **Pre-commit Vaulting:** Automated Gitleaks hooks to prevent accidental commits of API keys or credentials.
* **JWT Authentication:** Secure claims verification across protected endpoints (`/api/v1/candidates/me`, `/api/v1/gov-jobs/:id/apply`).
* **Tiered Rate Limiting:** Global rate limiting (100 req/min general, 5 req/min auth/login, 10 req/min job applications).

### 🐳 Ultra-Lightweight Multi-Stage Docker Build
* **Multi-Stage Optimization:** Reduced Docker container image size from **~850MB to 29.9MB** (97% size reduction) using `CGO_ENABLED=0` Go 1.25 static binary on Alpine Linux.
* **Resilient Healthchecks:** Uses `pg_isready` DB healthchecks in `docker-compose.yml` to prevent service initialization race conditions.

---

## 🏗️ System Architecture

                              +---------------------+
                              |   React Frontend    |
                              |     (Port 8080)     |
                              +----------+----------+
                                         |
                                 (HTTP / REST API)
                                         |
                              +----------v----------+
                              |     API Gateway     |
                              |     (Port 3001)     |
                              +----------+----------+
                                         |
     +-----------------------------------+-----------------------------------+
     |                                   |                                   |
+--------v---------+                +--------v---------+                +--------v---------+|  Backend (Go)    |                | AI Recommender   |                |  Auth Service    ||   (Port 8083)    |                |  (Port 8000)     |                |   (Port 8081)    |+--------+---------+                +------------------+                +------------------+|+----+----+-----------------------+|         |                       |+---v---+ +---v----+          +-------v-------+|  DB   | | Redis  |          | Scraper Engine|| Postgres| | Cache|          | (UPSC/SSC/etc)|+-------+ +--------+          +---------------+
---

## 📁 Repository Directory Structure

Rojgarsetu2.0/├── .github/workflows/         # GitHub Actions CI/CD workflows (Go 1.25 & golangci-lint)├── backend_go/                # Go Core Microservice & Scraper Engine│   ├── cmd/                   # Application entrypoints (server, diag CLI)│   ├── internal/              # Core domain logic│   │   ├── legal/             # Anti-fake validator & domain filter│   │   └── middleware/        # JWT authentication & rate limiting│   ├── migrations/            # PostgreSQL database migration scripts│   ├── docs/                  # Auto-generated Swagger 2.0 API documentation│   └── Dockerfile             # Optimized multi-stage Alpine build (~29.9 MB)├── services/│   └── crawler-go/            # Dedicated Scraper Engine & Private Parsers├── database/                  # SQL schema definitions & PostGIS extensions├── frontend/                  # React / Vite SPA Client├── mobile_app_flutter/        # Flutter Mobile Application├── deployment/                # Kubernetes manifests & Nginx configurations├── docker-compose.yml         # Production microservices orchestrator└── README.md                  # Project documentation
---

## 🏃 Quick Start (Development & Local Setup)

### Prerequisites
* **Docker Desktop** (v24.0+)
* **Go** (v1.25+)
* **Node.js** (v22+)
* **Flutter SDK** (v3.41.1+)

### 1. Clone & Environment Setup
```bash
git clone [https://github.com/dippy79/Rojgarsetu2.0.git](https://github.com/dippy79/Rojgarsetu2.0.git)
cd Rojgarsetu2.0
Create a .env file in the root directory:Code snippet# Database Credentials
POSTGRES_USER=rojgar_admin
POSTGRES_PASSWORD=YOUR_SECURE_POSTGRES_PASSWORD
POSTGRES_DB=rojgarsetu2
POSTGRES_PORT=5432

# Caching & JWT
REDIS_HOST=redis
REDIS_PORT=6379
JWT_SECRET=YOUR_RANDOM_256_BIT_SECRET_KEY

# API Ports
BACKEND_PORT=8083
FRONTEND_PORT=8080
API_GATEWAY_PORT=3001
AI_ENGINE_PORT=8000
AUTH_SERVICE_PORT=8081
2. Launch Services via Docker ComposeBash# Build and ignite all 7 microservices
docker compose up --build -d

# Verify container status
docker compose ps
3. Check Backend Logs & Swagger DocsBash# Stream live logs from backend
docker compose logs -f backend

# Healthcheck Verification
curl http://localhost:8083/health
# Response: {"status":"UP"}
Access Interactive Swagger UI:👉 http://localhost:8083/swagger/index.html📋 Primary API EndpointsPublic RoutesMethodEndpointDescriptionGET/healthLive service health checkPOST/api/v1/auth/registerRegister new user accountPOST/api/v1/auth/loginAuthenticate and obtain JWT tokenGET/api/v1/gov-jobsList verified government job listingsGET/api/v1/private-jobsList verified private job listingsGET/api/v1/searchFull-text tsvector search queryGET/swagger/*Interactive Swagger API documentationProtected Routes (JWT Required)MethodEndpointDescriptionGET/api/v1/candidates/meFetch active user candidate profilePUT/api/v1/candidates/meUpdate candidate resume & profile detailsPOST/api/v1/gov-jobs/:id/applyApply for a specific government job listingGET/api/v1/wsWebSocket connection for real-time notifications🔄 CI/CD Pipeline (GitHub Actions)The repository includes a GitHub Actions workflow (.github/workflows/ci.yml) configured for Go 1.25:YAML- Run Go Unit Tests (`go test -v ./...`)
- Install active Go 1.25 GolangCI-Lint
- Run GolangCI-Lint static code analysis
- Verify Docker Multi-Stage Image Compilation
⚖️ Legal Compliance & Bot PolicyRojgarSetu 2.0 operates in strict compliance with Indian IT Laws:IT Act 2000 Section 79 (Intermediary Guidelines)Official Source Attribution: Every aggregated job post strictly links back to the original government or corporate application portal.Content Integrity: Scraper engines do not modify official notice content or collect application fees.Web Scraper Ethics & DPDP Act 2023User-Agent Identification: Crawlers identify as RojgarSetuBot/2.0 (+https://rojgarsetu.in/bot-policy).Robots.txt Adherence: Scrapers check and respect robots.txt rate limits before fetching data.Takedown API: Takedown requests are accepted at POST /api/v1/legal/takedown or via legal@rojgarsetu.in.🤝 Contributing & SupportFork the repositoryCreate your feature branch (git checkout -b feature/NewFeature)Commit your changes (git commit -m 'feat: Add NewFeature')Push to the branch (git push origin feature/NewFeature)Open a Pull RequestRepository: https://github.com/dippy79/Rojgarsetu2.0License: EULA License
