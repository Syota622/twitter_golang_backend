-- name: CreateUser :one
INSERT INTO users (
    username,
    hashed_password,
    email
) VALUES (
    $1, $2, $3
) RETURNING id, username, hashed_password, email, created_at, updated_at;
