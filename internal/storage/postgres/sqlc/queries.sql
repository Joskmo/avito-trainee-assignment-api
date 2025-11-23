-- name: CreateUser :one
INSERT INTO users (user_id, username, is_active, team_name)
VALUES ($1, $2, $3, $4)
ON CONFLICT (user_id) DO UPDATE
SET 
    username = EXCLUDED.username,
    is_active = EXCLUDED.is_active,
    team_name = EXCLUDED.team_name
RETURNING *;

-- name: TeamExists :one
SELECT EXISTS (
  SELECT 1 FROM users WHERE team_name = $1
);

-- name: GetUser :one
SELECT * FROM users
WHERE user_id = $1;

-- name: GetTeam :many
SELECT * FROM users
WHERE team_name = $1;

-- name: SetUserActivity :one
UPDATE users
SET is_active = $2
WHERE user_id = $1
RETURNING *;

-- name: CreatePR :one
INSERT INTO pull_requests (pull_request_id, pull_request_name, author_id)
VALUES ($1, $2, $3)
ON CONFLICT (pull_request_id) DO NOTHING
RETURNING *;

-- name: AssignReviewer :one
INSERT INTO pr_reviewer_assignment (pr_id, reviewer_id)
VALUES ($1, $2)
RETURNING *;
