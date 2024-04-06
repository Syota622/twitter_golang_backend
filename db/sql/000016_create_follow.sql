-- name: CreateFollow :one
INSERT INTO follows (user_id, follow_id) VALUES ($1, $2) RETURNING *;
