// Package pr provides handlers and service logic for managing pull requests.
package pr

import (
	"net/http"

	"github.com/Joskmo/avito-trainee-assignment-api/internal/errors"
	"github.com/Joskmo/avito-trainee-assignment-api/internal/json"
	repo "github.com/Joskmo/avito-trainee-assignment-api/internal/storage/postgres/sqlc"
)

// CreatePR handles the creation of a new pull request.
func (h *Handler) CreatePR(w http.ResponseWriter, r *http.Request) {
	var req repo.CreatePRParams
	if err := json.Read(r, &req); err != nil {
		errors.WriteAppError(w, "invalid json in CreatePR", errors.ErrInvalidInput)
		return
	}

	response, err := h.service.CreatePR(r.Context(), req)
	if err != nil {
		errors.WriteAppError(w, "failed to create PR", err)
		return
	}

	json.Write(w, http.StatusCreated, response)
}

// MergePR handles the merging of a pull request.
func (h *Handler) MergePR(w http.ResponseWriter, r *http.Request) {
	var req struct {
		PullRequestID string `json:"pull_request_id"`
	}
	if err := json.Read(r, &req); err != nil {
		errors.WriteAppError(w, "invalid json in MergePR", errors.ErrInvalidInput)
		return
	}

	response, err := h.service.MergePR(r.Context(), req.PullRequestID)
	if err != nil {
		errors.WriteAppError(w, "failed to merge PR", err)
		return
	}

	json.Write(w, http.StatusOK, response)
}

// ReassignReviewer handles the reassignment of a reviewer for a pull request.
func (h *Handler) ReassignReviewer(w http.ResponseWriter, r *http.Request) {
	var req struct {
		PullRequestID string `json:"pull_request_id"`
		OldUserID     string `json:"old_user_id"`
	}
	if err := json.Read(r, &req); err != nil {
		errors.WriteAppError(w, "invalid json in ReassignReviewer", errors.ErrInvalidInput)
		return
	}

	response, err := h.service.ReassignReviewer(r.Context(), req.PullRequestID, req.OldUserID)
	if err != nil {
		errors.WriteAppError(w, "failed to reassign reviewer", err)
		return
	}

	json.Write(w, http.StatusOK, response)
}

// GetUserReviews handles the retrieval of pull requests assigned to a user for review.
func (h *Handler) GetUserReviews(w http.ResponseWriter, r *http.Request) {
	userID := r.URL.Query().Get("user_id")

	response, err := h.service.GetUserReviews(r.Context(), userID)
	if err != nil {
		errors.WriteAppError(w, "failed to get user reviews", err)
		return
	}

	json.Write(w, http.StatusOK, response)
}
