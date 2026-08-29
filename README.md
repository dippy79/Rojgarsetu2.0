# 🚀 RojgarSetu 2.0 Engine
### Enterprise-Grade Civic Tech Job & Course Aggregation Platform — Production Ready ✅

[![CI Build](https://github.com/dippy79/Rojgarsetu2.0/actions/workflows/production-deploy.yml/badge.svg)](https://github.com/dippy79/Rojgarsetu2.0/actions)
[![Go Version](https://img.shields.io/badge/Go-1.24-00ADD8.svg?style=flat&logo=go)](https://golang.org)
[![Next.js](https://img.shields.io/badge/Next.js-14.2-black.svg?style=flat&logo=next.js)](https://nextjs.org)
[![Docker Image Size](https://img.shields.io/badge/Docker%20Image-29.9%20MB-blue.svg?style=flat&logo=docker)](https://docker.com)
[![PostgreSQL](https://img.shields.io/badge/PostgreSQL-16-green.svg?style=flat&logo=postgresql)](https://postgresql.org)
[![Redis](https://img.shields.io/badge/Redis-7-red.svg?style=flat&logo=redis)](https://redis.io)

**RojgarSetu 2.0** is an enterprise-grade, microservices-driven job aggregation platform engineered for high-performance job scraping, anti-fake verification, real-time notifications, and AI-driven candidate recommendation across India.

---

## 📊 Production Health Audit (8/8 Containers Healthy)

| Service Name | Image / Component | Port Mapping | Health Status | Highlights |
| :--- | :--- | :--- | :--- | :--- |
| **`rojgar-backend`** | `rojgarsetu2-backend` | `8083:8083` | 🟢 Healthy | Go 1.24 multi-stage Alpine binary — **29.9 MB** |
| **`rojgar-frontend`** | `rojgarsetu2-frontend` | `8080:80` | 🟢 Healthy | Next.js 14 (Pages Router) + Tailwind CSS |
| **`rojgar-api-gateway`** | `rojgarsetu2-api-gateway` | `3001:3000` | 🟢 Healthy | Node.js rate-limited API gateway |
| **`rojgar-ai-engine`** | `rojgarsetu2-ai-engine` | `8000:8000` | 🟢 Healthy | Python FastAPI + Gemini 2.0 LLM |
| **`rojgar-auth-service`** | `rojgar-auth-service` | `8081:8081` | 🟢 Healthy | Java Spring Boot JWT Auth service |
| **`rojgar-crawler`** | `rojgar-crawler` | `8082:8080` | 🟢 Healthy | Go-based multi-source scheduler |
| **`rojgar-postgres`** | `postgres:16-alpine` | `5432:5432` | 🟢 Healthy | `tsvector` Search & 20+ core tables |
| **`rojgar-redis`** | `redis:7-alpine` | `6379:6379` | 🟢 Healthy | Centralized caching & rate-limiting |

---

## 🌟 Key Architecture & Engine Enhancements

### 🛡️ Anti-Fake & Verification Engine
* **Canonical Domain Strictness:** Restricts government job ingestion exclusively to official government domains (`.gov.in`, `.nic.in`).
* **MD5 Hash Deduplication:** Generates unique hash signatures `MD5(Company + Title + Location)` to eliminate duplicate job listings across multi-source scrapers.
* **Scam Keyword Filtering:** Real-time NLP-assisted scanning to reject fraudulent private job postings before writing to PostgreSQL.

### 🕷️ High-Performance Scraper Suite & Diagnostics
* **Multi-Source Crawlers:** Ingests public sector jobs (UPSC, SSC, RRB) and private job feeds (Indeed, LinkedIn, Glassdoor, etc.).
* **Browser Pool Memory Management:** Bounded `chromedp` browser pool prevents RAM spikes during deep headless crawls.
* **Aggregator Core:** Real-time synchronization with official gazettes and recruitment notifications.

### 🔐 Enterprise Security & Reliability
* **Unified JWT Secret:** Shared token validation across Go, Java, and Node.js microservices.
* **Tiered Rate Limiting:** Global rate limiting enforced at API Gateway (100 req/min general, 5 req/min auth).
* **Automated DB Backups:** Daily 2 AM cron job with 7-day retention policy and Gzip compression.
* **Full-Text Search:** Optimized PostgreSQL `tsvector` with GIN indexing for sub-millisecond keyword matching.

### 📈 Analytics & Real-Time Intelligence
* **Dashboard Analytics:** Live throughput tracking using **Recharts** (Application velocity, Status mix).
* **AI Neural Engine:** Hyper-personalized job matching based on candidate skill graphs.
* **WebSocket Notifications:** Real-time status updates and alerts via backend-integrated WS handlers.

---

## 🏗️ System Architecture

                              +---------------------+
                              |   Next.js Frontend  |
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
+--------v---------+                +--------v---------+                +--------v---------+
|  Backend (Go)    |                | AI Recommender   |                |  Auth Service    |
|   (Port 8083)    |                |  (Port 8000)     |                |   (Port 8081)    |
+--------+---------+                +------------------+                +------------------+
         |
+--------v---------+                +--------v---------+                +--------v---------+
|  DB (PostgreSQL) |                | Redis (Cache)    |                | Crawler (Go)     |
|   (Port 5432)    |                |  (Port 6379)     |                |   (Port 8082)    |
+------------------+                +------------------+                +------------------+

---

## 📁 Repository Directory Structure

```text
Rojgarsetu2.0/
├── backend_go/                # Go Core Microservice (Business Logic, API)
│   ├── cmd/                   # Application entrypoints (server, migrations)
│   ├── internal/              # Domain logic, handlers, services, workers
│   └── migrations/            # SQL migration scripts (00001 - 00023)
├── services/
│   ├── crawler-go/            # Dedicated multi-source scraper engine
│   ├── ai-engine-python/      # FastAPI + Gemini LLM matching service
│   ├── auth-java/             # Spring Boot authentication microservice
│   └── api-gateway-node/      # Express.js rate-limited API gateway
├── frontend/                  # Next.js 14 (Pages Router) + Tailwind CSS
├── deployment/                # Production manifests (Docker, Nginx, K8s)
├── scripts/                   # System verification & backup utilities
└── docker-compose.yml         # Local orchestration profile
```

---

## 🏃 Quick Start (Development & Local Setup)

### Prerequisites
* **Docker Desktop** (v24.0+)
* **Go** (v1.24+)
* **Node.js** (v20+)
* **Python** (v3.10+)

### 1. Clone & Environment Setup
```bash
git clone https://github.com/dippy79/Rojgarsetu2.0.git
cd Rojgarsetu2.0
cp .env.example .env
```
Update the `.env` file with your secure credentials.

### 🔐 HttpOnly Cookie Authentication
RojgarSetu 2.0 uses **HttpOnly, Secure, and SameSite: Strict** cookies for JWT storage. This eliminates XSS-based token theft.
* The frontend automatically includes credentials in API calls.
* No sensitive tokens are stored in `localStorage`.

### 2. Launch Production Stack
```bash
# Build and ignite all microservices
docker compose up --build -d

# Verify container health
docker compose ps
```

### 3. Master Verification
Run the automated system health check:
```powershell
# Windows
.\scripts\verify-all.ps1

# Linux/macOS
./scripts/verify-all.sh
```

---

## ⚖️ Legal Compliance & Bot Policy
RojgarSetu 2.0 operates in strict compliance with Indian IT Laws (**IT Act 2000 Section 79**):
* **Source Attribution:** Every aggregated post strictly links to the original official portal.
* **Content Integrity:** Scrapers do not modify notice content or collect application fees.
* **User-Agent:** Crawlers identify as `RojgarSetuBot/2.0`.
* **Takedown API:** Requests accepted at `POST /api/v1/legal/takedown`.

🤝 **Contributing:** Fork → Feature Branch → Pull Request.
🛡️ **Security:** Report vulnerabilities to `security@rojgarsetu.in`.
📜 **License:** EULA / MIT (See LICENSE).
