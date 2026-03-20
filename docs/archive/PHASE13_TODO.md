# PHASE13_TODO.md - SQLC QUERIER INTEGRATION

## Steps (Complete one by one, update status)

- [x] 1. Backup current files (database.go.bak2, models.go.bak2)
- [x] 2. Fix backend_go/sqlc.yaml (out: internal/db)
- [x] 3. `cd backend_go && sqlc generate`
- [x] 4. Verify generated: internal/db/models.go, queries.go, querier.go exist
- [x] 5. Add paginated/filtered queries to ALL queries/*.sql files (ListXXX, GetXXXCount, SearchXXX)
- [x] 6. `cd backend_go && sqlc generate` (regen with new queries)
- [ ] 4. Verify generated: internal/db/models.go, queries.go, querier.go exist
- [ ] 5. Add paginated/filtered queries to ALL queries/*.sql files (ListXXX, GetXXXCount, SearchXXX)
- [ ] 6. `cd backend_go && sqlc generate` (regen with new queries)
- [ ] 7. Rewrite backend_go/internal/db/database.go to use *db.Queries
- [ ] 8. Update services: user_service.go, candidate_service.go, company_service.go, company_job_service.go (job), gov_job_service.go, priv_job_service.go, etc. for param types
- [ ] 9. `cd backend_go && go mod tidy && go build ./cmd/server/main.go`
- [ ] 10. Fix build errors
- [ ] 11. Add schema_v3.sql to deployment/docker-compose postgres initdb.d
- [ ] 12. Full rebuild: cd deployment && docker-compose down -v && docker-compose build backend --no-cache && docker-compose up -d backend
- [ ] 13. Test: curl http://localhost:8083/health
- [ ] 14. Generate PHASE 13 REPORT

**Next step: 1/14**

