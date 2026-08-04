# ATS + JSON-LD Task TODO

## Task #1 — JSON-LD parser in company_pages.go (defensive fallback)
- [x] Add `parseJSONLDJobs` + `jsonLDJob` struct to `services/crawler-go/internal/sources/company_pages.go`
- [x] Call it in `parseHTMLJobs` as a fallback before regex patterns
- [x] Build check (`go build -mod=vendor ./cmd/scheduler`) → BUILD_SUCCESS
- [x] Single-run test against existing company list → confirmed 0-result no-op (TCS="Access Denied" block, Wipro=JS SPA, no JSON-LD in raw HTML)
- [x] Show diff + honest before/after → JSON-LD fallback WORKS on synthetic JSON-LD fixture (1 job), confirmed no-op on existing 15-company list

## Task #2 — ATS fetchers (separate sources)
- [ ] Create `services/crawler-go/internal/sources/greenhouse.go` (GreenhouseSource, orgs: stripe, gitlab)
- [ ] Create `services/crawler-go/internal/sources/lever.go` (LeverSource, org: lever)
- [ ] Wire both into `scheduler/run.go` privSources slice
- [ ] Build check
- [ ] Real crawl → show RunSummary fetch/save counts for greenhouse + lever

## Assessment
- [ ] Honest assessment of the 15 original CompanyList companies: JSON-LD/ATS-fixable vs bot-blocked (TCS/Infosys Akamai) vs SPA-only with no structured-data path
