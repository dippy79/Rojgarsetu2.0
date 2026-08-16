-- 000022_fulltext_search.down.sql
-- Rollback full-text search to trigger-based approach

-- Drop GENERATED columns and revert to trigger-based approach for jobs_government
ALTER TABLE jobs_government DROP COLUMN IF EXISTS search_vector;

-- Recreate search_vector column for manual updates
ALTER TABLE jobs_government ADD COLUMN IF NOT EXISTS search_vector tsvector;

-- Recreate function for jobs_government
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

-- Recreate trigger for jobs_government
CREATE TRIGGER trg_gov_jobs_search
    BEFORE INSERT OR UPDATE OF title, department, location, eligibility
    ON jobs_government
    FOR EACH ROW
    EXECUTE FUNCTION gov_jobs_search_update();

-- Recreate GIN index for jobs_government
CREATE INDEX IF NOT EXISTS idx_jobs_gov_search_vector ON jobs_government USING GIN(search_vector);

-- Drop GENERATED columns and revert to trigger-based approach for jobs_private
ALTER TABLE jobs_private DROP COLUMN IF EXISTS search_vector;

-- Recreate search_vector column for manual updates
ALTER TABLE jobs_private ADD COLUMN IF NOT EXISTS search_vector tsvector;

-- Recreate function for jobs_private
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

-- Recreate trigger for jobs_private
CREATE TRIGGER trg_priv_jobs_search
    BEFORE INSERT OR UPDATE OF title, company, location, description
    ON jobs_private
    FOR EACH ROW
    EXECUTE FUNCTION priv_jobs_search_update();

-- Recreate GIN index for jobs_private
CREATE INDEX IF NOT EXISTS idx_jobs_priv_search_vector ON jobs_private USING GIN(search_vector);

-- Backfill existing data
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