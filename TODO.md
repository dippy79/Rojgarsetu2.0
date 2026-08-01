# Phase 7.2/7.3/7.4 + Phase 6 Migration Fix — Implementation TODO

## Status — ALL DONE ✅
- [x] Step 1: Add `SanitizeString()` to `services/crawler-go/internal/sources/base.go` (goquery-based HTML stripping)
- [x] Step 2: Apply sanitization in `services/crawler-go/internal/parser/parser.go` `ParseJob()`
- [x] Step 3: Add JWT_SECRET length >= 32 validation in `backend_go/config/config.go`
- [x] Step 4: Create `monitoring/promtail-config.yml` (minimal promtail config)
- [x] Step 5: Create `backend_go/migrations/000010_company_case_insensitive.up.sql` (dedup cleanup + unique index)
- [x] Step 6: Create `backend_go/migrations/000010_company_case_insensitive.down.sql`
- [x] Step 7: Verify `go build ./...` in `services/crawler-go/` — PASSED
- [x] Step 8: Verify `go build ./...` in `backend_go/` — PASSED
- [x] Step 9: Python import check for `services/ai-engine-python/` — PASSED
