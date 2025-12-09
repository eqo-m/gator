-- name: CreateFeed :one
INSERT INTO feeds (id, created_at, updated_at, name,url,user_id)
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
f.name AS feed_name,
f.url,
u.name AS user_name
FROM 
feeds f 
INNER JOIN users u ON f.user_id = u.id;
-- name: CreateFeedFollow :one
WITH inserted_feed_follow AS (
    INSERT INTO feed_follow(
        id,
        created_at,
        updated_at,
        user_id,
        feed_id
    )  VALUES (
        $1,
        $2,
        $3,
        $4,
        $5
    )
    RETURNING *
)
SELECT inserted_feed_follow.*,
feeds.name AS feed_name,
users.name AS user_name
FROM inserted_feed_follow
INNER JOIN users ON inserted_feed_follow.user_id = users.id
INNER JOIN feeds ON inserted_feed_follow.feed_id = feeds.id;
-- name: GetFeedByURL :one
SELECT * FROM feeds
WHERE url = $1;
-- name: GetFeedFollowsForUser :many
SELECT 
ff.id,
ff.created_at,
f.id AS feed_id,
f.name AS feed_name,
f.url AS feed_url
FROM
feed_follow ff
INNER JOIN 
feeds f ON ff.feed_id = f.id
WHERE
ff.user_id = $1;
-- name: UnfollowFeed :one
DELETE FROM feed_follow
WHERE
user_id=$1 AND feed_id=$2
RETURNING *;
-- name: MarkFeedFetched :one
UPDATE feeds
SET last_fetched_at=NOW(),
updated_at=NOW()
WHERE id=$1
RETURNING * ;
-- name: GetNextFeedToFetch :one
SELECT
id,
created_at,
updated_at,
name,
url,
user_id 
FROM feeds
ORDER BY last_fetched_at
NULLS FIRST
LIMIT 1;

