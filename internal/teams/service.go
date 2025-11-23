package teams

import (
	"context"
	"math/rand"

	"github.com/Joskmo/avito-trainee-assignment-api/internal/domain"
	"github.com/Joskmo/avito-trainee-assignment-api/internal/errors"
	repo "github.com/Joskmo/avito-trainee-assignment-api/internal/storage/postgres/sqlc"
	"github.com/jackc/pgx/v5/pgtype"
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

func (s *svc) DeactivateUsers(ctx context.Context, userIDs []string) (DeactivateUsersResponse, error) {
	tx, err := s.db.Begin(ctx)
	if err != nil {
		return DeactivateUsersResponse{}, errors.InternalError
	}
	defer func() { _ = tx.Rollback(ctx) }()

	qtx := s.repo.WithTx(tx)

	deactivatedMap := make(map[string]bool)
	for _, uid := range userIDs {
		deactivatedMap[uid] = true
	}

	updatedPRsMap := make(map[string]struct{})

	for _, uid := range userIDs {
		// Get user to find team name
		user, err := qtx.GetUser(ctx, uid)
		if err != nil {
			return DeactivateUsersResponse{}, err
		}

		// 1. Deactivate user
		_, err = qtx.SetUserActivity(ctx, repo.SetUserActivityParams{
			UserID:   uid,
			IsActive: false,
		})
		if err != nil {
			return DeactivateUsersResponse{}, err
		}

		// 2. Get active assignments
		prs, err := qtx.GetPRsByReviewer(ctx, uid)
		if err != nil {
			return DeactivateUsersResponse{}, err
		}

		for _, pr := range prs {
			updatedPRsMap[pr.PullRequestID] = struct{}{}

			// 3. Find replacement
			candidates, err := qtx.GetActiveTeamMembersExcept(ctx, repo.GetActiveTeamMembersExceptParams{
				TeamName: user.TeamName,
				UserID:   pr.AuthorID,
			})
			if err != nil {
				return DeactivateUsersResponse{}, err
			}

			currentReviewers, err := qtx.GetPRReviewers(ctx, pr.PullRequestID)
			if err != nil {
				return DeactivateUsersResponse{}, err
			}
			currentReviewerMap := make(map[string]bool)
			for _, r := range currentReviewers {
				currentReviewerMap[r] = true
			}

			var validCandidates []string
			for _, c := range candidates {
				if deactivatedMap[c.UserID] {
					continue
				}
				if currentReviewerMap[c.UserID] {
					continue
				}
				if c.UserID == uid {
					continue
				}
				validCandidates = append(validCandidates, c.UserID)
			}

			if len(validCandidates) > 0 {
				idx := rand.Intn(len(validCandidates))
				newReviewerID := validCandidates[idx]

				_, err = qtx.ReplaceReviewer(ctx, repo.ReplaceReviewerParams{
					PrID:       pr.PullRequestID,
					ReviewerID: uid,
					ReplacedBy: pgtype.Text{String: newReviewerID, Valid: true},
				})
				if err != nil {
					return DeactivateUsersResponse{}, err
				}

				_, err = qtx.AssignReviewer(ctx, repo.AssignReviewerParams{
					PrID:       pr.PullRequestID,
					ReviewerID: newReviewerID,
				})
				if err != nil {
					return DeactivateUsersResponse{}, err
				}
			} else {
				err = qtx.DeleteReviewer(ctx, repo.DeleteReviewerParams{
					PrID:       pr.PullRequestID,
					ReviewerID: uid,
				})
				if err != nil {
					return DeactivateUsersResponse{}, err
				}
			}
		}
	}

	var response DeactivateUsersResponse
	for prID := range updatedPRsMap {
		pr, err := qtx.GetPR(ctx, prID)
		if err != nil {
			return DeactivateUsersResponse{}, err
		}

		reviewers, err := qtx.GetPRReviewers(ctx, prID)
		if err != nil {
			return DeactivateUsersResponse{}, err
		}

		status := "UNKNOWN"
		if pr.Status.Valid {
			status = string(pr.Status.PrStatusEnum)
		}

		response.UpdatedPRs = append(response.UpdatedPRs, domain.PRWithReviewers{
			PullRequestID:     pr.PullRequestID,
			PullRequestName:   pr.PullRequestName,
			AuthorID:          pr.AuthorID,
			Status:            status,
			AssignedReviewers: reviewers,
		})
	}

	if err := tx.Commit(ctx); err != nil {
		return DeactivateUsersResponse{}, errors.InternalError
	}

	return response, nil
}
