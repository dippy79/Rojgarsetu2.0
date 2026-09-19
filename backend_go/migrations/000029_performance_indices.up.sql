-- Add pg_trgm extension for fast ILIKE searches
CREATE EXTENSION IF NOT EXISTS pg_trgm;

-- GIN Trigram indexes for fast pattern matching on department and location
CREATE INDEX IF NOT EXISTS idx_gov_jobs_dept_trgm ON jobs_government USING gin (department gin_trgm_ops);
CREATE INDEX IF NOT EXISTS idx_gov_jobs_loc_trgm ON jobs_government USING gin (location gin_trgm_ops);

-- Composite index for the most common filter combination
CREATE INDEX IF NOT EXISTS idx_gov_jobs_filter_composite ON jobs_government (department, location, source) WHERE is_active = true;

-- Optimized index for private jobs as well
CREATE INDEX IF NOT EXISTS idx_priv_jobs_loc_trgm ON jobs_private USING gin (location gin_trgm_ops);
CREATE INDEX IF NOT EXISTS idx_priv_jobs_comp_trgm ON jobs_private USING gin (company gin_trgm_ops);
