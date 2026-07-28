-- name: CreateCompanyJob :one
INSERT INTO company_jobs (
  company_id, title, description, requirements, responsibilities, location,
  job_type, experience_min, experience_max, salary_min, salary_max,
  salary_currency, skills, benefits, application_url, application_email,
  is_remote, expires_at
) VALUES (
  $1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17, $18
) RETURNING *;

-- name: GetCompanyJobByID :one
SELECT * FROM company_jobs WHERE id = $1;

-- name: GetCompanyJobsByCompanyID :many
SELECT * FROM company_jobs WHERE company_id = $1 AND is_active = true ORDER BY created_at DESC LIMIT $2 OFFSET $3;

-- name: UpdateCompanyJob :one
UPDATE company_jobs SET
  title = $2, description = $3, requirements = $4, responsibilities = $5, location = $6,
  job_type = $7, experience_min = $8, experience_max = $9, salary_min = $10, salary_max = $11,
  salary_currency = $12, skills = $13, benefits = $14, application_url = $15, application_email = $16,
  is_remote = $17, expires_at = $18, updated_at = NOW()
WHERE id = $1 RETURNING *;

-- name: DeleteCompanyJob :exec
UPDATE company_jobs SET is_active = false WHERE id = $1;

-- name: ListActiveCompanyJobs :many
SELECT * FROM company_jobs WHERE is_active = true
  AND ($1::text = '' OR location ILIKE '%' || $1 || '%')
  AND ($2::text = '' OR job_type = $2)
ORDER BY created_at DESC LIMIT $3 OFFSET $4;

-- name: IncrementJobViews :exec
UPDATE company_jobs SET views_count = views_count + 1 WHERE id = $1;

