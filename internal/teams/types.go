// Package teams provides handlers and service logic for managing teams.
package teams

import (
	"context"

	"github.com/Joskmo/avito-trainee-assignment-api/internal/domain"
	repo "github.com/Joskmo/avito-trainee-assignment-api/internal/storage/postgres/sqlc"
	"github.com/jackc/pgx/v5/pgxpool"
)

// Service defines the interface for the teams service.
type Service interface {
	GetTeamByName(ctx context.Context, teamName string) ([]repo.User, error)
	CreateTeam(ctx context.Context, tempTeam tempTeamParams) ([]repo.User, error)
	DeactivateUsers(ctx context.Context, userIDs []string) (DeactivateUsersResponse, error)
}

// Handler handles HTTP requests for the teams service.
type Handler struct {
	service Service
}

// NewHandler creates a new teams handler.
func NewHandler(service Service) *Handler {
	return &Handler{
		service: service,
	}
}

type svc struct {
	repo *repo.Queries
	db   *pgxpool.Pool
}

// NewService creates a new teams service.
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

// DeactivateUsersRequest represents the request body for mass user deactivation.
type DeactivateUsersRequest struct {
	Users []string `json:"users"`
}

// DeactivateUsersResponse represents the response for mass user deactivation.
type DeactivateUsersResponse struct {
	UpdatedPRs []domain.PRWithReviewers `json:"updated_prs"`
}

// CreateTeamResponse represents the response for creating a team.
type CreateTeamResponse struct {
	Team TeamResponse `json:"team"`
}

// TeamMember represents a member of a team.
type TeamMember struct {
	UserID   string `json:"user_id"`
	Username string `json:"username"`
	IsActive bool   `json:"is_active"`
}

// TeamResponse represents the response for getting a team.
type TeamResponse struct {
	TeamName string       `json:"team_name"`
	Members  []TeamMember `json:"members"`
}
