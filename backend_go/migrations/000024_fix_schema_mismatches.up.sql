-- 000024_fix_schema_mismatches.up.sql

-- 1. Add job_hash to jobs_private (required for deduplication)
ALTER TABLE jobs_private
ADD COLUMN IF NOT EXISTS job_hash VARCHAR(64);

-- Create a unique constraint on job_hash for UPSERT operations
DO $$
BEGIN
    IF NOT EXISTS (SELECT 1 FROM pg_constraint WHERE conname = 'jobs_private_job_hash_key') THEN
        ALTER TABLE jobs_private ADD CONSTRAINT jobs_private_job_hash_key UNIQUE (job_hash);
    END IF;
END $$;

-- 2. Fix crawler_logs table to match the Go code expectations
-- Rename source_name to source if it exists
DO $$
BEGIN
    IF EXISTS (SELECT 1 FROM information_schema.columns WHERE table_name = 'crawler_logs' AND column_name = 'source_name') THEN
        ALTER TABLE crawler_logs RENAME COLUMN source_name TO source;
    END IF;
END $$;

-- Ensure other columns expected by Go code exist
ALTER TABLE crawler_logs ADD COLUMN IF NOT EXISTS jobs_saved INT DEFAULT 0;
ALTER TABLE crawler_logs ADD COLUMN IF NOT EXISTS errors TEXT;
ALTER TABLE crawler_logs ADD COLUMN IF NOT EXISTS started_at TIMESTAMP WITH TIME ZONE DEFAULT NOW();
ALTER TABLE crawler_logs ADD COLUMN IF NOT EXISTS completed_at TIMESTAMP WITH TIME ZONE DEFAULT NOW();

-- Rename existing columns to match code if necessary (though keeping duplicates_found etc is fine if code doesn't use them)
-- But the code uses: source, status, jobs_found, jobs_saved, errors, started_at, completed_at
-- Our table has: source (was source_name), status, jobs_found, jobs_added (should be jobs_saved?), duplicates_found, errors_count, error_message (should be errors?), execution_time_ms, created_at

-- Let's align exactly with what store.go SaveLog uses:
-- VALUES ('all_sources', $1, $2, $3, $4, NOW() - INTERVAL '1 second', NOW())
-- $1=status, $2=found, $3=saved, $4=errMsg

DO $$
BEGIN
    IF EXISTS (SELECT 1 FROM information_schema.columns WHERE table_name = 'crawler_logs' AND column_name = 'jobs_added') AND
       NOT EXISTS (SELECT 1 FROM information_schema.columns WHERE table_name = 'crawler_logs' AND column_name = 'jobs_saved') THEN
        ALTER TABLE crawler_logs RENAME COLUMN jobs_added TO jobs_saved;
    END IF;

    IF EXISTS (SELECT 1 FROM information_schema.columns WHERE table_name = 'crawler_logs' AND column_name = 'error_message') AND
       NOT EXISTS (SELECT 1 FROM information_schema.columns WHERE table_name = 'crawler_logs' AND column_name = 'errors') THEN
        ALTER TABLE crawler_logs RENAME COLUMN error_message TO errors;
    END IF;
END $$;
