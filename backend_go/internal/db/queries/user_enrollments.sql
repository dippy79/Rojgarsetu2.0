-- name: GetUserEnrollments :many
SELECT id, user_id, trade_id, status, enrolled_at, expires_at, completed_at, progress_pct, metadata, created_at, updated_at
FROM user_enrollments
WHERE user_id = $1
  AND ($2::text = '' OR status = $2)
ORDER BY enrolled_at DESC
LIMIT $3 OFFSET $4;

-- name: GetUserEnrollmentsCount :one
SELECT COUNT(*) FROM user_enrollments
WHERE user_id = $1
  AND ($2::text = '' OR status = $2);

-- name: GetUserEnrollmentByID :one
SELECT id, user_id, trade_id, status, enrolled_at, expires_at, completed_at, progress_pct, metadata, created_at, updated_at
FROM user_enrollments
WHERE id = $1;

-- name: GetUserEnrollmentByUserAndTrade :one
SELECT id, user_id, trade_id, status, enrolled_at, expires_at, completed_at, progress_pct, metadata, created_at, updated_at
FROM user_enrollments
WHERE user_id = $1 AND trade_id = $2 AND status = 'active';

-- name: GetExpiringEnrollments :many
SELECT id, user_id, trade_id, status, enrolled_at, expires_at, completed_at, progress_pct, metadata, created_at, updated_at
FROM user_enrollments
WHERE status = 'active'
  AND expires_at <= NOW() + INTERVAL '7 days'
  AND expires_at > NOW()
ORDER BY expires_at ASC;

-- name: GetExpiringEnrollmentsWithTrade :many
SELECT 
    ue.id, ue.user_id, ue.trade_id, ue.status, ue.enrolled_at, ue.expires_at, ue.completed_at, ue.progress_pct, ue.metadata, ue.created_at, ue.updated_at,
    jt.name as trade_name, jt.slug as trade_slug, jt.category_id,
    jc.name as category_name, jc.slug as category_slug, jc.color as category_color, jc.icon as category_icon
FROM user_enrollments ue
JOIN job_trades jt ON ue.trade_id = jt.id
JOIN job_categories jc ON jt.category_id = jc.id
WHERE ue.status = 'active'
  AND ue.expires_at <= NOW() + INTERVAL '7 days'
  AND ue.expires_at > NOW()
ORDER BY ue.expires_at ASC;

-- name: CreateUserEnrollment :one
INSERT INTO user_enrollments (user_id, trade_id, expires_at, metadata)
VALUES ($1, $2, $3, $4)
RETURNING id, user_id, trade_id, status, enrolled_at, expires_at, completed_at, progress_pct, metadata, created_at, updated_at;

-- name: UpdateUserEnrollment :one
UPDATE user_enrollments
SET status = $2, expires_at = COALESCE($3, expires_at), completed_at = $4, progress_pct = $5, metadata = $6, updated_at = NOW()
WHERE id = $1
RETURNING id, user_id, trade_id, status, enrolled_at, expires_at, completed_at, progress_pct, metadata, created_at, updated_at;

-- name: UpdateEnrollmentProgress :one
UPDATE user_enrollments
SET progress_pct = $2, updated_at = NOW()
WHERE id = $1
RETURNING id, user_id, trade_id, status, enrolled_at, expires_at, completed_at, progress_pct, metadata, created_at, updated_at;

-- name: CompleteEnrollment :one
UPDATE user_enrollments
SET status = 'completed', completed_at = NOW(), progress_pct = 100, updated_at = NOW()
WHERE id = $1
RETURNING id, user_id, trade_id, status, enrolled_at, expires_at, completed_at, progress_pct, metadata, created_at, updated_at;

-- name: CancelEnrollment :one
UPDATE user_enrollments
SET status = 'cancelled', updated_at = NOW()
WHERE id = $1
RETURNING id, user_id, trade_id, status, enrolled_at, expires_at, completed_at, progress_pct, metadata, created_at, updated_at;

-- name: DeleteUserEnrollment :exec
DELETE FROM user_enrollments WHERE id = $1;