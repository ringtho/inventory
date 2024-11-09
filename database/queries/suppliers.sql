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