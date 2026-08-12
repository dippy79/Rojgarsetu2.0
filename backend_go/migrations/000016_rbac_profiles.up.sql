-- 000016_rbac_profiles.up.sql
-- Add RBAC profile fields to users and candidate/company metadata.

ALTER TABLE users
    ADD COLUMN IF NOT EXISTS role VARCHAR(20) DEFAULT 'CANDIDATE' CHECK (role IN ('CANDIDATE', 'COMPANY', 'ADMIN')),
    ADD COLUMN IF NOT EXISTS google_id VARCHAR,
    ADD COLUMN IF NOT EXISTS linkedin_id VARCHAR,
    ADD COLUMN IF NOT EXISTS preferred_language VARCHAR(10) DEFAULT 'en';

ALTER TABLE candidates
    ADD COLUMN IF NOT EXISTS deleted_at TIMESTAMPTZ;

ALTER TABLE companies
    ADD COLUMN IF NOT EXISTS verified_status BOOLEAN DEFAULT false,
    ADD COLUMN IF NOT EXISTS job_credits INT DEFAULT 5;
