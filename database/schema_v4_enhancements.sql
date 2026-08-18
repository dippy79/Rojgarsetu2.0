-- RojgarSetu 2.0 - Schema v4.0 - ANTI-FAKE + PROFILE ENHANCEMENTS

-- ============================================
-- ENHANCE GOVERNMENT JOBS
-- ============================================
ALTER TABLE jobs_government ADD COLUMN IF NOT EXISTS is_verified BOOLEAN DEFAULT false;
ALTER TABLE jobs_government ADD COLUMN IF NOT EXISTS verification_meta JSONB;
ALTER TABLE jobs_government ADD COLUMN IF NOT EXISTS scam_score DOUBLE PRECISION DEFAULT 0.0;
ALTER TABLE jobs_government ADD COLUMN IF NOT EXISTS canonical_url TEXT;

-- ============================================
-- ENHANCE PRIVATE JOBS
-- ============================================
ALTER TABLE jobs_private ADD COLUMN IF NOT EXISTS is_verified BOOLEAN DEFAULT false;
ALTER TABLE jobs_private ADD COLUMN IF NOT EXISTS verification_meta JSONB;
ALTER TABLE jobs_private ADD COLUMN IF NOT EXISTS scam_score DOUBLE PRECISION DEFAULT 0.0;
ALTER TABLE jobs_private ADD COLUMN IF NOT EXISTS canonical_url TEXT;

-- ============================================
-- ENHANCE CANDIDATES (Profiles)
-- ============================================
ALTER TABLE candidates ADD COLUMN IF NOT EXISTS portfolio_links JSONB DEFAULT '[]';
ALTER TABLE candidates ADD COLUMN IF NOT EXISTS skill_graph JSONB DEFAULT '{}';
ALTER TABLE candidates ADD COLUMN IF NOT EXISTS education_history JSONB DEFAULT '[]';
ALTER TABLE candidates ADD COLUMN IF NOT EXISTS work_history JSONB DEFAULT '[]';
ALTER TABLE candidates ADD COLUMN IF NOT EXISTS certifications JSONB DEFAULT '[]';
ALTER TABLE candidates ADD COLUMN IF NOT EXISTS match_preferences JSONB DEFAULT '{}';

-- ============================================
-- ENHANCE COMPANIES (Profiles)
-- ============================================
ALTER TABLE companies ADD COLUMN IF NOT EXISTS employer_badge TEXT; -- 'verified', 'gold', 'premium'
ALTER TABLE companies ADD COLUMN IF NOT EXISTS ats_integration_meta JSONB DEFAULT '{}';
ALTER TABLE companies ADD COLUMN IF NOT EXISTS social_links JSONB DEFAULT '{}';
ALTER TABLE companies ADD COLUMN IF NOT EXISTS gallery_urls TEXT[] DEFAULT '{}';

-- ============================================
-- ENHANCE COMPANY JOBS (Analytics)
-- ============================================
ALTER TABLE company_jobs ADD COLUMN IF NOT EXISTS job_analytics JSONB DEFAULT '{"impressions": 0, "clicks": 0, "applications": 0}';
ALTER TABLE company_jobs ADD COLUMN IF NOT EXISTS auto_renewal BOOLEAN DEFAULT false;

-- ============================================
-- ANTI-SPAM LOGS
-- ============================================
CREATE TABLE IF NOT EXISTS scam_alerts (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    target_type TEXT NOT NULL, -- 'job_gov', 'job_priv', 'company'
    target_id UUID NOT NULL,
    alert_reason TEXT,
    severity TEXT, -- 'low', 'medium', 'high', 'critical'
    detected_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    is_resolved BOOLEAN DEFAULT false,
    resolved_by UUID REFERENCES users(id),
    resolved_at TIMESTAMP
);

CREATE INDEX idx_scam_alerts_target ON scam_alerts(target_type, target_id);
CREATE INDEX idx_scam_alerts_severity ON scam_alerts(severity);

-- ============================================
-- SKILLS REGISTRY (Global)
-- ============================================
CREATE TABLE IF NOT EXISTS skill_registry (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    name TEXT UNIQUE NOT NULL,
    category TEXT,
    demand_score INTEGER DEFAULT 0,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_skill_registry_name ON skill_registry(name);
