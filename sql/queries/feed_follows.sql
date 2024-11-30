-- name: CreateFeedFollow :one
INSERT INTO feed_follows (user_id, feed_id)
VALUES (
  (SELECT users.id FROM users WHERE users.name = $1),
  $2
)
RETURNING feed_follows.*, 
  (SELECT users.name FROM users 
  WHERE feed_follows.user_id = users.id ) AS user_name,
  (SELECT feeds.name FROM feeds
  WHERE feed_follows.feed_id = feeds.id) AS feed_name;


-- name: GetFeedFollowsForUser :many
SELECT feed_follows.*, users.name AS user_name, feeds.name AS feed_name
FROM feed_follows
JOIN users ON users.id = feed_follows.user_id
JOIN feeds ON feeds.id = feed_follows.feed_id
WHERE users.name = $1;

-- name: DeleteFeedFollowForUser :exec
DELETE FROM feed_follows 
WHERE user_id = (
  SELECT users.id FROM users
  WHERE users.name = $1
) AND feed_id = (
  SELECT feeds.id FROM feeds
  WHERE feeds.url = $2
);



