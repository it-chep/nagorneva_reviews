-- +goose Up
CREATE TABLE doctor_courses (
    doctor_id BIGINT NOT NULL REFERENCES doctors(id) ON DELETE CASCADE,
    course_id BIGINT NOT NULL REFERENCES courses(id) ON DELETE RESTRICT,
    PRIMARY KEY (doctor_id, course_id)
);

CREATE INDEX doctor_courses_course_idx ON doctor_courses (course_id, doctor_id);

-- +goose Down
DROP TABLE doctor_courses;
