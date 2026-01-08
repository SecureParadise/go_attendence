-- name: CreateAttendance :one
INSERT INTO attendance (
    student_id,
    subject_id,
    teacher_id,
    semester_id,
    date,
    check_in,
    check_out,
    status,
    method,
    remarks
) VALUES (
    $1, $2, $3, $4, $5, $6, $7, $8, $9, $10
) RETURNING *;

-- name: GetAttendance :one
SELECT * FROM attendance
WHERE id = $1 AND deleted_at IS NULL LIMIT 1;

-- name: ListAttendanceByStudent :many
SELECT * FROM attendance
WHERE student_id = $1 AND deleted_at IS NULL
ORDER BY date DESC;

-- name: ListAttendanceBySubject :many
SELECT * FROM attendance
WHERE subject_id = $1 AND deleted_at IS NULL
ORDER BY date DESC;

-- name: UpdateAttendance :one
UPDATE attendance
SET
    check_in = COALESCE(sqlc.narg(check_in), check_in),
    check_out = COALESCE(sqlc.narg(check_out), check_out),
    status = COALESCE(sqlc.narg(status), status),
    remarks = COALESCE(sqlc.narg(remarks), remarks),
    updated_at = NOW()
WHERE id = sqlc.arg(id) AND deleted_at IS NULL
RETURNING *;

-- name: SoftDeleteAttendance :exec
UPDATE attendance
SET deleted_at = NOW()
WHERE id = $1;

-- name: GetAttendanceByStudentSubjectDate :one
SELECT * FROM attendance
WHERE student_id = $1 AND subject_id = $2 AND date = $3 AND deleted_at IS NULL
LIMIT 1;

-- name: ListAttendanceForReport :many
SELECT 
    a.*, 
    s.first_name, s.last_name, s.roll_no,
    sub.name as subject_name,
    t.first_name as teacher_first_name, t.last_name as teacher_last_name
FROM attendance a
JOIN students s ON a.student_id = s.id
JOIN subjects sub ON a.subject_id = sub.id
JOIN teachers t ON a.teacher_id = t.id
WHERE a.semester_id = $1 
  AND a.date >= $2 
  AND a.date <= $3
  AND a.deleted_at IS NULL
ORDER BY a.date DESC, s.roll_no ASC;
