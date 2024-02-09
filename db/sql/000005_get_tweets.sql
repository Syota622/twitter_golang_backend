-- name: GetTweets :many
SELECT id, user_id, message, created_at, updated_at FROM tweets ORDER BY created_at DESC;
