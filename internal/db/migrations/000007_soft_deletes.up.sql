-- Add deleted_at to all relevant tables
ALTER TABLE users ADD COLUMN deleted_at TIMESTAMPTZ;
ALTER TABLE departments ADD COLUMN deleted_at TIMESTAMPTZ;
ALTER TABLE branches ADD COLUMN deleted_at TIMESTAMPTZ;
ALTER TABLE semesters ADD COLUMN deleted_at TIMESTAMPTZ;
ALTER TABLE students ADD COLUMN deleted_at TIMESTAMPTZ;
ALTER TABLE teachers ADD COLUMN deleted_at TIMESTAMPTZ;
ALTER TABLE subjects ADD COLUMN deleted_at TIMESTAMPTZ;
