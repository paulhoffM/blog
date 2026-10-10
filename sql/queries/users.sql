-- name: CreateUser :one
INSERT INTO users (id, created_at, updated_at, name)
VALUES (
    $1,
    $2,
    $3,
    $4
)
RETURNING *;

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

-- name: CreatePost :one
INSERT INTO posts (id, created_at, updated_at, title, url, description, published_at, feed_id)
VALUES (
    $1,
    $2,
    $3,
    $4,
    $5,
    $6,
    $7,
    $8
)
RETURNING *;

-- name: GetUser :one
SELECT * FROM users WHERE name = $1;

-- name: GetFeedByUrl :one
SELECT * FROM feeds WHERE url = $1;

-- name: GetPosts :many
SELECT * FROM posts ORDER BY published_at DESC LIMIT $1;

-- name: Reset :exec
DELETE FROM users;

-- name: List :many
SELECT name FROM users;

-- name: ListFeeds :many
SELECT 
    feeds.name AS feed_name,
    feeds.url, 
    users.name AS user_name
FROM 
    feeds
    INNER JOIN users ON feeds.user_id = users.id;

-- name: CreateFeedFollow :many
WITH inserted_feed_follow AS (
INSERT INTO feed_follows(id, created_at, updated_at, feed_id, user_id)
VALUES(
    $1,
    $2,
    $3,
    $4,
    $5
)
RETURNING *
)
SELECT
    inserted_feed_follow.*,
    feeds.name AS feed_name,
    users.name AS user_name
FROM inserted_feed_follow
INNER JOIN feeds ON inserted_feed_follow.feed_id = feeds.id
INNER JOIN users ON inserted_feed_follow.user_id = users.id; 

-- name: GetFeedFollowsForUser :many
SELECT
    feed_follows.*,
    feeds.name AS feed_name,
    users.name AS user_name
FROM
    feed_follows
     INNER JOIN users ON feed_follows.user_id = users.id 
     INNER JOIN feeds ON feed_follows.feed_id = feeds.id
WHERE
    feed_follows.user_id = $1;


-- name: DeleteFollow :exec
DELETE FROM 
    feed_follows
USING
    feeds
WHERE
    feed_follows.feed_id = feeds.id
    AND feeds.url = $1
    AND feed_follows.user_id = $2;

-- name: MarkFeedFetched :one
UPDATE feeds
SET updated_at = NOW(), last_fetched_at = NOW()
WHERE id = $1
RETURNING *;

-- name: GetNextFeedToFetch :one
SELECT *
FROM feeds
ORDER BY last_fetched_at NULLS FIRST
LIMIT 1;