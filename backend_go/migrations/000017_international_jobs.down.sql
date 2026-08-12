-- 000017_international_jobs.down.sql
-- Roll back international jobs tables and company_job_postings columns.

ALTER TABLE company_job_postings
    DROP COLUMN IF EXISTS currency_code,
    DROP COLUMN IF EXISTS work_location_type,
    DROP COLUMN IF EXISTS visa_sponsorship;

DROP TABLE IF EXISTS candidate_internal_ratings CASCADE;
DROP TABLE IF EXISTS company_reviews CASCADE;
DROP TABLE IF EXISTS job_reports CASCADE;
