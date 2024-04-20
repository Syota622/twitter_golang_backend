-- name: DeactivateUser :exec
UPDATE users SET is_deleted = true WHERE id = $1;