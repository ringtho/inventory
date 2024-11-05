-- +goose Up
CREATE TABLE categories (
    id UUID PRIMARY KEY,
    created_at TIMESTAMP NOT NULL,
    updated_at TIMESTAMP NOT NULL,
    name VARCHAR(100) NOT NULL UNIQUE,
    description TEXT
);


-- +goose Down
DROP TABLE categories;