-- name: GetGroupMessages :many
SELECT * FROM messages
WHERE group_id = $1
ORDER BY created_at DESC;