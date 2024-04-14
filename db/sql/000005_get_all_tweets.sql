-- name: GetAllTweets :many
SELECT 
    tweets.*,
    users.username,
    COALESCE(retweet_counts.count, 0) AS retweet_count,
    COALESCE(like_counts.count, 0) AS like_count,
    CASE WHEN bookmarks.tweet_id IS NOT NULL THEN TRUE ELSE FALSE END AS is_bookmarked
FROM 
    tweets 
JOIN 
    users ON tweets.user_id = users.id
LEFT JOIN 
    (SELECT tweet_id, COUNT(*) AS count FROM retweets GROUP BY tweet_id) AS retweet_counts
    ON tweets.id = retweet_counts.tweet_id
LEFT JOIN 
    (SELECT tweet_id, COUNT(*) AS count FROM likes GROUP BY tweet_id) AS like_counts
    ON tweets.id = like_counts.tweet_id
LEFT JOIN 
    bookmarks ON tweets.id = bookmarks.tweet_id AND bookmarks.user_id = $3
ORDER BY 
    tweets.updated_at DESC
LIMIT $1 OFFSET $2;