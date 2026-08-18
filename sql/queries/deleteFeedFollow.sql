-- name: DeleteFeedFollow :exec
DELETE FROM feed_follows
WHERE
  $1 = user_id
AND
  $2 = feed_id;
