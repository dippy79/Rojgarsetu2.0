-- Search queries using full-text search and pg_trgm

-- name: UnifiedSearch :many
-- Unified search across government and private jobs using plainto_tsquery
SELECT 
    'gov' as job_type,
    id,
    title,
    department,
    location,
    apply_url,
    last_date,
    source,
    eligibility,
    vacancy_count,
    salary,
    exam_date,
    created_at,
    ts_rank(search_vector, plainto_tsquery('english', $1)) as rank
FROM jobs_government
WHERE 
    is_active = true
    AND search_vector @@ plainto_tsquery('english', $1)

UNION ALL

SELECT 
    'private' as job_type,
    id,
    company as department,
    location,
    url as apply_url,
    posted_at as last_date,
    source,
    description as eligibility,
    0 as vacancy_count,
    salary,
    NULL as exam_date,
    created_at,
    ts_rank(search_vector, plainto_tsquery('english', $1)) as rank
FROM jobs_private
WHERE 
    is_active = true
    AND search_vector @@ plainto_tsquery('english', $1)

ORDER BY rank DESC
LIMIT $2 OFFSET $3;

-- name: UnifiedSearchCount :one
SELECT COUNT(*) FROM (
    SELECT 1 FROM jobs_government
    WHERE is_active = true AND search_vector @@ plainto_tsquery('english', $1)
    UNION ALL
    SELECT 1 FROM jobs_private
    WHERE is_active = true AND search_vector @@ plainto_tsquery('english', $1)
) as combined;

-- name: SearchCompanyJobs :many
SELECT 
    cj.id,
    cj.title,
    cj.description,
    cj.location,
    cj.job_type,
    cj.salary_min,
    cj.salary_max,
    cj.skills,
    cj.is_remote,
    cj.created_at,
    c.name as company_name,
    ts_rank(cj.search_vector, plainto_tsquery('english', $1)) as rank
FROM company_jobs cj
LEFT JOIN companies c ON cj.company_id = c.id
WHERE 
    cj.is_active = true
    AND (
        cj.search_vector @@ plainto_tsquery('english', $1)
        OR cj.title ILIKE '%' || $2 || '%'
        OR cj.description ILIKE '%' || $2 || '%'
        OR cj.location ILIKE '%' || $2 || '%'
        OR c.name ILIKE '%' || $2 || '%'
    )
ORDER BY rank DESC, cj.created_at DESC
LIMIT $3 OFFSET $4;

-- name: SearchCompanyJobsCount :one
SELECT COUNT(*)
FROM company_jobs cj
LEFT JOIN companies c ON cj.company_id = c.id
WHERE 
    cj.is_active = true
    AND (
        cj.search_vector @@ plainto_tsquery('english', $1)
        OR cj.title ILIKE '%' || $2 || '%'
        OR cj.description ILIKE '%' || $2 || '%'
        OR cj.location ILIKE '%' || $2 || '%'
        OR c.name ILIKE '%' || $2 || '%'
    );

-- name: SearchGovJobs :many
SELECT 
    id,
    title,
    department,
    location,
    apply_url,
    last_date,
    source,
    eligibility,
    vacancy_count,
    salary,
    exam_date,
    created_at,
    ts_rank(search_vector, plainto_tsquery('english', $1)) as rank
FROM jobs_government
WHERE 
    is_active = true
    AND (
        search_vector @@ plainto_tsquery('english', $1)
        OR title ILIKE '%' || $2 || '%'
        OR department ILIKE '%' || $2 || '%'
        OR location ILIKE '%' || $2 || '%'
        OR eligibility ILIKE '%' || $2 || '%'
    )
ORDER BY rank DESC, created_at DESC
LIMIT $3 OFFSET $4;

-- name: SearchGovJobsCount :one
SELECT COUNT(*)
FROM jobs_government
WHERE 
    is_active = true
    AND (
        search_vector @@ plainto_tsquery('english', $1)
        OR title ILIKE '%' || $2 || '%'
        OR department ILIKE '%' || $2 || '%'
        OR location ILIKE '%' || $2 || '%'
        OR eligibility ILIKE '%' || $2 || '%'
    );

-- name: SearchPrivJobs :many
SELECT 
    id,
    company,
    title,
    location,
    url,
    salary,
    experience,
    job_type,
    skills,
    description,
    source,
    posted_at,
    created_at,
    ts_rank(search_vector, plainto_tsquery('english', $1)) as rank
FROM jobs_private
WHERE 
    is_active = true
    AND (
        search_vector @@ plainto_tsquery('english', $1)
        OR title ILIKE '%' || $2 || '%'
        OR company ILIKE '%' || $2 || '%'
        OR location ILIKE '%' || $2 || '%'
        OR description ILIKE '%' || $2 || '%'
    )
ORDER BY rank DESC, created_at DESC
LIMIT $3 OFFSET $4;

-- name: SearchPrivJobsCount :one
SELECT COUNT(*)
FROM jobs_private
WHERE 
    is_active = true
    AND (
        search_vector @@ plainto_tsquery('english', $1)
        OR title ILIKE '%' || $2 || '%'
        OR company ILIKE '%' || $2 || '%'
        OR location ILIKE '%' || $2 || '%'
        OR description ILIKE '%' || $2 || '%'
    );

