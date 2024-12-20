package sharedhttp

import (
	"encoding/json"
	"errors"
	"io"
	"log/slog"
	"net/http"
	"sync"
	"time"

	"github.com/Macaquit0/Tropical-BFF/pkg/logger"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"github.com/go-chi/httplog/v2"
	"github.com/go-playground/validator/v10"
	"github.com/rs/zerolog/log"
)

type routerOpts struct {
	ServerEnableTracing bool
	ServerVersion       string
	ServerName          string
	ServerEnv           string
}

type Router struct {
	Engine   *chi.Mux
	Opt      routerOpts
	log      *logger.Logger
	once     sync.Once
	validate *validator.Validate
}

func NewRouter(log *logger.Logger, opts routerOpts) *Router {
	middleware.RequestIDHeader = "x-request-id"

	logger := httplog.NewLogger(opts.ServerName, httplog.Options{
		JSON:             true,
		LogLevel:         slog.LevelDebug,
		Concise:          true,
		RequestHeaders:   true,
		MessageFieldName: "message",
		TimeFieldFormat:  time.RFC850,
		Tags: map[string]string{
			"version": opts.ServerVersion,
			"env":     opts.ServerEnv,
		},
		QuietDownRoutes: []string{
			"/",
			"/status",
			"/ready",
		},
		QuietDownPeriod: 10 * time.Second,
	})
	chiRouter := chi.NewRouter()

	router := Router{
		log:    log,
		Engine: chiRouter,
	}

	chiRouter.Use(middleware.Compress(5))
	chiRouter.Use(middleware.Heartbeat("/status"))
	chiRouter.Use(httplog.RequestLogger(logger))
	chiRouter.Use(middleware.AllowContentType("application/json", "multipart/form-data"))

	chiRouter.Use(middleware.Recoverer)
	chiRouter.Use(middleware.RequestID)
	chiRouter.Use(middleware.RealIP)
	chiRouter.Use(middleware.Logger)
	chiRouter.Use(middleware.Timeout(30 * time.Second))

	chiRouter.Use(cors.Handler(cors.Options{
		AllowedOrigins: []string{"*"},
		AllowedHeaders: []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		AllowedMethods: []string{"GET", "POST", "PUT", "DELETE", "OPTIONS", "PATCH"},
	}))

	// chiRouter.Get("/ready", func(w http.ResponseWriter, r *http.Request) {
	// 	w.Write([]byte("ok"))
	// })

	return &router
}

func (r *Router) Decoder(rr io.Reader, v interface{}) error {
	defer io.Copy(io.Discard, rr)
	err := json.NewDecoder(rr).Decode(v)
	if err != nil {
		return err
	}

	if err := r.validateStruct(v); err != nil {
		if errors.Is(err, validator.ValidationErrors{}) {
			valErrs := err.(validator.ValidationErrors)
			return valErrs
		}

		if errors.Is(err, &validator.InvalidValidationError{}) {
			log.Error().Err(err).Msg("InvalidValidationError")
		}

		return err
	}

	return nil
}

func (v *Router) validateStruct(obj any) error {
	v.lazyinit()
	return v.validate.Struct(obj)
}

func (v *Router) Validator() any {
	v.lazyinit()
	return v.validate
}

func (v *Router) lazyinit() {
	v.once.Do(func() {
		v.validate = validator.New()
		v.validate.SetTagName("binding")
	})
}

func InjectHeaders() func(http.Handler) http.Handler {
	f := func(h http.Handler) http.Handler {
		fn := func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			h.ServeHTTP(w, r)
		}
		return http.HandlerFunc(fn)
	}
	return f
}
