<<<<<<< HEAD
🚀 RojgarSetu 2.0 EngineEnterprise-Grade Civic Tech Job & Course Aggregation Platform — Production Ready ✅RojgarSetu 2.0 is an enterprise-grade, microservices-driven job aggregation platform engineered for high-performance job scraping, anti-fake verification, real-time notifications, and AI-driven candidate recommendation across India.📊 Production Health & Verification Matrix (8/8 Services Verified)Service NameImage / ComponentPort MappingHealth / Test StatusKey Architectural Hardeningrojgar-backendrojgarsetu2-backend8083:8083🟢 Healthy / VerifiedGo 1.24 multi-stage Alpine binary — 29.9 MB, GIN indexed tsvector FTSrojgar-frontendrojgarsetu2-frontend8080:80🟢 Healthy / VerifiedNext.js 14 + Tailwind CSS, Webpack peer-dependency resolution hardenedrojgar-api-gatewayrojgarsetu2-api-gateway3001:3000🟢 Healthy / VerifiedNode 22 API Gateway with CSRF protection & In-Memory Redis fallbackrojgar-ai-enginerojgarsetu2-ai-engine8000:8000🟢 Healthy / VerifiedPython FastAPI + Gemini LLM with ThreadPoolExecutor parallel DB queriesrojgar-auth-servicerojgar-auth-service8081:8081🟢 Healthy / VerifiedJava Spring Boot + Go Admin MFA Claim (mfa_verified) verification middlewarerojgar-crawlerrojgar-crawler8082:8080🟢 Healthy / VerifiedGo crawler with semaphore bounded worker pools (MAX_WORKERS)rojgar-postgrespostgres:16-alpine5432:5432🟢 Healthy / Verifiedtsvector Search, row-level policies & 20+ core tablesrojgar-redisredis:7-alpine6379:6379🟢 Healthy / VerifiedDistributed caching & rate-limiting with graceful gateway degradation🌟 Key Architecture & Engine Enhancements🛡️ Security, Auth & Anti-Fake EngineAdmin MFA Enforcement: Dedicated Go middleware verifying the mfa_verified claim in JWT payloads, enforcing 403 Forbidden checks on all /api/admin/* routes.CSRF & Dynamic Rate Limiting: Mandatory CSRF token verification for state-changing HTTP requests, paired with a dynamic fallback to local in-memory rate-limiting whenever Redis is unreachable.Canonical Domain Strictness: Restricts government job ingestion exclusively to official domains (.gov.in, .nic.in).MD5 Hash Deduplication: Generates unique hash signatures MD5(Company + Title + Location) to eliminate duplicate job listings across multi-source scrapers.Scam Keyword Filtering: Real-time NLP-assisted scanning to reject fraudulent job postings before database writes.🕷️ High-Performance Scraper Suite & Worker PoolsBounded Crawler Pool: Semaphore-driven Go concurrency mechanism regulating chromedp instances via MAX_WORKERS to guarantee zero memory leaks or RAM spikes during deep scans.Multi-Source Crawlers: Scrapes public sector jobs (UPSC, SSC, RRB) and private job feeds concurrently.Aggregator Core: Real-time synchronization with official gazettes and recruitment notifications.⚡ AI Neural Recommender & Multi-ThreadingAsync Parallel Database Execution: Refactored Python AI matching engine using ThreadPoolExecutor and asyncio.gather to query company_jobs, jobs_private, and jobs_government concurrently, reducing processing latency by ~3x (1.01s total execution).Gemini LLM Matching: Hyper-personalized job matching based on candidate skill graphs and contextual relevance.🔐 Enterprise Infrastructure & CI/CD ResilienceSynchronized Runtimes: Standardized across Go 1.24 and Node.js 22.x environments.Resilient CI/CD Pipelines: GitHub Actions workflow stabilized with Webpack --legacy-peer-deps resolution, Go linting, and Trivy security scanning configured with Maven Central rate-limit protection.Unified JWT Security: Shared token validation across Go, Java, and Node.js microservices with HttpOnly cookie handling.Automated DB Backups: Daily cron job with 7-day retention policy and Gzip compression.🏗️ System ArchitecturePlaintext                              +---------------------+
                              |   Next.js Frontend  |
                              |     (Port 8080)     |
                              +----------+----------+
                                         |
                                  (HTTP / REST API)
                                         |
                              +----------v----------+
                              |     API Gateway     |
                              |  (CSRF & Rate Limit)|
                              |     (Port 3001)     |
                              +----------+----------+
                                         |
     +-----------------------------------+-----------------------------------+
     |                                   |                                   |
+--------v---------+                +--------v---------+                +--------v---------+
|  Backend (Go)    |                | AI Recommender   |                |  Auth Service    |
| (MFA Middleware) |                |(Parallel Queries)|                |  (Spring Boot)   |
|   (Port 8083)    |                |   (Port 8000)    |                |   (Port 8081)    |
+--------+---------+                +------------------+                +------------------+
         |
+--------v---------+                +--------v---------+                +--------v---------+
| DB (PostgreSQL)  |                | Redis (Cache)    |                | Crawler Pool (Go)|
|   (Port 5432)    |                |   (Port 6379)    |                | (Worker Bounded) |
+------------------+                +------------------+                |   (Port 8082)    |
                                                                        +------------------+
📁 Repository Directory StructurePlaintextRojgarsetu2.0/
├── backend_go/                 # Go 1.24 Core Microservice (Business Logic, Admin MFA, API)
│   ├── cmd/                    # Application entrypoints (server, migrations)
│   ├── internal/               # Domain logic, handlers, security middleware, workers
│   └── migrations/             # SQL migration scripts (00001 - 00023)
├── services/
│   ├── crawler-go/             # Scraper engine with semaphore worker pools
│   ├── ai-engine-python/       # FastAPI + Gemini LLM engine with parallel DB query executor
│   ├── auth-java/              # Spring Boot authentication microservice
│   └── api-gateway-node/       # Express.js / Node 22 gateway with CSRF & Redis fallback
├── frontend/                   # Next.js 14 + Tailwind CSS (Hardened Webpack pipeline)
├── deployment/                 # Production manifests (Docker, Nginx, K8s)
├── scripts/                    # System verification, test automation & backup utilities
└── docker-compose.yml          # Local & production orchestration profile
🏃 Quick Start (Development & Local Setup)PrerequisitesDocker Desktop (v24.0+)Go (v1.24+)Node.js (v22+)Python (v3.10+)1. Clone & Environment SetupBashgit clone https://github.com/dippy79/Rojgarsetu2.0.git
cd Rojgarsetu2.0
cp .env.example .env
Update the .env file with your secure credentials.🔐 HttpOnly Cookie & MFA AuthenticationRojgarSetu 2.0 uses HttpOnly, Secure, and SameSite: Strict cookies for JWT storage alongside strict MFA claims for administrative routes (/api/admin/*).Frontend automatically attaches credentials in API calls.Zero sensitive tokens stored in localStorage.2. Launch Production StackBash# Build and ignite all microservices
docker compose up --build -d

# Verify container health
docker compose ps
3. Master Automated Health & Fix VerificationRun the comprehensive automated verification matrix across Gateway, Auth, AI Engine, and Scraper modules:PowerShell# Windows
.\scripts\verify-all.ps1

# Linux/macOS
./scripts/verify-all.sh
⚖️ Legal Compliance & Bot PolicyRojgarSetu 2.0 operates in strict compliance with Indian IT Laws (IT Act 2000 Section 79):Source Attribution: Every aggregated post strictly links to the original official portal.Content Integrity: Scrapers do not modify notice content or collect application fees.User-Agent: Crawlers identify as RojgarSetuBot/2.0.Takedown API: Requests accepted at POST /api/v1/legal/takedown.🤝 Contributing: Fork → Feature Branch → Pull Request.🛡️ Security: Report vulnerabilities to security@rojgarsetu.in.📜 License: EULA / MIT (See LICENSE).
=======
# RojgarSetu 2.0

RojgarSetu is a service-oriented job aggregation platform for discovering public and private-sector opportunities, indexing job data, and exposing search and recommendation capabilities through a web application.

The repository contains the application services, local Docker Compose stack, database schemas, deployment manifests, monitoring configuration, and CI workflows.

## Repository Status

This README documents the repository's local development and deployment layout. Run the verification commands in this file against the current checkout before treating a build as production-ready.

## Architecture

The default `docker-compose.yml` stack contains:

| Service | Purpose | Host port |
| --- | --- | ---: |
| `postgres` | PostgreSQL database | `5432` |
| `redis` | Cache and rate-limit storage | `6379` |
| `backend` | Go API and core business logic | `8083` |
| `auth-service` | Spring Boot authentication service | `8081` |
| `ai-engine` | FastAPI recommendation and AI service | `8000` |
| `crawler` | Go job crawler and scheduled ingestion | `8082` |
| `api-gateway` | Node.js gateway and service proxy | `3001` |
| `frontend` | Next.js web application | `8080` |

Optional monitoring services are enabled with the `monitoring` Compose profile:

| Service | Host port |
| --- | ---: |
| Prometheus | `9090` |
| Grafana | `3002` |

Request flow is generally:

```text
Browser -> Frontend -> API Gateway -> Backend/Auth/AI/Crawler
                                      |
                                 PostgreSQL + Redis
```

## Repository Layout

```text
backend_go/                 Go backend API, migrations, and tests
services/api-gateway-node/  Node.js gateway
services/auth-java/         Spring Boot authentication service
services/ai-engine-python/  FastAPI AI/recommendation service
services/crawler-go/        Go crawler and ingestion workers
frontend/                   Next.js frontend
database/                   SQL schema snapshots
deployment/                 Production Docker, Nginx, and Kubernetes files
monitoring/                 Prometheus, Grafana, Loki, and Promtail config
scripts/                    Verification and backup scripts
.github/workflows/          CI and deployment workflows
```

## Prerequisites

- Docker Desktop with Docker Compose
- Go 1.25 or newer
- Node.js 20 or newer
- Java 17
- Python 3.10 or newer
- Maven for the Java service

Docker Compose is the recommended way to run the complete stack. Individual service commands are useful for focused development.

## Local Setup

1. Clone the repository and enter it:

   ```bash
   git clone https://github.com/dippy79/Rojgarsetu2.0.git
   cd Rojgarsetu2.0
   ```

2. Create a private environment file from the safe template:

   ```bash
   cp .env.example .env
   ```

   On Windows PowerShell, use:

   ```powershell
   Copy-Item .env.example .env
   ```

3. Edit `.env` locally and replace every placeholder with a value from your secret manager or local development setup. Never commit `.env`.

4. Start the application stack:

   ```bash
   docker compose up --build -d
   docker compose ps
   ```

5. Start optional monitoring:

   ```bash
   docker compose --profile monitoring up -d
   ```

Useful local URLs include:

- Frontend: `http://localhost:8080`
- API Gateway: `http://localhost:3001`
- Go backend health endpoint: `http://localhost:8083/health`
- Auth service health endpoint: `http://localhost:8081/actuator/health`
- AI engine health endpoint: `http://localhost:8000/health`
- Crawler health endpoint: `http://localhost:8082/health`

## Environment Configuration

`.env.example` lists the configuration names required by the Compose stack. It contains placeholders only. Depending on the service profile, the following values may be required:

- Database: `POSTGRES_USER`, `POSTGRES_PASSWORD`, `POSTGRES_DB`
- Authentication: `JWT_SECRET`, `REFRESH_TOKEN_KEY`, `CSRF_SECRET`
- AI and crawler integrations: `GEMINI_API_KEY`, `AI_ENGINE_API_KEY`, `YOUTUBE_API_KEY`, `CRAWLER_ADMIN_KEY`
- Networking: `ALLOWED_ORIGINS`, `API_GATEWAY_URL`
- Monitoring: `GRAFANA_ADMIN_USER`, `GRAFANA_ADMIN_PASSWORD`

Use long, randomly generated values for secrets. Keep development, CI, staging, and production credentials separate. Do not place secrets in source code, Dockerfiles, workflow YAML, SQL files, screenshots, test fixtures, issue reports, or log output.

## Development Commands

### Frontend

```bash
cd frontend
npm install
npm run dev
npm run type-check
npm run build
```

### Go backend

```bash
cd backend_go
go mod download
go build ./...
go test ./... --timeout=60s
```

### Go crawler

```bash
cd services/crawler-go
go mod download
go build ./...
go test ./... --timeout=60s
```

### Java authentication service

```bash
cd services/auth-java
mvn dependency:resolve --no-transfer-progress
mvn clean package -DskipTests --batch-mode
```

### Python AI engine

```bash
cd services/ai-engine-python
python -m venv .venv
python -m pip install -r requirements.txt
```

Run the service using the entrypoint defined in that service's source and Dockerfile, with configuration supplied through environment variables.

## Verification

The repository includes a PowerShell verification script:

```powershell
.\scripts\verify-all.ps1
```

At minimum, validate the Compose configuration before starting services:

```bash
docker compose config --quiet
```

CI also runs service builds and a Trivy filesystem vulnerability scan. The scan intentionally excludes dependency-heavy local directories such as `services/auth-java`, `node_modules`, `frontend/.next`, and `.git`; this does not replace dependency review or application security testing.

## Security and Secret Handling

- `.env`, local environment overrides, credentials, and private key files must remain untracked.
- Review `git diff --cached` before every commit.
- Do not print environment variables or full connection strings in CI logs.
- Use GitHub Actions Secrets for CI and deployment credentials; reference secret names, never their values, in workflow files.
- Rotate a credential immediately if it appears in a commit, log, screenshot, artifact, or pull request. Removing the text from a later commit does not invalidate the exposed credential.
- Report suspected vulnerabilities privately using the process in [SECURITY.md](SECURITY.md), not through a public issue.

Before opening a pull request, inspect tracked files for accidental credentials and confirm that only placeholder values appear in documentation and examples.

## Deployment

Deployment assets are under `deployment/` and include production Compose, Nginx, and Kubernetes configuration. Review the target environment's secret-management, TLS, database backup, network, and rollback requirements before deployment. Do not copy local `.env` files into an image or repository artifact.

## Legal and Crawling Guidelines

Use the crawler only against sources and APIs that permit automated access. Respect robots instructions, terms of service, rate limits, attribution requirements, and takedown requests. Do not collect credentials or unnecessary personal data, and do not represent aggregated listings as an official government publication.

## Contributing

1. Create a focused branch.
2. Keep secrets out of changes, test data, logs, and commit messages.
3. Run the relevant service checks and `docker compose config --quiet`.
4. Open a pull request describing behavior changes and security considerations.

See [SECURITY.md](SECURITY.md) for vulnerability reporting and security expectations. The repository license is available in [LICENSE](LICENSE).
>>>>>>> 30793fe8 (Readme.md updated)
