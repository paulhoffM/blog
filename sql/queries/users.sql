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

-- name: GetFeed :one
SELECT * FROM feeds WHERE name = $1;

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