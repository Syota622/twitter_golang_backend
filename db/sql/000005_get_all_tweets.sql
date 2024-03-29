-- name: GetAllTweets :many
SELECT 
    tweets.*,
    COUNT(retweets.id) AS retweet_count
FROM 
    tweets 
LEFT JOIN 
    retweets ON tweets.id = retweets.tweet_id
GROUP BY 
    tweets.id
ORDER BY 
    tweets.updated_at DESC
LIMIT $1 OFFSET $2;
