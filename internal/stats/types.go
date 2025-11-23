package stats

import (
	"context"

	repo "github.com/Joskmo/avito-trainee-assignment-api/internal/storage/postgres/sqlc"
	"github.com/jackc/pgx/v5/pgxpool"
)

// Service defines the interface for the stats service.
type Service interface {
	GetStats(ctx context.Context) (Response, error)
}

// Handler handles HTTP requests for the stats service.
type Handler struct {
	service Service
}

// NewHandler creates a new stats handler.
func NewHandler(service Service) *Handler {
	return &Handler{
		service: service,
	}
}

type svc struct {
	repo *repo.Queries
	db   *pgxpool.Pool
}

// NewService creates a new stats service.
func NewService(repo *repo.Queries, db *pgxpool.Pool) Service {
	return &svc{
		repo: repo,
		db:   db,
	}
}

// ReviewerStat represents statistics for a single reviewer.
type ReviewerStat struct {
	ReviewerID      string `json:"reviewer_id"`
	AssignmentCount int64  `json:"assignment_count"`
}

// PRStatusStat represents statistics for PR statuses.
type PRStatusStat struct {
	Status string `json:"status"`
	Count  int64  `json:"count"`
}

// Response represents the response for the stats endpoint.
type Response struct {
	TopReviewers     []ReviewerStat `json:"top_reviewers"`
	PRStatus         []PRStatusStat `json:"pr_status_distribution"`
	TotalActiveUsers int64          `json:"total_active_users"`
}
