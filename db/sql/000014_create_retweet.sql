-- name: CreateRetweet :one
INSERT INTO retweets (user_id, tweet_id) VALUES ($1, $2) RETURNING *;