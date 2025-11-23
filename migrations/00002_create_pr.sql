-- +goose Up
-- +goose StatementBegin
CREATE SEQUENCE IF NOT EXISTS pr_sequence;
DROP TYPE IF EXISTS pr_status_enum;
CREATE TYPE pr_status_enum AS ENUM ('OPEN', 'MERGED');
CREATE TABLE IF NOT EXISTS pull_requests (
    pull_request_id TEXT PRIMARY KEY DEFAULT ('pr-' || nextval('pr_sequence')),
    pull_request_name TEXT NOT NULL,
    author_id TEXT NOT NULL REFERENCES users(user_id),
    status pr_status_enum DEFAULT 'OPEN',
    merged_at TIMESTAMPTZ
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS pull_requests;
DROP TYPE IF EXISTS pr_status_enum;
DROP SEQUENCE IF EXISTS pr_sequence;
-- +goose StatementEnd
