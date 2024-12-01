-- name: GetFeed :one
SELECT * FROM feeds WHERE url = $1;