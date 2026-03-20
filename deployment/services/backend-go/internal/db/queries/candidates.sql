-- name: CreateCandidate :one
INSERT INTO candidates (
  user_id, phone, resume_url, resume_parsed, skills, experience_years,
  current_company, current_position, location, linkedin_url, github_url,
  portfolio_url, bio, is_open_to_work, expected_salary, preferred_job_type,
  preferred_location
) VALUES (
  $1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17
) RETURNING *;

-- name: GetCandidateByUserID :one
SELECT * FROM candidates WHERE user_id = $1;

-- name: GetCandidateByID :one
SELECT * FROM candidates WHERE id = $1;

-- name: UpdateCandidate :one
UPDATE candidates SET 
  phone = $2, resume_url = $3, skills = $4, experience_years = $5,
  current_company = $6, current_position = $7, location = $8,
  linkedin_url = $9, github_url = $10, portfolio_url = $11, bio = $12,
  is_open_to_work = $13, expected_salary = $14, preferred_job_type = $15,
  preferred_location = $16, updated_at = NOW()
WHERE id = $1 RETURNING *;

-- name: ListCandidates :many
SELECT * FROM candidates WHERE is_active = true
  AND ($1::text = '' OR location ILIKE '%' || $1 || '%')
ORDER BY created_at DESC LIMIT $2 OFFSET $3;

-- name: GetCandidatesCount :one
SELECT COUNT(*) FROM candidates WHERE is_active = true
  AND ($1::text = '' OR location ILIKE '%' || $1 || '%');

