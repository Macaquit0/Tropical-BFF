package handlers

import (
	"encoding/json"
	"net/http"

	registerservices "github.com/Macaquit0/Tropical-BFF/internal/services"
)

func (h *Handler) CreateUser(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var req registerservices.CreateUserRequestParams
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.log.Error(ctx).Msg("[register-handler] - error on parse register request body. error: %s", err.Error())
		h.router.ErrorHandler(w, r, err)
		return
	}

	registerResponse, err := h.services.CreateUser(ctx, req)
	if err != nil {
		h.router.ErrorHandler(w, r, err)
		return
	}

	response := map[string]string{
		"message": registerResponse.Message,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	if err := json.NewEncoder(w).Encode(response); err != nil {
		h.log.Error(ctx).Msg("[register-handler] - error encoding response. error: %s", err.Error())
	}
}
