-- +goose Up
CREATE TABLE items(
    id UUID PRIMARY KEY,
    name TEXT NOT NULL UNIQUE,
    quantity INTEGER NOT NULL,
    cost INTEGER NOT NULL
);

-- +goose Down
DROP TABLE items;