-- 000009_search_index.down.sql
-- Rollback full-text search changes

DROP TRIGGER IF EXISTS trg_company_jobs_search ON company_jobs;
DROP TRIGGER IF EXISTS trg_gov_jobs_search ON jobs_government;
DROP TRIGGER IF EXISTS trg_priv_jobs_search ON jobs_private;

DROP FUNCTION IF EXISTS company_jobs_search_update();
DROP FUNCTION IF EXISTS gov_jobs_search_update();
DROP FUNCTION IF EXISTS priv_jobs_search_update();

DROP INDEX IF EXISTS idx_company_jobs_search_vector;
DROP INDEX IF EXISTS idx_jobs_gov_search_vector;
DROP INDEX IF EXISTS idx_jobs_priv_search_vector;

DROP INDEX IF EXISTS idx_company_jobs_title_trgm;
DROP INDEX IF EXISTS idx_company_jobs_location_trgm;
DROP INDEX IF EXISTS idx_company_jobs_description_trgm;

DROP INDEX IF EXISTS idx_jobs_gov_title_trgm;
DROP INDEX IF EXISTS idx_jobs_gov_department_trgm;
DROP INDEX IF EXISTS idx_jobs_gov_location_trgm;

DROP INDEX IF EXISTS idx_jobs_priv_title_trgm;
DROP INDEX IF EXISTS idx_jobs_priv_company_trgm;
DROP INDEX IF EXISTS idx_jobs_priv_location_trgm;

ALTER TABLE company_jobs DROP COLUMN IF EXISTS search_vector;
ALTER TABLE jobs_government DROP COLUMN IF EXISTS search_vector;
ALTER TABLE jobs_private DROP COLUMN IF EXISTS search_vector;

DROP TABLE IF EXISTS jobs_private;
