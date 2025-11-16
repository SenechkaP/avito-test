package users

import (
	"encoding/json"
	"net/http"

	"github.com/SenechkaP/avito-test/internal/handler"
	"github.com/SenechkaP/avito-test/internal/model"
	"github.com/SenechkaP/avito-test/pkg/response"
	"github.com/go-chi/chi/v5"
)

type UserHandler struct {
	Service UserServiceInterface
}

func NewUserHandler(s UserServiceInterface) *UserHandler {
	return &UserHandler{Service: s}
}

func (h *UserHandler) RegisterRoutes(r chi.Router) {
	r.Post("/users/setIsActive", h.setIsActive)
	r.Get("/users/getReview", h.getReviews)
}

func (h *UserHandler) setIsActive(w http.ResponseWriter, r *http.Request) {
	var req model.SetUserActiveRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid JSON", http.StatusBadRequest)
		return
	}
	u, err := h.Service.SetUserActive(r.Context(), req.UserID, req.IsActive)
	if err != nil {
		handler.HandleServiceError(w, err)
		return
	}
	response.ResponseJSON(w, http.StatusOK, u)
}

func (h *UserHandler) getReviews(w http.ResponseWriter, r *http.Request) {
	userID := r.URL.Query().Get("user_id")
	if userID == "" {
		http.Error(w, "user_id required", http.StatusBadRequest)
		return
	}
	prs, err := h.Service.GetPRsByReviewer(r.Context(), userID)
	if err != nil {
		handler.HandleServiceError(w, err)
		return
	}
	response.ResponseJSON(w, http.StatusOK, map[string]any{"user_id": userID, "pull_requests": prs})
}
