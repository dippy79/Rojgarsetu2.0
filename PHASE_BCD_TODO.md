# RojgarSetu 2.0 — Phase B/C/D Execution TODO

## Phase B — Enterprise-Grade Job Ingestion (VERIFIED)
- [x] JSON-LD fallback in company_pages.go (defensive, no regression)
- [x] ATS sources: greenhouse.go (stripe+gitlab), lever.go (lever)
- [x] Wired into scheduler/run.go privSources
- [x] Real crawl: greenhouse fetched=733 saved=733, lever 0/0; no regression

## Phase C — EdTech Course Source Expansion (VERIFIED)
- [x] geeksforgeeks.go CourseFetcher (stdlib only)
- [x] Registered in scheduler/run.go courseSources
- [x] Schema compat with 000012 (courses table)
- [x] Real crawl: geeksforgeeks fetched=7 saved=7; coursera skills fix 60/60
- [x] DB counts verified (courses=67)

## Phase D — Multi-Language Support
### D1. Migration (APPLIED + VERIFIED)
- [x] 000013_add_language_columns.up/down.sql (all 5 tables + partial indexes)
- [x] Applied to Postgres (amitsharma / rojgarsetu2)
- [x] Verified: \d company_jobs + \d jobs_private — language column + partial index
- [x] 000013_backfill_language.up.sql applied (ascii() fix)
- [x] Backfill: jobs_gov=6, jobs_private=817, company_jobs=0, courses=67 (2 ur), yt=0

### D2. Crawler pipeline
- [x] internal/lang/lang.go + lang_test.go
- [ ] Wire lang.Detect into store.go Save* functions
- [ ] Add language to crawler INSERT statements
- [ ] Build crawler

### D3. Backend + AI filter
- [ ] Add language filter to search_service.go
- [ ] Add language filter to sqlc queries
- [ ] Update AI engine service.py
- [ ] sqlc generate + rebuild backend

### D4. Frontend filter UI
- [ ] Add language filter to FilterBar.jsx
- [ ] Add language to 4 pages
- [ ] Rebuild frontend

### D5. Verification
- [ ] Rebuild crawler + run real crawl
- [ ] Show fetch/save counts with language tagging
- [ ] Verify DB language distribution
- [ ] Consolidated test report (B + C + D)
