-- name: SaveJob :one
INSERT INTO saved_jobs (user_id, gov_job_id, priv_job_id)
VALUES ($1, $2, $3)
RETURNING *;

-- name: GetSavedJobsByUser :many
SELECT sj.*, g.title as gov_title, p.title as priv_title, p.company as priv_company
FROM saved_jobs sj
LEFT JOIN jobs_government g ON sj.gov_job_id = g.id
LEFT JOIN jobs_private p ON sj.priv_job_id = p.id
WHERE sj.user_id = $1
ORDER BY sj.created_at DESC;

-- name: DeleteSavedJob :exec
DELETE FROM saved_jobs WHERE id = $1 AND user_id = $2;
