-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS data.users (
    id SERIAL PRIMARY KEY,
    telegram_id BIGINT,
    chat_id BIGINT,
    is_running BOOL NOT NULL,
    location varchar(127),
    window_floor TIME NOT NULL,
    window_ceil TIME NOT NULL CHECK (window_ceil > window_floor),
    is_deleted BOOL
);

-- +goose StatementEnd
-- +goose Down
-- +goose StatementBegin
DROP TABLE datga.users;

-- +goose StatementEnd
