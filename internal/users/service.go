package users

import (
	"context"
	"errors"

	apperrors "github.com/Joskmo/avito-trainee-assignment-api/internal/errors"
	repo "github.com/Joskmo/avito-trainee-assignment-api/internal/storage/postgres/sqlc"
	"github.com/jackc/pgx/v5"
)

func (s *svc) SetUserActivity(ctx context.Context, userActivityParams repo.SetUserActivityParams) (repo.User, error) {
	// validation
	if userActivityParams.UserID == "" {
		return repo.User{}, apperrors.ErrInvalidInput
	}

	user, err := s.repo.SetUserActivity(ctx, userActivityParams)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return repo.User{}, apperrors.ErrNotFound
		}
		return repo.User{}, err
	}

	return user, nil
}
