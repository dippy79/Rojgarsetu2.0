-- Align candidates table with Go models.go and fix registration crash
ALTER TABLE candidates ADD COLUMN IF NOT EXISTS resume_parsed JSONB DEFAULT '{}'::jsonb;

-- Rename experience to experience_years if needed to match code expectations,
-- or add it as a new column to avoid breaking old data.
DO $$
BEGIN
    IF NOT EXISTS (SELECT 1 FROM information_schema.columns WHERE table_name='candidates' AND column_name='experience_years') THEN
        ALTER TABLE candidates ADD COLUMN experience_years INTEGER DEFAULT 0;
    END IF;
END $$;
