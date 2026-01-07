CREATE TABLE IF NOT EXISTS students (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),

    -- Academic identity
    roll_no VARCHAR(50) NOT NULL UNIQUE,

    -- Personal info
    first_name VARCHAR(50) NOT NULL,
    middle_name VARCHAR(50),
    last_name VARCHAR(50) NOT NULL,
    image VARCHAR(255),

    -- Academic metadata
    batch VARCHAR(50),

    -- Relations
    user_id UUID NOT NULL UNIQUE,
    branch_id UUID NOT NULL,
    current_semester_id UUID,

    -- RFID & Biometric
    rfid_tag_id VARCHAR(255) ,
    fingerprint_hash VARCHAR(255),

    -- Timestamps
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),

    -- Foreign keys
    CONSTRAINT fk_students_user
        FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,

    CONSTRAINT fk_students_branch
        FOREIGN KEY (branch_id) REFERENCES branches(id),

    CONSTRAINT fk_students_current_semester
        FOREIGN KEY (current_semester_id) REFERENCES semesters(id)
);


-- Teacher
CREATE TABLE IF NOT EXISTS teachers (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),

    -- Academic identity
    card_no VARCHAR(50) NOT NULL UNIQUE,

    -- Personal info
    first_name VARCHAR(50) NOT NULL,
    middle_name VARCHAR(50),
    last_name VARCHAR(50) NOT NULL,
    image VARCHAR(255),

    -- Relations
    user_id UUID NOT NULL UNIQUE,
    department_id UUID NOT NULL,

    -- RFID & Biometric
    rfid_tag_id VARCHAR(255),
    fingerprint_hash VARCHAR(255),

    -- Timestamps
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),

    -- Foreign keys
    CONSTRAINT fk_teachers_user
        FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,

    CONSTRAINT fk_teachers_department
        FOREIGN KEY (department_id) REFERENCES departments(id)
);
