package loginhandler

import (
	"encoding/json"
	"net/http"

	loginservices "github.com/bff-cognito/internal/services"
)

type Handler struct {
	services loginservices.UserService
	log      loginservices.Logger
}

func NewHandler(services loginservices.UserService, log loginservices.Logger) *Handler {
	return &Handler{
		services: services,
		log:      log,
	}
}

func (h *Handler) Login(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var req loginservices.LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.log.Error(ctx).Msgf("[login-handler] - error on parse login request body. error: %s", err.Error())
		h.services.ErrorHandler(w, r, err)
		return
	}

	token, err := h.services.Login(ctx, req.Email, req.Password)
	if err != nil {
		h.services.ErrorHandler(w, r, err)
		return
	}

	response := map[string]string{
		"token": token,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(response); err != nil {
		h.log.Error(ctx).Msgf("[login-handler] - error encoding response. error: %s", err.Error())
	}
}
