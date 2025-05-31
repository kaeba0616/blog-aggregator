
-- name: CreateFeed :one
INSERT INTO feeds (id, created_at, updated_at, name, url, user_id, last_fetched_at)
VALUES (
    $1,
    $2,
    $3,
    $4,
    $5,
    $6,
    $7
)
RETURNING *;

-- name: GetFeeds :many
SELECT * FROM feeds;

-- name: GetFeedByURL :one
SELECT * FROM feeds
WHERE url = $1;

-- name: UpdateFeedFetched :one
UPDATE feeds
SET last_fetched_at = NOW(),
    updated_at = NOW()
WHERE id = $1
RETURNING *;

-- name: GetNextFeedFetched :one
SELECT * FROM feeds
ORDER BY last_fetched_at ASC NULLS FIRST
LIMIT 1;
