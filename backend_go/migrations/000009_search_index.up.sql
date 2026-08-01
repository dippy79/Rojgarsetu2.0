-- 000009_search_index.up.sql
-- Enable full-text search on jobs tables using pg_trgm for fuzzy matching
-- and tsvector for full-text search

CREATE EXTENSION IF NOT EXISTS pg_trgm;

-- Add tsvector columns for full-text search on company_jobs
ALTER TABLE company_jobs ADD COLUMN IF NOT EXISTS search_vector tsvector;

-- Add tsvector columns for full-text search on jobs_government
ALTER TABLE jobs_government ADD COLUMN IF NOT EXISTS search_vector tsvector;

-- Add tsvector columns for full-text search on jobs_private
ALTER TABLE jobs_private ADD COLUMN IF NOT EXISTS search_vector tsvector;

-- Create GIN indexes for fast full-text search
CREATE INDEX IF NOT EXISTS idx_company_jobs_search_vector ON company_jobs USING GIN(search_vector);
CREATE INDEX IF NOT EXISTS idx_jobs_gov_search_vector ON jobs_government USING GIN(search_vector);
CREATE INDEX IF NOT EXISTS idx_jobs_priv_search_vector ON jobs_private USING GIN(search_vector);

-- Create trigram GIN indexes for fuzzy LIKE/ILIKE matching on key text columns
CREATE INDEX IF NOT EXISTS idx_company_jobs_title_trgm ON company_jobs USING GIN(title gin_trgm_ops);
CREATE INDEX IF NOT EXISTS idx_company_jobs_location_trgm ON company_jobs USING GIN(location gin_trgm_ops);
CREATE INDEX IF NOT EXISTS idx_company_jobs_description_trgm ON company_jobs USING GIN(description gin_trgm_ops);

CREATE INDEX IF NOT EXISTS idx_jobs_gov_title_trgm ON jobs_government USING GIN(title gin_trgm_ops);
CREATE INDEX IF NOT EXISTS idx_jobs_gov_department_trgm ON jobs_government USING GIN(department gin_trgm_ops);
CREATE INDEX IF NOT EXISTS idx_jobs_gov_location_trgm ON jobs_government USING GIN(location gin_trgm_ops);

CREATE INDEX IF NOT EXISTS idx_jobs_priv_title_trgm ON jobs_private USING GIN(title gin_trgm_ops);
CREATE INDEX IF NOT EXISTS idx_jobs_priv_company_trgm ON jobs_private USING GIN(company gin_trgm_ops);
CREATE INDEX IF NOT EXISTS idx_jobs_priv_location_trgm ON jobs_private USING GIN(location gin_trgm_ops);

-- Create or replace function to update search_vector on company_jobs
CREATE OR REPLACE FUNCTION company_jobs_search_update() RETURNS trigger AS $$
BEGIN
    NEW.search_vector := to_tsvector('english', 
        COALESCE(NEW.title, '') || ' ' || 
        COALESCE(NEW.description, '') || ' ' || 
        COALESCE(NEW.location, '') || ' ' ||
        COALESCE(array_to_string(NEW.skills, ' '), '')
    );
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

-- Create or replace function to update search_vector on jobs_government
CREATE OR REPLACE FUNCTION gov_jobs_search_update() RETURNS trigger AS $$
BEGIN
    NEW.search_vector := to_tsvector('english', 
        COALESCE(NEW.title, '') || ' ' || 
        COALESCE(NEW.department, '') || ' ' || 
        COALESCE(NEW.location, '') || ' ' ||
        COALESCE(NEW.eligibility, '')
    );
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

-- Create or replace function to update search_vector on jobs_private
CREATE OR REPLACE FUNCTION priv_jobs_search_update() RETURNS trigger AS $$
BEGIN
    NEW.search_vector := to_tsvector('english', 
        COALESCE(NEW.title, '') || ' ' || 
        COALESCE(NEW.company, '') || ' ' || 
        COALESCE(NEW.location, '') || ' ' ||
        COALESCE(NEW.description, '')
    );
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

-- Drop existing triggers if they exist
DROP TRIGGER IF EXISTS trg_company_jobs_search ON company_jobs;
DROP TRIGGER IF EXISTS trg_gov_jobs_search ON jobs_government;
DROP TRIGGER IF EXISTS trg_priv_jobs_search ON jobs_private;

-- Create triggers to automatically update search_vector on INSERT or UPDATE
CREATE TRIGGER trg_company_jobs_search
    BEFORE INSERT OR UPDATE OF title, description, location, skills
    ON company_jobs
    FOR EACH ROW
    EXECUTE FUNCTION company_jobs_search_update();

CREATE TRIGGER trg_gov_jobs_search
    BEFORE INSERT OR UPDATE OF title, department, location, eligibility
    ON jobs_government
    FOR EACH ROW
    EXECUTE FUNCTION gov_jobs_search_update();

CREATE TRIGGER trg_priv_jobs_search
    BEFORE INSERT OR UPDATE OF title, company, location, description
    ON jobs_private
    FOR EACH ROW
    EXECUTE FUNCTION priv_jobs_search_update();

-- Backfill existing rows
UPDATE company_jobs SET search_vector = to_tsvector('english', 
    COALESCE(title, '') || ' ' || 
    COALESCE(description, '') || ' ' || 
    COALESCE(location, '') || ' ' ||
    COALESCE(array_to_string(skills, ' '), '')
) WHERE search_vector IS NULL;

UPDATE jobs_government SET search_vector = to_tsvector('english', 
    COALESCE(title, '') || ' ' || 
    COALESCE(department, '') || ' ' || 
    COALESCE(location, '') || ' ' ||
    COALESCE(eligibility, '')
) WHERE search_vector IS NULL;

UPDATE jobs_private SET search_vector = to_tsvector('english', 
    COALESCE(title, '') || ' ' || 
    COALESCE(company, '') || ' ' || 
    COALESCE(location, '') || ' ' ||
    COALESCE(description, '')
) WHERE search_vector IS NULL;

COMMENT ON TABLE company_jobs IS 'Private job postings with full-text search support';
COMMENT ON TABLE jobs_government IS 'Government job postings with full-text search support';
COMMENT ON TABLE jobs_private IS 'Private job listings with full-text search support';

