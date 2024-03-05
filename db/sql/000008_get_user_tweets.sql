-- name: GetUserTweets :many
SELECT id, user_id, message, created_at, updated_at, image_url 
FROM tweets WHERE user_id = $1 ORDER BY created_at DESC;
