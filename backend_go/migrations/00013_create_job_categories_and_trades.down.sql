-- 00013_create_job_categories_and_trades.down.sql
-- Drops job_categories, job_trades, user_enrollments, and user_notification_logs tables

BEGIN;

-- Drop triggers first
DROP TRIGGER IF EXISTS update_user_enrollments_updated_at ON user_enrollments;
DROP TRIGGER IF EXISTS update_job_trades_updated_at ON job_trades;
DROP TRIGGER IF EXISTS update_job_categories_updated_at ON job_categories;

-- Drop tables in reverse dependency order
DROP TABLE IF EXISTS user_notification_logs;
DROP TABLE IF EXISTS user_enrollments;
DROP TABLE IF EXISTS job_trades;
DROP TABLE IF EXISTS job_categories;

COMMIT;