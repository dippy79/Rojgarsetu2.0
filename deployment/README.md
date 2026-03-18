# RojgarSetu 2.0 - Docker Deployment Guide

## Overview

RojgarSetu 2.0 is a full-stack job portal application with microservices architecture. This guide covers the Docker-based deployment setup.

## Architecture

```
┌─────────────────────────────────────────────────────────────────┐
│                        Docker Network                           │
│                      (rojgar-network)                           │
├─────────────────────────────────────────────────────────────────┤
│                                                                 │
│  ┌──────────────┐    ┌──────────────┐    ┌──────────────┐     │
│  │ API Gateway  │───▶│ Auth Service │    │ AI Engine    │     │
│  │  (Node.js)   │    │   (Java)     │    │  (Python)    │     │
│  │   :3000      │    │   :8081      │    │   :8000      │     │
│  └──────┬───────┘    └──────────────┘    └──────────────┘     │
│         │                                                      │
│         ▼                                                      │
│  ┌──────────────┐    ┌──────────────┐    ┌──────────────┐     │
│  │   Crawler    │    │   Postgres   │    │    Redis     │     │
│  │    (Go)      │    │   :5432      │    │   :6379      │     │
│  │   :8082      │    └──────────────┘    └──────────────┘     │
│  └──────────────┘                                              │
│                                                                 │
│  ┌──────────────┐    ┌──────────────┐                          │
│  │  Prometheus  │    │   Grafana    │                          │
│  │   :9090      │◀───│   :3002      │                          │
│  └──────────────┘    └──────────────┘                          │
│                                                                 │
└─────────────────────────────────────────────────────────────────┘
```

## Prerequisites

- Docker Engine 20.10+
- Docker Compose v2.0+
- At least 4GB RAM available
- Ports 3000, 3002, 5432, 6379, 8000, 8081, 8082, 9090 available

## Services

| Service        | Port | Description                          |
|---------------|------|--------------------------------------|
| API Gateway   | 3000 | Node.js Express API Gateway         |
| Auth Service  | 8081 | Java Spring Boot Authentication     |
| Crawler       | 8082 | Go Job Listing Crawler              |
| AI Engine     | 8000 | Python FastAPI Recommendation Engine|
| Postgres      | 5432 | PostgreSQL Database                  |
| Redis         | 6379 | Redis Cache                         |
| Prometheus   | 9090 | Metrics Collection                  |
| Grafana       | 3002 | Monitoring Dashboard                |

## Quick Start

### 1. Clone and Navigate

```bash
cd rojgarsetu2
```

### 2. Copy Environment File

```bash
cp .env.example .env
```

### 3. Build and Start All Services

```bash
cd deployment
docker compose build
docker compose up -d
```

### 4. Verify Containers

```bash
docker ps
```

All containers should show status as "Up" with healthy status for services with healthchecks.

## Run Instructions

### Start All Services
```bash
cd deployment
docker compose up -d
```

### Start Specific Service
```bash
docker compose up -d api-gateway
docker compose up -d postgres
```

### Stop All Services
```bash
docker compose down
```

### Stop and Remove Volumes (Complete Reset)
```bash
docker compose down -v
```

## Build Instructions

### Build All Services
```bash
cd deployment
docker compose build
```

### Build Specific Service
```bash
docker compose build api-gateway
docker compose build ai-engine
docker compose build crawler-service
docker compose build auth-service
```

### Rebuild Without Cache
```bash
docker compose build --no-cache
```

## Logging Commands

### View All Logs
```bash
docker compose logs -f
```

### View Specific Service Logs
```bash
docker compose logs -f api-gateway
docker compose logs -f postgres
docker compose logs -f redis
docker compose logs -f ai-engine
docker compose logs -f crawler-service
docker compose logs -f auth-service
docker compose logs -f prometheus
docker compose logs -f grafana
```

### View Last N Lines
```bash
docker compose logs --tail=100 api-gateway
```

### Filter Logs by Level
```bash
docker compose logs | grep ERROR
```

## Container Restart Instructions

