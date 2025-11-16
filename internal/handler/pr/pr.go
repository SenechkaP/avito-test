package pr

import (
	"encoding/json"
	"net/http"

	"github.com/SenechkaP/avito-test/internal/handler"
	"github.com/SenechkaP/avito-test/internal/model"
	"github.com/SenechkaP/avito-test/pkg/response"
	"github.com/go-chi/chi/v5"
)

type PRHandler struct {
	Service PRServiceInterface
}

func NewPRHandler(s PRServiceInterface) *PRHandler {
	return &PRHandler{Service: s}
}

func (h *PRHandler) RegisterRoutes(r chi.Router) {
	r.Post("/pullRequest/create", h.createPR)
	r.Post("/pullRequest/merge", h.mergePR)
	r.Post("/pullRequest/reassign", h.reassignReviewer)
	r.Get("/pullRequest/stats", h.getStats)
}

func (h *PRHandler) createPR(w http.ResponseWriter, r *http.Request) {
	var req model.CreatePRRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid JSON", http.StatusBadRequest)
		return
	}
	pr, err := h.Service.CreatePullRequest(r.Context(), &req)
	if err != nil {
		handler.HandleServiceError(w, err)
		return
	}
	response.ResponseJSON(w, http.StatusCreated, map[string]any{"pr": pr})
}

func (h *PRHandler) mergePR(w http.ResponseWriter, r *http.Request) {
	var req model.MergePRRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid JSON", http.StatusBadRequest)
		return
	}
	pr, err := h.Service.MergePullRequest(r.Context(), req.PullRequestID)
	if err != nil {
		handler.HandleServiceError(w, err)
		return
	}
	response.ResponseJSON(w, http.StatusOK, map[string]any{"pr": pr})
}

func (h *PRHandler) reassignReviewer(w http.ResponseWriter, r *http.Request) {
	var req model.ReassignRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid JSON", http.StatusBadRequest)
		return
	}
	newUID, pr, err := h.Service.ReassignReviewer(r.Context(), req.PullRequestID, req.OldUserID)
	if err != nil {
		handler.HandleServiceError(w, err)
		return
	}
	response.ResponseJSON(w, http.StatusOK, map[string]any{"pr": pr, "replaced_by": newUID})
}

func (h *PRHandler) getStats(w http.ResponseWriter, r *http.Request) {
	groupBy := r.URL.Query().Get("group_by")
	if groupBy == "" {
		groupBy = "user"
	}

	switch groupBy {
	case "user":
		stats, err := h.Service.GetAssignmentsCountByUser(r.Context())
		if err != nil {
			handler.HandleServiceError(w, err)
			return
		}
		response.ResponseJSON(w, http.StatusOK, map[string]any{"group_by": "user", "stats": stats})
		return
	case "pr":
		stats, err := h.Service.GetAssignmentsCountByPR(r.Context())
		if err != nil {
			handler.HandleServiceError(w, err)
			return
		}
		response.ResponseJSON(w, http.StatusOK, map[string]any{"group_by": "pr", "stats": stats})
		return
	default:
		http.Error(w, "invalid group_by, allowed: user|pr", http.StatusBadRequest)
		return
	}
}
