-- name: CreateUser :one
INSERT INTO users (
  email,
  password_hash,
  user_role
) VALUES (
  $1, $2,$3 
)
RETURNING *;
-- name: GetUserByEmail :one
SELECT * FROM users
WHERE email = $1 AND deleted_at IS NULL LIMIT 1;

-- name: UpdateUserProfileCompleted :one
UPDATE users
SET is_profile_completed = $2
WHERE id = $1
RETURNING *;
