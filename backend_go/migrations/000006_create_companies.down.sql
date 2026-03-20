-- 000006_create_companies.down.sql

DROP TRIGGER IF EXISTS update_companies_updated_at ON companies;
DROP TABLE IF EXISTS companies;
