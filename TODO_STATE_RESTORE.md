# Task State Restoration — CRITICAL AUDIT CATCH

Date: 2026-08-07

## Goal
Restore local crawler codebase to full parity with the running container
(21 sources / ~1045 items baseline) and complete Task D Batch 1 (vertical
slice restructuring).

## Audit Finding
Local disk had diverged from the running container. Missing logic:
- Task Group A: `deriveJobRegion` heuristic + `JobRegion` fields + migration 000015
- Task Group B: `tutorialspoint.go`, `w3schools.go` (course fetchers)
- Task Group C: `youtube.go` native XML RSS strategy
- Scheduler: 21-source registration

## Restoration Status

### Task Group A — Job Region classification ✅ DONE
- [x] `services/crawler-go/internal/store/store.go`:
      Added `deriveJobRegion(title, location)` heuristic returning
      `india` / `overseas` / `global_remote`.
- [x] Wired `deriveJobRegion` into `SaveGovJob` and `SavePrivJob` INSERT/UPDATE
      statements (adds `job_region` column + `job_region = EXCLUDED.job_region`).
- [x] `backend_go/migrations/000015_add_job_region.up.sql` created
      (ADD COLUMN IF NOT EXISTS + index on jobs_government & jobs_private).
- [x] `backend_go/migrations/000015_add_job_region.down.sql` created
      (DROP COLUMN IF EXISTS).
- [x] Verified `go build -mod=vendor ./...` passes (scheduler.exe produced).

### Task Group B — Course fetchers ✅ ALREADY PRESENT
- [x] `internal/courses/tutorialspoint.go` present
- [x] `internal/courses/w3schools.go` present
- [x] Wired in scheduler (8 course sources)

### Task Group C — YouTube native XML RSS ✅ ALREADY PRESENT
- [x] `internal/videos/youtube.go` present
- [x] Wired in scheduler (1 video source)

### Scheduler — 21-source registration ✅ ALREADY PRESENT
- [x] `internal/scheduler/run.go` wires exactly 21 sources:
      5 gov + 6 priv + 1 naukri + 8 course + 1 video
- [x] Matches crawl_summary.json baseline (21 sources / 1045 items)

### Task D Batch 1 — Vertical slice restructuring ✅ ALREADY PRESENT
- [x] `internal/scheduler/run.go` rewritten to feature-based vertical slices
      (`jobs`→gov/priv, `courses`, `videos`) via packages:
      shared/, jobs/gov/, jobs/priv/, courses/, videos/
- [x] Extraction logic unchanged (behavior-preserving)

## Verification
- [x] `go build -mod=vendor ./...` → success (no compiler errors)
- [x] `go vet -mod=vendor ./...` → success
- [x] `go build -o scheduler.exe ./cmd/scheduler` → success (24.6 MB binary)
- [x] `go test -mod=vendor ./internal/store/...` → no test files

## Pending / Follow-up
- [ ] Rebuild + redeploy crawler container (no-cache) to pick up job_region writes
- [ ] Trigger test crawl and confirm `job_region` populated in DB
- [ ] Confirm live crawl matches 21-source / ~1045-item baseline
