package users

import (
	"net/http"

	"github.com/Joskmo/avito-trainee-assignment-api/internal/errors"
	"github.com/Joskmo/avito-trainee-assignment-api/internal/json"
	repo "github.com/Joskmo/avito-trainee-assignment-api/internal/storage/postgres/sqlc"
)

func (h *handler) SetUserActivity(w http.ResponseWriter, r *http.Request) {
	var req SetUserActivityRequest
	if err := json.Read(r, &req); err != nil {
		errors.WriteAppError(w, "invalid json in SetUserActivity", errors.ErrInvalidInput)
		return
	}

	if req.IsActive == nil {
		errors.WriteAppError(w, "is_active field is required", errors.ErrInvalidInput)
		return
	}

	params := repo.SetUserActivityParams{
		UserID:   req.UserID,
		IsActive: *req.IsActive,
	}
	user, err := h.service.SetUserActivity(r.Context(), params)
	if err != nil {
		errors.WriteAppError(w, "failed to set user activity", err)
		return
	}
	response := SetUserActivityResponse{User: user}
	json.Write(w, http.StatusOK, response)
}
