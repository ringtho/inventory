-- name: CreateProduct :one
INSERT INTO products(
    id,
    name,
    description,
    price,
    stock_level,
    category_id,
    supplier_id,
    sku,
    created_at,
    updated_at
)
VALUES($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
RETURNING *;

-- name: GetProducts :many
SELECT * FROM products;

-- name: GetProduct :one
SELECT * FROM products WHERE id = $1;

-- name: DeleteProduct :exec
DELETE FROM products WHERE id = $1;

-- name: UpdateProduct :one
UPDATE products
SET
name = $2,
description = $3,
price = $4,
stock_level = $5,
category_id = $6,
supplier_id = $7,
sku = $8,
updated_at = $9
WHERE id = $1
RETURNING *;