-- name: GetCommentsByTweetID :many
SELECT * FROM comments WHERE tweet_id = $1 ORDER BY created_at DESC;