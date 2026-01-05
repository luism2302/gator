-- name: CreateFeedFollow :one
WITH inserted_feed_follows AS(
    INSERT INTO feed_follows(id, created_at, updated_at, user_id, feed_id) VALUES (
	$1,
	$2,
	$3,
	$4,
	$5
    )
    RETURNING *
)
SELECT 
    inserted_feed_follows.id AS follow_id,
    inserted_feed_follows.created_at AS created_at, 
    inserted_feed_follows.updated_at AS updated_at,
    inserted_feed_follows.user_id AS user_id,
    users.name AS user_name,
    inserted_feed_follows.feed_id AS feed_id,
    feeds.name AS feed_name
FROM inserted_feed_follows
    JOIN users ON users.id = inserted_feed_follows.user_id
    JOIN feeds ON feeds.id = inserted_feed_follows.feed_id
;

-- name: GetFeedFollowsForUser :many
SELECT feeds.name 
FROM feed_follows 
    JOIN feeds ON feeds.id = feed_follows.feed_id
WHERE feed_follows.user_id = $1;

-- name: DeleteFeedFollow :exec
DELETE FROM feed_follows WHERE user_id = $1 AND feed_id = $2;
