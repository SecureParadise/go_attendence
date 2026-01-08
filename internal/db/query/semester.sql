-- name: CreateSemester :one
INSERT INTO semesters(
    number,
    name,
    branch_id
)VALUES(
    $1,$2,$3
)
RETURNING *;

-- name: GetSemesterByNumberAndBranch :one
SELECT * FROM semesters
WHERE number = $1 AND branch_id = $2 AND deleted_at IS NULL
LIMIT 1;