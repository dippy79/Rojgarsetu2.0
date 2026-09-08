-- Migration 20: Harmonizing with 19 (which now handles base crawler tables correctly)
-- Ensuring indices and constraints that might be missing or need specific naming

-- Note: Tables crawler_sources, crawled_jobs, crawler_logs are now correctly defined in 19.

-- Add missing status column if not present
DO $$
BEGIN
    IF NOT EXISTS (SELECT 1 FROM information_schema.columns WHERE table_name='crawled_jobs' AND column_name='status') THEN
        ALTER TABLE crawled_jobs ADD COLUMN status VARCHAR(20) DEFAULT 'ACTIVE';
    END IF;
END $$;

-- Add specific indices for performance if not already there
CREATE INDEX IF NOT EXISTS idx_crawled_jobs_status ON crawled_jobs(status) WHERE status = 'ACTIVE';
CREATE INDEX IF NOT EXISTS idx_crawled_jobs_source_id ON crawled_jobs(source_id);
