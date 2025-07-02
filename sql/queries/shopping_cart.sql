-- name: GetShoppingCart :many
SELECT * FROM shopping_cart
WHERE user_id = $1;

-- name: AddItemInCart :exec
INSERT INTO shopping_cart(item_id, user_id, quantity, cost, item_name)
VALUES (
    $1,
    $2,
    $3,
    $4,
    $5
);