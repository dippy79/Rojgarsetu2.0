-- Fix schema to match sqlc generated code precisely
DROP TABLE IF EXISTS email_queue CASCADE;
DROP TABLE IF EXISTS platform_stats CASCADE;

-- Create email_queue table
CREATE TABLE IF NOT EXISTS email_queue (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    to_email TEXT NOT NULL,
    subject TEXT NOT NULL,
    body TEXT NOT NULL,
    status TEXT DEFAULT 'pending', -- pending, sent, failed
    attempts INTEGER DEFAULT 0,
    last_attempt_at TIMESTAMP WITH TIME ZONE,
    sent_at TIMESTAMP WITH TIME ZONE,
    error_message TEXT,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- Create platform_stats table (Single row singleton)
CREATE TABLE IF NOT EXISTS platform_stats (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    total_jobs BIGINT DEFAULT 0,
    total_candidates BIGINT DEFAULT 0,
    total_companies BIGINT DEFAULT 0,
    total_placements BIGINT DEFAULT 0,
    total_applications BIGINT DEFAULT 0,
    visits_today BIGINT DEFAULT 0,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- Insert initial row
INSERT INTO platform_stats (total_jobs, total_candidates, total_companies, total_placements, total_applications, visits_today)
VALUES (0, 0, 0, 0, 0, 0)
ON CONFLICT DO NOTHING;
