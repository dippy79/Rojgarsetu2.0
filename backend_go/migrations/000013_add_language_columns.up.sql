-- 000013_add_language_columns.up.sql
-- Multi-Language Support (Phase D)
--
-- Adds a `language` column across ALL job/education tables: jobs_government,
-- jobs_private, company_jobs, courses, and youtube_videos. Stores an ISO
-- 639-1 code ('en', 'hi', 'ta', 'te', 'bn', 'gu', 'kn', 'ml', 'pa', 'ur',
-- 'or', ...).
--
-- company_jobs is included because it is the 5th job table (employer-posted-
-- jobs flow) and appears in the AI recommendation engine's UNION query.
--
-- Default 'en' makes the migration safe for existing rows (the overwhelmingly
-- English corpus) and for crawler rows that fail language detection (empty
-- description → detector returns 'en'). We add a partial index so the common
-- `language = 'en'` filter stays cheap and the far rarer non-English filter
-- is also covered.
--
-- The column is NOT NOT NULL because some rows may be inserted by code paths
-- that predate the column (e.g. raw admin inserts); COALESCE in queries
-- treats NULL as 'en'. The crawler, however, always writes a concrete value.
--
-- NOTE: This migration only adds the columns with DEFAULT 'en'. A separate
-- one-time backfill (see 000013_backfill_language.up.sql) re-runs real
-- language detection on the ~800 existing rows so pre-existing non-English
-- content is not silently mislabeled. New crawls always write a detected
-- value via lang.Detect.

BEGIN;

-- ============================================
-- jobs_government
-- ============================================
ALTER TABLE jobs_government
    ADD COLUMN IF NOT EXISTS language TEXT NOT NULL DEFAULT 'en';

CREATE INDEX IF NOT EXISTS idx_jobs_gov_language
    ON jobs_government(language) WHERE language <> 'en';

-- ============================================
-- jobs_private
-- ============================================
ALTER TABLE jobs_private
    ADD COLUMN IF NOT EXISTS language TEXT NOT NULL DEFAULT 'en';

CREATE INDEX IF NOT EXISTS idx_jobs_priv_language
    ON jobs_private(language) WHERE language <> 'en';

-- ============================================
-- company_jobs (5th job table — employer-posted flow)
-- ============================================
ALTER TABLE company_jobs
    ADD COLUMN IF NOT EXISTS language TEXT NOT NULL DEFAULT 'en';

CREATE INDEX IF NOT EXISTS idx_company_jobs_language
    ON company_jobs(language) WHERE language <> 'en';

-- ============================================
-- courses
-- ============================================
ALTER TABLE courses
    ADD COLUMN IF NOT EXISTS language TEXT NOT NULL DEFAULT 'en';

CREATE INDEX IF NOT EXISTS idx_courses_language
    ON courses(language) WHERE language <> 'en';

-- ============================================
-- youtube_videos
-- ============================================
ALTER TABLE youtube_videos
    ADD COLUMN IF NOT EXISTS language TEXT NOT NULL DEFAULT 'en';

CREATE INDEX IF NOT EXISTS idx_yt_videos_language
    ON youtube_videos(language) WHERE language <> 'en';

COMMIT;
