-- name: CreateDepartment :one
INSERT INTO departments (
    name,
    hod_name,
    dhod_name
) VALUES (
$1, $2,$3
)
RETURNING *;

-- name: GetDepartmentByName :one
SELECT * FROM departments
WHERE name = $1 AND deleted_at IS NULL LIMIT 1;
