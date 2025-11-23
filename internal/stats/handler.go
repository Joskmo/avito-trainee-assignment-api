// Package stats provides handlers and service logic for statistics.
package stats

import (
	"net/http"

	"github.com/Joskmo/avito-trainee-assignment-api/internal/errors"
	"github.com/Joskmo/avito-trainee-assignment-api/internal/json"
)

// GetStats handles the retrieval of system statistics.
func (h *Handler) GetStats(w http.ResponseWriter, r *http.Request) {
	stats, err := h.service.GetStats(r.Context())
	if err != nil {
		errors.WriteAppError(w, "failed to get stats", err)
		return
	}

	json.Write(w, http.StatusOK, stats)
}
