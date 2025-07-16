-- name: GetAllUsers :many
SELECT * FROM users;

-- name: CreateNewUser :one
INSERT INTO users (id, created_at, updated_at, email, hashed_password)
VALUES (
    gen_random_uuid(),
    NOW(),
    NOW(),
    $1,
    $2
)
RETURNING *;

-- name: GetUserPassword :one
SELECT hashed_password, id FROM users
WHERE email = $1;

-- name: GetUserEmail :one
SELECT email FROM users
WHERE id = $1;

-- name: IsEmailExists :one
SELECT EXISTS(
    SELECT * FROM users
    WHERE email = $1
);