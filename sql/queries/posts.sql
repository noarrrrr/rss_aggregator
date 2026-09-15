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

-- name: GetPostsForUser :many
SELECT * FROM posts
INNER JOIN feeds ON feed_id = feeds.id
INNER JOIN users ON $1 = users.id
ORDER BY published_at
LIMIT $2;