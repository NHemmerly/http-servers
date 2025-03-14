-- +goose Up
-- +goose StatementBegin
ALTER TABLE users
ADD is_chirpy_red BOOLEAN NOT NULL
DEFAULT FALSE;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
SELECT 'down SQL query';
-- +goose StatementEnd
