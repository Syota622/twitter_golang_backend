// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.25.0
// source: 000007_update_profile.sql

package generated

import (
	"context"
	"database/sql"
)

const updateUserProfile = `-- name: UpdateUserProfile :exec
UPDATE users
SET username = $2, email = $3, bio = $4, profile_image_url = $5, background_image_url = $6
WHERE id = $1
`

type UpdateUserProfileParams struct {
	ID                 int32          `json:"id"`
	Username           string         `json:"username"`
	Email              string         `json:"email"`
	Bio                sql.NullString `json:"bio"`
	ProfileImageUrl    sql.NullString `json:"profile_image_url"`
	BackgroundImageUrl sql.NullString `json:"background_image_url"`
}

func (q *Queries) UpdateUserProfile(ctx context.Context, arg UpdateUserProfileParams) error {
	_, err := q.db.ExecContext(ctx, updateUserProfile,
		arg.ID,
		arg.Username,
		arg.Email,
		arg.Bio,
		arg.ProfileImageUrl,
		arg.BackgroundImageUrl,
	)
	return err
}
