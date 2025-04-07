-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS data.reminders (
    id SERIAL PRIMARY KEY,
    user_id INT NOT NULL,
    text TEXT NOT NULL,
    tag VARCHAR(255),
    prompt TEXT,
    frequency INTERVAL NOT NULL,
    next_reminder TIMESTAMP
    WITH
        TIME ZONE NOT NULL,
        is_deleted BOOL,
        FOREIGN KEY (user_id) REFERENCES data.users (id)
);

-- +goose StatementEnd
-- +goose Down
-- +goose StatementBegin
DROP TABLE data.reminders;

-- +goose StatementEnd
