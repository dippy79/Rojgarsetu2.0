-- Migration 26: Add job_hash to jobs_government to support crawler deduplication
ALTER TABLE jobs_government ADD COLUMN IF NOT EXISTS job_hash VARCHAR(64);

-- Populate job_hash for existing records if any
UPDATE jobs_government SET job_hash = md5(title || department || apply_url) WHERE job_hash IS NULL;

-- Add unique constraint
CREATE UNIQUE INDEX IF NOT EXISTS idx_jobs_gov_job_hash_unique ON jobs_government(job_hash);
