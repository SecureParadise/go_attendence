CREATE TYPE attendance_status AS ENUM ('present', 'absent', 'late', 'excused');
CREATE TYPE attendance_method AS ENUM ('manual', 'qr', 'face', 'rfid', 'fingerprint');

CREATE TABLE attendance (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    student_id UUID NOT NULL REFERENCES students(id),
    subject_id UUID NOT NULL REFERENCES subjects(id),
    teacher_id UUID NOT NULL REFERENCES teachers(id),
    semester_id UUID NOT NULL REFERENCES semesters(id),
    date DATE NOT NULL DEFAULT CURRENT_DATE,
    check_in TIMESTAMPTZ,
    check_out TIMESTAMPTZ,
    status attendance_status NOT NULL DEFAULT 'present',
    method attendance_method NOT NULL DEFAULT 'manual',
    remarks TEXT,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    deleted_at TIMESTAMPTZ,
    
    -- Ensure a student can only have one attendance record per subject per day
    UNIQUE (student_id, subject_id, date)
);

CREATE INDEX ON attendance (student_id);
CREATE INDEX ON attendance (subject_id);
CREATE INDEX ON attendance (teacher_id);
CREATE INDEX ON attendance (date);
