-- +goose Up
-- +goos StatementBegin
CREATE SCHEMA IF NOT EXISTS data;

-- +goos StatementEnd
-- +goose Down
-- +goose StatementBegin
DROP SCHEMA IF EXISTS data;

-- +goose StatementEnd
