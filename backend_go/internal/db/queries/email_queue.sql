-- name: EnqueueEmail :one
INSERT INTO email_queue (to_email, subject, body)
VALUES ($1, $2, $3)
RETURNING *;

-- name: GetPendingEmails :many
SELECT * FROM email_queue
WHERE status = 'pending' AND attempts < 3
ORDER BY created_at ASC
LIMIT $1;

-- name: UpdateEmailStatus :one
UPDATE email_queue
SET status = $2, attempts = attempts + 1, last_attempt_at = NOW(), sent_at = CASE WHEN $2 = 'sent' THEN NOW() ELSE sent_at END, error_message = $3
WHERE id = $1
RETURNING *;
