-- 000023_add_forms_videos_uploads.up.sql

-- Government Forms Table
CREATE TABLE IF NOT EXISTS government_forms (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    title TEXT NOT NULL,
    department TEXT,
    category TEXT,
    pdf_url TEXT,
    official_apply_url TEXT,
    deadline TIMESTAMPTZ,
    created_at TIMESTAMPTZ DEFAULT NOW(),
    updated_at TIMESTAMPTZ DEFAULT NOW()
);

-- Videos Table (if not already existing or to enhance)
-- Note: youtube_videos already exists in schema v3, but we might want to unify or enhance it.
-- Let's stick to enhancing what exists or adding a general videos table if needed.
-- The prompt asked for a 'videos' table.

CREATE TABLE IF NOT EXISTS platform_videos (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    title TEXT NOT NULL,
    description TEXT,
    thumbnail_url TEXT,
    video_url TEXT NOT NULL,
    category TEXT,
    duration TEXT,
    views_count INT8 DEFAULT 0,
    created_at TIMESTAMPTZ DEFAULT NOW(),
    updated_at TIMESTAMPTZ DEFAULT NOW()
);

-- File Uploads Table
CREATE TABLE IF NOT EXISTS file_uploads (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    file_type TEXT NOT NULL, -- 'resume', 'avatar', 'document'
    file_url TEXT NOT NULL,
    original_name TEXT NOT NULL,
    file_size INT8,
    created_at TIMESTAMPTZ DEFAULT NOW()
);

-- Indexes for performance
CREATE INDEX IF NOT EXISTS idx_jobs_gov_cat_loc ON jobs_government(department, location);
CREATE INDEX IF NOT EXISTS idx_jobs_priv_cat_loc ON jobs_private(company, location);
CREATE INDEX IF NOT EXISTS idx_gov_forms_deadline ON government_forms(deadline);
CREATE INDEX IF NOT EXISTS idx_file_uploads_user_id ON file_uploads(user_id);
