-- 000008_create_job_applications.down.sql

DROP TRIGGER IF EXISTS update_job_applications_updated_at ON job_applications;
DROP TABLE IF EXISTS job_applications;
