package teams

import (
	"context"

	"github.com/Joskmo/avito-trainee-assignment-api/internal/errors"
	repo "github.com/Joskmo/avito-trainee-assignment-api/internal/storage/postgres/sqlc"
)

func (s *svc) GetTeamByName(ctx context.Context, teamName string) ([]repo.User, error) {
	return s.repo.GetTeam(ctx, teamName)
}

func (s *svc) CreateTeam(ctx context.Context, tempTeam tempTeamParams) ([]repo.User, error) {
	// validation
	if tempTeam.TeamName == "" {
		return nil, errors.ErrInvalidInput
	}
	if len(tempTeam.Members) == 0 {
		return nil, errors.ErrInvalidInput
	}
	for _, u := range tempTeam.Members {
		if u.Username == "" {
			return nil, errors.ErrInvalidInput
		}
	}

	tx, err := s.db.Begin(ctx)
	if err != nil {
		return nil, errors.InternalError
	}
	defer func() { _ = tx.Rollback(ctx) }()

	qtx := s.repo.WithTx(tx)

	exists, err := qtx.TeamExists(ctx, tempTeam.TeamName)
	if err != nil {
		return nil, errors.InternalError
	}
	if exists {
		return nil, errors.ErrTeamExists
	}

	var createdUsers []repo.User

	for _, tempUser := range tempTeam.Members {
		user, err := qtx.CreateUser(ctx, repo.CreateUserParams{
			UserID:   tempUser.UserID,
			Username: tempUser.Username,
			IsActive: tempUser.IsActive,
			TeamName: tempTeam.TeamName,
		})
		if err != nil {
			return nil, errors.ErrTeamExists
		}
		createdUsers = append(createdUsers, user)
	}

	if err := tx.Commit(ctx); err != nil {
		return nil, errors.InternalError
	}

	return createdUsers, nil
}
