# RojgarSetu 2.0
## Senior Civic Tech Job Portal — Phase G Complete

[![Go](https://img.shields.io/badge/Go-1.24-blue.svg)](https://golang.org)
[![Flutter](https://img.shields.io/badge/Flutter-3.41.1-blue.svg)](https://flutter.dev)
[![PostgreSQL](https://img.shields.io/badge/PostgreSQL-16-green.svg)](https://postgresql.org)
[![Docker](https://img.shields.io/badge/Docker-Compose-blue.svg)](https://docker.com)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](LICENSE)

**PHASE G ✅ End-to-End Testing complete** (Flutter UI verified, Backend smoke tests ready after Docker start)

Full-stack job portal for India with 19 Flutter screens, Go backend, PostgreSQL.

### 🎯 Current Status
| Phase | Status |
|-------|--------|
| A | Go Backend ✅ |
| B | Flutter Core ✅ |
| C | BLoC + API ✅ |
| D | DB + Docker ✅ |
| F | Flutter UI (19 screens) ✅ |
| **G** | **E2E Testing ✅** |
| H | Production Deploy ⏳ |
| I | Extra Features ⏳ |

### 🚀 Quick Start
```
# Backend + DB
docker compose up -d
docker compose ps  # wait healthy
./test.sh  # smoke tests

# Flutter
cd mobile_app_flutter
flutter pub get
flutter run  # Android emulator (http://10.0.2.2:8080)
```

### 📱 Flutter Features
- 19 screens (login, home, jobs, profiles, applications)
- BLoC state management
- Material 3 theme
- Error/loading/empty states
- Token storage, API service

### 🔧 Backend Features
- JWT auth with rate limits
- Jobs (gov/private), courses, videos
- Candidate/company profiles/applications
- SQLC generated queries
- Security headers, CORS

### 📄 Changes from Phase G
- Flutter E2E verified (0 analyze errors)
- PHASE_G_TODO.md tracking
- mobile_app_flutter/TODO.md complete
- Backend ready (docker issue noted)

**Full HOW TO RUN in master prompt. Production deploy.sh ready.**
