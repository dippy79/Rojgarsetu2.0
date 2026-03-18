-- migrations/00001_initial.up.sql
-- From schema_v2.sql (abridged)
CREATE EXTENSION IF NOT EXISTS \"uuid-ossp\";

CREATE TABLE jobs_government (
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

-- Add other tables: jobs_private, courses, youtube_videos...
-- Full schema in database/schema_v2.sql
