-- +goose Up
CREATE TABLE posts(
    id UUID PRIMARY KEY NOT NULL,
    created_at TIMESTAMPTZ NOT NULL,
    updated_at TIMESTAMPTZ NOT NULL,
    title TEXT,
    url TEXT UNIQUE NOT NULL,
    description TEXT,
    published_at TIMESTAMPTZ,
    feed_id UUID NOT NULL,
    CONSTRAINT fk_feed_id
    FOREIGN KEY (feed_id)
    REFERENCES feeds(id) ON DELETE CASCADE
);

-- +goose Down
DROP TABLE posts;