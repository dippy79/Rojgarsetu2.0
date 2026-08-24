-- name: CreateFileUpload :one
INSERT INTO file_uploads (user_id, file_type, file_url, original_name, file_size)
VALUES ($1, $2, $3, $4, $5)
RETURNING *;

-- name: GetFileUploadsByUserID :many
SELECT * FROM file_uploads WHERE user_id = $1 ORDER BY created_at DESC;

-- name: DeleteFileUpload :exec
DELETE FROM file_uploads WHERE id = $1 AND user_id = $2;
