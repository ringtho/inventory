-- name: CreateCategory :one
INSERT INTO categories (
    id, name, description, created_at, updated_at
)
VALUES ($1,$2,$3,$4,$5)
RETURNING *;