-- name: CreateComment :one
INSERT INTO comments (tweet_id, comment)
VALUES ($1, $2)
RETURNING *;