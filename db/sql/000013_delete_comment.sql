-- name: DeleteComment :exec
DELETE FROM comments WHERE id = $1;