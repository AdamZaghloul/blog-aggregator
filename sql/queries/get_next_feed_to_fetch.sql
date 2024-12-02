-- name: GetNextFeedToFetch :one
SELECT * from feeds ORDER BY last_fetched_at NULLS FIRST LIMIT 1;