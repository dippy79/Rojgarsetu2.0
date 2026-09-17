-- Drop and Recreate to fix schema mismatch with sqlc queries
DROP TABLE IF EXISTS email_queue;
DROP TABLE IF EXISTS platform_stats;

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

-- Create platform_stats table
CREATE TABLE IF NOT EXISTS platform_stats (
    id SERIAL PRIMARY KEY,
    total_jobs INTEGER DEFAULT 0,
    total_candidates INTEGER DEFAULT 0,
    total_companies INTEGER DEFAULT 0,
    total_placements INTEGER DEFAULT 0,
    total_applications INTEGER DEFAULT 0,
    visits_today INTEGER DEFAULT 0,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- Insert initial row if not exists
INSERT INTO platform_stats (id, total_jobs, total_candidates, total_companies, total_placements, total_applications, visits_today)
VALUES (1, 0, 0, 0, 0, 0, 0)
ON CONFLICT (id) DO NOTHING;
