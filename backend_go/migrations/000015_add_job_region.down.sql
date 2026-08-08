-- 000015_add_job_region.down.sql
-- Removes the job_region column from jobs_government and jobs_private.

BEGIN;

ALTER TABLE jobs_private
    DROP COLUMN IF EXISTS job_region;

ALTER TABLE jobs_government
    DROP COLUMN IF EXISTS job_region;

COMMIT;
