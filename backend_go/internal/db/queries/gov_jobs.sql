-- name: GetGovJobs :many
SELECT id, title, department, location, apply_url, last_date, 
       source, eligibility, vacancy_count, salary, exam_date, created_at
FROM jobs_government 
WHERE is_active = true
  AND ($1::text = '' OR department ILIKE '%' || $1 || '%')
  AND ($2::text = '' OR location ILIKE '%' || $2 || '%')
  AND ($3::text = '' OR source = $3)
ORDER BY created_at DESC 
LIMIT $4 OFFSET $5;

-- name: GetGovJobsCount :one
SELECT COUNT(*) FROM jobs_government WHERE is_active = true
  AND ($1::text = '' OR department ILIKE '%' || $1 || '%')
  AND ($2::text = '' OR location ILIKE '%' || $2 || '%')
  AND ($3::text = '' OR source = $3);

-- name: GetGovJobByID :one
SELECT id, title, department, location, apply_url, last_date, 
       source, eligibility, vacancy_count, salary, exam_date, created_at
FROM jobs_government 
WHERE id = $1 AND is_active = true;

