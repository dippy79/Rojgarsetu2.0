 # Phase D — Multi-Language Support — Execution TODO

## Status Legend
- [x] done
- [ ] pending

## Phase B — Enterprise-Grade Job Ingestion (VERIFIED)
- [x] JSON-LD fallback in company_pages.go (defensive, no regression)
- [x] ATS sources: greenhouse.go (orgs stripe+gitlab), lever.go (org lever)
- [x] Wired into scheduler/run.go privSources
- [x] Real crawl run → greenhouse fetched=733 saved=733, lever fetched=0 saved=0
- [x] DB counts verified (jobs_private=779, greenhouse=739)
- [x] No regression: upsc 6/6, naukri 20/20, company_pages 1/1

## Phase C — EdTech Course Source Expansion (VERIFIED)
- [x] geeksforgeeks.go implemented as CourseFetcher (stdlib only)
- [x] Registered in scheduler/run.go courseSources
- [x] Schema compatibility verified with 000012 (courses table)
- [x] Real crawl run → geeksforgeeks fetched=7 saved=7
- [x] coursera skills NOT NULL fix → fetched=60 saved=60
- [x] DB counts verified (courses=67)

## Phase D — Multi-Language Support
### D1. Migration
- [x] Drafted 000013_add_language_columns.up.sql (adds `language` to jobs_government, jobs_private, courses, youtube_videos)
- [x] Drafted 000013_add_language_columns.down.sql
- [ ] APPLY migration to Postgres (needs explicit confirmation)
- [ ] Verify columns exist

### D2. Language detection pipeline (crawler)
- [x] internal/lang/lang.go (dependency-free Unicode/heuristic detector)
- [x] internal/lang/lang_test.go (empty-description → 'en' verified PASS)
- [ ] Wire lang.Detect into store.go SaveGovJob / SavePrivJob / SaveCourse / SaveVideo
- [ ] Add `language` to crawler save INSERT statements
- [ ] Build crawler (go build -mod=vendor ./cmd/scheduler)

### D3. Backend search_service.go + AI recommendation filter
- [ ] Add language filter to search_service.go (searchGovJobs / searchPrivJobs)
- [ ] Add language filter to priv_jobs.sql / gov_jobs.sql / courses.sql / videos.sql queries
- [ ] Update AI engine service.py to filter by language
- [ ] sqlc generate + rebuild backend

### D4. Frontend filter UI
- [ ] Add language filter to FilterBar.jsx
- [ ] Add language to PrivateJobsPage.jsx / CoursesPage.jsx / GovJobsPage.jsx / VideosPage.jsx
- [ ] Rebuild frontend

### D5. Verification
- [ ] Rebuild crawler + run real crawl
- [ ] Show fetch/save counts with language tagging
- [ ] Verify DB: SELECT language, COUNT(*) FROM ... GROUP BY language
- [ ] Consolidated test report
