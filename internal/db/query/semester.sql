-- name: CreateSemester :one
INSERT INTO semesters(
    number,
    name,
    branch_id
)VALUES(
    $1,$2,$3
)
RETURNING *;