-- name: InsertNewRefreshToken :exec
INSERT INTO refresh_tokens(token, user_id, expires_at)
VALUES(
    $1,
    $2,
    $3
);

-- name: GetRefreshToken :one
SELECT * FROM refresh_tokens
WHERE token = $1;