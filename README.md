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
