package teams

import (
	"context"

	repo "github.com/Joskmo/avito-trainee-assignment-api/internal/storage/postgres/sqlc"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Service interface {
	GetTeamByName(ctx context.Context, teamName string) ([]repo.User, error)
	CreateTeam(ctx context.Context, tempTeam tempTeamParams) ([]repo.User, error)
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

type tempUserParams struct {
	UserID   string `json:"user_id"`
	Username string `json:"username"`
	IsActive bool   `json:"is_active"`
}

type tempTeamParams struct {
	TeamName string           `json:"team_name"`
	Members  []tempUserParams `json:"members"`
}

type CreateTeamResponse struct {
	Team TeamResponse `json:"team"`
}

type TeamMember struct {
	UserID   string `json:"user_id"`
	Username string `json:"username"`
	IsActive bool   `json:"is_active"`
}

type TeamResponse struct {
	TeamName string       `json:"team_name"`
	Members  []TeamMember `json:"members"`
}
