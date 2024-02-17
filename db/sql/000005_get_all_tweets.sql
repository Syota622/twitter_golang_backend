-- name: GetTweets :many
SELECT id, user_id, message, image_url, created_at, updated_at FROM tweets ORDER BY created_at DESC;
