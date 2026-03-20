-- name: GetPrivJobs :many
SELECT id, company, title, location, url, salary, experience, 
       job_type, skills, description, source, posted_at, created_at
FROM jobs_private 
WHERE is_active = true
  AND ($1::text = '' OR company ILIKE '%' || $1 || '%')
  AND ($2::text = '' OR location ILIKE '%' || $2 || '%')
  AND ($3::text = '' OR source = $3)
  AND ($4::text = '' OR job_type = $4)
ORDER BY created_at DESC 
LIMIT $5 OFFSET $6;

-- name: GetPrivJobsCount :one
SELECT COUNT(*) FROM jobs_private WHERE is_active = true
  AND ($1::text = '' OR company ILIKE '%' || $1 || '%')
  AND ($2::text = '' OR location ILIKE '%' || $2 || '%')
  AND ($3::text = '' OR source = $3)
  AND ($4::text = '' OR job_type = $4);

-- name: GetPrivJobByID :one
SELECT id, company, title, location, url, salary, experience, 
       job_type, skills, description, source, posted_at, created_at
FROM jobs_private 
WHERE id = $1 AND is_active = true;

