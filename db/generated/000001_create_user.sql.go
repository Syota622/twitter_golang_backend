// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.25.0
// source: 000001_create_user.sql

package generated

import (
	"context"
	"database/sql"
)

const createUser = `-- name: CreateUser :one
INSERT INTO users (
    username,
    hashed_password,
    email,
    confirmation_token
) VALUES (
    $1, $2, $3, $4
) RETURNING id, username, hashed_password, email, created_at, updated_at, confirmation_token, is_confirmed
`

type CreateUserParams struct {
	Username          string         `json:"username"`
	HashedPassword    string         `json:"hashed_password"`
	Email             string         `json:"email"`
	ConfirmationToken sql.NullString `json:"confirmation_token"`
}

func (q *Queries) CreateUser(ctx context.Context, arg CreateUserParams) (User, error) {
	row := q.db.QueryRowContext(ctx, createUser,
		arg.Username,
		arg.HashedPassword,
		arg.Email,
		arg.ConfirmationToken,
	)
	var i User
	err := row.Scan(
		&i.ID,
		&i.Username,
		&i.HashedPassword,
		&i.Email,
		&i.CreatedAt,
		&i.UpdatedAt,
		&i.ConfirmationToken,
		&i.IsConfirmed,
	)
	return i, err
}
