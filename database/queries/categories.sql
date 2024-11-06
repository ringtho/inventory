-- name: CreateCategory :one
INSERT INTO categories (
    id, name, description, created_at, updated_at, created_by
)
VALUES ($1,$2,$3,$4,$5, $6)
RETURNING *;

-- name: GetCategories :many
SELECT * FROM categories;

-- name: GetCategoryById :one
SELECT * FROM categories WHERE id = $1;

-- name: UpdateCategory :one
UPDATE categories
SET
name = $2,
description = $3,
updated_at = $4
WHERE id = $1
RETURNING *;

-- name: DeleteCategory :exec
DELETE FROM categories WHERE id = $1;