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

-- name: GetUser :one
SELECT * FROM users WHERE name = $1;

-- name: GetFeedByUrl :one
SELECT * FROM feeds WHERE url = $1;

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
