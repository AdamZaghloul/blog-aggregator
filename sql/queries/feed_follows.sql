-- name: CreateFeedFollow :one
WITH inserted_feed_follow AS (
    INSERT INTO feed_follows (id, created_at, updated_at, user_id, feed_id)
    VALUES (
        $1,
        $2,
        $3,
        $4,
        $5
    )
    RETURNING *
)

SELECT inserted_feed_follow.*, f.name as feed_name, u.name as user_name FROM inserted_feed_follow INNER JOIN feeds f ON inserted_feed_follow.feed_id = f.id INNER JOIN users u ON inserted_feed_follow.user_id = u.id;