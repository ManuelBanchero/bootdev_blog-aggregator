-- name: GetUserById :one
SELECT * FROM users
WHERE $1 = id;
