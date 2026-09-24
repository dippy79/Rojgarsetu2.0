ALTER TABLE crawler_logs DROP COLUMN IF EXISTS source;
ALTER TABLE crawler_logs DROP COLUMN IF EXISTS errors;
ALTER TABLE crawler_logs DROP COLUMN IF EXISTS jobs_saved;
ALTER TABLE crawler_logs DROP COLUMN IF EXISTS started_at;
ALTER TABLE crawler_logs DROP COLUMN IF EXISTS completed_at;
