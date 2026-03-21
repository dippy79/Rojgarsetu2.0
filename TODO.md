# Migration Trigger Fix TODO

## Plan Steps
- [x] Edit backend_go/migrations/00003_refresh_tokens.up.sql (1 trigger)

- [x] Edit backend_go/migrations/000004_create_users.up.sql (2 triggers)

- [x] Edit backend_go/migrations/000005_create_candidates.up.sql (1 trigger)
- [x] Edit backend_go/migrations/000006_create_companies.up.sql (1 trigger)

- [x] Edit backend_go/migrations/000007_create_company_jobs.up.sql (1 trigger)

- [x] Edit backend_go/migrations/000008_create_job_applications.up.sql (1 trigger)

- [x] Verify: search_files "CREATE TRIGGER IF NOT EXISTS" → 0 results

- [x] Test: docker compose down -v && docker compose up -d && docker compose logs migrate (all pass ✅)


## Rules
- Exact string replace only
- No column changes
- No docker-compose changes

