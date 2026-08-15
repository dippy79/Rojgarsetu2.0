-- 00013_create_job_categories_and_trades.up.sql
-- Creates job_categories, job_trades, user_enrollments, and user_notification_logs tables
-- for the Trade Categories & Enrollment Notification System

BEGIN;

-- Function used by updated_at triggers (defensive recreation)
CREATE OR REPLACE FUNCTION update_updated_at_column()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = NOW();
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

-- ============================================
-- JOB CATEGORIES TABLE (6 Sectors)
-- ============================================
CREATE TABLE IF NOT EXISTS job_categories (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    name TEXT NOT NULL UNIQUE,
    slug TEXT NOT NULL UNIQUE,
    description TEXT,
    icon TEXT,
    color TEXT,
    display_order INT4 NOT NULL DEFAULT 0,
    is_active BOOLEAN NOT NULL DEFAULT true,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_job_categories_slug ON job_categories(slug);
CREATE INDEX IF NOT EXISTS idx_job_categories_is_active ON job_categories(is_active);
CREATE INDEX IF NOT EXISTS idx_job_categories_display_order ON job_categories(display_order);

-- ============================================
-- JOB TRADES TABLE (Specific trades within categories)
-- ============================================
CREATE TABLE IF NOT EXISTS job_trades (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    category_id UUID NOT NULL REFERENCES job_categories(id) ON DELETE CASCADE,
    name TEXT NOT NULL,
    slug TEXT NOT NULL,
    description TEXT,
    qualification_req TEXT,
    min_salary INT4,
    max_salary INT4,
    demand_level TEXT NOT NULL DEFAULT 'NORMAL' CHECK (demand_level IN ('CRITICAL', 'HIGH', 'MEDIUM', 'NORMAL')),
    icon TEXT,
    is_active BOOLEAN NOT NULL DEFAULT true,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    UNIQUE(category_id, slug)
);

CREATE INDEX IF NOT EXISTS idx_job_trades_category_id ON job_trades(category_id);
CREATE INDEX IF NOT EXISTS idx_job_trades_slug ON job_trades(slug);
CREATE INDEX IF NOT EXISTS idx_job_trades_demand_level ON job_trades(demand_level);
CREATE INDEX IF NOT EXISTS idx_job_trades_is_active ON job_trades(is_active);

-- ============================================
-- USER ENROLLMENTS TABLE
-- ============================================
CREATE TABLE IF NOT EXISTS user_enrollments (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    trade_id UUID NOT NULL REFERENCES job_trades(id) ON DELETE CASCADE,
    status TEXT NOT NULL DEFAULT 'active' CHECK (status IN ('active', 'completed', 'expired', 'cancelled')),
    enrolled_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    expires_at TIMESTAMPTZ NOT NULL,
    completed_at TIMESTAMPTZ,
    progress_pct INT4 NOT NULL DEFAULT 0 CHECK (progress_pct >= 0 AND progress_pct <= 100),
    metadata JSONB NOT NULL DEFAULT '{}',
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    UNIQUE(user_id, trade_id, status)
);

CREATE INDEX IF NOT EXISTS idx_user_enrollments_user_id ON user_enrollments(user_id);
CREATE INDEX IF NOT EXISTS idx_user_enrollments_trade_id ON user_enrollments(trade_id);
CREATE INDEX IF NOT EXISTS idx_user_enrollments_status ON user_enrollments(status);
CREATE INDEX IF NOT EXISTS idx_user_enrollments_expires_at ON user_enrollments(expires_at);
CREATE INDEX IF NOT EXISTS idx_user_enrollments_user_status_expires ON user_enrollments(user_id, status, expires_at);

-- ============================================
-- USER NOTIFICATION LOGS TABLE (for 2/day limit enforcement)
-- ============================================
CREATE TABLE IF NOT EXISTS user_notification_logs (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    enrollment_id UUID REFERENCES user_enrollments(id) ON DELETE SET NULL,
    notification_type TEXT NOT NULL CHECK (notification_type IN ('expiry_warning', 'expiry_final', 'enrollment_reminder', 'course_update')),
    channel TEXT NOT NULL DEFAULT 'in_app' CHECK (channel IN ('in_app', 'email', 'push', 'sms')),
    title TEXT NOT NULL,
    message TEXT NOT NULL,
    payload JSONB NOT NULL DEFAULT '{}',
    sent_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    read_at TIMESTAMPTZ,
    clicked_at TIMESTAMPTZ,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_user_notification_logs_user_id ON user_notification_logs(user_id);
CREATE INDEX IF NOT EXISTS idx_user_notification_logs_enrollment_id ON user_notification_logs(enrollment_id);
CREATE INDEX IF NOT EXISTS idx_user_notification_logs_sent_at ON user_notification_logs(sent_at DESC);
-- ✅ UPDATED: Composite IMMUTABLE Index replacing DATE(sent_at)
CREATE INDEX IF NOT EXISTS idx_user_notification_logs_user_sent_at ON user_notification_logs(user_id, sent_at DESC);
CREATE INDEX IF NOT EXISTS idx_user_notification_logs_type ON user_notification_logs(notification_type);

-- ============================================
-- UPDATED_AT TRIGGERS
-- ============================================
DROP TRIGGER IF EXISTS update_job_categories_updated_at ON job_categories;
CREATE TRIGGER update_job_categories_updated_at
    BEFORE UPDATE ON job_categories
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

DROP TRIGGER IF EXISTS update_job_trades_updated_at ON job_trades;
CREATE TRIGGER update_job_trades_updated_at
    BEFORE UPDATE ON job_trades
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

DROP TRIGGER IF EXISTS update_user_enrollments_updated_at ON user_enrollments;
CREATE TRIGGER update_user_enrollments_updated_at
    BEFORE UPDATE ON user_enrollments
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

-- ============================================
-- SEED DATA: 6 SECTORS (Job Categories)
-- ============================================
INSERT INTO job_categories (name, slug, description, icon, color, display_order) VALUES
    ('Skilled Trades', 'skilled-trades', 'Hands-on technical trades requiring specialized training and certification', 'wrench', '#10B981', 1),
    ('Construction', 'construction', 'Building and infrastructure trades including carpentry, masonry, and electrical work', 'building', '#F59E0B', 2),
    ('Transport', 'transport', 'Vehicle operation, logistics, and transportation management roles', 'truck', '#3B82F6', 3),
    ('Service', 'service', 'Customer-facing service roles in hospitality, retail, and personal care', 'users', '#EC4899', 4),
    ('Healthcare', 'healthcare', 'Medical and healthcare support roles requiring certification', 'heart-pulse', '#EF4444', 5),
    ('Creative', 'creative', 'Design, media, arts, and creative technology roles', 'palette', '#8B5CF6', 6)
ON CONFLICT (slug) DO UPDATE SET
    name = EXCLUDED.name,
    description = EXCLUDED.description,
    icon = EXCLUDED.icon,
    color = EXCLUDED.color,
    display_order = EXCLUDED.display_order,
    updated_at = NOW();

-- ============================================
-- SEED DATA: Job Trades for each category
-- ============================================

-- Skilled Trades
INSERT INTO job_trades (category_id, name, slug, description, qualification_req, min_salary, max_salary, demand_level, icon) 
SELECT id, 'Electrician', 'electrician', 'Install, maintain, and repair electrical systems', 'ITI Electrician / Diploma in Electrical Engineering', 250000, 600000, 'CRITICAL', 'zap'
FROM job_categories WHERE slug = 'skilled-trades'
ON CONFLICT (category_id, slug) DO UPDATE SET
    name = EXCLUDED.name,
    description = EXCLUDED.description,
    qualification_req = EXCLUDED.qualification_req,
    min_salary = EXCLUDED.min_salary,
    max_salary = EXCLUDED.max_salary,
    demand_level = EXCLUDED.demand_level,
    icon = EXCLUDED.icon,
    updated_at = NOW();

INSERT INTO job_trades (category_id, name, slug, description, qualification_req, min_salary, max_salary, demand_level, icon) 
SELECT id, 'Plumber', 'plumber', 'Install and repair water, gas, and drainage systems', 'ITI Plumber / Diploma in Plumbing Technology', 200000, 500000, 'HIGH', 'wrench'
FROM job_categories WHERE slug = 'skilled-trades'
ON CONFLICT (category_id, slug) DO UPDATE SET
    name = EXCLUDED.name,
    description = EXCLUDED.description,
    qualification_req = EXCLUDED.qualification_req,
    min_salary = EXCLUDED.min_salary,
    max_salary = EXCLUDED.max_salary,
    demand_level = EXCLUDED.demand_level,
    icon = EXCLUDED.icon,
    updated_at = NOW();

INSERT INTO job_trades (category_id, name, slug, description, qualification_req, min_salary, max_salary, demand_level, icon) 
SELECT id, 'Welder', 'welder', 'Join metal parts using various welding techniques', 'ITI Welder / Certificate in Welding Technology', 180000, 450000, 'HIGH', 'flame'
FROM job_categories WHERE slug = 'skilled-trades'
ON CONFLICT (category_id, slug) DO UPDATE SET
    name = EXCLUDED.name,
    description = EXCLUDED.description,
    qualification_req = EXCLUDED.qualification_req,
    min_salary = EXCLUDED.min_salary,
    max_salary = EXCLUDED.max_salary,
    demand_level = EXCLUDED.demand_level,
    icon = EXCLUDED.icon,
    updated_at = NOW();

INSERT INTO job_trades (category_id, name, slug, description, qualification_req, min_salary, max_salary, demand_level, icon) 
SELECT id, 'HVAC Technician', 'hvac-technician', 'Install and maintain heating, ventilation, and air conditioning systems', 'ITI Refrigeration & AC / Diploma in HVAC', 220000, 550000, 'MEDIUM', 'fan'
FROM job_categories WHERE slug = 'skilled-trades'
ON CONFLICT (category_id, slug) DO UPDATE SET
    name = EXCLUDED.name,
    description = EXCLUDED.description,
    qualification_req = EXCLUDED.qualification_req,
    min_salary = EXCLUDED.min_salary,
    max_salary = EXCLUDED.max_salary,
    demand_level = EXCLUDED.demand_level,
    icon = EXCLUDED.icon,
    updated_at = NOW();

-- Construction
INSERT INTO job_trades (category_id, name, slug, description, qualification_req, min_salary, max_salary, demand_level, icon) 
SELECT id, 'Carpenter', 'carpenter', 'Construct, install, and repair wooden structures', 'ITI Carpenter / Certificate in Carpentry', 180000, 450000, 'HIGH', 'hammer'
FROM job_categories WHERE slug = 'construction'
ON CONFLICT (category_id, slug) DO UPDATE SET
    name = EXCLUDED.name,
    description = EXCLUDED.description,
    qualification_req = EXCLUDED.qualification_req,
    min_salary = EXCLUDED.min_salary,
    max_salary = EXCLUDED.max_salary,
    demand_level = EXCLUDED.demand_level,
    icon = EXCLUDED.icon,
    updated_at = NOW();

INSERT INTO job_trades (category_id, name, slug, description, qualification_req, min_salary, max_salary, demand_level, icon) 
SELECT id, 'Mason', 'mason', 'Build structures using bricks, concrete blocks, and stone', 'ITI Mason / Certificate in Masonry', 160000, 400000, 'MEDIUM', 'brick'
FROM job_categories WHERE slug = 'construction'
ON CONFLICT (category_id, slug) DO UPDATE SET
    name = EXCLUDED.name,
    description = EXCLUDED.description,
    qualification_req = EXCLUDED.qualification_req,
    min_salary = EXCLUDED.min_salary,
    max_salary = EXCLUDED.max_salary,
    demand_level = EXCLUDED.demand_level,
    icon = EXCLUDED.icon,
    updated_at = NOW();

INSERT INTO job_trades (category_id, name, slug, description, qualification_req, min_salary, max_salary, demand_level, icon) 
SELECT id, 'Construction Supervisor', 'construction-supervisor', 'Oversee construction projects and manage site operations', 'Diploma in Civil Engineering / B.Tech Civil', 350000, 800000, 'CRITICAL', 'hard-hat'
FROM job_categories WHERE slug = 'construction'
ON CONFLICT (category_id, slug) DO UPDATE SET
    name = EXCLUDED.name,
    description = EXCLUDED.description,
    qualification_req = EXCLUDED.qualification_req,
    min_salary = EXCLUDED.min_salary,
    max_salary = EXCLUDED.max_salary,
    demand_level = EXCLUDED.demand_level,
    icon = EXCLUDED.icon,
    updated_at = NOW();

INSERT INTO job_trades (category_id, name, slug, description, qualification_req, min_salary, max_salary, demand_level, icon) 
SELECT id, 'Scaffolder', 'scaffolder', 'Erect and dismantle scaffolding for construction work', 'Certificate in Scaffolding / ITI', 180000, 420000, 'MEDIUM', 'building-2'
FROM job_categories WHERE slug = 'construction'
ON CONFLICT (category_id, slug) DO UPDATE SET
    name = EXCLUDED.name,
    description = EXCLUDED.description,
    qualification_req = EXCLUDED.qualification_req,
    min_salary = EXCLUDED.min_salary,
    max_salary = EXCLUDED.max_salary,
    demand_level = EXCLUDED.demand_level,
    icon = EXCLUDED.icon,
    updated_at = NOW();

-- Transport
INSERT INTO job_trades (category_id, name, slug, description, qualification_req, min_salary, max_salary, demand_level, icon) 
SELECT id, 'Heavy Vehicle Driver', 'heavy-vehicle-driver', 'Operate trucks, trailers, and heavy commercial vehicles', 'Commercial Driving License (CDL) / HMV License', 250000, 550000, 'CRITICAL', 'truck'
FROM job_categories WHERE slug = 'transport'
ON CONFLICT (category_id, slug) DO UPDATE SET
    name = EXCLUDED.name,
    description = EXCLUDED.description,
    qualification_req = EXCLUDED.qualification_req,
    min_salary = EXCLUDED.min_salary,
    max_salary = EXCLUDED.max_salary,
    demand_level = EXCLUDED.demand_level,
    icon = EXCLUDED.icon,
    updated_at = NOW();

INSERT INTO job_trades (category_id, name, slug, description, qualification_req, min_salary, max_salary, demand_level, icon) 
SELECT id, 'Light Vehicle Driver', 'light-vehicle-driver', 'Drive cars, vans, and light commercial vehicles', 'Valid Driving License (LMV)', 150000, 350000, 'HIGH', 'car'
FROM job_categories WHERE slug = 'transport'
ON CONFLICT (category_id, slug) DO UPDATE SET
    name = EXCLUDED.name,
    description = EXCLUDED.description,
    qualification_req = EXCLUDED.qualification_req,
    min_salary = EXCLUDED.min_salary,
    max_salary = EXCLUDED.max_salary,
    demand_level = EXCLUDED.demand_level,
    icon = EXCLUDED.icon,
    updated_at = NOW();

INSERT INTO job_trades (category_id, name, slug, description, qualification_req, min_salary, max_salary, demand_level, icon) 
SELECT id, 'Logistics Coordinator', 'logistics-coordinator', 'Coordinate supply chain and transportation logistics', 'Diploma in Logistics / Supply Chain Management', 250000, 600000, 'MEDIUM', 'package'
FROM job_categories WHERE slug = 'transport'
ON CONFLICT (category_id, slug) DO UPDATE SET
    name = EXCLUDED.name,
    description = EXCLUDED.description,
    qualification_req = EXCLUDED.qualification_req,
    min_salary = EXCLUDED.min_salary,
    max_salary = EXCLUDED.max_salary,
    demand_level = EXCLUDED.demand_level,
    icon = EXCLUDED.icon,
    updated_at = NOW();

INSERT INTO job_trades (category_id, name, slug, description, qualification_req, min_salary, max_salary, demand_level, icon) 
SELECT id, 'Forklift Operator', 'forklift-operator', 'Operate forklifts and material handling equipment', 'Forklift Operator Certificate', 180000, 400000, 'MEDIUM', 'forklift'
FROM job_categories WHERE slug = 'transport'
ON CONFLICT (category_id, slug) DO UPDATE SET
    name = EXCLUDED.name,
    description = EXCLUDED.description,
    qualification_req = EXCLUDED.qualification_req,
    min_salary = EXCLUDED.min_salary,
    max_salary = EXCLUDED.max_salary,
    demand_level = EXCLUDED.demand_level,
    icon = EXCLUDED.icon,
    updated_at = NOW();

-- Service
INSERT INTO job_trades (category_id, name, slug, description, qualification_req, min_salary, max_salary, demand_level, icon) 
SELECT id, 'Hotel Front Desk', 'hotel-front-desk', 'Manage guest check-in, check-out, and front desk operations', 'Diploma in Hotel Management / Hospitality', 180000, 400000, 'HIGH', 'concierge-bell'
FROM job_categories WHERE slug = 'service'
ON CONFLICT (category_id, slug) DO UPDATE SET
    name = EXCLUDED.name,
    description = EXCLUDED.description,
    qualification_req = EXCLUDED.qualification_req,
    min_salary = EXCLUDED.min_salary,
    max_salary = EXCLUDED.max_salary,
    demand_level = EXCLUDED.demand_level,
    icon = EXCLUDED.icon,
    updated_at = NOW();

INSERT INTO job_trades (category_id, name, slug, description, qualification_req, min_salary, max_salary, demand_level, icon) 
SELECT id, 'Retail Sales Associate', 'retail-sales-associate', 'Assist customers and manage sales in retail environments', '12th Pass / Diploma in Retail Management', 140000, 300000, 'MEDIUM', 'shopping-bag'
FROM job_categories WHERE slug = 'service'
ON CONFLICT (category_id, slug) DO UPDATE SET
    name = EXCLUDED.name,
    description = EXCLUDED.description,
    qualification_req = EXCLUDED.qualification_req,
    min_salary = EXCLUDED.min_salary,
    max_salary = EXCLUDED.max_salary,
    demand_level = EXCLUDED.demand_level,
    icon = EXCLUDED.icon,
    updated_at = NOW();

INSERT INTO job_trades (category_id, name, slug, description, qualification_req, min_salary, max_salary, demand_level, icon) 
SELECT id, 'Security Guard', 'security-guard', 'Protect property and ensure safety of premises', 'Security Guard Training Certificate', 150000, 300000, 'HIGH', 'shield'
FROM job_categories WHERE slug = 'service'
ON CONFLICT (category_id, slug) DO UPDATE SET
    name = EXCLUDED.name,
    description = EXCLUDED.description,
    qualification_req = EXCLUDED.qualification_req,
    min_salary = EXCLUDED.min_salary,
    max_salary = EXCLUDED.max_salary,
    demand_level = EXCLUDED.demand_level,
    icon = EXCLUDED.icon,
    updated_at = NOW();

INSERT INTO job_trades (category_id, name, slug, description, qualification_req, min_salary, max_salary, demand_level, icon) 
SELECT id, 'Beauty & Wellness Therapist', 'beauty-wellness-therapist', 'Provide beauty treatments and wellness services', 'Diploma in Cosmetology / Beauty Therapy', 160000, 400000, 'MEDIUM', 'sparkles'
FROM job_categories WHERE slug = 'service'
ON CONFLICT (category_id, slug) DO UPDATE SET
    name = EXCLUDED.name,
    description = EXCLUDED.description,
    qualification_req = EXCLUDED.qualification_req,
    min_salary = EXCLUDED.min_salary,
    max_salary = EXCLUDED.max_salary,
    demand_level = EXCLUDED.demand_level,
    icon = EXCLUDED.icon,
    updated_at = NOW();

-- Healthcare
INSERT INTO job_trades (category_id, name, slug, description, qualification_req, min_salary, max_salary, demand_level, icon) 
SELECT id, 'Nursing Assistant', 'nursing-assistant', 'Provide basic patient care under nursing supervision', 'ANM / GNM / Certificate in Nursing Assistant', 200000, 450000, 'CRITICAL', 'heart-pulse'
FROM job_categories WHERE slug = 'healthcare'
ON CONFLICT (category_id, slug) DO UPDATE SET
    name = EXCLUDED.name,
    description = EXCLUDED.description,
    qualification_req = EXCLUDED.qualification_req,
    min_salary = EXCLUDED.min_salary,
    max_salary = EXCLUDED.max_salary,
    demand_level = EXCLUDED.demand_level,
    icon = EXCLUDED.icon,
    updated_at = NOW();

INSERT INTO job_trades (category_id, name, slug, description, qualification_req, min_salary, max_salary, demand_level, icon) 
SELECT id, 'Medical Lab Technician', 'medical-lab-technician', 'Perform laboratory tests and analyze samples', 'DMLT / B.Sc MLT', 220000, 500000, 'HIGH', 'microscope'
FROM job_categories WHERE slug = 'healthcare'
ON CONFLICT (category_id, slug) DO UPDATE SET
    name = EXCLUDED.name,
    description = EXCLUDED.description,
    qualification_req = EXCLUDED.qualification_req,
    min_salary = EXCLUDED.min_salary,
    max_salary = EXCLUDED.max_salary,
    demand_level = EXCLUDED.demand_level,
    icon = EXCLUDED.icon,
    updated_at = NOW();

INSERT INTO job_trades (category_id, name, slug, description, qualification_req, min_salary, max_salary, demand_level, icon) 
SELECT id, 'Pharmacy Assistant', 'pharmacy-assistant', 'Assist pharmacists in dispensing medications', 'D.Pharm / Certificate in Pharmacy', 180000, 400000, 'MEDIUM', 'pill'
FROM job_categories WHERE slug = 'healthcare'
ON CONFLICT (category_id, slug) DO UPDATE SET
    name = EXCLUDED.name,
    description = EXCLUDED.description,
    qualification_req = EXCLUDED.qualification_req,
    min_salary = EXCLUDED.min_salary,
    max_salary = EXCLUDED.max_salary,
    demand_level = EXCLUDED.demand_level,
    icon = EXCLUDED.icon,
    updated_at = NOW();

INSERT INTO job_trades (category_id, name, slug, description, qualification_req, min_salary, max_salary, demand_level, icon) 
SELECT id, 'Physiotherapy Assistant', 'physiotherapy-assistant', 'Assist physiotherapists in patient rehabilitation', 'Diploma in Physiotherapy', 200000, 450000, 'MEDIUM', 'activity'
FROM job_categories WHERE slug = 'healthcare'
ON CONFLICT (category_id, slug) DO UPDATE SET
    name = EXCLUDED.name,
    description = EXCLUDED.description,
    qualification_req = EXCLUDED.qualification_req,
    min_salary = EXCLUDED.min_salary,
    max_salary = EXCLUDED.max_salary,
    demand_level = EXCLUDED.demand_level,
    icon = EXCLUDED.icon,
    updated_at = NOW();

-- Creative
INSERT INTO job_trades (category_id, name, slug, description, qualification_req, min_salary, max_salary, demand_level, icon) 
SELECT id, 'Graphic Designer', 'graphic-designer', 'Create visual concepts for digital and print media', 'Diploma in Graphic Design / B.Des', 250000, 700000, 'HIGH', 'palette'
FROM job_categories WHERE slug = 'creative'
ON CONFLICT (category_id, slug) DO UPDATE SET
    name = EXCLUDED.name,
    description = EXCLUDED.description,
    qualification_req = EXCLUDED.qualification_req,
    min_salary = EXCLUDED.min_salary,
    max_salary = EXCLUDED.max_salary,
    demand_level = EXCLUDED.demand_level,
    icon = EXCLUDED.icon,
    updated_at = NOW();

INSERT INTO job_trades (category_id, name, slug, description, qualification_req, min_salary, max_salary, demand_level, icon) 
SELECT id, 'Video Editor', 'video-editor', 'Edit and assemble video content for various platforms', 'Diploma in Video Editing / Certificate Course', 200000, 600000, 'MEDIUM', 'film'
FROM job_categories WHERE slug = 'creative'
ON CONFLICT (category_id, slug) DO UPDATE SET
    name = EXCLUDED.name,
    description = EXCLUDED.description,
    qualification_req = EXCLUDED.qualification_req,
    min_salary = EXCLUDED.min_salary,
    max_salary = EXCLUDED.max_salary,
    demand_level = EXCLUDED.demand_level,
    icon = EXCLUDED.icon,
    updated_at = NOW();

INSERT INTO job_trades (category_id, name, slug, description, qualification_req, min_salary, max_salary, demand_level, icon) 
SELECT id, 'UI/UX Designer', 'ui-ux-designer', 'Design user interfaces and experiences for digital products', 'Diploma in UI/UX / B.Des / Certification', 300000, 800000, 'HIGH', 'layout'
FROM job_categories WHERE slug = 'creative'
ON CONFLICT (category_id, slug) DO UPDATE SET
    name = EXCLUDED.name,
    description = EXCLUDED.description,
    qualification_req = EXCLUDED.qualification_req,
    min_salary = EXCLUDED.min_salary,
    max_salary = EXCLUDED.max_salary,
    demand_level = EXCLUDED.demand_level,
    icon = EXCLUDED.icon,
    updated_at = NOW();

INSERT INTO job_trades (category_id, name, slug, description, qualification_req, min_salary, max_salary, demand_level, icon) 
SELECT id, 'Content Writer', 'content-writer', 'Create written content for digital and print media', 'Bachelor in Journalism / English / Mass Comm', 200000, 550000, 'MEDIUM', 'pen-tool'
FROM job_categories WHERE slug = 'creative'
ON CONFLICT (category_id, slug) DO UPDATE SET
    name = EXCLUDED.name,
    description = EXCLUDED.description,
    qualification_req = EXCLUDED.qualification_req,
    min_salary = EXCLUDED.min_salary,
    max_salary = EXCLUDED.max_salary,
    demand_level = EXCLUDED.demand_level,
    icon = EXCLUDED.icon,
    updated_at = NOW();

COMMIT;