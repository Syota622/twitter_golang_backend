-- name: CreateUser :one
INSERT INTO users (
    username,
    hashed_password,
    email,
    confirmation_token
) VALUES (
    $1, $2, $3, $4
) RETURNING id, username, hashed_password, email, created_at, updated_at, confirmation_token, is_confirmed;