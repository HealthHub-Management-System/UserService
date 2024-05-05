-- +goose Up
-- +goose StatementBegin
ALTER TABLE users ALTER COLUMN password SET DATA TYPE BYTEA USING (password::BYTEA);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE users ALTER COLUMN password SET DATA TYPE VARCHAR(255) USING (password::VARCHAR(255));
-- +goose StatementEnd
