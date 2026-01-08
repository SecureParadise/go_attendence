-- name: GetTeacherByUserID :one
SELECT * FROM teachers
WHERE user_id = $1 AND deleted_at IS NULL LIMIT 1;

-- name: ListTeachersByDepartment :many
SELECT * FROM teachers
WHERE department_id = $1 AND deleted_at IS NULL;

-- name: UpdateTeacherDepartment :one
UPDATE teachers
SET department_id = $2, updated_at = NOW()
WHERE id = $1
RETURNING *;
