// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.25.0

package generated

import (
	"database/sql"
	"time"
)

type Comment struct {
	ID        int32        `json:"id"`
	TweetID   int32        `json:"tweet_id"`
	Comment   string       `json:"comment"`
	CreatedAt sql.NullTime `json:"created_at"`
	UpdatedAt sql.NullTime `json:"updated_at"`
}

type Like struct {
	ID        int32     `json:"id"`
	UserID    int32     `json:"user_id"`
	TweetID   int32     `json:"tweet_id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type Retweet struct {
	ID        int32     `json:"id"`
	UserID    int32     `json:"user_id"`
	TweetID   int32     `json:"tweet_id"`
	CreatedAt time.Time `json:"created_at"`
}

type Tweet struct {
	ID        int32          `json:"id"`
	UserID    int32          `json:"user_id"`
	Message   string         `json:"message"`
	CreatedAt sql.NullTime   `json:"created_at"`
	UpdatedAt sql.NullTime   `json:"updated_at"`
	ImageUrl  sql.NullString `json:"image_url"`
}

type User struct {
	ID                 int32          `json:"id"`
	Username           string         `json:"username"`
	HashedPassword     string         `json:"hashed_password"`
	Email              string         `json:"email"`
	CreatedAt          sql.NullTime   `json:"created_at"`
	UpdatedAt          sql.NullTime   `json:"updated_at"`
	ConfirmationToken  sql.NullString `json:"confirmation_token"`
	IsConfirmed        sql.NullBool   `json:"is_confirmed"`
	Bio                sql.NullString `json:"bio"`
	ProfileImageUrl    sql.NullString `json:"profile_image_url"`
	BackgroundImageUrl sql.NullString `json:"background_image_url"`
}
