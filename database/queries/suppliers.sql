-- name: CreateSupplier :one
INSERT INTO suppliers(
    id, name, email, description, phone, country, created_at, updated_at
) 
VALUES($1, $2, $3, $4, $5, $6, $7, $8)
RETURNING *;