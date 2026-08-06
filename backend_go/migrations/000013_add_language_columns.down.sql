-- 000013_add_language_columns.down.sql
-- Reverts the Multi-Language Support (Phase D) migration.
-- Drops the partial indexes and the language columns added in the up migration.

BEGIN;

DROP INDEX IF EXISTS idx_jobs_gov_language;
DROP INDEX IF EXISTS idx_jobs_priv_language;
DROP INDEX IF EXISTS idx_company_jobs_language;
DROP INDEX IF EXISTS idx_courses_language;
DROP INDEX IF EXISTS idx_yt_videos_language;

ALTER TABLE jobs_government DROP COLUMN IF EXISTS language;
ALTER TABLE jobs_private DROP COLUMN IF EXISTS language;
ALTER TABLE company_jobs DROP COLUMN IF EXISTS language;
ALTER TABLE courses DROP COLUMN IF EXISTS language;
ALTER TABLE youtube_videos DROP COLUMN IF EXISTS language;

COMMIT;
