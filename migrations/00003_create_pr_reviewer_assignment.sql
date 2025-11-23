-- +goose Up
-- +goose StatementBegin
CREATE SEQUENCE IF NOT EXISTS assignment_seq;
CREATE TABLE IF NOT EXISTS pr_reviewer_assignment (
    assignment_id TEXT PRIMARY KEY DEFAULT ('a' || nextval('assignment_seq')),
    pr_id TEXT NOT NULL REFERENCES pull_requests(pull_request_id),
    reviewer_id TEXT NOT NULL REFERENCES users(user_id),
    assigned_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    replaced_by TEXT REFERENCES users(user_id),
    UNIQUE(pr_id, reviewer_id)
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS pr_reviewer_assignment;
DROP SEQUENCE IF EXISTS assignment_seq;
-- +goose StatementEnd
