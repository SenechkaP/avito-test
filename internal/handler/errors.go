package handler

import (
	"errors"
	"net/http"

	"github.com/SenechkaP/avito-test/internal/service"
	"github.com/SenechkaP/avito-test/pkg/response"
)

type errorBody struct {
	Error errorObject `json:"error"`
}

type errorObject struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

func writeErrorJSON(w http.ResponseWriter, status int, code, message string) {
	response.ResponseJSON(w, status, errorBody{
		Error: errorObject{
			Code:    code,
			Message: message,
		},
	})
}

func HandleServiceError(w http.ResponseWriter, err error) {
	if err == nil {
		return
	}

	status := http.StatusInternalServerError
	code := "INTERNAL_ERROR"
	msg := "internal server error"

	switch {
	case errors.Is(err, service.ErrInvalidParams):
		status = http.StatusBadRequest
		code = "INVALID_PARAMS"
		msg = "invalid parameters"
	case errors.Is(err, service.ErrNotFound):
		status = http.StatusNotFound
		code = "NOT_FOUND"
		msg = "not found"
	case errors.Is(err, service.ErrTeamExists):
		status = http.StatusConflict
		code = "TEAM_EXISTS"
		msg = "provided team_name already exists"
	case errors.Is(err, service.ErrPRExists):
		status = http.StatusConflict
		code = "PR_EXISTS"
		msg = "pull request already exists"
	case errors.Is(err, service.ErrPRMerged):
		status = http.StatusConflict
		code = "PR_MERGED"
		msg = "cannot modify merged PR"
	case errors.Is(err, service.ErrNotAssigned):
		status = http.StatusConflict
		code = "NOT_ASSIGNED"
		msg = "reviewer is not assigned to this PR"
	case errors.Is(err, service.ErrNoCandidate):
		status = http.StatusConflict
		code = "NO_CANDIDATE"
		msg = "no active replacement candidate in team"
	default:
	}

	writeErrorJSON(w, status, code, msg)
}
