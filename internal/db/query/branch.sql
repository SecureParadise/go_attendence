-- name: CreateBranch :one
INSERT INTO branches (
    name,
    code,
    department_id
) VALUES (
    $1, $2, $3
)
RETURNING *;

-- name: GetBranchByCode :one
SELECT * FROM branches
WHERE code = $1 AND deleted_at IS NULL LIMIT 1;