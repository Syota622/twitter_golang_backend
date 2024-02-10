-- name: CreateTweet :one
INSERT INTO tweets (user_id, message, image_url) VALUES ($1, $2, $3) RETURNING id, user_id, message, image_url, created_at, updated_at;
