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

-- name: GetItemById :one
SELECT name, cost, quantity FROM items
WHERE id = $1;

-- name: UpdateItemQuantity :exec
UPDATE items
SET quantity = $1
WHERE id = $2;