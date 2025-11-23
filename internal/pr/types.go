// Package pr provides handlers and service logic for managing pull requests.
package pr

import (
	"context"

	"github.com/Joskmo/avito-trainee-assignment-api/internal/domain"
	repo "github.com/Joskmo/avito-trainee-assignment-api/internal/storage/postgres/sqlc"
	"github.com/jackc/pgx/v5/pgxpool"
)

// Service defines the interface for the PR service.
type Service interface {
	CreatePR(ctx context.Context, createPRParams repo.CreatePRParams) (CreatePRResponse, error)
	MergePR(ctx context.Context, prID string) (Response, error)
	ReassignReviewer(ctx context.Context, prID, oldUserID string) (ReassignResponse, error)
	GetUserReviews(ctx context.Context, userID string) (UserReviewsResponse, error)
}

// Handler handles HTTP requests for the PR service.
type Handler struct {
	service Service
}

// NewHandler creates a new PR handler.
func NewHandler(service Service) *Handler {
	return &Handler{
		service: service,
	}
}

type svc struct {
	repo *repo.Queries
	db   *pgxpool.Pool
}

// NewService creates a new PR service.
func NewService(repo *repo.Queries, db *pgxpool.Pool) Service {
	return &svc{
		repo: repo,
		db:   db,
	}
}

// WithReviewers represents a PR with its assigned reviewers.
type WithReviewers = domain.PRWithReviewers

// CreatePRResponse represents the response for creating a PR.
type CreatePRResponse struct {
	PR WithReviewers `json:"pr"`
}

// Response represents the response for getting a PR.
type Response struct {
	PR WithReviewers `json:"pr"`
}

// ReassignResponse represents the response for reassigning a reviewer.
type ReassignResponse struct {
	PR         WithReviewers `json:"pr"`
	ReplacedBy string        `json:"replaced_by"`
}

// Short represents a short version of a PR.
type Short struct {
	PullRequestID   string `json:"pull_request_id"`
	PullRequestName string `json:"pull_request_name"`
	AuthorID        string `json:"author_id"`
	Status          string `json:"status"`
}

// UserReviewsResponse represents the response for getting user reviews.
type UserReviewsResponse struct {
	UserID       string  `json:"user_id"`
	PullRequests []Short `json:"pull_requests"`
}
