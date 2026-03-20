-- migrations/00001_initial.down.sql
DROP TABLE IF EXISTS youtube_videos CASCADE;
DROP TABLE IF EXISTS courses CASCADE;
DROP TABLE IF EXISTS jobs_private CASCADE;
DROP TABLE IF EXISTS jobs_government CASCADE;
DROP TABLE IF EXISTS crawler_logs CASCADE;
DROP FUNCTION IF EXISTS update_updated_at_column();

