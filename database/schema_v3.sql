-- RojgarSetu 2.0 - Schema v3.0 - AUTH + PRIVATE JOBS + USERS
-- Compatible with schema_v2.sql - adds users, candidates, companies, applications
-- PostgreSQL optimized schema

-- Enable extensions (already in v2)
CREATE EXTENSION IF NOT EXISTS \"uuid-ossp\";

-- ============================================
-- USERS TABLE (core auth)
-- ============================================
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
    last_login TIMESTAMP,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_users_email ON users(email);
CREATE INDEX idx_users_role ON users(role);
CREATE INDEX idx_users_is_active ON users(is_active);
CREATE INDEX idx_users_created_at ON users(created_at DESC);

-- ============================================
-- CANDIDATES TABLE
-- ============================================
CREATE TABLE IF NOT EXISTS candidates (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    phone TEXT,
    resume_url TEXT,
    resume_parsed JSONB,
    skills TEXT[],
    experience_years INTEGER DEFAULT 0,
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
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    UNIQUE(user_id)
);

CREATE INDEX idx_candidates_user_id ON candidates(user_id);
CREATE INDEX idx_candidates_skills ON candidates USING GIN(skills);
CREATE INDEX idx_candidates_location ON candidates(location);
CREATE INDEX idx_candidates_is_open_to_work ON candidates(is_open_to_work);

-- ============================================
-- COMPANIES TABLE
-- ============================================
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
    founded_year INTEGER,
    is_verified BOOLEAN DEFAULT false,
    is_active BOOLEAN DEFAULT true,
    approved_by UUID REFERENCES users(id),
    approved_at TIMESTAMP,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    UNIQUE(user_id)
);

CREATE INDEX idx_companies_user_id ON companies(user_id);
CREATE INDEX idx_companies_industry ON companies(industry);
CREATE INDEX idx_companies_is_verified ON companies(is_verified);

-- ============================================
-- COMPANY JOBS (private jobs posted by companies)
-- ============================================
CREATE TABLE IF NOT EXISTS company_jobs (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    company_id UUID NOT NULL REFERENCES companies(id) ON DELETE CASCADE,
    title TEXT NOT NULL,
    description TEXT,
    requirements TEXT,
    responsibilities TEXT,
    location TEXT,
    job_type TEXT, -- full-time, part-time, contract, internship
    experience_min INTEGER DEFAULT 0,
    experience_max INTEGER,
    salary_min INTEGER,
    salary_max INTEGER,
    salary_currency TEXT DEFAULT 'INR',
    skills TEXT[],
    benefits TEXT[],
    application_url TEXT,
    application_email TEXT,
    is_remote BOOLEAN DEFAULT false,
    is_active BOOLEAN DEFAULT true,
    views_count INTEGER DEFAULT 0,
    applications_count INTEGER DEFAULT 0,
    posted_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    expires_at TIMESTAMP,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_company_jobs_company_id ON company_jobs(company_id);
CREATE INDEX idx_company_jobs_location ON company_jobs(location);
CREATE INDEX idx_company_jobs_job_type ON company_jobs(job_type);
CREATE INDEX idx_company_jobs_skills ON company_jobs USING GIN(skills);
CREATE INDEX idx_company_jobs_is_active ON company_jobs(is_active);
CREATE INDEX idx_company_jobs_posted_at ON company_jobs(posted_at DESC);

-- ============================================
-- JOB APPLICATIONS
-- ============================================
CREATE TABLE IF NOT EXISTS job_applications (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    job_id UUID NOT NULL REFERENCES company_jobs(id) ON DELETE CASCADE,
    candidate_id UUID NOT NULL REFERENCES candidates(id) ON DELETE CASCADE,
    status TEXT DEFAULT 'pending' CHECK (status IN ('pending', 'reviewing', 'shortlisted', 'rejected', 'hired')),
    cover_letter TEXT,
    resume_url TEXT,
    match_score DOUBLE PRECISION,
    notes TEXT,
    applied_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    UNIQUE(job_id, candidate_id)
);

CREATE INDEX idx_job_applications_job_id ON job_applications(job_id);
CREATE INDEX idx_job_applications_candidate_id ON job_applications(candidate_id);
CREATE INDEX idx_job_applications_status ON job_applications(status);
CREATE INDEX idx_job_applications_applied_at ON job_applications(applied_at DESC);

-- ============================================
-- UPDATE TIMESTAMP TRIGGER (shared with v2)
-- ============================================
CREATE OR REPLACE FUNCTION update_updated_at_column()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = CURRENT_TIMESTAMP;
    RETURN NEW;
END;
$$ language 'plpgsql';

-- Apply to new tables
DO $$
BEGIN
    IF NOT EXISTS (SELECT 1 FROM pg_trigger WHERE tgname = 'update_users_updated_at') THEN
        CREATE TRIGGER update_users_updated_at 
        BEFORE UPDATE ON users 
        FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();
    END IF;
    
    IF NOT EXISTS (SELECT 1 FROM pg_trigger WHERE tgname = 'update_candidates_updated_at') THEN
        CREATE TRIGGER update_candidates_updated_at 
        BEFORE UPDATE ON candidates 
        FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();
    END IF;
    
    IF NOT EXISTS (SELECT 1 FROM pg_trigger WHERE tgname = 'update_companies_updated_at') THEN
        CREATE TRIGGER update_companies_updated_at 
        BEFORE UPDATE ON companies 
        FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();
    END IF;
    
    IF NOT EXISTS (SELECT 1 FROM pg_trigger WHERE tgname = 'update_company_jobs_updated_at') THEN
        CREATE TRIGGER update_company_jobs_updated_at 
        BEFORE UPDATE ON company_jobs 
        FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();
    END IF;
    
    IF NOT EXISTS (SELECT 1 FROM pg_trigger WHERE tgname = 'update_job_applications_updated_at') THEN
        CREATE TRIGGER update_job_applications_updated_at 
        BEFORE UPDATE ON job_applications 
        FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();
    END IF;
END $$;

-- Views for active records
CREATE OR REPLACE VIEW v_active_candidates AS
SELECT c.*, u.name, u.email, u.role FROM candidates c JOIN users u ON c.user_id = u.id WHERE c.is_open_to_work = true;

CREATE OR REPLACE VIEW v_active_company_jobs AS
SELECT cj.*, c.name as company_name, c.logo_url, c.is_verified FROM company_jobs cj 
JOIN companies c ON cj.company_id = c.id WHERE cj.is_active = true AND (cj.expires_at IS NULL OR cj.expires_at > CURRENT_TIMESTAMP);

-- Comments
COMMENT ON TABLE users IS 'Central user authentication table';
COMMENT ON TABLE candidates IS 'Candidate profiles linked to users';
COMMENT ON TABLE companies IS 'Company profiles linked to users';
COMMENT ON TABLE company_jobs IS 'Private job postings by verified companies';
COMMENT ON TABLE job_applications IS 'Job applications with AI match scoring';

-- Grants (for app user)
-- GRANT ALL ON ALL TABLES IN SCHEMA public TO rojgarsetu;
-- GRANT USAGE, SELECT ON ALL SEQUENCES IN SCHEMA public TO rojgarsetu;

