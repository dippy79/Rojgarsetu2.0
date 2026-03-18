-- Fixed schema_v3 for sqlc compatibility
-- All TIMESTAMP -> TIMESTAMPTZ, INT consistent

CREATE EXTENSION IF NOT EXISTS \"uuid-ossp\";

-- users
CREATE TABLE IF NOT EXISTS users (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    name TEXT NOT NULL,
    email TEXT UNIQUE NOT NULL,
    password_hash TEXT NOT NULL,
    role TEXT NOT NULL CHECK (role IN ('candidate', 'company', 'admin', 'superadmin')) DEFAULT 'candidate',
    phone TEXT,
    avatar_url TEXT,
    is_active BOOLEAN DEFAULT true,
    is_verified BOOLEAN DEFAULT false,
    last_login TIMESTAMPTZ,
    created_at TIMESTAMPTZ DEFAULT NOW(),
    updated_at TIMESTAMPTZ DEFAULT NOW()
);

-- candidates
CREATE TABLE IF NOT EXISTS candidates (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    phone TEXT,
    resume_url TEXT,
    resume_parsed JSONB,
    skills TEXT[],
    experience_years INT4 DEFAULT 0,
    current_company TEXT,
    current_position TEXT,
    location TEXT,
    linkedin_url TEXT,
    github_url TEXT,
    portfolio_url TEXT,
    bio TEXT,
    is_open_to_work BOOLEAN DEFAULT false,
    expected_salary TEXT,
    preferred_job_type TEXT[],
    preferred_location TEXT[],
    created_at TIMESTAMPTZ DEFAULT NOW(),
    updated_at TIMESTAMPTZ DEFAULT NOW()
);

-- companies
CREATE TABLE IF NOT EXISTS companies (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    name TEXT NOT NULL,
    industry TEXT,
    company_size TEXT,
    website TEXT,
    logo_url TEXT,
    description TEXT,
    headquarters TEXT,
    founded_year INT4,
    is_verified BOOLEAN DEFAULT false,
    is_active BOOLEAN DEFAULT true,
    approved_by UUID REFERENCES users(id),
    approved_at TIMESTAMPTZ,
    created_at TIMESTAMPTZ DEFAULT NOW(),
    updated_at TIMESTAMPTZ DEFAULT NOW()
);

-- company_jobs
CREATE TABLE IF NOT EXISTS company_jobs (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    company_id UUID NOT NULL REFERENCES companies(id) ON DELETE CASCADE,
    title TEXT NOT NULL,
    description TEXT,
    requirements TEXT,
    responsibilities TEXT,
    location TEXT,
    job_type TEXT,
    experience_min INT4 DEFAULT 0,
    experience_max INT4,
    salary_min INT4,
    salary_max INT4,
    salary_currency TEXT DEFAULT 'INR',
    skills TEXT[],
    benefits TEXT[],
    application_url TEXT,
    application_email TEXT,
    is_remote BOOLEAN DEFAULT false,
    is_active BOOLEAN DEFAULT true,
    views_count INT4 DEFAULT 0,
    applications_count INT4 DEFAULT 0,
    posted_at TIMESTAMPTZ DEFAULT NOW(),
    expires_at TIMESTAMPTZ,
    created_at TIMESTAMPTZ DEFAULT NOW(),
    updated_at TIMESTAMPTZ DEFAULT NOW()
);

-- job_applications
CREATE TABLE IF NOT EXISTS job_applications (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    job_id UUID NOT NULL REFERENCES company_jobs(id) ON DELETE CASCADE,
    candidate_id UUID NOT NULL REFERENCES candidates(id) ON DELETE CASCADE,
    status TEXT DEFAULT 'pending' CHECK (status IN ('pending', 'reviewing', 'shortlisted', 'rejected', 'hired')),
    cover_letter TEXT,
    resume_url TEXT,
    match_score DOUBLE PRECISION,
    notes TEXT,
    applied_at TIMESTAMPTZ DEFAULT NOW(),
    updated_at TIMESTAMPTZ DEFAULT NOW()
);

-- Trigger function
CREATE OR REPLACE FUNCTION update_updated_at_column()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = NOW();
    RETURN NEW;
END;
$$ language 'plpgsql';

-- Triggers (only if not exists)
-- ... (same as original)

-- Views (same as original)

