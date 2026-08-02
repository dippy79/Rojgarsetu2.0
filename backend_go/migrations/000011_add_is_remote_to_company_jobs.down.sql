-- 000011_add_is_remote_to_company_jobs.down.sql
-- Rollback: remove is_remote column and its index

DROP INDEX IF EXISTS idx_company_jobs_is_remote;
ALTER TABLE company_jobs DROP COLUMN IF EXISTS is_remote;
