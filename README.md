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
