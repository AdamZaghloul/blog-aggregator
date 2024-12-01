-- name: GetFeeds :many
SELECT f.name, f.url, u.name as user_name FROM feeds f LEFT JOIN users u ON f.user_id = u.id;