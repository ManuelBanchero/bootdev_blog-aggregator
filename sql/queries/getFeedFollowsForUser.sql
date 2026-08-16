-- name: GetFeedFollowsForUser :many
SELECT * FROM feed_follows
WHERE $1 = user_id;
