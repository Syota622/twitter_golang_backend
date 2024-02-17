-- name: GetAllTweets :many
SELECT id, user_id, message, created_at, updated_at
FROM tweets
ORDER BY created_at DESC
LIMIT $1 OFFSET $2;
