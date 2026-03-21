# Backend Phase A Completion TODO

## Phase 1: main.go - Port 8080, CORS env, wire middleware/services/handlers/routes [COMPLETED: port/CORS/middleware/services wired]
- Change serverPort = 8080
- CORS from ALLOWED_ORIGIN env
- Wire ALL middleware (add missing)
- Init ALL services/handlers
- Wire ALL public/protected routes exactly per specs

## Phase 2: services/ - Full impl no stubs [PENDING]
- user_service.go: Complete methods
- token_service.go: Verify complete
- candidate_service.go: Complete
- company_service.go: Full impl
- job_service.go: Full impl
- application_service.go: Full impl

## Phase 3: handlers/ - Full impl [PENDING]
- auth_handler.go: Register/Login/Logout/Refresh
- user_handler.go: ListUsers/GetUser
- candidate_handler.go: List/Get/MyProfile/UpdateMyProfile
- company_handler.go: List/Get/MyCompany/UpdateMyCompany
- job_handler.go: All methods
- application_handler.go: All methods

## Phase 4: middleware/ - Verify all [PENDING]
- auth.go JWT
- cors.go env origin
- ratelimit.go global+login
- body_limit.go
- secure_headers.go

## Phase 5: Routes wiring (in main.go) [PENDING]

## Verification [PENDING]
- go build ./...
- docker compose up -d
- curl tests for register/login/401/429
- Flutter analyze

**Next step: Phase 1 main.go**
