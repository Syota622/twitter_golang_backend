// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.25.0
// source: 0000017_get_follow.sql

package generated

import (
	"context"
)

const isFollowing = `-- name: IsFollowing :one
SELECT EXISTS(
  SELECT 1
  FROM follows
  WHERE user_id = $1 AND follow_id = $2
) AS is_following
`

type IsFollowingParams struct {
	UserID   int32 `json:"user_id"`
	FollowID int32 `json:"follow_id"`
}

func (q *Queries) IsFollowing(ctx context.Context, arg IsFollowingParams) (bool, error) {
	row := q.db.QueryRowContext(ctx, isFollowing, arg.UserID, arg.FollowID)
	var is_following bool
	err := row.Scan(&is_following)
	return is_following, err
}
