package users

import (
	"context"

	repo "github.com/Joskmo/avito-trainee-assignment-api/internal/storage/postgres/sqlc"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Service interface {
	SetUserActivity(ctx context.Context, userActivityParams repo.SetUserActivityParams) (repo.User, error)
}

type handler struct {
	service Service
}

func NewHandler(service Service) *handler {
	return &handler{
		service: service,
	}
}

type svc struct {
	repo *repo.Queries
	db   *pgxpool.Pool
}

func NewService(repo *repo.Queries, db *pgxpool.Pool) Service {
	return &svc{
		repo: repo,
		db:   db,
	}
}

type SetUserActivityRequest struct {
	UserID   string `json:"user_id"`
	IsActive *bool  `json:"is_active"`
}

type SetUserActivityResponse struct {
	User repo.User `json:"user"`
}
