-- +goose Up
-- +goose StatementBegin
ALTER TABLE users
ADD hashed_password TEXT NOT NULL
DEFAULT 'unset';
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
SELECT 'down SQL query';
-- +goose StatementEnd
