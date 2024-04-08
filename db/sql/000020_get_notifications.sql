-- name: GetNotificationsByUserID :many
SELECT * FROM notifications WHERE user_id = $1 ORDER BY created_at DESC;