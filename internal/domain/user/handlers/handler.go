package handlers

import (
	"context"

	"github.com/Macaquit0/Tropical-BFF/internal/domain/user/services"
	sharedhttp "github.com/Macaquit0/Tropical-BFF/pkg/http"
	"github.com/Macaquit0/Tropical-BFF/pkg/logger"
	"github.com/go-chi/chi/v5"
)

type Services interface {
	CreateUser(ctx context.Context, req services.CreateUserRequestParams) (services.CreateUserResponse, error)
	Login(ctx context.Context, req services.LoginRequestParams) (services.LoginResponse, error)
	//GetUserProfile(ctx context.Context, userID string) (services.GetUserProfileResponse, error)
	//UpdateUserProfile(ctx context.Context, userID string, req services.UpdateUserProfileRequest) error
	//DeleteUser(ctx context.Context, userID string) error
}

type Handler struct {
	log        *logger.Logger
	router     *sharedhttp.Router
	services   Services
}

func NewHandler(log *logger.Logger, router *sharedhttp.Router, services Services) *Handler {
	return &Handler{
		log:        log,
		router:     router,
		services:   services,	
	}
}

func (h *Handler) Routes() {
	h.router.Engine.Route("/api/v1/users", func(r chi.Router) {
		r.Use(sharedhttp.InjectHeaders())
		r.Post("/register", h.CreateUser)
		r.Post("/login", h.Login)
		//r.Get("/{user_id}", h.GetUserProfile)
		//r.Put("/{user_id}", h.UpdateUserProfile)
		//r.Delete("/{user_id}", h.DeleteUser)
	})
}
