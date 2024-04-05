-- name: CreateLikesTable
CREATE TABLE likes (
    id SERIAL PRIMARY KEY,
    user_id INTEGER NOT NULL,
    tweet_id INTEGER NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);