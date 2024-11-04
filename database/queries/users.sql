-- name: CreateUser :one
INSERT INTO users(
    id, username, email, name, password, role, profile_picture_url, created_at, updated_at
)
VALUES($1, $2, $3, $4, $5, $6, $7, $8, $9)
RETURNING id, username, email, name, role, profile_picture_url, created_at, updated_at;

-- name: GetUserById :one
SELECT * FROM users WHERE id = $1;

-- name: GetUserByEmail :one
SELECT * FROM users WHERE email = $1;

-- name: GetAllUsers :many
SELECT 
    id, username, email, name, role, profile_picture_url, created_at, updated_at
FROM users;

-- name: DeleteUser :exec
DELETE FROM users WHERE id = $1 AND role != 'admin';
