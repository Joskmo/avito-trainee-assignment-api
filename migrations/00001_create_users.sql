-- +goose Up
-- +goose StatementBegin
CREATE SEQUENCE IF NOT EXISTS user_seq;
CREATE TABLE IF NOT EXISTS users (
    user_id TEXT PRIMARY KEY DEFAULT ('u' || nextval('user_seq')),
    username TEXT NOT NULL,
    is_active BOOLEAN NOT NULL DEFAULT TRUE,
    team_name TEXT NOT NULL
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS users;
DROP SEQUENCE IF EXISTS user_seq;
-- +goose StatementEnd
