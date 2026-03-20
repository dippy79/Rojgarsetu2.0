-- name: GetUserByEmail :one
SELECT id, name, email, password_hash, role, phone, avatar_url, is_active, is_verified, last_login, created_at, updated_at
FROM users
WHERE email = $1 LIMIT 1;

-- name: CreateUser :one
INSERT INTO users (id, name, email, password_hash, role, phone)
VALUES ($1, $2, $3, $4, $5, $6)
RETURNING id, name, email, password_hash, role, phone, avatar_url, is_active, is_verified, last_login, created_at, updated_at;

-- name: EmailExists :one
SELECT EXISTS (SELECT 1 FROM users WHERE email = $1);

-- name: GetUserByID :one
SELECT id, name, email, password_hash, role, phone, avatar_url, is_active, is_verified, last_login, created_at, updated_at
FROM users
WHERE id = $1 LIMIT 1;

-- name: UpdateLastLogin :one
UPDATE users SET last_login = NOW() WHERE id = $1 RETURNING *;

-- name: ListUsers :many
SELECT id, name, email, role, phone, avatar_url, is_active, is_verified, created_at, updated_at
FROM users
WHERE role = $1 OR $1 = ''
ORDER BY created_at DESC
LIMIT $2 OFFSET $3;

