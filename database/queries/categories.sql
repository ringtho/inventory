-- name: CreateCategory :one
INSERT INTO categories (
    id, name, description, created_at, updated_at, created_by
)
VALUES ($1,$2,$3,$4,$5, $6)
RETURNING *;

-- name: GetCategories :many
SELECT * FROM categories;