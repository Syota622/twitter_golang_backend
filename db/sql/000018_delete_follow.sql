-- name: Unfollow :one
DELETE FROM follows
WHERE user_id = $1 AND follow_id = $2
RETURNING *;