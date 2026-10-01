-- +goose Up
-- The initial user is created by the application from ADMIN_EMAIL/ADMIN_PASSWORD.
-- Password is bcrypt-hashed; never insert a plaintext password in a migration.

-- +goose Down
SELECT 1;
