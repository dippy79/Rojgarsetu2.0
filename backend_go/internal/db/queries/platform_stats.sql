-- name: GetPlatformStats :one
SELECT * FROM platform_stats LIMIT 1;

-- name: UpdatePlatformStats :one
UPDATE platform_stats
SET total_jobs = $1, total_candidates = $2, total_companies = $3, total_placements = $4, total_applications = $5, visits_today = $6, updated_at = NOW()
WHERE id = (SELECT id FROM platform_stats LIMIT 1)
RETURNING *;

-- name: IncrementVisits :exec
UPDATE platform_stats SET visits_today = visits_today + 1 WHERE id = (SELECT id FROM platform_stats LIMIT 1);
