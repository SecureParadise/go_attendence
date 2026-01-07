ALTER TABLE semesters ADD CONSTRAINT unique_semester_number_branch UNIQUE (number, branch_id);
