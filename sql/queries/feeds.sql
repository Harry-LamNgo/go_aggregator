-- name: CreateFeed :one
INSERT INTO feeds (id, created_at, updated_at, name, url, user_id)
VALUES (
    $1,
    $2,
    $3,
    $4,
    $5,
    $6
)
RETURNING *;

-- name: GetFeeds :many
SELECT
    feeds.*,
    COUNT(feed_follows.user_id) AS number_users_followed,
    STRING_AGG(users.name, ', ') AS user_followed
FROM feeds
LEFT JOIN feed_follows ON feeds.id = feed_follows.feed_id
LEFT JOIN users ON feed_follows.user_id = users.id
GROUP BY feeds.id, feeds.name, feeds.url;

-- name: GetFeedByURL :one
SELECT * FROM feeds
WHERE url = $1;

-- name: MarkFeedFetched :exec
UPDATE feeds
    SET updated_at = NOW(), last_fetched_at = NOW()
    WHERE id = $1;

-- name: GetNextFeedToFetch :one
SELECT * FROM feeds
    ORDER BY last_fetched_at ASC NULLS FIRST
    LIMIT 1;