-- name: CreateCompany :one
INSERT INTO companies (
  user_id, name, industry, company_size, website, logo_url, description,
  headquarters, founded_year
) VALUES (
  $1, $2, $3, $4, $5, $6, $7, $8, $9
) RETURNING *;

-- name: GetCompanyByUserID :one
SELECT * FROM companies WHERE user_id = $1;

-- name: GetCompanyByID :one
SELECT * FROM companies WHERE id = $1;

-- name: UpdateCompany :one
UPDATE companies SET 
  name = $2, industry = $3, company_size = $4, website = $5, logo_url = $6,
  description = $7, headquarters = $8, founded_year = $9, updated_at = NOW()
WHERE id = $1 RETURNING *;

-- name: ListCompanies :many
SELECT * FROM companies WHERE is_active = true ORDER BY created_at DESC LIMIT $1 OFFSET $2;

