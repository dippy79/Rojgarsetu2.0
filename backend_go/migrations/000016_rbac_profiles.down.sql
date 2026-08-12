-- 000016_rbac_profiles.down.sql
-- Roll back RBAC profile fields added in 000016_rbac_profiles.up.sql.

ALTER TABLE users
    DROP COLUMN IF EXISTS role,
    DROP COLUMN IF EXISTS google_id,
    DROP COLUMN IF EXISTS linkedin_id,
    DROP COLUMN IF EXISTS preferred_language;

ALTER TABLE candidates
    DROP COLUMN IF EXISTS deleted_at;

ALTER TABLE companies
    DROP COLUMN IF EXISTS verified_status,
    DROP COLUMN IF EXISTS job_credits;
