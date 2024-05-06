-- +goose Up
-- +goose StatementBegin
CREATE TYPE role AS ENUM ('patient', 'doctor', 'admin');
ALTER TABLE users ADD COLUMN role role NOT NULL DEFAULT 'patient';
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE users DROP COLUMN IF EXISTS role;
DROP TYPE IF EXISTS role;
-- +goose StatementEnd
