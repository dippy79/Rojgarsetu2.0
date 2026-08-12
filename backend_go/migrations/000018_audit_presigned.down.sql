-- 000018_audit_presigned.down.sql
-- Roll back audit and presigned URL log tables and fuzzy search index.

DROP INDEX IF EXISTS idx_candidates_skills_trgm;

DROP TABLE IF EXISTS presigned_url_log CASCADE;
DROP TABLE IF EXISTS audit_logs CASCADE;
