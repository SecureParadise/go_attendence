-- name: CreateStudent :one
INSERT INTO students(
    roll_no,
    first_name,
    middle_name ,
    last_name,
    image ,
    batch,
    user_id ,
    branch_id ,
    current_semester_id
) VALUES (
$1,$2,$3,$4,$5,$6,$7,$8,$9
)
RETURNING *;

-- name: GetStudentByRollNo :one
SELECT * FROM students
WHERE roll_no = $1 AND deleted_at IS NULL
LIMIT 1;