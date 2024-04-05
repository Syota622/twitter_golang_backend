-- name: CreateNotification :one
INSERT INTO notifications (user_id, notified_by_id, type, post_id, comment_id)
VALUES ($1, $2, $3, $4, $5)
RETURNING *;