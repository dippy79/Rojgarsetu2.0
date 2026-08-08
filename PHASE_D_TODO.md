# Task Group D — Crawler package restructure (behavior-preserving)

Goal: Restructure the flat `services/crawler-go/internal/sources/` into clean
packages — `shared/`, `jobs/gov/`, `jobs/priv/`, `courses/`, `videos/` — while
preserving every struct/interface/orchestrator pattern byte-for-byte. Any
behavior difference after restructuring is a BUG.

## Verified baseline (19 sources wired in scheduler/run.go)

- Gov (5): upsc, ssc, ncs, rrb, employment_news
- Priv (6): indeed, linkedin, google_jobs, company_pages, greenhouse, lever
- Legacy (1): naukri  (saved as job_source, converted to PrivJobSource)
- Course (6): nptel, swayam, nsdc, coursera, udemy, geeksforgeeks
- Video (1): youtube  (OfficialChannels / gov channels via YouTube API)
= 19 sources. "21" in the brief is stale/aspirational — target is 19.
Add tutorialspoint/w3schools/region ONLY after restructure is verified stable
(separate follow-up work, NOT part of this task).

## Approved target structure

```
services/crawler-go/internal/
  shared/
    source.go      (package shared  — Source interface, JobSource, BaseSource.Name())
    types.go       (package shared  — GovJobSource, PrivJobSource, CourseSource,
                     YouTubeVideoSource, RSS structs, fetcher interfaces, all
                     validation/normalization helpers)
    base.go        (package shared  — BaseSource, UserAgentRotator, DomainLimiter,
                     SanitizeString, DoRequest, robots-txt helpers)
  jobs/gov/        (package gov  — upsc.go, ssc.go, ncs.go, rrb.go, employment_news.go)
  jobs/priv/       (package priv — indeed.go, linkedin.go, google_jobs.go,
                     company_pages.go, company_pages_test.go, greenhouse.go,
                     lever.go, naukri.go)
  courses/         (package courses — nptel.go, swayam.go, nsdc.go, coursera.go,
                     udemy.go, geeksforgeeks.go)
  videos/          (package videos — youtube.go)
  store/store.go   (change import `sources` -> `shared`; uses shared.GovJobSource
                     etc. — all bodies unchanged)
  scheduler/run.go (change imports -> gov, priv, courses, videos, shared;
                     fully-qualified constructor calls, e.g.
                     gov.NewUPSCSource(), priv.NewIndeedSource(), etc.)
```

## Source-mode note

All moved files keep their CONTENT byte-identical. Only the `package X` clause
and import paths change. Cross-package callers are updated to the new paths.
Every struct/interface/controller body is preserved.

## Batching plan (go build + go test after each batch)

- [ ] Batch 1: create `shared/` (source.go, types.go, base.go) — copy content,
      change package to `shared`. `go build ./...` + `go test ./...`.
      (Interim: old sources/ still present; shared adds no duplicate symbols
      because it occupies a new package path. sources/ still imports itself,
      so no break yet.)
- [ ] Batch 2: create `jobs/gov/` (5 gov files) with `package gov`. Build+test.
      Keep old copies in sources/ until scheduler+store switch (final rename
      step removes them), OR delete old sources/ files here. Decide per batch
      to keep intermediate builds green.
- [ ] Batch 3: create `jobs/priv/` (7 priv + test) with `package priv`.
      Build+test.
- [ ] Batch 4: create `courses/` (6 course files) with `package courses`.
      Build+test.
- [ ] Batch 5: create `videos/` (youtube.go) with `package videos`. Build+test.
- [ ] Batch 6: switch `store/store.go` import sources->shared. Build+test.
- [ ] Batch 7: rewrite `scheduler/run.go` to import shared/gov/priv/courses/
      videos and call the new fully-qualified constructors. Build+test.
- [ ] Batch 8: remove the now-orphaned old `internal/sources/` dir. Build+test.
- [ ] Batch 9: final `go build ./...` + `go vet ./...` + `go test ./...`.
      Confirm 19-source wiring intact (grep run.go for all 19 constructor
      calls + naukri).

## Post-restructure verification (proof nothing broke)

- [ ] Run the real crawler (scheduler binary / go run cmd/scheduler).
- [ ] Confirm `sources_run == 19` (or the same count the pre-restructure
      baseline produced) and per-source fetch/save counts match the current
      real 19-source crawl exactly (compare against crawl_summary.json /
      crawler_logs.txt committed baseline).

## EXECUTION RULE

Move small batches. After each batch run:
   cd services/crawler-go && go build -mod=vendor ./... && go test -mod=vendor ./...
If a batch fails the build, fix it before proceeding. Do NOT skip a batch.

