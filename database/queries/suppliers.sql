-- name: CreateSupplier :one
INSERT INTO suppliers(
    id, name, email, description, phone, country, created_at, updated_at
) 
VALUES($1, $2, $3, $4, $5, $6, $7, $8)
RETURNING *;

-- name: GetAllSuppliers :many
SELECT * FROM suppliers;

-- name: GetSupplierById :one
SELECT * FROM suppliers WHERE id=$1;

-- name: DeleteSupplier :exec
DELETE FROM suppliers WHERE id=$1;

-- name: UpdateSupplier :one
UPDATE suppliers
SET 
name = $2,
email = $3,
description = $4,
phone = $5,
country = $6,
updated_at = $7
WHERE id = $1
RETURNING *;
