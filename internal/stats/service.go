package stats

import (
	"context"
)

func (s *svc) GetStats(ctx context.Context) (Response, error) {
	// Get top reviewers
	reviewerStats, err := s.repo.GetReviewerStats(ctx)
	if err != nil {
		return Response{}, err
	}

	topReviewers := make([]ReviewerStat, len(reviewerStats))
	for i, stat := range reviewerStats {
		topReviewers[i] = ReviewerStat{
			ReviewerID:      stat.ReviewerID,
			AssignmentCount: stat.AssignmentCount,
		}
	}

	// Get PR status stats
	prStats, err := s.repo.GetPRStatusStats(ctx)
	if err != nil {
		return Response{}, err
	}

	prStatusStats := make([]PRStatusStat, len(prStats))
	for i, stat := range prStats {
		status := "UNKNOWN"
		if stat.Status.Valid {
			status = string(stat.Status.PrStatusEnum)
		}
		prStatusStats[i] = PRStatusStat{
			Status: status,
			Count:  stat.Count,
		}
	}

	// Get total active users
	totalActiveUsers, err := s.repo.GetTotalActiveUsers(ctx)
	if err != nil {
		return Response{}, err
	}

	return Response{
		TopReviewers:     topReviewers,
		PRStatus:         prStatusStats,
		TotalActiveUsers: totalActiveUsers,
	}, nil
}
