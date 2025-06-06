-- name: CreateFeed :one
INSERT INTO feeds(id, created_at, updated_at, user_id, name, url) VALUES (
	$1,
	$2,
	$3,
	$4,
	$5,
	$6
)
RETURNING *;

-- name: GetFeeds :many
SELECT
	feeds.name, 
	feeds.url,
	users.name as user_name
FROM feeds 
	INNER JOIN users ON users.id = feeds.user_id
;

-- name: GetFeed :one
SELECT * FROM feeds WHERE url = $1;

-- name: MarkedFeedFetched :exec
UPDATE feeds SET last_fetched_at = $1, updated_at = $2 WHERE id = $3;

-- name: GetNextFeedToFetch :one
SELECT *
FROM feeds
ORDER BY last_fetched_at ASC NULLS FIRST;