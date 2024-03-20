// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.25.0
// source: 000010_delete_tweet.sql

package generated

import (
	"context"
)

const deleteTweet = `-- name: DeleteTweet :exec
DELETE FROM tweets WHERE id = $1
`

func (q *Queries) DeleteTweet(ctx context.Context, id int32) error {
	_, err := q.db.ExecContext(ctx, deleteTweet, id)
	return err
}
