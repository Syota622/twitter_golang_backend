// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.25.0
// source: 000024_get_group_messages.sql

package generated

import (
	"context"
	"database/sql"
)

const getGroupMessages = `-- name: GetGroupMessages :many
SELECT id, group_id, user_id, message, created_at FROM group_messages
WHERE group_id = $1
ORDER BY created_at DESC
`

func (q *Queries) GetGroupMessages(ctx context.Context, groupID sql.NullInt32) ([]GroupMessage, error) {
	rows, err := q.db.QueryContext(ctx, getGroupMessages, groupID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []GroupMessage
	for rows.Next() {
		var i GroupMessage
		if err := rows.Scan(
			&i.ID,
			&i.GroupID,
			&i.UserID,
			&i.Message,
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
