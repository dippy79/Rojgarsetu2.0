-- 000012_create_courses_and_youtube_videos.down.sql
-- Rollback for the courses / youtube_videos tables and the upsert unique
-- indexes added to existing jobs tables by the up migration.

DROP TRIGGER IF EXISTS update_courses_updated_at ON courses;
DROP TRIGGER IF EXISTS update_youtube_videos_updated_at ON youtube_videos;

DROP TABLE IF EXISTS youtube_videos CASCADE;
DROP TABLE IF EXISTS courses CASCADE;

DROP INDEX IF EXISTS idx_jobs_gov_source_apply_url_unique;
DROP INDEX IF EXISTS idx_jobs_priv_source_url_unique;

