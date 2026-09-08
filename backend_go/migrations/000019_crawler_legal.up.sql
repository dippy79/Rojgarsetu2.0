-- Crawler source registry
CREATE TABLE IF NOT EXISTS crawler_sources (
    id          SERIAL PRIMARY KEY,
    name        VARCHAR(100) NOT NULL UNIQUE,
    category    VARCHAR(50) NOT NULL DEFAULT 'GOVT_JOB',
    source_type VARCHAR(50) NOT NULL,
    base_url    TEXT NOT NULL,
    robots_txt_url TEXT,
    is_active   BOOLEAN DEFAULT true,
    last_crawled_at TIMESTAMP WITH TIME ZONE,
    created_at  TIMESTAMPTZ DEFAULT now()
);

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
    status              VARCHAR(50) DEFAULT 'ACTIVE',
    is_taken_down       BOOLEAN DEFAULT false,
    created_at          TIMESTAMPTZ DEFAULT now(),
    updated_at          TIMESTAMPTZ DEFAULT now()
);

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

-- Crawler telemetry
CREATE TABLE IF NOT EXISTS crawler_logs (
    id               SERIAL PRIMARY KEY,
    source_id        INT REFERENCES crawler_sources(id),
    jobs_found       INT DEFAULT 0,
    jobs_added       INT DEFAULT 0,
    duplicates_found INT DEFAULT 0,
    status           VARCHAR(50) NOT NULL,
    error_message    TEXT,
    execution_time_ms INT DEFAULT 0,
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

-- Seed initial sources
INSERT INTO crawler_sources (name, category, source_type, base_url, robots_txt_url) VALUES
('UPSC', 'GOVT_JOB', 'html_upsc', 'https://www.upsc.gov.in', 'https://www.upsc.gov.in/robots.txt'),
('SSC', 'GOVT_JOB', 'html_ssc', 'https://ssc.gov.in', 'https://ssc.gov.in/robots.txt'),
('Railway RRB', 'GOVT_JOB', 'html_rrb', 'https://www.rrbapply.gov.in', 'https://www.rrbapply.gov.in/robots.txt'),
('NCS Portal', 'GOVT_JOB', 'html_ncs', 'https://www.ncs.gov.in', 'https://www.ncs.gov.in/robots.txt'),
('Adzuna API', 'PRIVATE_JOB', 'api_adzuna', 'https://api.adzuna.com/v1/api/jobs/in/search', NULL),
('Jooble API', 'PRIVATE_JOB', 'api_jooble', 'https://jooble.org/api', NULL)
ON CONFLICT (name) DO NOTHING;
