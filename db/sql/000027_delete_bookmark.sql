-- name: DeleteBookmark :exec
DELETE FROM bookmarks WHERE user_id = $1 AND tweet_id = $2;