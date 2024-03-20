-- name: DeleteTweet :exec
DELETE FROM tweets WHERE id = $1;