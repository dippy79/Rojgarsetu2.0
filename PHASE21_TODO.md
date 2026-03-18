# PHASE 21 IMPLEMENTATION TODO
## Approved Plan: Advanced Auth Security + Token Rotation

### Step 1: DB Layer & Migration [ ]
- [ ] Create `backend_go/internal/db/models.go` → add RefreshToken struct
- [ ] Create `backend_go/internal/db/queries/refresh_tokens.sql` → sqlc queries (Create/Get/Revoke/RevokeAll/CleanupExpired)
- [ ] Create `backend_go/migrations/00003_refresh_tokens.up.sql` → table + indexes (idx_user_revoked)
- [ ] Create `backend_go/migrations/00003_refresh_tokens.down.sql` → DROP TABLE
- [ ] Run `cd backend_go && sqlc generate`

### Step 2: Config Updates [ ]
- [ ] Edit `backend_go/config/config.go` → JWT.AccessTokenExpiry=15m, RefreshTokenExpiry=30d; LoginRateLimit=5

### Step 3: Auth Package [ ]
- [ ] Create `backend_go/internal/auth/jwt.go` → token generate/validate functions

### Step 4: Services [ ]
- [ ] Create `backend_go/internal/services/auth_service.go` → Login/Refresh/Logout/LogoutAll (+CleanupExpiredTokens goroutine)

### Step 5: Handlers [ ]
- [ ] Create `backend_go/internal/handlers/auth_handler.go` → POST /login, /refresh, /logout, /logout-all

### Step 6: Middleware [ ]
- [ ] Edit `backend_go/internal/middleware/auth.go` → integrate cfg.JWT.Secret
- [ ] Edit `backend_go/internal/middleware/ratelimit.go` → add LoginRateLimitMiddleware("/login", 5/min)

### Step 7: Main Integration [ ]
- [ ] Edit `backend_go/cmd/server/main.go` → init AuthService/Handler, add /api/v1/auth routes w/ login limiter

### Step 8: Build & Deploy [ ]
- [ ] `cd backend_go && go mod tidy && go build ./cmd/server/main.go`
- [ ] `cd deployment && docker-compose build backend --no-cache`
- [ ] `docker-compose down && docker-compose up -d postgres backend`
- [ ] Test: login → get tokens, refresh → new pair + old revoked, logout → revoke, logout-all → all revoked

### Step 9: Verification [ ]
- [ ] Check logs/DB for flows (tokens issued/revoked, IP/UA binding, rate limit)
- [ ] Phase 21 REPORT

**Notes**: Apply improvements - migration indexes + DEFAULT NOW(), Cleanup goroutine, LOGIN_RATE_LIMIT env.

