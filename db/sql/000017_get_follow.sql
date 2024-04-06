-- name: IsFollowing :one
SELECT EXISTS(
  SELECT 1
  FROM follows
  WHERE user_id = $1 AND follow_id = $2
) AS is_following;