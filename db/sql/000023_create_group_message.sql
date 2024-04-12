-- name: CreateGroupMessage :one
INSERT INTO group_messages (group_id, user_id, message)
VALUES ($1, $2, $3)
RETURNING *;