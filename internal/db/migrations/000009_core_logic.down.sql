DROP TABLE IF EXISTS attendance_records;
DROP TABLE IF EXISTS enrollments;
DROP TABLE IF EXISTS class_sessions;

ALTER TABLE users DROP COLUMN IF EXISTS department_id;
