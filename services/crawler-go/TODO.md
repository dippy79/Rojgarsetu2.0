# Crawler Wiring TODO

Goal: Wire all implemented source scrapers into the scheduler so crawled data
actually reaches the database.

## Steps

- [x] 0. Research: read store.go, source constructors/Fetch signatures, browser/proxy
         packages, migrations, and confirm courses/youtube_videos tables are missing.
- [x] 1. Add migration `backend_go/migrations/000012_create_courses_and_youtube_videos.up.sql`
         (+ down) creating `courses`, `youtube_videos`, and unique upsert constraints.
- [x] 2. Add `SaveGovJob`, `SavePrivJob`, `SaveCourse`, `SaveVideo` to
         `services/crawler-go/internal/store/store.go` (leave `SaveJob` as dead code).
- [x] 3. Add `services/crawler-go/internal/scheduler/run.go` orchestrator that runs all
         16 sources with per-source timeouts, routes results to the correct tables,
         and continues past per-source failures.
- [x] 4. Rewire `services/crawler-go/cmd/scheduler/main.go`: keep `/health`, add
         `POST /trigger` + `/crawl`, and a `time.Ticker` scheduler driven by
         `CRAWL_INTERVAL_HOURS` (default 6).
- [x] 5. Build (`go build ./...` from services/crawler-go).

## Post-rebuild review fixes (store.go date safety)

- [x] 6. Confirmed `SaveVideo()` does NOT insert `YouTubeVideoSource.Source` — the
         INSERT references only channel..category. Migration 000012 left unchanged
         (no `source` column needed on youtube_videos).
- [x] 7. Added `parseDateToTime()` helper in `services/crawler-go/internal/store/store.go`
         that converts scraped `*string` dates to `*time.Time` (RFC3339, ISO, dd/mm/yyyy,
         etc.) and returns `nil` (→ NULL) for empty/unparseable input.
- [x] 8. `SaveCourse()` now stores `parseDateToTime(course.StartDate/EndDate)` instead of
         raw strings into the TIMESTAMPTZ `start_date`/`end_date` columns.
- [x] 9. `SaveGovJob()` now stores `parseDateToTime(job.LastDate/ExamDate)` instead of raw
         strings into the TIMESTAMPTZ `last_date`/`exam_date` columns (live path today).
- [x] 10. Build verified: `go build -mod=vendor ./...` passes.

