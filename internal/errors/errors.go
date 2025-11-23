// Package errors defines application-specific errors and helper functions for handling them.
package errors

import (
	"errors"
	"log/slog"
	"net/http"

	"github.com/Joskmo/avito-trainee-assignment-api/internal/json"
)

// AppError represents an application error with a code, message, and HTTP status.
type AppError struct {
	Code       string `json:"code"`
	Message    string `json:"message"`
	HTTPStatus int    `json:"-"`
}

func (e *AppError) Error() string {
	return e.Message
}

// NewAppError creates a new AppError.
func NewAppError(code, message string, status int) *AppError {
	return &AppError{
		Code:       code,
		Message:    message,
		HTTPStatus: status,
	}
}

// ErrorResponse represents the JSON response for an error.
type ErrorResponse struct {
	Error *AppError `json:"error"`
}

var (
	// ErrTeamExists indicates that the team already exists.
	ErrTeamExists = NewAppError("TEAM_EXISTS", "team_name already exists", http.StatusBadRequest)
	// ErrPRExists indicates that the PR already exists.
	ErrPRExists = NewAppError("PR_EXISTS", "PR id already exists", http.StatusConflict)
	// ErrPRMerged indicates that the PR is already merged.
	ErrPRMerged = NewAppError("PR_MERGED", "cannot reassign on merged PR", http.StatusConflict)
	// ErrNotAssigned indicates that the reviewer is not assigned to the PR.
	ErrNotAssigned = NewAppError("NOT_ASSIGNED", "reviewer not assigned", http.StatusConflict)
	// ErrNoCandidate indicates that no suitable candidate was found for assignment.
	ErrNoCandidate = NewAppError("NO_CANDIDATE", "no suitable candidate found", http.StatusConflict)
	// ErrNotFound indicates that the requested resource was not found.
	ErrNotFound = NewAppError("NOT_FOUND", "resource not found", http.StatusNotFound)

	// ErrInvalidInput indicates that the input data is invalid.
	ErrInvalidInput = NewAppError("INVALID_INPUT", "input data is invalid", http.StatusBadRequest)
	// InternalError indicates an internal server error.
	InternalError = NewAppError("INTERNAL_ERROR", "internal server error", http.StatusInternalServerError)
)

// WriteAppError writes an error response to the ResponseWriter.
func WriteAppError(w http.ResponseWriter, logMsg string, err error) {
	var appErr *AppError
	if errors.As(err, &appErr) {
		slog.Error(logMsg, "error", err)
		json.Write(w, appErr.HTTPStatus, ErrorResponse{
			Error: appErr,
		})
		return
	}

	slog.Error(logMsg, "error", err)
	json.Write(w, http.StatusInternalServerError, ErrorResponse{
		Error: InternalError,
	})
}
