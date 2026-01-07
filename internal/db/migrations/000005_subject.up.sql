-- subject
CREATE TABLE IF NOT EXISTS subjects (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name VARCHAR(255) NOT NULL,
    code VARCHAR(10) NOT NULL,
    is_lab BOOLEAN NOT NULL DEFAULT FALSE,
    credits INTEGER,

    -- Relations
    branch_id UUID NOT NULL,
    semester_id UUID NOT NULL,
    teacher_id UUID NOT NULL,

    -- Timestamps
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),

    -- Foreign keys
    CONSTRAINT fk_subjects_branch
        FOREIGN KEY (branch_id) REFERENCES branches(id),

    CONSTRAINT fk_subjects_semester
        FOREIGN KEY (semester_id) REFERENCES semesters(id),

    CONSTRAINT fk_subjects_teacher
        FOREIGN KEY (teacher_id) REFERENCES teachers(id)
);
