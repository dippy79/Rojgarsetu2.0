# Backend Security Audit Fixes - TODO

## Approved Plan Steps (Phase 2)

1. ✅ FIX1: backend_go/internal/middleware/auth.go - JWT exp validated (redundant check noted)
2. ✅ FIX2: backend_go/cmd/server/main.go - JWT_SECRET len <32 fatal check
3. ✅ FIX3: backend_go/internal/services/token_service.go - SHA256 hashing implemented
4. ✅ [ ] FIX4: backend_go/cmd/server/main.go - Wire GlobalRateLimit middleware
5. ✅ [ ] FIX5: backend_go/cmd/server/main.go - Wire LoginRateLimit on POST /auth/login
6. ✅ FIX6: backend_go/internal/middleware/body_limit.go created
7. ✅ [ ] FIX7: backend_go/config/config.go - Add AllowedOrigin + update main.go CORS
8. ✅ [ ] FIX8: backend_go/cmd/server/main.go - Add /robots.txt route
9. ✅ [ ] FIX9: backend_go/cmd/server/main.go - DB SSL warning
10.❌ [ ] VERIFY: cd backend_go &amp;&amp; go mod tidy &amp;&amp; go build ./...
11.❌ [ ] FINAL CHECKLIST print

**Progress: 0/11 complete**
**Next: Start with FIX1 after reading auth_handler.go if needed.**
