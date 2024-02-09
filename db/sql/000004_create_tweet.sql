-- name: CreateTweet :one
INSERT INTO tweets (user_id, message) VALUES ($1, $2) RETURNING id, user_id, message, created_at, updated_at;
