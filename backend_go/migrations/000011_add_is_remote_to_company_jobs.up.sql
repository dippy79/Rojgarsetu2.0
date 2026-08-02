-- 000011_add_is_remote_to_company_jobs.up.sql
-- Add is_remote flag to company_jobs, required by search_service.go's
-- full-text search query which already selects this column.

ALTER TABLE company_jobs ADD COLUMN IF NOT EXISTS is_remote BOOLEAN NOT NULL DEFAULT false;

CREATE INDEX IF NOT EXISTS idx_company_jobs_is_remote ON company_jobs(is_remote);
