-- name: GetTweetByID :one
SELECT id, user_id, message, image_url, created_at, updated_at
FROM tweets
WHERE id = $1;
