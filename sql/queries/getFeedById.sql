-- name: GetFeedById :one
SELECT * FROM feeds
WHERE $1 = id;
