-- name: ListBookmarks :many
SELECT * FROM bookmarks WHERE user_id = $1 ORDER BY created_at DESC;