-- Crawler source registry
CREATE TABLE IF NOT EXISTS crawler_sources (
    id          SERIAL PRIMARY KEY,
    name        VARCHAR(100) NOT NULL,
    category    VARCHAR(50) NOT NULL,
    source_type VARCHAR(50) NOT NULL,
    base_url    TEXT NOT NULL,
    robots_txt_url TEXT,
    is_active   BOOLEAN DEFAULT true,
    created_at  TIMESTAMPTZ DEFAULT now()
);

-- Add missing columns to crawler_sources if table already exists
DO $$
BEGIN
    IF NOT EXISTS (SELECT 1 FROM information_schema.columns WHERE table_name = 'crawler_sources' AND column_name = 'category') THEN
        ALTER TABLE crawler_sources ADD COLUMN category VARCHAR(50) NOT NULL DEFAULT 'GOVT_JOB';
    END IF;
    IF NOT EXISTS (SELECT 1 FROM information_schema.columns WHERE table_name = 'crawler_sources' AND column_name = 'robots_txt_url') THEN
        ALTER TABLE crawler_sources ADD COLUMN robots_txt_url TEXT;
    END IF;
END $$;

-- Unified crawled jobs (links to existing jobs_government/jobs_private)
CREATE TABLE IF NOT EXISTS crawled_jobs (
    id                  SERIAL PRIMARY KEY,
    source_id           INT REFERENCES crawler_sources(id),
    job_type            VARCHAR(20) DEFAULT 'GOVT',
    title               VARCHAR(255) NOT NULL,
    company_or_dept     VARCHAR(255) NOT NULL,
    location            VARCHAR(255),
    qualification_req   TEXT,
    salary_or_pay_scale VARCHAR(100),
    apply_url           TEXT NOT NULL,
    source_attribution  VARCHAR(255) NOT NULL,
    hash_checksum       VARCHAR(64) UNIQUE NOT NULL,
    is_taken_down       BOOLEAN DEFAULT false,
    created_at          TIMESTAMPTZ DEFAULT now()
);

-- Add missing columns to crawled_jobs if table already exists
DO $$
BEGIN
    IF NOT EXISTS (SELECT 1 FROM information_schema.columns WHERE table_name = 'crawled_jobs' AND column_name = 'job_type') THEN
        ALTER TABLE crawled_jobs ADD COLUMN job_type VARCHAR(20) DEFAULT 'GOVT';
    END IF;
    IF NOT EXISTS (SELECT 1 FROM information_schema.columns WHERE table_name = 'crawled_jobs' AND column_name = 'qualification_req') THEN
        ALTER TABLE crawled_jobs ADD COLUMN qualification_req TEXT;
    END IF;
    IF NOT EXISTS (SELECT 1 FROM information_schema.columns WHERE table_name = 'crawled_jobs' AND column_name = 'salary_or_pay_scale') THEN
        ALTER TABLE crawled_jobs ADD COLUMN salary_or_pay_scale VARCHAR(100);
    END IF;
    IF NOT EXISTS (SELECT 1 FROM information_schema.columns WHERE table_name = 'crawled_jobs' AND column_name = 'apply_url') THEN
        ALTER TABLE crawled_jobs ADD COLUMN apply_url TEXT;
    END IF;
    IF NOT EXISTS (SELECT 1 FROM information_schema.columns WHERE table_name = 'crawled_jobs' AND column_name = 'source_attribution') THEN
        ALTER TABLE crawled_jobs ADD COLUMN source_attribution VARCHAR(255);
    END IF;
    IF NOT EXISTS (SELECT 1 FROM information_schema.columns WHERE table_name = 'crawled_jobs' AND column_name = 'is_taken_down') THEN
        ALTER TABLE crawled_jobs ADD COLUMN is_taken_down BOOLEAN DEFAULT false;
    END IF;
    -- Rename company to company_or_dept if it exists
    IF EXISTS (SELECT 1 FROM information_schema.columns WHERE table_name = 'crawled_jobs' AND column_name = 'company') AND
       NOT EXISTS (SELECT 1 FROM information_schema.columns WHERE table_name = 'crawled_jobs' AND column_name = 'company_or_dept') THEN
        ALTER TABLE crawled_jobs RENAME COLUMN company TO company_or_dept;
    END IF;
