-- name: CreateLike :one
INSERT INTO likes (user_id, tweet_id) VALUES ($1, $2) RETURNING *;