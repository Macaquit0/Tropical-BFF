package sharedhttp

import (
	"net/http"
	"time"

	"github.com/go-chi/chi/v5/middleware"
	"github.com/rs/zerolog"
)

func LoggerMiddleware(logger zerolog.Logger) func(http.Handler) http.Handler {
	return func(h http.Handler) http.Handler {
		middlewareFunction := func(rw http.ResponseWriter, r *http.Request) {
			wrapper := middleware.NewWrapResponseWriter(rw, r.ProtoMajor)
			start := time.Now()

			defer func() {
				logger.Info().
					Str("x-request-id", middleware.GetReqID(r.Context())).
					Str("method", r.Method).
					Str("path", r.URL.Path).
					Str("query", r.URL.RawQuery).
					Str("user-agent", r.UserAgent()).
					Str("x-account-id", r.Header.Get("x-account-id")).
					Str("x-customer-id", r.Header.Get("x-customer-id")).
					Int("status", wrapper.Status()).
					Dur("latency", time.Since(start)).
					Msg("request finished")
			}()

			h.ServeHTTP(wrapper, r)
		}

		return http.HandlerFunc(middlewareFunction)
	}
}
