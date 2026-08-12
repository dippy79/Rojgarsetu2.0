-- Complete schema for sqlc generate from database/schema_v3_fixed.sql + refresh_tokens
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

-- users table
CREATE TABLE users (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    name TEXT NOT NULL,
    email TEXT UNIQUE NOT NULL,
    password_hash TEXT NOT NULL,
    role TEXT NOT NULL,
    phone TEXT,
    avatar_url TEXT,
    is_active BOOLEAN,
    is_verified BOOLEAN,
    last_login TIMESTAMPTZ,
    created_at TIMESTAMPTZ DEFAULT NOW(),
    updated_at TIMESTAMPTZ DEFAULT NOW()
);

-- candidates table
CREATE TABLE candidates (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    user_id UUID NOT NULL REFERENCES users(id),
    phone TEXT,
    resume_url TEXT,
    resume_parsed JSONB,
    skills TEXT[],
    experience_years INT4,
    current_company TEXT,
    current_position TEXT,
    location TEXT,
    linkedin_url TEXT,
    github_url TEXT,
    portfolio_url TEXT,
    bio TEXT,
    is_open_to_work BOOLEAN,
    expected_salary TEXT,
    preferred_job_type TEXT[],
    preferred_location TEXT[],
    created_at TIMESTAMPTZ DEFAULT NOW(),
    updated_at TIMESTAMPTZ DEFAULT NOW()
);

-- companies table
CREATE TABLE companies (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    user_id UUID NOT NULL REFERENCES users(id),
    name TEXT NOT NULL,
    industry TEXT,
    company_size TEXT,
    website TEXT,
    logo_url TEXT,
    description TEXT,
    headquarters TEXT,
    founded_year INT4,
    is_verified BOOLEAN,
    is_active BOOLEAN,
    approved_by UUID REFERENCES users(id),
    approved_at TIMESTAMPTZ,
    created_at TIMESTAMPTZ DEFAULT NOW(),
    updated_at TIMESTAMPTZ DEFAULT NOW()
);

-- company_jobs table
CREATE TABLE company_jobs (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    company_id UUID NOT NULL REFERENCES companies(id),
    title TEXT NOT NULL,
    description TEXT,
    requirements TEXT,
    responsibilities TEXT,
    location TEXT,
    job_type TEXT,
    experience_min INT4,
    experience_max INT4,
    salary_min INT4,
    salary_max INT4,
    salary_currency TEXT,
    skills TEXT[],
    benefits TEXT[],
    application_url TEXT,
    application_email TEXT,
    is_remote BOOLEAN,
    is_active BOOLEAN,
    views_count INT4,
    applications_count INT4,
    posted_at TIMESTAMPTZ,
    expires_at TIMESTAMPTZ,
    created_at TIMESTAMPTZ DEFAULT NOW(),
    updated_at TIMESTAMPTZ DEFAULT NOW()
);

-- job_applications table
CREATE TABLE job_applications (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    job_id UUID NOT NULL REFERENCES company_jobs(id),
    candidate_id UUID NOT NULL REFERENCES candidates(id),
    status TEXT,
    cover_letter TEXT,
    resume_url TEXT,
    match_score DOUBLE PRECISION,
    notes TEXT,
    applied_at TIMESTAMPTZ,
    updated_at TIMESTAMPTZ DEFAULT NOW()
);

-- jobs_government table
CREATE TABLE jobs_government (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    title TEXT NOT NULL,
    department TEXT,
    location TEXT,
    apply_url TEXT,
    last_date TIMESTAMPTZ,
    source TEXT NOT NULL,
    eligibility TEXT,
    vacancy_count INT4,
    salary TEXT,
    exam_date TIMESTAMPTZ,
    notification_pdf_url TEXT,
    is_active BOOLEAN,
    created_at TIMESTAMPTZ DEFAULT NOW(),
    updated_at TIMESTAMPTZ DEFAULT NOW()
);

-- jobs_private table
CREATE TABLE jobs_private (
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
    posted_at TIMESTAMPTZ,
    is_active BOOLEAN,
    created_at TIMESTAMPTZ DEFAULT NOW(),
    updated_at TIMESTAMPTZ DEFAULT NOW()
);

-- courses table
CREATE TABLE courses (
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
    is_free BOOLEAN,
    source TEXT NOT NULL,
    start_date TIMESTAMPTZ,
    end_date TIMESTAMPTZ,
    enrollment_count INT4,
    rating TEXT,
    is_active BOOLEAN,
    created_at TIMESTAMPTZ DEFAULT NOW(),
    updated_at TIMESTAMPTZ DEFAULT NOW()
);

-- youtube_videos table
CREATE TABLE youtube_videos (
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
    is_active BOOLEAN,
    created_at TIMESTAMPTZ DEFAULT NOW(),
    updated_at TIMESTAMPTZ DEFAULT NOW()
);

-- job_categories table
CREATE TABLE job_categories (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    name TEXT NOT NULL UNIQUE,
    slug TEXT NOT NULL UNIQUE,
    description TEXT,
    icon TEXT,
    color TEXT,
    display_order INT4 NOT NULL DEFAULT 0,
    is_active BOOLEAN NOT NULL DEFAULT true,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- job_trades table
CREATE TABLE job_trades (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    category_id UUID NOT NULL REFERENCES job_categories(id) ON DELETE CASCADE,
    name TEXT NOT NULL,
    slug TEXT NOT NULL,
    description TEXT,
    qualification_req TEXT,
    min_salary INT4,
    max_salary INT4,
    demand_level TEXT NOT NULL DEFAULT 'NORMAL' CHECK (demand_level IN ('CRITICAL', 'HIGH', 'MEDIUM', 'NORMAL')),
    icon TEXT,
    is_active BOOLEAN NOT NULL DEFAULT true,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    UNIQUE(category_id, slug)
);

-- user_enrollments table
CREATE TABLE user_enrollments (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    trade_id UUID NOT NULL REFERENCES job_trades(id) ON DELETE CASCADE,
    status TEXT NOT NULL DEFAULT 'active' CHECK (status IN ('active', 'completed', 'expired', 'cancelled')),
    enrolled_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    expires_at TIMESTAMPTZ NOT NULL,
    completed_at TIMESTAMPTZ,
    progress_pct INT4 NOT NULL DEFAULT 0 CHECK (progress_pct >= 0 AND progress_pct <= 100),
    metadata JSONB NOT NULL DEFAULT '{}',
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    UNIQUE(user_id, trade_id, status)
);

-- user_notification_logs table
CREATE TABLE user_notification_logs (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    enrollment_id UUID REFERENCES user_enrollments(id) ON DELETE SET NULL,
    notification_type TEXT NOT NULL CHECK (notification_type IN ('expiry_warning', 'expiry_final', 'enrollment_reminder', 'course_update')),
    channel TEXT NOT NULL DEFAULT 'in_app' CHECK (channel IN ('in_app', 'email', 'push', 'sms')),
    title TEXT NOT NULL,
    message TEXT NOT NULL,
    payload JSONB NOT NULL DEFAULT '{}',
    sent_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    read_at TIMESTAMPTZ,
    clicked_at TIMESTAMPTZ,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- refresh_tokens table
CREATE TABLE refresh_tokens (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    ip_hash TEXT,
    ua_hash TEXT,
    expires_at TIMESTAMPTZ NOT NULL,
    revoked BOOLEAN NOT NULL DEFAULT false,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

