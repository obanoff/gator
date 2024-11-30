-- +goose Up
CREATE TABLE feed_follows (
  id SERIAL PRIMARY KEY,
  user_id UUID REFERENCES users (id) ON DELETE CASCADE,
  feed_id INT REFERENCES feeds (id) ON DELETE CASCADE,
  created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  UNIQUE (user_id, feed_id)
);

-- +goose Down
DROP TABLE feed_follows;
