-- name: CreatePost :one
INSERT INTO posts (
    id,
    created_at,
    updated_at,
    title,
    url,
    description,
    published,
    feed_id
)
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
-- name: GetPosts :many
SELECT * FROM posts
ORDER BY updated_at DESC
LIMIT($1);
-- name: GetPostsUser :many
SELECT
p.id,
p.created_at,
p.updated_at,
p.title,
p.url,
p.description,
p.published
FROM posts p
INNER JOIN feed_follow ff ON p.feed_id=ff.feed_id
WHERE
ff.user_id = $1
ORDER BY p.updated_at DESC
LIMIT $2;