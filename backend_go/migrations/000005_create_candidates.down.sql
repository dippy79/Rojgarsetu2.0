-- 000005_create_candidates.down.sql

DROP TRIGGER IF EXISTS update_candidates_updated_at ON candidates;
DROP TABLE IF EXISTS candidates;
