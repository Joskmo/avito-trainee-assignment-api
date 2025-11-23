package errors

import (
	"errors"
	"log/slog"
	"net/http"

	"github.com/Joskmo/avito-trainee-assignment-api/internal/json"
)

type AppError struct {
	Code       string `json:"code"`
	Message    string `json:"message"`
	HTTPStatus int    `json:"-"`
}

func (e *AppError) Error() string {
	return e.Message
}

func NewAppError(code, message string, status int) *AppError {
	return &AppError{
		Code:       code,
		Message:    message,
		HTTPStatus: status,
	}
}

type ErrorResponse struct {
	Error *AppError `json:"error"`
}

var (
	ErrTeamExists  = NewAppError("TEAM_EXISTS", "team_name already exists", http.StatusBadRequest)
	ErrPRExists    = NewAppError("PR_EXISTS", "PR id already exists", http.StatusConflict)
	ErrPRMerged    = NewAppError("PR_MERGED", "cannot reassign on merged PR", http.StatusConflict)
	ErrNotAssigned = NewAppError("NOT_ASSIGNED", "reviewer not assigned", http.StatusConflict)
	ErrNoCandidate = NewAppError("NO_CANDIDATE", "no suitable candidate found", http.StatusConflict)
	ErrNotFound    = NewAppError("NOT_FOUND", "resource not found", http.StatusNotFound)

	// not declared in openapi
	ErrInvalidInput = NewAppError("INVALID_INPUT", "input data is invalid", http.StatusBadRequest)
	InternalError   = NewAppError("INTERNAL_ERROR", "internal server error", http.StatusInternalServerError)
)

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
