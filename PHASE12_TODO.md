# PHASE 12 TODO - SQLC + Schema v3

## Completed ✅
- [x] Created database/schema_v3.sql with users/candidates/companies/company_jobs/job_applications tables, FKs, indexes, triggers

## Pending ⏳
1. Update backend_go/sqlc.yaml: Add schema_v3.sql to schema list
2. Add missing queries/*.sql: candidates.sql, companies.sql, jobs.sql, applications.sql
3. cd backend_go && go get github.com/sqlc-dev/sqlc/v2/cmd/sqlc@latest github.com/jackc/pgx/v5/stdlib
4. go mod tidy
5. sqlc generate
6. Refactor backend_go/internal/db/database.go: Remove manual SQL, embed sqlc Queries, implement Querier
7. Fix services/*_service.go: Match generated CreateUserParams etc. types
8. go build ./cmd/server/main.go → fix errors
9. Update docker-compose.yml: Add schema_v3.sql to postgres init
10. docker-compose up --build && test /health

## Phase 12 Report [After build success]

