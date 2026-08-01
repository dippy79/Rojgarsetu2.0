-- 000010_company_case_insensitive.up.sql
-- Phase 6: Company case-insensitive dedup + unique index
--
-- This migration:
-- 1. Finds all case-duplicate company groups (same LOWER(name))
-- 2. For each group, picks the earliest-created row as canonical
-- 3. Re-points all jobs (company_jobs table) referencing non-canonical
--    company IDs to the canonical company ID
-- 4. Deletes orphaned duplicate company rows
-- 5. Creates a UNIQUE index on LOWER(name) to prevent future duplicates
--
-- WARNING: This migration will DELETE duplicate company rows. Ensure
-- you have a backup before running on production.

BEGIN;

-- Step 1: Identify canonical company per case-insensitive group
-- Canonical = earliest created (oldest created_at), tie-breaking by ID
WITH canonical AS (
    SELECT DISTINCT ON (LOWER(name))
        id AS canonical_id,
        LOWER(name) AS normalized_name
    FROM companies
    ORDER BY LOWER(name), created_at ASC, id ASC
),
-- Step 2: Identify all non-canonical company IDs that need re-pointing
non_canonical AS (
    SELECT c.id AS duplicate_id, cn.canonical_id
    FROM companies c
    JOIN canonical cn ON LOWER(c.name) = cn.normalized_name
    WHERE c.id <> cn.canonical_id
),
-- Step 3: Update company_jobs to point to canonical company
updated_jobs AS (
    UPDATE company_jobs
    SET company_id = nc.canonical_id
    FROM non_canonical nc
    WHERE company_jobs.company_id = nc.duplicate_id
    RETURNING company_jobs.id
),
-- Step 4: Update the jobs table (crawler-scraped jobs) to re-point
-- to canonical company IDs (if they reference companies directly)
updated_jobs_table AS (
    UPDATE jobs
    SET company_id = nc.canonical_id
    FROM non_canonical nc
    WHERE jobs.company_id = nc.duplicate_id
    RETURNING jobs.id
)
-- Step 5: Delete the non-canonical duplicate company rows
DELETE FROM companies
WHERE id IN (SELECT duplicate_id FROM non_canonical);

-- Step 6: Create the unique index on LOWER(name) to prevent future duplicates
-- Using a unique index instead of a UNIQUE constraint because it allows
-- easier online creation and can be used by ON CONFLICT in INSERT statements.
CREATE UNIQUE INDEX IF NOT EXISTS idx_companies_name_lower_unique
    ON companies (LOWER(name));

COMMENT ON INDEX idx_companies_name_lower_unique IS
    'Prevents case-insensitive duplicate company names';

COMMIT;
