-- 000022_fulltext_search.up.sql
-- Convert full-text search to use GENERATED columns instead of triggers
-- This provides better performance and automatic maintenance

-- Drop existing triggers for government jobs
DROP TRIGGER IF EXISTS trg_gov_jobs_search ON jobs_government;
DROP FUNCTION IF EXISTS gov_jobs_search_update();

-- Convert jobs_government search_vector to GENERATED column
ALTER TABLE jobs_government 
  DROP COLUMN IF EXISTS search_vector;

ALTER TABLE jobs_government
  ADD COLUMN IF NOT EXISTS search_vector tsvector
  GENERATED ALWAYS AS (
    to_tsvector('english', 
      coalesce(title,'') || ' ' ||
      coalesce(department,'') || ' ' || 
      coalesce(eligibility,'') || ' ' ||
      coalesce(location,'')
    )
  ) STORED;

-- Recreate GIN index for jobs_government
DROP INDEX IF EXISTS idx_jobs_gov_search_vector;
CREATE INDEX IF NOT EXISTS idx_gov_jobs_fts
  ON jobs_government USING gin(search_vector);

-- Drop existing triggers for private jobs
DROP TRIGGER IF EXISTS trg_priv_jobs_search ON jobs_private;
DROP FUNCTION IF EXISTS priv_jobs_search_update();

-- Convert jobs_private search_vector to GENERATED column
ALTER TABLE jobs_private
  DROP COLUMN IF EXISTS search_vector;

ALTER TABLE jobs_private
  ADD COLUMN IF NOT EXISTS search_vector tsvector
  GENERATED ALWAYS AS (
    to_tsvector('english', 
      coalesce(title,'') || ' ' ||
      coalesce(location,'') || ' ' ||
      coalesce(company,'') || ' ' ||
      coalesce(description,'')
    )
  ) STORED;

-- Recreate GIN index for jobs_private
DROP INDEX IF EXISTS idx_jobs_priv_search_vector;
CREATE INDEX IF NOT EXISTS idx_priv_jobs_fts
  ON jobs_private USING gin(search_vector);

-- Note: company_jobs search_vector will continue using triggers for now
-- as it includes skills array which is more complex for GENERATED columns