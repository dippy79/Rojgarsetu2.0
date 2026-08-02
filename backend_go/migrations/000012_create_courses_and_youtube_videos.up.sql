-- 000012_create_courses_and_youtube_videos.up.sql
-- Adds the courses and youtube_videos tables.
--
-- These tables were referenced by the backend's sqlc queries (courses.sql,
-- videos.sql) and by the crawler sources (nptel.go, swayam.go, nsdc.go,
-- coursera.go, udemy.go, youtube.go) but were never actually created by any
-- migration. 00001_initial.up.sql only created jobs_government and even
-- carried the comment "Add other tables: jobs_private, courses,
-- youtube_videos..." without doing so.
--
-- This migration also adds UNIQUE indexes required for idempotent crawler
-- upserts (ON CONFLICT ... DO UPDATE) on:
--   * jobs_government (source, apply_url)
--   * jobs_private    (source, url)
--   * courses         (source, url)
--   * youtube_videos  (video_id)
--
-- Column types mirror backend_go/internal/db/sqlc_schema.sql and were
-- cross-checked against the crawler structs in
-- services/crawler-go/internal/sources/types.go.

BEGIN;

-- Function used by updated_at triggers. Recreated defensively because the
-- legacy definition lives in database/schema_v2.sql (not a migration) while
-- 000007 already references it.
CREATE OR REPLACE FUNCTION update_updated_at_column()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = NOW();
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

-- ============================================
-- COURSES TABLE
-- ============================================
CREATE TABLE IF NOT EXISTS courses (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    provider TEXT NOT NULL,
    title TEXT NOT NULL,
    url TEXT NOT NULL,
    duration TEXT,
    mode TEXT,
    level TEXT,
    skills TEXT[] NOT NULL DEFAULT '{}',
    description TEXT,
    thumbnail_url TEXT,
    price TEXT,
    is_free BOOLEAN NOT NULL DEFAULT true,
    source TEXT NOT NULL,
    start_date TIMESTAMPTZ,
    end_date TIMESTAMPTZ,
    enrollment_count INT4,
    rating TEXT,
    is_active BOOLEAN NOT NULL DEFAULT true,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_courses_provider ON courses(provider);
CREATE INDEX IF NOT EXISTS idx_courses_mode ON courses(mode);
CREATE INDEX IF NOT EXISTS idx_courses_level ON courses(level);
CREATE INDEX IF NOT EXISTS idx_courses_source ON courses(source);
CREATE INDEX IF NOT EXISTS idx_courses_skills ON courses USING GIN(skills);
CREATE INDEX IF NOT EXISTS idx_courses_created_at ON courses(created_at DESC);
CREATE INDEX IF NOT EXISTS idx_courses_is_active ON courses(is_active);

-- Unique per (source, url) so crawler upserts are idempotent
CREATE UNIQUE INDEX IF NOT EXISTS idx_courses_source_url_unique ON courses(source, url);

-- ============================================
-- YOUTUBE VIDEOS TABLE
-- ============================================
CREATE TABLE IF NOT EXISTS youtube_videos (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    channel TEXT NOT NULL,
    channel_id TEXT,
    title TEXT NOT NULL,
    url TEXT NOT NULL,
    thumbnail TEXT,
    description TEXT,
    video_id TEXT NOT NULL,
    published_at TIMESTAMPTZ,
    duration TEXT,
    view_count BIGINT,
    like_count BIGINT,
    category TEXT,
    is_active BOOLEAN NOT NULL DEFAULT true,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_yt_videos_channel ON youtube_videos(channel);
CREATE INDEX IF NOT EXISTS idx_yt_videos_channel_id ON youtube_videos(channel_id);
CREATE INDEX IF NOT EXISTS idx_yt_videos_published_at ON youtube_videos(published_at DESC);
CREATE INDEX IF NOT EXISTS idx_yt_videos_category ON youtube_videos(category);
CREATE INDEX IF NOT EXISTS idx_yt_videos_created_at ON youtube_videos(created_at DESC);
CREATE INDEX IF NOT EXISTS idx_yt_videos_is_active ON youtube_videos(is_active);

-- Unique per video_id so crawler upserts are idempotent
CREATE UNIQUE INDEX IF NOT EXISTS idx_yt_videos_video_id_unique ON youtube_videos(video_id);

-- ============================================
-- UNIQUE INDEXES FOR EXISTING TABLES (upsert support)
-- ============================================
CREATE UNIQUE INDEX IF NOT EXISTS idx_jobs_gov_source_apply_url_unique
    ON jobs_government(source, apply_url);

CREATE UNIQUE INDEX IF NOT EXISTS idx_jobs_priv_source_url_unique
    ON jobs_private(source, url);

-- ============================================
-- UPDATED_AT TRIGGERS
-- ============================================
DROP TRIGGER IF EXISTS update_courses_updated_at ON courses;
CREATE TRIGGER update_courses_updated_at
    BEFORE UPDATE ON courses
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

DROP TRIGGER IF EXISTS update_youtube_videos_updated_at ON youtube_videos;
CREATE TRIGGER update_youtube_videos_updated_at
    BEFORE UPDATE ON youtube_videos
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

COMMIT;

