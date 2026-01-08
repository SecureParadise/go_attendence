-- name: CreateClassSession :one
INSERT INTO class_sessions (
    subject_id,
    teacher_id,
    semester_id,
    scheduled_start
) VALUES (
    $1, $2, $3, $4
) RETURNING *;

-- name: GetClassSession :one
SELECT * FROM class_sessions
WHERE id = $1 AND deleted_at IS NULL LIMIT 1;

-- name: GetActiveSessionByTeacher :one
SELECT * FROM class_sessions
WHERE teacher_id = $1 
  AND actual_start <= NOW() 
  AND actual_start + INTERVAL '90 minutes' >= NOW()
  AND deleted_at IS NULL
LIMIT 1;

-- name: CreateAttendanceRecord :one
INSERT INTO attendance_records (
    student_id,
    session_id,
    scan_time,
    score,
    status,
    method
) VALUES (
    $1, $2, $3, $4, $5, $6
) RETURNING *;

-- name: UpdateAttendanceRecord :one
UPDATE attendance_records
SET
    scan_time = COALESCE(sqlc.narg(scan_time), scan_time),
    score = COALESCE(sqlc.narg(score), score),
    status = COALESCE(sqlc.narg(status), status),
    method = COALESCE(sqlc.narg(method), method),
    updated_at = NOW()
WHERE id = sqlc.arg(id) AND deleted_at IS NULL
RETURNING *;

-- name: GetAttendanceRecordByStudentAndSession :one
SELECT * FROM attendance_records
WHERE student_id = $1 AND session_id = $2 AND deleted_at IS NULL
LIMIT 1;

-- name: GetActiveSessionBySubject :one
SELECT * FROM class_sessions
WHERE subject_id = $1 
  AND actual_start <= NOW() 
  AND actual_start + INTERVAL '90 minutes' >= NOW()
  AND deleted_at IS NULL
LIMIT 1;

-- name: CreateEnrollment :one
INSERT INTO enrollments (
    student_id,
    branch_id,
    semester_id,
    academic_year
) VALUES (
    $1, $2, $3, $4
) RETURNING *;

-- name: GetStudentAttendancePercentage :one
SELECT 
    sub.name as subject_name,
    SUM(ar.score) as total_score,
    COUNT(cs.id) as total_sessions,
    (SUM(ar.score) / COUNT(cs.id)) * 100 as percentage
FROM attendance_records ar
JOIN class_sessions cs ON ar.session_id = cs.id
JOIN subjects sub ON cs.subject_id = sub.id
WHERE ar.student_id = $1 
  AND cs.semester_id = $2
  AND ar.deleted_at IS NULL
GROUP BY sub.name;

-- name: ListAttendanceRecordsBySession :many
SELECT 
    ar.*, 
    s.first_name, s.last_name, s.roll_no
FROM attendance_records ar
JOIN students s ON ar.student_id = s.id
WHERE ar.session_id = $1 AND ar.deleted_at IS NULL;

-- name: GetActiveSessionForStudent :one
SELECT cs.* FROM class_sessions cs
JOIN enrollments e ON cs.semester_id = e.semester_id
WHERE e.student_id = $1 
  AND e.is_active = TRUE
  AND cs.actual_start <= NOW() 
  AND cs.actual_start + INTERVAL '90 minutes' >= NOW()
  AND cs.deleted_at IS NULL
LIMIT 1;
