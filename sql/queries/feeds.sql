-- name: CreateFeed :one
INSERT INTO feeds(id, user_id, name, url) VALUES (
	$1,
	$2,
	$3,
	$4 
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
