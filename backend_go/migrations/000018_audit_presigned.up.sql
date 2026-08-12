-- 000018_audit_presigned.up.sql
-- Add audit log and presigned URL log tables, plus fuzzy search index for candidates.

CREATE TABLE IF NOT EXISTS audit_logs (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID REFERENCES users(id),
    action VARCHAR(100) NOT NULL,
    ip_address VARCHAR(50),
    metadata JSONB,
    created_at TIMESTAMPTZ DEFAULT now()
);

CREATE INDEX IF NOT EXISTS idx_audit_user ON audit_logs(user_id, created_at DESC);

CREATE TABLE IF NOT EXISTS presigned_url_log (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID REFERENCES users(id),
    file_key TEXT NOT NULL,
    expires_at TIMESTAMPTZ,
    used BOOLEAN DEFAULT false,
    created_at TIMESTAMPTZ DEFAULT now()
);

CREATE EXTENSION IF NOT EXISTS pg_trgm;
CREATE INDEX IF NOT EXISTS idx_candidates_skills_trgm
    ON candidates USING gin(bio gin_trgm_ops);
