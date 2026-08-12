-- name: GetJobCategories :many
SELECT id, name, slug, description, icon, color, display_order, is_active, created_at, updated_at
FROM job_categories
WHERE is_active = true
ORDER BY display_order ASC;

-- name: GetJobCategoryBySlug :one
SELECT id, name, slug, description, icon, color, display_order, is_active, created_at, updated_at
FROM job_categories
WHERE slug = $1 AND is_active = true;

-- name: GetJobCategoryByID :one
SELECT id, name, slug, description, icon, color, display_order, is_active, created_at, updated_at
FROM job_categories
WHERE id = $1 AND is_active = true;

-- name: CreateJobCategory :one
INSERT INTO job_categories (name, slug, description, icon, color, display_order)
VALUES ($1, $2, $3, $4, $5, $6)
RETURNING id, name, slug, description, icon, color, display_order, is_active, created_at, updated_at;

-- name: UpdateJobCategory :one
UPDATE job_categories
SET name = $2, slug = $3, description = $4, icon = $5, color = $6, display_order = $7, is_active = $8, updated_at = NOW()
WHERE id = $1
RETURNING id, name, slug, description, icon, color, display_order, is_active, created_at, updated_at;

-- name: DeleteJobCategory :exec
DELETE FROM job_categories WHERE id = $1;