-- +goose Up
CREATE TABLE users (
    id BIGSERIAL PRIMARY KEY,
    email TEXT NOT NULL UNIQUE,
    password TEXT NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE TABLE cities (
    id BIGSERIAL PRIMARY KEY,
    name TEXT NOT NULL UNIQUE,
    lat DOUBLE PRECISION NOT NULL,
    lon DOUBLE PRECISION NOT NULL
);

CREATE TABLE specialties (
    id BIGSERIAL PRIMARY KEY,
    name TEXT NOT NULL UNIQUE
);

CREATE TABLE doctors (
    id BIGSERIAL PRIMARY KEY,
    name TEXT NOT NULL,
    full_name TEXT NOT NULL,
    photo TEXT NOT NULL DEFAULT '',
    personal_data_consent BOOLEAN NOT NULL DEFAULT false,
    city_id BIGINT NOT NULL REFERENCES cities(id) ON DELETE RESTRICT,
    specialty_id BIGINT NOT NULL REFERENCES specialties(id) ON DELETE RESTRICT,
    lat DOUBLE PRECISION NOT NULL,
    lon DOUBLE PRECISION NOT NULL,
    is_active BOOLEAN NOT NULL DEFAULT true,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE TABLE courses (
    id BIGSERIAL PRIMARY KEY,
    name TEXT NOT NULL UNIQUE
);

CREATE TABLE reviews (
    id BIGSERIAL PRIMARY KEY,
    doctor_id BIGINT NOT NULL REFERENCES doctors(id) ON DELETE CASCADE,
    rating SMALLINT NOT NULL CHECK (rating BETWEEN 1 AND 5),
    comment TEXT NOT NULL DEFAULT '',
    course_id BIGINT NOT NULL REFERENCES courses(id) ON DELETE RESTRICT,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    is_active BOOLEAN NOT NULL DEFAULT true
);

CREATE INDEX doctors_map_idx ON doctors (is_active, lat, lon);
CREATE INDEX reviews_doctor_idx ON reviews (doctor_id, is_active, created_at DESC);
CREATE INDEX reviews_course_idx ON reviews (course_id, is_active, created_at DESC);

-- +goose Down
DROP TABLE reviews;
DROP TABLE courses;
DROP TABLE doctors;
DROP TABLE specialties;
DROP TABLE cities;
DROP TABLE users;
