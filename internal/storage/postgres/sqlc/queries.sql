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
RETURNING reviewer_id;

-- name: PRExists :one
SELECT EXISTS (
  SELECT 1 FROM pull_requests WHERE pull_request_id = $1
);

-- name: GetPR :one
SELECT * FROM pull_requests
WHERE pull_request_id = $1;

-- name: GetPRReviewers :many
SELECT reviewer_id FROM pr_reviewer_assignment
WHERE pr_id = $1 AND replaced_by IS NULL;

-- name: GetActiveTeamMembersExcept :many
SELECT * FROM users
WHERE team_name = $1 AND is_active = true AND user_id != $2;

-- name: MergePR :one
UPDATE pull_requests
SET status = 'MERGED', merged_at = COALESCE(merged_at, now())
WHERE pull_request_id = $1
RETURNING *;

-- name: CheckReviewerAssignment :one
SELECT EXISTS (
  SELECT 1 FROM pr_reviewer_assignment 
  WHERE pr_id = $1 AND reviewer_id = $2 AND replaced_by IS NULL
);

-- name: ReplaceReviewer :one
UPDATE pr_reviewer_assignment
SET replaced_by = $3
WHERE pr_id = $1 AND reviewer_id = $2 AND replaced_by IS NULL
RETURNING *;

-- name: GetPRsByReviewer :many
SELECT DISTINCT pr.* FROM pull_requests pr
JOIN pr_reviewer_assignment pra ON pr.pull_request_id = pra.pr_id
WHERE pra.reviewer_id = $1 AND pra.replaced_by IS NULL;
