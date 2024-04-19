-- name: CreateBookmark :exec
INSERT INTO bookmarks (user_id, tweet_id) VALUES ($1, $2) RETURNING *;