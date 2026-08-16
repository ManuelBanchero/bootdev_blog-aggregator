-- name: CreateFeedFollow :one
WITH inserted_feed_follow AS (
  INSERT INTO feed_follows (
    id, created_at, updated_at, user_id, feed_id
  ) VALUES (
  $1, $2, $3, $4, $5
  )
  RETURNING *
)
SELECT 
  inserted_feed_follow.*,
  F.name AS feed_name,
  U.name AS user_name,
  F.url AS feed_url
FROM 
  inserted_feed_follow 
INNER JOIN
  users U ON inserted_feed_follow.user_id = U.id
INNER JOIN
  feeds F ON inserted_feed_follow.feed_id = F.id;
  
