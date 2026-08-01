-- 000010_company_case_insensitive.down.sql
-- Drop the case-insensitive unique index on companies(name)
DROP INDEX IF EXISTS idx_companies_name_lower_unique;
