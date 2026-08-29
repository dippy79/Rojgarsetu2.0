-- Remove performance indexes
DROP INDEX IF EXISTS idx_gov_jobs_source;
DROP INDEX IF EXISTS idx_gov_jobs_created;
DROP INDEX IF EXISTS idx_priv_jobs_type;
DROP INDEX IF EXISTS idx_priv_jobs_salary;
DROP INDEX IF EXISTS idx_applications_candidate;
DROP INDEX IF EXISTS idx_crawled_hash;