END $$;

CREATE INDEX IF NOT EXISTS idx_crawled_jobs_hash ON crawled_jobs(hash_checksum);
CREATE INDEX IF NOT EXISTS idx_crawled_jobs_type ON crawled_jobs(job_type, is_taken_down);

-- Govt forms & admit cards
CREATE TABLE IF NOT EXISTS gov_forms_info (
    id                SERIAL PRIMARY KEY,
    source_id         INT REFERENCES crawler_sources(id),
    title             VARCHAR(255) NOT NULL,
    conducting_body   VARCHAR(255) NOT NULL,
    form_type         VARCHAR(50) NOT NULL,
    official_website  TEXT NOT NULL,
    notification_pdf  TEXT,
    hash_checksum     VARCHAR(64) UNIQUE NOT NULL,
    is_taken_down     BOOLEAN DEFAULT false,
    created_at        TIMESTAMPTZ DEFAULT now()
);

-- Crawler telemetry (extends existing analytics)
CREATE TABLE IF NOT EXISTS crawler_logs (
    id               SERIAL PRIMARY KEY,
    source_id        INT REFERENCES crawler_sources(id),
    jobs_found       INT DEFAULT 0,
    jobs_added       INT DEFAULT 0,
    duplicates_found INT DEFAULT 0,
    status           VARCHAR(50) NOT NULL,
    error_message    TEXT,
    created_at       TIMESTAMPTZ DEFAULT now()
);

-- Takedown requests (IT Act Section 79)
CREATE TABLE IF NOT EXISTS takedown_requests (
    id           UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    job_id       INT,
    form_id      INT,
    requester    TEXT NOT NULL,
    reason       TEXT NOT NULL,
    status       VARCHAR(20) DEFAULT 'PENDING',
    created_at   TIMESTAMPTZ DEFAULT now(),
    resolved_at  TIMESTAMPTZ
);

-- Seed initial sources (using ON CONFLICT DO NOTHING)
INSERT INTO crawler_sources (name, category, source_type, base_url, robots_txt_url)
SELECT 'UPSC', 'GOVT_JOB', 'html_upsc', 'https://www.upsc.gov.in', 'https://www.upsc.gov.in/robots.txt'
WHERE NOT EXISTS (SELECT 1 FROM crawler_sources WHERE name = 'UPSC');

INSERT INTO crawler_sources (name, category, source_type, base_url, robots_txt_url)
SELECT 'SSC', 'GOVT_JOB', 'html_ssc', 'https://ssc.gov.in', 'https://ssc.gov.in/robots.txt'
WHERE NOT EXISTS (SELECT 1 FROM crawler_sources WHERE name = 'SSC');

INSERT INTO crawler_sources (name, category, source_type, base_url, robots_txt_url)
SELECT 'Railway RRB', 'GOVT_JOB', 'html_rrb', 'https://www.rrbapply.gov.in', 'https://www.rrbapply.gov.in/robots.txt'
WHERE NOT EXISTS (SELECT 1 FROM crawler_sources WHERE name = 'Railway RRB');

INSERT INTO crawler_sources (name, category, source_type, base_url, robots_txt_url)
SELECT 'NCS Portal', 'GOVT_JOB', 'html_ncs', 'https://www.ncs.gov.in', 'https://www.ncs.gov.in/robots.txt'
WHERE NOT EXISTS (SELECT 1 FROM crawler_sources WHERE name = 'NCS Portal');

INSERT INTO crawler_sources (name, category, source_type, base_url, robots_txt_url)
SELECT 'Adzuna API', 'PRIVATE_JOB', 'api_adzuna', 'https://api.adzuna.com/v1/api/jobs/in/search', NULL
WHERE NOT EXISTS (SELECT 1 FROM crawler_sources WHERE name = 'Adzuna API');

INSERT INTO crawler_sources (name, category, source_type, base_url, robots_txt_url)
SELECT 'Jooble API', 'PRIVATE_JOB', 'api_jooble', 'https://jooble.org/api', NULL
WHERE NOT EXISTS (SELECT 1 FROM crawler_sources WHERE name = 'Jooble API');
