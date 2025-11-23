// Package domain contains shared domain models.
package domain

// PRWithReviewers represents a PR with its assigned reviewers.
type PRWithReviewers struct {
	PullRequestID     string   `json:"pull_request_id"`
	PullRequestName   string   `json:"pull_request_name"`
	AuthorID          string   `json:"author_id"`
	Status            string   `json:"status"`
	AssignedReviewers []string `json:"assigned_reviewers"`
}
