-- name: CreateUser :one
INSERT INTO users(id, username, email, name, password, role, profile_picture_url, created_at, updated_at)
VALUES($1, $2, $3, $4, $5, $6, $7, $8, $9)
RETURNING id, username, email, name, role, profile_picture_url, created_at, updated_at;
