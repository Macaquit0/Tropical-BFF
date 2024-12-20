package handlers

import (
	"encoding/json"
	"net/http"

	loginservices "github.com/Macaquit0/Tropical-BFF/internal/services"
)

func (h *Handler) Login(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var req loginservices.LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.log.Error(ctx).Msg("[login-handler] - error on parse login request body. error: %s", err.Error())
		h.router.ErrorHandler(w, r, err)
		return
	}

	token, err := h.services.Login(ctx, req.Email, req.Password)
	if err != nil {
		h.router.ErrorHandler(w, r, err)
		return
	}

	response := map[string]string{
		"token": token,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(response); err != nil {
		h.log.Error(ctx).Msg("[login-handler] - error encoding response. error: %s", err.Error())
	}
}
