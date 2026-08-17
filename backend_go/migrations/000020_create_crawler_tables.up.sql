-- 1. Crawler Sources Table
CREATE TABLE IF NOT EXISTS crawler_sources (
    id SERIAL PRIMARY KEY,
    name VARCHAR(100) NOT NULL UNIQUE,
    source_type VARCHAR(50) NOT NULL,
    base_url TEXT NOT NULL,
    is_active BOOLEAN DEFAULT true,
    last_crawled_at TIMESTAMP WITH TIME ZONE,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

INSERT INTO crawler_sources (name, source_type, base_url, is_active) VALUES
('NCS', 'API', 'https://www.ncs.gov.in', true),
('UPSC', 'RSS', 'https://upsc.gov.in', true),
('SSC', 'RSS', 'https://ssc.gov.in', true),
('Railway', 'HTML', 'https://indianrailways.gov.in', true)
ON CONFLICT (name) DO NOTHING;

-- 2. Crawled Jobs Table
CREATE TABLE IF NOT EXISTS crawled_jobs (
    id BIGSERIAL PRIMARY KEY,
    source_id INT REFERENCES crawler_sources(id) ON DELETE SET NULL,
    external_job_id VARCHAR(255) NOT NULL,
    job_hash VARCHAR(64) UNIQUE NOT NULL,
    title VARCHAR(255) NOT NULL,
    organization VARCHAR(255) NOT NULL,
    job_type VARCHAR(50) DEFAULT 'GOVT',
    category_id INT REFERENCES job_categories(id) ON DELETE SET NULL,
    trade_id INT REFERENCES trades(id) ON DELETE SET NULL,
    qualification_required VARCHAR(255),
    total_vacancies INT DEFAULT 1,
    salary_range VARCHAR(100),
    job_location VARCHAR(255),
    official_notification_url TEXT,
    apply_url TEXT NOT NULL,
    published_at TIMESTAMP WITH TIME ZONE,
    application_deadline TIMESTAMP WITH TIME ZONE,
    status VARCHAR(50) DEFAULT 'ACTIVE',
    raw_payload JSONB,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX IF NOT EXISTS idx_crawled_jobs_hash ON crawled_jobs(job_hash);
CREATE INDEX IF NOT EXISTS idx_crawled_jobs_status ON crawled_jobs(status);
CREATE INDEX IF NOT EXISTS idx_crawled_jobs_trade ON crawled_jobs(trade_id);

-- 3. Crawler Logs Table
CREATE TABLE IF NOT EXISTS crawler_logs (
    id BIGSERIAL PRIMARY KEY,
    source_name VARCHAR(100) NOT NULL,
    status VARCHAR(50) NOT NULL,
    jobs_found INT DEFAULT 0,
    jobs_added INT DEFAULT 0,
    duplicates_found INT DEFAULT 0,
    errors_count INT DEFAULT 0,
    error_message TEXT,
    execution_time_ms INT DEFAULT 0,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);