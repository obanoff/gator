-- name: CreateFeed :one
INSERT INTO feeds (name, url, user_id)
VALUES ($1, $2, (
  SELECT users.id FROM users WHERE users.name = $3
))
RETURNING *;

-- name: GetAllFeeds :many
SELECT feeds.name, feeds.url, users.name AS owner
FROM feeds
JOIN users ON users.id = feeds.user_id;

-- name: GetFeedByUrl :one
SELECT * FROM feeds WHERE feeds.url = $1;

-- name: MarkFeedFetched :exec
UPDATE feeds 
SET updated_at = NOW(),
  last_fetched_at = NOW()
WHERE id = $1;

-- name: GetNextFeedToFetch :one
SELECT id, url FROM feeds
ORDER BY last_fetched_at ASC NULLS FIRST
LIMIT 1;

