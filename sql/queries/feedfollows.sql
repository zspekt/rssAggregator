-- name: CreateFeedFollow :one
INSERT INTO feedfollow (id, created_at, updated_at, feed_id, user_id)
    VALUES ($1, $2, $3, $4, $5)
RETURNING
    *;

-- name: DeleteFeedFollow :exec
DELETE FROM feedfollow
WHERE id = $1;

-- name: GetFeedFollowsByUser :many
SELECT
    *
FROM
    feedfollow
WHERE
    user_id = $1;
