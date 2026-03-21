-- 000007_create_company_jobs.up.sql
-- Company jobs table

CREATE TABLE IF NOT EXISTS company_jobs (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    company_id UUID NOT NULL REFERENCES companies(id) ON DELETE CASCADE,
    title TEXT NOT NULL,
    description TEXT NOT NULL DEFAULT '',
    location TEXT,
    job_type TEXT NOT NULL DEFAULT 'full_time',
    salary_min INT4,
    salary_max INT4,
    skills TEXT[] NOT NULL DEFAULT '{}',
    is_active BOOLEAN NOT NULL DEFAULT true,
    views INT4 NOT NULL DEFAULT 0,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_company_jobs_company_id ON company_jobs(company_id);
CREATE INDEX IF NOT EXISTS idx_company_jobs_location ON company_jobs(location);
CREATE INDEX IF NOT EXISTS idx_company_jobs_job_type ON company_jobs(job_type);
CREATE INDEX IF NOT EXISTS idx_company_jobs_skills ON company_jobs USING GIN(skills);
CREATE INDEX IF NOT EXISTS idx_company_jobs_is_active ON company_jobs(is_active);
CREATE INDEX IF NOT EXISTS idx_company_jobs_created_at ON company_jobs(created_at DESC);

DROP TRIGGER IF EXISTS update_company_jobs_updated_at ON company_jobs;
CREATE TRIGGER update_company_jobs_updated_at
    BEFORE UPDATE ON company_jobs
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

COMMENT ON TABLE company_jobs IS 'Private job postings by verified companies';
