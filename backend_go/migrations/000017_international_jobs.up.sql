-- 000017_international_jobs.up.sql
-- Add tables for job reports, company reviews, and candidate internal ratings.

CREATE TABLE IF NOT EXISTS job_reports (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    job_id UUID NOT NULL,
    job_type VARCHAR(20) NOT NULL,
    reporter_user_id UUID REFERENCES users(id) ON DELETE CASCADE,
    reason TEXT NOT NULL,
    status VARCHAR(20) DEFAULT 'PENDING',
    created_at TIMESTAMPTZ DEFAULT now()
);

CREATE TABLE IF NOT EXISTS company_reviews (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    company_id UUID NOT NULL REFERENCES companies(id) ON DELETE CASCADE,
    candidate_id UUID NOT NULL REFERENCES candidates(id) ON DELETE CASCADE,
    rating INT NOT NULL CHECK (rating BETWEEN 1 AND 5),
    review_text TEXT,
    created_at TIMESTAMPTZ DEFAULT now(),
    UNIQUE(company_id, candidate_id)
);

CREATE TABLE IF NOT EXISTS candidate_internal_ratings (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    company_id UUID NOT NULL REFERENCES companies(id) ON DELETE CASCADE,
    candidate_id UUID NOT NULL REFERENCES candidates(id) ON DELETE CASCADE,
    private_rating INT CHECK (private_rating BETWEEN 1 AND 5),
    recruiter_notes TEXT,
    created_at TIMESTAMPTZ DEFAULT now(),
    UNIQUE(company_id, candidate_id)
);

ALTER TABLE company_jobs
    ADD COLUMN IF NOT EXISTS currency_code VARCHAR(3) DEFAULT 'INR',
    ADD COLUMN IF NOT EXISTS work_location_type VARCHAR(20) DEFAULT 'ON_SITE' CHECK (work_location_type IN ('REMOTE', 'HYBRID', 'ON_SITE')),
    ADD COLUMN IF NOT EXISTS visa_sponsorship BOOLEAN DEFAULT false;
