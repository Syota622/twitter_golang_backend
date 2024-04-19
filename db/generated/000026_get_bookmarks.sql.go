// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.25.0
// source: 000026_get_bookmarks.sql

package generated

import (
	"context"
)

const listBookmarks = `-- name: ListBookmarks :many
SELECT id, user_id, tweet_id, created_at FROM bookmarks WHERE user_id = $1 ORDER BY created_at DESC
`

func (q *Queries) ListBookmarks(ctx context.Context, userID int32) ([]Bookmark, error) {
	rows, err := q.db.QueryContext(ctx, listBookmarks, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []Bookmark
	for rows.Next() {
		var i Bookmark
		if err := rows.Scan(
			&i.ID,
			&i.UserID,
			&i.TweetID,
			&i.CreatedAt,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}
