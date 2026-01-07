-- name: CreateUser :one
INSERT INTO users (
  email,
  password_hash,
  user_role
) VALUES (
  $1, $2,$3 
)
RETURNING *;