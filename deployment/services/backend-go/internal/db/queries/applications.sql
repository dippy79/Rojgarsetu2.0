-- name: CreateJobApplication :one
INSERT INTO job_applications (job_id, candidate_id, cover_letter, resume_url)
VALUES ($1, $2, $3, $4) RETURNING *;

-- name: ApplyJob :one
INSERT INTO job_applications (job_id, candidate_id, cover_letter, resume_url)
VALUES ($1, $2, $3, $4) RETURNING *;

-- name: GetJobApplicationByID :one
SELECT * FROM job_applications WHERE id = $1;

-- name: ListJobApplicationsByJobID :many
SELECT * FROM job_applications WHERE job_id = $1 ORDER BY applied_at DESC LIMIT $2 OFFSET $3;

-- name: GetJobApplicationsCountByJobID :one
SELECT COUNT(*) FROM job_applications WHERE job_id = $1;

-- name: ListJobApplicationsByCandidateID :many
SELECT * FROM job_applications WHERE candidate_id = $1 ORDER BY applied_at DESC LIMIT $2 OFFSET $3;

-- name: UpdateJobApplicationStatus :one
UPDATE job_applications SET status = $2, notes = $3, updated_at = NOW() WHERE id = $1 RETURNING *;

