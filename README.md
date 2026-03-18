# RojgarSetu 2.0
## Senior Civic Tech Job Portal

[![Go](https://img.shields.io/badge/Go-1.21-blue.svg)](https://golang.org)
[![Kubernetes](https://img.shields.io/badge/Kubernetes-1.28-green.svg)](https://kubernetes.io)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](LICENSE)

RojgarSetu 2.0 is a production-ready job portal platform for government and private sector jobs, courses, and videos. Built with modern microservices, monitoring, and disaster recovery.

### 🚀 Tech Stack
- **Backend**: Go (Gin, SQLC, Zerolog)
- **Database**: PostgreSQL 15
- **Cache**: Redis 7
- **Auth**: JWT (HS256 with rotation), Refresh Tokens
- **API Gateway**: Node.js Express
- **Auth Service**: Java Spring Boot
- **Crawler**: Go (ChromeDP for job scraping)
- **AI Engine**: Python FastAPI (recommendations)
- **Frontend**: React with modern UI
- **Orchestration**: Docker Compose, Kubernetes (HPA, PDB)
- **Monitoring**: Prometheus, Grafana, Loki
- **Storage**: S3 backups
- **Security**: Rate Limiting, SQLi Protection, TLS

### 🔒 Security Features
- JWT token rotation & validation
- SQL injection protection (prepared statements)
- Rate limiting (Redis-based)
- TLS encryption (Let's Encrypt ready)
- Secrets via env vars & K8s Secrets
- Input validation & sanitization
- Deadlock retry logic
- Crawler staleness alerts

### 🛠 Local Setup
1. Clone repo:
   ```
   git clone https://github.com/dippy79/rojgarsetu2.git
   cd rojgarsetu2
   ```

2. Copy & edit env template:
   ```
   cp .env.example .env
   # Edit .env with your secrets (DB_PASSWORD, JWT_SECRET, etc.)
   ```

3. Edit `.env` with your secrets (DB_PASSWORD, JWT_SECRET, etc.)

4. Start:
   ```
   cd deployment
   docker compose up -d
   ```

5. Access:
   - Frontend: http://localhost:3001
   - API: http://localhost:3000
Grafana: http://localhost:3002 (admin/admin **local dev only** - change in production)
   - Prometheus: http://localhost:9090

### ☁️ Production Deployment
1. **Kubernetes**:
   ```
   kubectl create namespace rojgarsetu
   kubectl create secret generic rojgarsetu-secrets --from-env-file=.env -n rojgarsetu
   kubectl apply -f deployment/k8s/
   ```

2. **Blue/Green**:
   ```
   deployment/blue-green-deploy.ps1
   ```

3. **GitHub Actions**: Auto-deploy on main (see `.github/workflows/production-deploy.yml`)

### 🆘 Disaster Recovery
1. **Daily S3 Backups**: scripts/db-backup-s3.ps1
2. **Dry-run Restore**: scripts/db-dry-run-restore.ps1
3. **Crawler Block Runbook**: docs/ops/RUNBOOK_CRAWLER_BLOCK.md
4. **DB Restore Runbook**: docs/ops/RUNBOOK_DB_RESTORE.md

### 📊 Monitoring
- Prometheus dashboards for DB pool, TLS health, containers
- Grafana pre-configured (import dashboards)
- Staleness alerts for crawler

### 🔧 Development
```
# Install deps
cd backend_go && go mod tidy
cd api_gateway_node && npm ci
cd frontend && npm ci

# Local services
deployment/run_local_services.ps1

# Load test
deployment/loadtest/k6.js
```

### 📄 License
MIT - See LICENSE

### 🤝 Contributing
Fork → Branch → PR to main

**Production Ready - No hardcoded secrets!**

