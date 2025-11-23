// Package users provides handlers and service logic for managing users.
package users

import (
	"context"

	repo "github.com/Joskmo/avito-trainee-assignment-api/internal/storage/postgres/sqlc"
	"github.com/jackc/pgx/v5/pgxpool"
)

// Service defines the interface for the users service.
type Service interface {
	SetUserActivity(ctx context.Context, userActivityParams repo.SetUserActivityParams) (repo.User, error)
}

// Handler handles HTTP requests for the users service.
type Handler struct {
	service Service
}

// NewHandler creates a new users handler.
func NewHandler(service Service) *Handler {
	return &Handler{
		service: service,
	}
}

type svc struct {
	repo *repo.Queries
	db   *pgxpool.Pool
}

// NewService creates a new users service.
func NewService(repo *repo.Queries, db *pgxpool.Pool) Service {
	return &svc{
		repo: repo,
		db:   db,
	}
}

// SetUserActivityRequest represents the request for setting user activity.
type SetUserActivityRequest struct {
	UserID   string `json:"user_id"`
	IsActive *bool  `json:"is_active"`
}

// SetUserActivityResponse represents the response for setting user activity.
type SetUserActivityResponse struct {
	User repo.User `json:"user"`
}
