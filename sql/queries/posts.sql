-- name: CreatePost :exec
INSERT INTO posts (title, url, description, feed_id, published_at)
VALUES ($1, $2, $3, $4, $5)
ON CONFLICT (url)
DO UPDATE
SET title = EXCLUDED.title,
  description = EXCLUDED.description,
  feed_id = EXCLUDED.feed_id,
  published_at = EXCLUDED.published_at,
  updated_at = CURRENT_TIMESTAMP;

-- name: GetPostsByUser :many
SELECT posts.title, posts.description, posts.url, posts.published_at FROM posts
JOIN feed_follows ON feed_follows.feed_id = posts.feed_id
JOIN users ON users.id = feed_follows.user_id
WHERE users.name = $1
ORDER BY posts.published_at DESC
LIMIT $2;
