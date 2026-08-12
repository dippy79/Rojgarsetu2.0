-- name: GetUserNotificationLogs :many
SELECT id, user_id, enrollment_id, notification_type, channel, title, message, payload, sent_at, read_at, clicked_at, created_at
FROM user_notification_logs
WHERE user_id = $1
  AND ($2::text = '' OR notification_type = $2)
ORDER BY sent_at DESC
LIMIT $3 OFFSET $4;

-- name: GetUserNotificationLogsCount :one
SELECT COUNT(*) FROM user_notification_logs
WHERE user_id = $1
  AND ($2::text = '' OR notification_type = $2);

-- name: GetNotificationLogByID :one
SELECT id, user_id, enrollment_id, notification_type, channel, title, message, payload, sent_at, read_at, clicked_at, created_at
FROM user_notification_logs
WHERE id = $1;

-- name: GetDailyNotificationCount :one
SELECT COUNT(*) FROM user_notification_logs
WHERE user_id = $1
  AND DATE(sent_at) = CURRENT_DATE;

-- name: GetDailyNotificationCountByType :one
SELECT COUNT(*) FROM user_notification_logs
WHERE user_id = $1
  AND DATE(sent_at) = CURRENT_DATE
  AND notification_type = $2;

-- name: CreateNotificationLog :one
INSERT INTO user_notification_logs (user_id, enrollment_id, notification_type, channel, title, message, payload)
VALUES ($1, $2, $3, $4, $5, $6, $7)
RETURNING id, user_id, enrollment_id, notification_type, channel, title, message, payload, sent_at, read_at, clicked_at, created_at;

-- name: MarkNotificationRead :one
UPDATE user_notification_logs
SET read_at = NOW()
WHERE id = $1
RETURNING id, user_id, enrollment_id, notification_type, channel, title, message, payload, sent_at, read_at, clicked_at, created_at;

-- name: MarkNotificationClicked :one
UPDATE user_notification_logs
SET clicked_at = NOW()
WHERE id = $1
RETURNING id, user_id, enrollment_id, notification_type, channel, title, message, payload, sent_at, read_at, clicked_at, created_at;

-- name: DeleteNotificationLog :exec
DELETE FROM user_notification_logs WHERE id = $1;