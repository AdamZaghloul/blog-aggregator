-- name: GetPostsForUser :many
WITH feedfollows AS (SELECT ff.feed_id FROM feed_follows ff INNER JOIN users u ON ff.user_id = u.id WHERE u.id = $1)
SELECT posts.*, f.name FROM posts LEFT JOIN feeds f ON f.id = posts.feed_id WHERE posts.feed_id IN (SELECT feed_id FROM feedfollows) ORDER BY posts.published_at DESC LIMIT $2;