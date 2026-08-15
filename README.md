# RojgarSetu 2.0
## Complete Senior Civic Tech Job Portal — All Phases ✅

[![Go](https://img.shields.io/badge/Go-1.24-blue.svg)](https://golang.org)
[![Flutter](https://img.shields.io/badge/Flutter-3.41.1-blue.svg)](https://flutter.dev)
[![PostgreSQL](https://img.shields.io/badge/PostgreSQL-16-green.svg)](https://postgresql.org)
[![Docker](https://img.shields.io/badge/Docker-Compose-blue.svg)](https://docker.com)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](LICENSE)
[![Railway](https://img.shields.io/badge/Railway-Deploy-brightgreen.svg)](https://railway.app)

**Production-Ready Job Portal for India** — Backend API, Flutter Mobile App, Docker Deploy, VPS/Cloud Ready.

## 🚀 Features
- **Auth**: JWT + Refresh Tokens + Rate Limiting (5/min login)
- **Jobs**: Government + Private jobs, Company jobs, Applications
- **Profiles**: Candidate + Company full CRUD
- **Content**: Courses, Videos public API
- **Security**: CORS env, 1MB Body Limit, CSP/HSTS Headers
- **Mobile**: 19 Flutter Screens (auth, home, jobs, profiles, applications)
- **Infra**: PostgreSQL, Redis, Docker Compose, K8s ready
- **Crawler**: Automated job crawler with deduplication and scheduled execution

## 📊 Phase Status
| Phase | Description | Status |
|-------|-------------|--------|
| A | Go Backend API (auth/jobs/profiles) | ✅ |
| B | Flutter Core + Models | ✅ |
| C | BLoC + API Integration | ✅ |
| D | DB Migrations + Docker | ✅ |
| F | Flutter UI 19 Screens | ✅ |
| G | End-to-End Testing | ✅ |
| **H** | **Production Deploy** | **✅** |
| **I** | **Automated Job Crawler Module** | **✅** |

## 🏃 Quick Start (Local)
```
git clone https://github.com/dippy79/Rojgarsetu2.0.git
cd Rojgarsetu2.0
docker compose up -d  # postgres, redis, backend
curl localhost:8080/health  # {"status":"ok"}
curl -X POST localhost:8080/api/v1/auth/register -H "Content-Type: application/json" -d '{"email":"test@example.com","password":"Test@1234","role":"candidate","name":"Test"}'
```

**Flutter**:
```
cd mobile_app_flutter
flutter pub get
flutter run  # Android: use http://10.0.2.2:8080
```

## ☁️ Production Deploy
### 1. Railway.app (Free Tier)
```
railway login
railway init
railway up
```
- Auto detects docker-compose.yml
- Free Postgres/Redis
- Custom domain free

### 2. DigitalOcean VPS ($12/mo)
```
./deploy.sh  # VPS IP as param
# Auto: docker-compose up, nginx.conf, certbot SSL
```

**Flutter APK**:
```
cd mobile_app_flutter
flutter build apk --release --dart-define API_BASE_URL=https://yourdomain.com/api/v1
```

## 📋 API Routes
### Public
| Method | Path | Desc |
|--------|------|------|
| POST | /api/v1/auth/register | Create user |
| POST | /api/v1/auth/login | JWT tokens |
| POST | /api/v1/auth/refresh | New access token |
| GET | /api/v1/jobs | Active jobs |
| GET | /api/v1/companies | List companies |
| GET | /api/v1/candidates | List candidates |

### Protected (JWT)
| Method | Path | Desc |
|--------|------|------|
| POST | /api/v1/auth/logout | Revoke tokens |
| GET | /api/v1/users | List users |
| GET | /api/v1/candidates/me | My profile |
| PUT | /api/v1/candidates/me | Update profile |
| POST | /api/v1/jobs | Create job |
| POST | /api/v1/jobs/:id/apply | Apply job |

## 🕷️ Automated Job Crawler & Intelligence Engine (`backend_go/internal/crawler`)

RojgarSetu 2.0 features an automated, multi-source job crawler capable of ingesting jobs from public portals and APIs with built-in deduplication and scheduled execution.

### Architecture Features
- **Deduplication Engine**: Uses SHA-256 content hashing (`hash_checksum`) to prevent duplicate job insertion.
- **HTML & JSON Parsers**: Native parsing support for standard web pages (DOM extraction) and REST JSON responses.
- **Cron Scheduler**: Background ticker executing periodic crawls every 6 hours (`scheduler.go`).
- **Telemetry & Logging**: All runs are tracked in `crawler_logs` with telemetry metrics (found, added, duplicates, error count, execution time).

### API Endpoints (`/api/v1/crawler`)
| Method | Endpoint | Description |
| :--- | :--- | :--- |
| `POST` | `/api/v1/crawler/crawl` | Manually trigger crawler execution (optional query param: `source_id`) |
| `GET` | `/api/v1/crawler/stats` | Fetch aggregated crawler statistics (total crawled, unique, duplicates) |
| `GET` | `/api/v1/crawler/health` | System health check & 24h error log status |

## 🔮 Phase I — Future (Optional)
- FCM Push Notifications
- Resume Upload (S3/MinIO)
- Email Verification (SMTP)
- Password Reset Flow
- Admin Panel (users/jobs moderation)
- Private Jobs Flutter Screen
- Full Text Search (pg_trgm)
- Analytics Dashboard (Grafana)

## 📝 License
MIT — Free for commercial use.

**⭐ Star if useful!** https://github.com/dippy79/Rojgarsetu2.0

## ⚖️ Legal Compliance

RojgarSetu operates as a legally compliant job aggregation platform in accordance with Indian IT laws:

### IT Act 2000 Section 79 Compliance
- **Intermediary Platform**: RojgarSetu acts as an intermediary, sourcing job listings from official government portals (public records) and authorized API partners
- **Content Attribution**: All job listings include source attribution and links to original official portals
- **No Content Modification**: We do not modify or host original content; users are redirected to official sources

### Crawler Best Practices
- **Bot Identification**: Crawler identifies itself via User-Agent: `RojgarSetuBot/2.0 (+https://rojgarsetu.in/bot-policy; support@rojgarsetu.in)`
- **Robots.txt Respect**: All crawls check and respect robots.txt directives before fetching content
- **Rate Limiting**: Polite 2-second delay between requests to avoid overloading source servers
- **SHA-256 Deduplication**: Prevents duplicate listings through content hashing

### Data Protection (DPDP Act 2023)
- User data processed solely for job matching and recruitment purposes
- Users may request data deletion or access by contacting privacy@rojgarsetu.in
- Minimal data collection with explicit consent mechanisms

### Takedown Rights
- **Takedown API**: `POST /api/v1/legal/takedown` for content removal requests
- **Response Time**: Takedown requests processed as per IT Act guidelines
- **Contact**: legal@rojgarsetu.in for legal inquiries
- **Bot Policy**: https://rojgarsetu.in/bot-policy

### Legal Endpoints
|| Method | Endpoint | Description |
|| :--- | :--- | :--- |
|| `GET` | `/api/v1/legal/disclaimer` | Returns IT Act compliance disclaimer, privacy policy, and terms of service |
|| `POST` | `/api/v1/legal/takedown` | Submit takedown requests for content removal |
|| `GET` | `/api/v1/crawler/forms` | Access government forms and admit card information |
