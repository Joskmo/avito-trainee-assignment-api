// Package teams provides handlers and service logic for managing teams.
package teams

import (
	"net/http"

	"github.com/Joskmo/avito-trainee-assignment-api/internal/errors"
	"github.com/Joskmo/avito-trainee-assignment-api/internal/json"
)

// GetTeamByName handles the retrieval of a team by its name.
func (h *Handler) GetTeamByName(w http.ResponseWriter, r *http.Request) {
	teamName := r.URL.Query().Get("team_name")

	users, err := h.service.GetTeamByName(r.Context(), teamName)
	if err != nil || len(users) == 0 {
		errors.WriteAppError(w, "team not found", errors.ErrNotFound)
		return
	}

	members := make([]TeamMember, len(users))
	for i, user := range users {
		members[i] = TeamMember{
			UserID:   user.UserID,
			Username: user.Username,
			IsActive: user.IsActive,
		}
	}

	response := TeamResponse{
		TeamName: teamName,
		Members:  members,
	}

	json.Write(w, http.StatusOK, response)
}

// CreateTeam handles the creation of a new team.
func (h *Handler) CreateTeam(w http.ResponseWriter, r *http.Request) {
	var req tempTeamParams
	if err := json.Read(r, &req); err != nil {
		errors.WriteAppError(w, "invalid json in CreateTeam", errors.ErrInvalidInput)
		return
	}

	users, err := h.service.CreateTeam(r.Context(), req)
	if err != nil {
		errors.WriteAppError(w, "error during creating team", err)
		return
	}

	members := make([]TeamMember, len(users))
	for i, user := range users {
		members[i] = TeamMember{
			UserID:   user.UserID,
			Username: user.Username,
			IsActive: user.IsActive,
		}
	}

	response := struct {
		Team TeamResponse `json:"team"`
	}{
		Team: TeamResponse{
			TeamName: req.TeamName,
			Members:  members,
		},
	}
	json.Write(w, http.StatusCreated, response)

}