### Restart All Services
```bash
docker compose restart
```

### Restart Specific Service
```bash
docker compose restart api-gateway
docker compose restart ai-engine
docker compose restart crawler-service
docker compose restart auth-service
docker compose restart postgres
docker compose restart redis
```

### Recreate Container (Preserves Data)
```bash
docker compose up -d --force-recreate api-gateway
```

## Debugging Commands

### Check Container Status
```bash
docker compose ps
```

### Check Container Health
```bash
docker inspect --format='{{.State.Health.Status}}' rojgar-api-gateway
docker inspect --format='{{.State.Health.Status}}' rojgar-postgres
docker inspect --format='{{.State.Health.Status}}' rojgar-redis
```

### Access Container Shell
```bash
docker exec -it rojgar-api-gateway sh
docker exec -it rojgar-postgres psql -U rojgarsetu
docker exec -it rojgar-redis redis-cli
docker exec -it rojgar-ai-engine sh
docker exec -it rojgar-crawler-service sh
docker exec -it rojgar-auth-service sh
```

### Check Resource Usage
```bash
docker stats
docker stats --no-stream rojgar-api-gateway
```

### View Container Logs with Timestamps
```bash
docker compose logs -t api-gateway
```

### Check Network Configuration
```bash
docker network inspect rojgarsetu2_rojgar-network
```

### View Environment Variables
```bash
docker exec rojgar-api-gateway env
```

## Database Commands

### Connect to PostgreSQL
```bash
docker exec -it rojgar-postgres psql -U rojgarsetu -d rojgarsetu
```

### Common PostgreSQL Commands
```sql
-- List tables
\dt

-- Exit
\q
```

### Backup Database
```bash
docker exec rojgar-postgres pg_dump -U rojgarsetu rojgarsetu > backup.sql
```

### Restore Database
```bash
docker exec -i rojgar-postgres psql -U rojgarsetu rojgarsetu < backup.sql
```

## Redis Commands

### Connect to Redis
```bash
docker exec -it rojgar-redis redis-cli
```

### Common Redis Commands
```bash
# Check connection
PING

# View all keys
KEYS *

# View info
INFO

# Exit
EXIT
```

## Monitoring

### Access Prometheus
- URL: http://localhost:9090
- Check targets: http://localhost:9090/targets

### Access Grafana
- URL: http://localhost:3002
- Default credentials: admin / admin123
- Add Prometheus data source: http://prometheus:9090

### Check Prometheus Metrics
```bash
curl http://localhost:9090/api/v1/query?query=up
```

## Troubleshooting

### Container Won't Start
```bash
# Check logs
docker compose logs <service-name>

# Check configuration
docker compose config

# Check port conflicts
netstat -ano | findstr "3000"
```

### Database Connection Issues
```bash
# Verify postgres is running
docker exec rojgar-postgres pg_isready

# Check network connectivity
docker exec rojgar-api-gateway ping postgres
```

### Permission Issues
```bash
# Fix volume permissions
docker compose down
sudo chown -R $USER:$USER .
docker compose up -d
```

### Clear All Data and Start Fresh
```bash
docker compose down -v
docker system prune -f
docker compose up -d
```

## Environment Variables

Create a `.env` file in the project root:

```env
# Required
DB_PASSWORD=your_secure_password
JWT_SECRET=your_jwt_secret_key
GRAFANA_PASSWORD=your_grafana_password

# Optional (have defaults)
# REDIS_PASSWORD=
# NODE_ENV=production
# CRAWLER_WORKERS=5
```

## Volumes

Data persists in Docker volumes:
- `postgres_data` - Database files
- `redis_data` - Redis persistence
- `ai_models` - ML model files
- `prometheus_data` - Metrics history
- `grafana_data` - Dashboard configurations

## Networks

All services communicate via the `rojgar-network` bridge network using service names as hostnames.

## Development Notes

- Services use healthchecks to ensure proper startup order
- Logs are accessible via `docker compose logs`
- Use `docker compose exec` to run commands inside containers
- All service-to-service communication uses Docker DNS (service names)

