-- name: GetUser :one
SELECT * FROM users WHERE name = $1;