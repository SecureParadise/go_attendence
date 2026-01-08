-- Create class_sessions table
CREATE TABLE class_sessions (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    subject_id UUID NOT NULL REFERENCES subjects(id),
    teacher_id UUID NOT NULL REFERENCES teachers(id),
    semester_id UUID NOT NULL REFERENCES semesters(id),
    scheduled_start TIMESTAMPTZ NOT NULL,
    actual_start TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    deleted_at TIMESTAMPTZ
);

-- Create enrollments table to track students across academic years/semesters
CREATE TABLE enrollments (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    student_id UUID NOT NULL REFERENCES students(id),
    branch_id UUID NOT NULL REFERENCES branches(id),
    semester_id UUID NOT NULL REFERENCES semesters(id),
    academic_year VARCHAR(10) NOT NULL, -- e.g. "2080"
    is_active BOOLEAN NOT NULL DEFAULT TRUE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    deleted_at TIMESTAMPTZ,
    UNIQUE (student_id, academic_year, semester_id)
);

-- Create attendance_records table with scoring
CREATE TABLE attendance_records (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    student_id UUID NOT NULL REFERENCES students(id),
    session_id UUID NOT NULL REFERENCES class_sessions(id),
    scan_time TIMESTAMPTZ,
    score DECIMAL(3, 2) NOT NULL DEFAULT 0.0,
    status attendance_status NOT NULL DEFAULT 'absent',
    method attendance_method NOT NULL DEFAULT 'manual',
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    deleted_at TIMESTAMPTZ,
    UNIQUE (student_id, session_id)
);

-- Add department_id to users to support HOD logic
ALTER TABLE users ADD COLUMN department_id UUID REFERENCES departments(id);

-- Create indexes for performance
CREATE INDEX ON class_sessions (subject_id);
CREATE INDEX ON class_sessions (teacher_id);
CREATE INDEX ON class_sessions (scheduled_start);
CREATE INDEX ON attendance_records (student_id);
CREATE INDEX ON attendance_records (session_id);
CREATE INDEX ON enrollments (student_id);
CREATE INDEX ON enrollments (academic_year);
