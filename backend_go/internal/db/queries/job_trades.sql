-- name: GetJobTrades :many
SELECT id, category_id, name, slug, description, qualification_req, min_salary, max_salary, demand_level, icon, is_active, created_at, updated_at
FROM job_trades
WHERE is_active = true
  AND ($1::uuid = '00000000-0000-0000-0000-000000000000'::uuid OR category_id = $1)
  AND ($2::text = '' OR demand_level = $2)
ORDER BY category_id, name ASC;

-- name: GetJobTradesCount :one
SELECT COUNT(*) FROM job_trades
WHERE is_active = true
  AND ($1::uuid = '00000000-0000-0000-0000-000000000000'::uuid OR category_id = $1)
  AND ($2::text = '' OR demand_level = $2);

-- name: GetJobTradeByID :one
SELECT id, category_id, name, slug, description, qualification_req, min_salary, max_salary, demand_level, icon, is_active, created_at, updated_at
FROM job_trades
WHERE id = $1 AND is_active = true;

-- name: GetJobTradeBySlug :one
SELECT id, category_id, name, slug, description, qualification_req, min_salary, max_salary, demand_level, icon, is_active, created_at, updated_at
FROM job_trades
WHERE slug = $1 AND category_id = $2 AND is_active = true;

-- name: GetJobTradesByCategory :many
SELECT id, category_id, name, slug, description, qualification_req, min_salary, max_salary, demand_level, icon, is_active, created_at, updated_at
FROM job_trades
WHERE category_id = $1 AND is_active = true
ORDER BY name ASC;

-- name: CreateJobTrade :one
INSERT INTO job_trades (category_id, name, slug, description, qualification_req, min_salary, max_salary, demand_level, icon)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
RETURNING id, category_id, name, slug, description, qualification_req, min_salary, max_salary, demand_level, icon, is_active, created_at, updated_at;

-- name: UpdateJobTrade :one
UPDATE job_trades
SET category_id = $2, name = $3, slug = $4, description = $5, qualification_req = $6, min_salary = $7, max_salary = $8, demand_level = $9, icon = $10, is_active = $11, updated_at = NOW()
WHERE id = $1
RETURNING id, category_id, name, slug, description, qualification_req, min_salary, max_salary, demand_level, icon, is_active, created_at, updated_at;

-- name: DeleteJobTrade :exec
DELETE FROM job_trades WHERE id = $1;