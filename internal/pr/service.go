// Package pr provides handlers and service logic for managing pull requests.
package pr

import (
	"context"
	"errors"
	"math/rand"
	"time"

	apperrors "github.com/Joskmo/avito-trainee-assignment-api/internal/errors"
	repo "github.com/Joskmo/avito-trainee-assignment-api/internal/storage/postgres/sqlc"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
)

func (s *svc) CreatePR(ctx context.Context, createPRParams repo.CreatePRParams) (CreatePRResponse, error) {
	// validate input
	if createPRParams.AuthorID == "" {
		return CreatePRResponse{}, apperrors.ErrInvalidInput
	}
	if createPRParams.PullRequestName == "" {
		return CreatePRResponse{}, apperrors.ErrInvalidInput
	}
	if createPRParams.PullRequestID == "" {
		return CreatePRResponse{}, apperrors.ErrInvalidInput
	}

	// check if PR already exists
	prExists, err := s.repo.PRExists(ctx, createPRParams.PullRequestID)
	if err != nil {
		return CreatePRResponse{}, err
	}
	if prExists {
		return CreatePRResponse{}, apperrors.ErrPRExists
	}

	// get author and validate exists
	author, err := s.repo.GetUser(ctx, createPRParams.AuthorID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return CreatePRResponse{}, apperrors.ErrNotFound
		}
		return CreatePRResponse{}, err
	}

	// get active team members (excluding author)
	teamMembers, err := s.repo.GetActiveTeamMembersExcept(ctx, repo.GetActiveTeamMembersExceptParams{
		TeamName: author.TeamName,
		UserID:   createPRParams.AuthorID,
	})
	if err != nil {
		return CreatePRResponse{}, err
	}

	tx, err := s.db.Begin(ctx)
	if err != nil {
		return CreatePRResponse{}, apperrors.InternalError
	}
	defer func() { _ = tx.Rollback(ctx) }()

	qtx := s.repo.WithTx(tx)

	// create PR
	pr, err := qtx.CreatePR(ctx, createPRParams)
	if err != nil {
		return CreatePRResponse{}, err
	}

	// select up to 2 random reviewers
	reviewers := selectRandomReviewers(teamMembers, 2)

	// assign reviewers
	for _, reviewerID := range reviewers {
		_, err := qtx.AssignReviewer(ctx, repo.AssignReviewerParams{
			PrID:       createPRParams.PullRequestID,
			ReviewerID: reviewerID,
		})
		if err != nil {
			return CreatePRResponse{}, err
		}
	}

	if err := tx.Commit(ctx); err != nil {
		return CreatePRResponse{}, apperrors.InternalError
	}

	// convert status to string
	status := "OPEN"
	if pr.Status.Valid {
		status = string(pr.Status.PrStatusEnum)
	}

	return CreatePRResponse{
		PR: WithReviewers{
			PullRequestID:     pr.PullRequestID,
			PullRequestName:   pr.PullRequestName,
			AuthorID:          pr.AuthorID,
			Status:            status,
			AssignedReviewers: reviewers,
		},
	}, nil
}

// selectRandomReviewers randomly selects up to maxReviewers from the user list
func selectRandomReviewers(users []repo.User, maxReviewers int) []string {
	if len(users) == 0 {
		return []string{}
	}

	// determine actual count
	count := maxReviewers
	if len(users) < count {
		count = len(users)
	}

	// create random generator
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))

	// shuffle users
	shuffled := make([]repo.User, len(users))
	copy(shuffled, users)
	rng.Shuffle(len(shuffled), func(i, j int) {
		shuffled[i], shuffled[j] = shuffled[j], shuffled[i]
	})

	// take first count users
	reviewers := make([]string, count)
	for i := 0; i < count; i++ {
		reviewers[i] = shuffled[i].UserID
	}

	return reviewers
}

func (s *svc) MergePR(ctx context.Context, prID string) (Response, error) {
	// validate input
	if prID == "" {
		return Response{}, apperrors.ErrInvalidInput
	}

	// merge PR (idempotent - already merged PR will just return current state)
	pr, err := s.repo.MergePR(ctx, prID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return Response{}, apperrors.ErrNotFound
		}
		return Response{}, err
	}

	// get current reviewers
	reviewerIDs, err := s.repo.GetPRReviewers(ctx, prID)
	if err != nil {
		return Response{}, err
	}

	// convert status
	status := "OPEN"
	if pr.Status.Valid {
		status = string(pr.Status.PrStatusEnum)
	}

	return Response{
		PR: WithReviewers{
			PullRequestID:     pr.PullRequestID,
			PullRequestName:   pr.PullRequestName,
			AuthorID:          pr.AuthorID,
			Status:            status,
			AssignedReviewers: reviewerIDs,
		},
	}, nil
}

