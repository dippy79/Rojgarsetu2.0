-- RojgarSetu 2.0 - New Database Schema v2
-- PostgreSQL with optimized design for 4 major data categories
-- [FULL COMPLETE CONTENT from read_file]

CREATE EXTENSION IF NOT EXISTS \"uuid-ossp\";

CREATE TABLE IF NOT EXISTS jobs_government (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    title TEXT NOT NULL,
    department TEXT,
    location TEXT,
    apply_url TEXT,
    last_date DATE,
    source TEXT NOT NULL,
    eligibility TEXT,
    vacancy_count INTEGER,
    salary TEXT,
    exam_date DATE,
    notification_pdf_url TEXT,
    is_active BOOLEAN DEFAULT true,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- ALL indexes
CREATE INDEX idx_gov_jobs_department ON jobs_government(department);
CREATE INDEX idx_gov_jobs_location ON jobs_government(location);
CREATE INDEX idx_gov_jobs_last_date ON jobs_government(last_date);
CREATE INDEX idx_gov_jobs_source ON jobs_government(source);
CREATE INDEX idx_gov_jobs_created_at ON jobs_government(created_at DESC);
CREATE INDEX idx_gov_jobs_is_active ON jobs_government(is_active);

CREATE TABLE IF NOT EXISTS jobs_private (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    company TEXT NOT NULL,
    title TEXT NOT NULL,
    location TEXT,
    url TEXT,
    salary TEXT,
    experience TEXT,
    job_type TEXT,
    skills TEXT[],
    description TEXT,
    source TEXT NOT NULL,
    posted_at TIMESTAMP,
    is_active BOOLEAN DEFAULT true,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_priv_jobs_company ON jobs_private(company);
CREATE INDEX idx_priv_jobs_location ON jobs_private(location);
CREATE INDEX idx_priv_jobs_source ON jobs_private(source);
CREATE INDEX idx_priv_jobs_title ON jobs_private(title);
CREATE INDEX idx_priv_jobs_posted_at ON jobs_private(posted_at DESC);
CREATE INDEX idx_priv_jobs_is_active ON jobs_private(is_active);
CREATE INDEX idx_priv_jobs_skills ON jobs_private USING GIN(skills);

CREATE TABLE IF NOT EXISTS courses (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    provider TEXT NOT NULL,
    title TEXT NOT NULL,
    url TEXT NOT NULL,
    duration TEXT,
    mode TEXT,
    level TEXT,
    skills TEXT[],
    description TEXT,
    thumbnail_url TEXT,
    price TEXT,
    is_free BOOLEAN DEFAULT true,
    source TEXT NOT NULL,
    start_date DATE,
    end_date DATE,
    enrollment_count INTEGER,
    rating DECIMAL(3,2),
    is_active BOOLEAN DEFAULT true,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_courses_provider ON courses(provider);
CREATE INDEX idx_courses_mode ON courses(mode);
CREATE INDEX idx_courses_level ON courses(level);
CREATE INDEX idx_courses_source ON courses(source);
CREATE INDEX idx_courses_skills ON courses USING GIN(skills);
CREATE INDEX idx_courses_created_at ON courses(created_at DESC);
CREATE INDEX idx_courses_is_active ON courses(is_active);

CREATE TABLE IF NOT EXISTS youtube_videos (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    channel TEXT NOT NULL,
    channel_id TEXT,
    title TEXT NOT NULL,
    url TEXT NOT NULL,
    thumbnail TEXT,
    description TEXT,
    video_id TEXT NOT NULL,
    published_at TIMESTAMP,
    duration TEXT,
    view_count BIGINT,
    like_count BIGINT,
    category TEXT,
    is_active BOOLEAN DEFAULT true,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_yt_videos_channel ON youtube_videos(channel);
CREATE INDEX idx_yt_videos_channel_id ON youtube_videos(channel_id);
CREATE INDEX idx_yt_videos_published_at ON youtube_videos(published_at DESC);
CREATE INDEX idx_yt_videos_category ON youtube_videos(category);
CREATE INDEX idx_yt_videos_created_at ON youtube_videos(created_at DESC);
CREATE INDEX idx_yt_videos_is_active ON youtube_videos(is_active);
CREATE INDEX idx_yt_videos_video_id ON youtube_videos(video_id) UNIQUE;

CREATE TABLE IF NOT EXISTS crawler_logs (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    source TEXT NOT NULL,
    status TEXT NOT NULL,
    jobs_found INTEGER DEFAULT 0,
    jobs_saved INTEGER DEFAULT 0,
    errors TEXT,
    started_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    completed_at TIMESTAMP,
    duration_seconds INTEGER
);

CREATE INDEX idx_crawler_logs_source ON crawler_logs(source);
CREATE INDEX idx_crawler_logs_status ON crawler_logs(status);
CREATE INDEX idx_crawler_logs_started_at ON crawler_logs(started_at DESC);

CREATE OR REPLACE FUNCTION update_updated_at_column()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = CURRENT_TIMESTAMP;
    RETURN NEW;
END;
$$ language 'plpgsql';

CREATE TRIGGER update_jobs_government_updated_at BEFORE UPDATE ON jobs_government FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();
CREATE TRIGGER update_jobs_private_updated_at BEFORE UPDATE ON jobs_private FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();
CREATE TRIGGER update_courses_updated_at BEFORE UPDATE ON courses FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();
CREATE TRIGGER update_youtube_videos_updated_at BEFORE UPDATE ON youtube_videos FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

CREATE OR REPLACE VIEW v_active_gov_jobs AS SELECT id, title, department, location, apply_url, last_date, source, eligibility, vacancy_count, salary, exam_date, notification_pdf_url, created_at FROM jobs_government WHERE is_active = true ORDER BY created_at DESC;

CREATE OR REPLACE VIEW v_active_priv_jobs AS SELECT id, company, title, location, url, salary, experience, job_type, skills, description, source, posted_at, created_at FROM jobs_private WHERE is_active = true ORDER BY created_at DESC;

CREATE OR REPLACE VIEW v_active_courses AS SELECT id, provider, title, url, duration, mode, level, skills, description, thumbnail_url, price, is_free, source, start_date, end_date, enrollment_count, rating, created_at FROM courses WHERE is_active = true ORDER BY created_at DESC;

CREATE OR REPLACE VIEW v_active_yt_videos AS SELECT id, channel, channel_id, title, url, thumbnail, description, video_id, published_at, duration, view_count, like_count, category, created_at FROM youtube_videos WHERE is_active = true ORDER BY published_at DESC;

CREATE SEQUENCE IF NOT EXISTS jobs_government_id_seq;
CREATE SEQUENCE IF NOT EXISTS jobs_private_id_seq;
CREATE SEQUENCE IF NOT EXISTS courses_id_seq;
CREATE SEQUENCE IF NOT EXISTS youtube_videos_id_seq;

COMMENT ON TABLE jobs_government IS 'Government job postings from authentic sources (NCS, SSC, UPSC, State PSC, RRB, Employment News)';
COMMENT ON TABLE jobs_private IS 'Private job postings from verified sources (LinkedIn, Indeed, Google Jobs, Company career pages)';
COMMENT ON TABLE courses IS 'Courses from government and private providers (NPTEL, SWAYAM, NSDC, Coursera, Udemy)';
COMMENT ON TABLE youtube_videos IS 'YouTube videos from official channels for job seekers and learners';
COMMENT ON TABLE crawler_logs IS 'Logs for monitoring crawler performance and debugging';

