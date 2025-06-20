// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.29.0
// source: feeds.sql

package database

import (
	"context"
	"database/sql"
	"time"

	"github.com/google/uuid"
)

const createFeed = `-- name: CreateFeed :one
INSERT INTO feeds(id, created_at, updated_at, user_id, name, url) VALUES (
	$1,
	$2,
	$3,
	$4,
	$5,
	$6
)
RETURNING id, created_at, updated_at, user_id, name, url, last_fetched_at
`

type CreateFeedParams struct {
	ID        uuid.UUID
	CreatedAt time.Time
	UpdatedAt time.Time
	UserID    uuid.UUID
	Name      string
	Url       string
}

func (q *Queries) CreateFeed(ctx context.Context, arg CreateFeedParams) (Feed, error) {
	row := q.db.QueryRowContext(ctx, createFeed,
		arg.ID,
		arg.CreatedAt,
		arg.UpdatedAt,
		arg.UserID,
		arg.Name,
		arg.Url,
	)
	var i Feed
	err := row.Scan(
		&i.ID,
		&i.CreatedAt,
		&i.UpdatedAt,
		&i.UserID,
		&i.Name,
		&i.Url,
		&i.LastFetchedAt,
	)
	return i, err
}

const getFeed = `-- name: GetFeed :one
SELECT id, created_at, updated_at, user_id, name, url, last_fetched_at FROM feeds WHERE url = $1
`

func (q *Queries) GetFeed(ctx context.Context, url string) (Feed, error) {
	row := q.db.QueryRowContext(ctx, getFeed, url)
	var i Feed
	err := row.Scan(
		&i.ID,
		&i.CreatedAt,
		&i.UpdatedAt,
		&i.UserID,
		&i.Name,
		&i.Url,
		&i.LastFetchedAt,
	)
	return i, err
}

const getFeeds = `-- name: GetFeeds :many
SELECT
	feeds.name, 
	feeds.url,
	users.name as user_name
FROM feeds 
	INNER JOIN users ON users.id = feeds.user_id
`

type GetFeedsRow struct {
	Name     string
	Url      string
	UserName string
}

func (q *Queries) GetFeeds(ctx context.Context) ([]GetFeedsRow, error) {
	rows, err := q.db.QueryContext(ctx, getFeeds)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []GetFeedsRow
	for rows.Next() {
		var i GetFeedsRow
		if err := rows.Scan(&i.Name, &i.Url, &i.UserName); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const getNextFeedToFetch = `-- name: GetNextFeedToFetch :one
SELECT id, created_at, updated_at, user_id, name, url, last_fetched_at
FROM feeds
ORDER BY last_fetched_at ASC NULLS FIRST
`

func (q *Queries) GetNextFeedToFetch(ctx context.Context) (Feed, error) {
	row := q.db.QueryRowContext(ctx, getNextFeedToFetch)
	var i Feed
	err := row.Scan(
		&i.ID,
		&i.CreatedAt,
		&i.UpdatedAt,
		&i.UserID,
		&i.Name,
		&i.Url,
		&i.LastFetchedAt,
	)
	return i, err
}

const markedFeedFetched = `-- name: MarkedFeedFetched :exec
UPDATE feeds SET last_fetched_at = $1, updated_at = $2 WHERE id = $3
`

type MarkedFeedFetchedParams struct {
	LastFetchedAt sql.NullTime
	UpdatedAt     time.Time
	ID            uuid.UUID
}

func (q *Queries) MarkedFeedFetched(ctx context.Context, arg MarkedFeedFetchedParams) error {
	_, err := q.db.ExecContext(ctx, markedFeedFetched, arg.LastFetchedAt, arg.UpdatedAt, arg.ID)
	return err
}
