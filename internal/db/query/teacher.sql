-- name: CreateTeacher :one
INSERT INTO teachers(
    card_no,
    first_name,
    middle_name ,
    last_name,
    image ,
    user_id ,
    department_id 
) VALUES (
$1,$2,$3,$4,$5,$6,$7
)
RETURNING *;

-- name: GetTeacherByCardNo :one
SELECT * FROM teachers
WHERE card_no = $1 AND deleted_at IS NULL
LIMIT 1;
