-- name: CreateRefreshToken :one
INSERT INTO refresh_tokens (
  user_id, 
  ip_hash, 
  ua_hash, 
  expires_at
) VALUES (
  $1, $2, $3, $4
) RETURNING *;

-- name: GetRefreshToken :one
SELECT * FROM refresh_tokens 
WHERE id = $1 FOR UPDATE;

-- name: GetRefreshTokenByToken :one
SELECT * FROM refresh_tokens 
WHERE id = $1 AND NOT revoked AND expires_at > NOW() FOR UPDATE;

-- name: RevokeRefreshToken :exec
UPDATE refresh_tokens SET revoked = true WHERE id = $1;

-- name: RevokeAllTokensForUser :exec
UPDATE refresh_tokens SET revoked = true WHERE user_id = $1;

-- name: CleanupExpiredTokens :many
DELETE FROM refresh_tokens 
WHERE revoked = true OR expires_at < NOW() 
RETURNING *;

