-- +goose Up
CREATE TABLE feeds (
	id UUID PRIMARY KEY,
	user_id UUID REFERENCES users(id) ON DELETE CASCADE NOT NULL,
	name TEXT NOT NULL,
	url TEXT UNIQUE NOT NULL
);

-- +goose Down
DROP TABLE feeds;

	

