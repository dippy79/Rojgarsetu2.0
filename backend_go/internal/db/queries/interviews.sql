-- name: CreateInterview :one
INSERT INTO interviews (application_id, candidate_id, company_id, scheduled_at, room_url, meeting_id)
VALUES ($1, $2, $3, $4, $5, $6)
RETURNING *;

-- name: GetInterviewByID :one
SELECT * FROM interviews WHERE id = $1;

-- name: GetInterviewsByCandidate :many
SELECT * FROM interviews WHERE candidate_id = $1 ORDER BY scheduled_at DESC;

-- name: GetInterviewsByCompany :many
SELECT * FROM interviews WHERE company_id = $1 ORDER BY scheduled_at DESC;

-- name: UpdateInterviewStatus :one
UPDATE interviews SET status = $2, updated_at = NOW() WHERE id = $1 RETURNING *;
