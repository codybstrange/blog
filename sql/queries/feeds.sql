-- name: AddFeed :one
INSERT INTO feeds(id, created_at, updated_at, name, url, user_id)
VALUES (
  $1, $2, $3, $4, $5, $6
)
RETURNING *;

-- name: GetAllFeeds :many
SELECT name, url, user_id FROM feeds;

-- name: GetFeedByURL :one
SELECT * FROM feeds
  WHERE feeds.url = $1;

-- name: MarkFeedFetched :exec
UPDATE feeds
SET last_fetched_at = $2, updated_at = $2
WHERE feeds.ID = $1;

-- name: GetNextFeedToFetch :one
SELECT * FROM feeds
ORDER BY last_fetched_at ASC NULLS FIRST;
