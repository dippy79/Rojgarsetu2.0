-- 000007_create_company_jobs.down.sql

DROP TRIGGER IF EXISTS update_company_jobs_updated_at ON company_jobs;
DROP TABLE IF EXISTS company_jobs;
