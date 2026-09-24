DROP INDEX IF EXISTS idx_jobs_gov_job_hash_unique;
ALTER TABLE jobs_government DROP COLUMN IF EXISTS job_hash;
