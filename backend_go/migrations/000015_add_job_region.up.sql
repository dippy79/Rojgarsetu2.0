-- 000015_add_job_region.up.sql
-- Adds a job_region column to jobs_government and jobs_private.
--
-- The LIVE database already has these columns applied (from an earlier
-- Task Group A execution: verified counts 319 India / 398 Overseas). This
-- migration is recreated on local disk so the schema change is not lost.
-- It uses ADD COLUMN IF NOT EXISTS so it is safe to re-run.
--
-- Values classify where a job is located:
--   * india          - domestic Indian job
--   * overseas       - job located outside India
--   * global_remote  - remote-first job open to candidates anywhere
--
-- The crawler populates this via the deriveJobRegion(title, location)
-- heuristic in services/crawler-go/internal/store/store.go.

BEGIN;

-- ============================================
-- JOBS_GOVERNMENT
-- ============================================
ALTER TABLE jobs_government
    ADD COLUMN IF NOT EXISTS job_region TEXT NOT NULL DEFAULT 'india';

CREATE INDEX IF NOT EXISTS idx_jobs_gov_job_region
    ON jobs_government(job_region);

-- ============================================
-- JOBS_PRIVATE
-- ============================================
ALTER TABLE jobs_private
    ADD COLUMN IF NOT EXISTS job_region TEXT NOT NULL DEFAULT 'india';

CREATE INDEX IF NOT EXISTS idx_jobs_priv_job_region
    ON jobs_private(job_region);

COMMIT;