func (s *svc) ReassignReviewer(ctx context.Context, prID, oldUserID string) (ReassignResponse, error) {
	// validate input
	if prID == "" || oldUserID == "" {
		return ReassignResponse{}, apperrors.ErrInvalidInput
	}

	// check PR exists and get it
	pr, err := s.repo.GetPR(ctx, prID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return ReassignResponse{}, apperrors.ErrNotFound
		}
		return ReassignResponse{}, err
	}

	// check PR is not merged
	if pr.Status.Valid && pr.Status.PrStatusEnum == repo.PrStatusEnumMERGED {
		return ReassignResponse{}, apperrors.ErrPRMerged
	}

	// check old reviewer is actually assigned
	isAssigned, err := s.repo.CheckReviewerAssignment(ctx, repo.CheckReviewerAssignmentParams{
		PrID:       prID,
		ReviewerID: oldUserID,
	})
	if err != nil {
		return ReassignResponse{}, err
	}
	if !isAssigned {
		return ReassignResponse{}, apperrors.ErrNotAssigned
	}

	// get old reviewer to find their team
	oldReviewer, err := s.repo.GetUser(ctx, oldUserID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return ReassignResponse{}, apperrors.ErrNotFound
		}
		return ReassignResponse{}, err
	}

	// get active team members except old reviewer and PR author
	teamMembers, err := s.repo.GetActiveTeamMembersExcept(ctx, repo.GetActiveTeamMembersExceptParams{
		TeamName: oldReviewer.TeamName,
		UserID:   oldUserID,
	})
	if err != nil {
		return ReassignResponse{}, err
	}

	// get current reviewers to exclude them from candidates
	currentReviewers, err := s.repo.GetPRReviewers(ctx, prID)
	if err != nil {
		return ReassignResponse{}, err
	}

	currentReviewersMap := make(map[string]bool)
	for _, id := range currentReviewers {
		currentReviewersMap[id] = true
	}

	// filter out PR author and current reviewers from candidates
	var candidates []repo.User
	for _, member := range teamMembers {
		if member.UserID != pr.AuthorID && !currentReviewersMap[member.UserID] {
			candidates = append(candidates, member)
		}
	}

	// check if there are any candidates
	if len(candidates) == 0 {
		return ReassignResponse{}, apperrors.ErrNoCandidate
	}

	// select random new reviewer
	newReviewers := selectRandomReviewers(candidates, 1)
	if len(newReviewers) == 0 {
		return ReassignResponse{}, apperrors.ErrNoCandidate
	}
	newReviewerID := newReviewers[0]

	tx, err := s.db.Begin(ctx)
	if err != nil {
		return ReassignResponse{}, apperrors.InternalError
	}
	defer func() { _ = tx.Rollback(ctx) }()

	qtx := s.repo.WithTx(tx)

	// mark old reviewer as replaced
	_, err = qtx.ReplaceReviewer(ctx, repo.ReplaceReviewerParams{
		PrID:       prID,
		ReviewerID: oldUserID,
		ReplacedBy: pgtype.Text{String: newReviewerID, Valid: true},
	})
	if err != nil {
		return ReassignResponse{}, err
	}

	// assign new reviewer
	_, err = qtx.AssignReviewer(ctx, repo.AssignReviewerParams{
		PrID:       prID,
		ReviewerID: newReviewerID,
	})
	if err != nil {
		return ReassignResponse{}, err
	}

	if err := tx.Commit(ctx); err != nil {
		return ReassignResponse{}, apperrors.InternalError
	}

	// get updated reviewers list
	reviewerIDs, err := s.repo.GetPRReviewers(ctx, prID)
	if err != nil {
		return ReassignResponse{}, err
	}

	// convert status
	status := "OPEN"
	if pr.Status.Valid {
		status = string(pr.Status.PrStatusEnum)
	}

	return ReassignResponse{
		PR: WithReviewers{
			PullRequestID:     pr.PullRequestID,
			PullRequestName:   pr.PullRequestName,
			AuthorID:          pr.AuthorID,
			Status:            status,
			AssignedReviewers: reviewerIDs,
		},
		ReplacedBy: newReviewerID,
	}, nil
}

func (s *svc) GetUserReviews(ctx context.Context, userID string) (UserReviewsResponse, error) {
	// validate input
	if userID == "" {
		return UserReviewsResponse{}, apperrors.ErrInvalidInput
	}

	// get PRs where user is reviewer
	prs, err := s.repo.GetPRsByReviewer(ctx, userID)
	if err != nil {
		return UserReviewsResponse{}, err
	}

	// convert to short format
	prShorts := make([]Short, len(prs))
	for i, pr := range prs {
		status := "OPEN"
		if pr.Status.Valid {
			status = string(pr.Status.PrStatusEnum)
		}

		prShorts[i] = Short{
			PullRequestID:   pr.PullRequestID,
			PullRequestName: pr.PullRequestName,
			AuthorID:        pr.AuthorID,
			Status:          status,
		}
	}

	return UserReviewsResponse{
		UserID:       userID,
		PullRequests: prShorts,
	}, nil
}
