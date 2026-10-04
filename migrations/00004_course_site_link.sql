-- +goose Up
ALTER TABLE courses
    ADD COLUMN site_link TEXT NOT NULL DEFAULT '';

-- +goose Down
ALTER TABLE courses
    DROP COLUMN site_link;
