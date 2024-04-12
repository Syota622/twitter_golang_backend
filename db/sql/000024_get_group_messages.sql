-- name: GetGroupMessages :many
SELECT * FROM group_messages
WHERE group_id = $1
ORDER BY created_at DESC;