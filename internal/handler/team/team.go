package team

import (
	"encoding/json"
	"net/http"

	"github.com/SenechkaP/avito-test/internal/handler"
	"github.com/SenechkaP/avito-test/internal/model"
	"github.com/SenechkaP/avito-test/pkg/response"
	"github.com/go-chi/chi/v5"
)

type TeamHandler struct {
	Service TeamServiceInterface
}

func NewTeamHandler(s TeamServiceInterface) *TeamHandler {
	return &TeamHandler{Service: s}
}

func (h *TeamHandler) RegisterRoutes(r chi.Router) {
	r.Post("/team/add", h.createTeam)
	r.Get("/team/get", h.getTeam)
}

func (h *TeamHandler) createTeam(w http.ResponseWriter, r *http.Request) {
	var req model.CreateTeamRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid JSON", http.StatusBadRequest)
		return
	}
	if err := h.Service.CreateOrUpdateTeam(r.Context(), &req); err != nil {
		handler.HandleServiceError(w, err)
		return
	}
	response.ResponseJSON(w, http.StatusCreated, map[string]any{"team": req})
}

func (h *TeamHandler) getTeam(w http.ResponseWriter, r *http.Request) {
	teamName := r.URL.Query().Get("team_name")
	t, err := h.Service.GetTeamByName(r.Context(), teamName)
	if err != nil {
		handler.HandleServiceError(w, err)
		return
	}
	response.ResponseJSON(w, http.StatusOK, map[string]any{"team": t})
}
