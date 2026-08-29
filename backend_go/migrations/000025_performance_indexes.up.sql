-- Jobs & Applications performance indexes
CREATE INDEX IF NOT EXISTS idx_gov_jobs_source ON jobs_government(source, created_at DESC);
CREATE INDEX IF NOT EXISTS idx_gov_jobs_created ON jobs_government(created_at DESC);
CREATE INDEX IF NOT EXISTS idx_priv_jobs_type ON jobs_private(job_type, created_at DESC);
CREATE INDEX IF NOT EXISTS idx_priv_jobs_salary ON jobs_private(salary);
CREATE INDEX IF NOT EXISTS idx_applications_candidate ON job_applications(candidate_id, status, created_at DESC);
CREATE INDEX IF NOT EXISTS idx_crawled_hash ON crawled_jobs(hash_checksum);
