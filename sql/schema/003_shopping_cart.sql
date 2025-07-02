-- +goose Up
CREATE TABLE shopping_cart(
    item_id UUID NOT NULL REFERENCES items (id),
    user_id UUID NOT NULL REFERENCES users (id) ON DELETE CASCADE,
    quantity INTEGER NOT NULL,
    cost INTEGER NOT NULL,
    item_name TEXT NOT NULL REFERENCES items (name)
);

-- +goose Down
DROP TABLE shopping_cart;