-- name: GetAllItems :many
SELECT * FROM items;

-- name: InsertItem :one
INSERT INTO items(id, name, quantity, cost)
VALUES(
    gen_random_uuid(),
    $1,
    $2,
    $3
)
RETURNING *;

-- name: GetItemNameById :one
SELECT name FROM items
WHERE id = $1;