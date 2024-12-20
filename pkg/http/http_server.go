package sharedhttp

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/backend/bff-cognito/pkg/logger"
	"github.com/go-chi/chi/v5"

	"github.com/backend/bff-cognito/pkg/metric"
)

type Server struct {
	Server *http.Server
	Router *Router
}

type ServerOpts struct {
	ServerPort          string
	ServerEnableTracing bool
	ServerEnv           string
	ServerName          string
	ServerVersion       string
}

func NewServer(log *logger.Logger, opts ServerOpts) *Server {
	metric.NewDefault()

	routerOptions := routerOpts{
		ServerEnableTracing: opts.ServerEnableTracing,
		ServerEnv:           opts.ServerEnv,
		ServerName:          opts.ServerName,
		ServerVersion:       opts.ServerVersion,
	}

	router := NewRouter(log, routerOptions)

	server := &http.Server{
		Addr:              fmt.Sprintf(":%s", opts.ServerPort),
		ReadTimeout:       30 * time.Second,
		WriteTimeout:      30 * time.Second,
		IdleTimeout:       30 * time.Second,
		ReadHeaderTimeout: 3 * time.Second,
		Handler:           router.Engine,
	}

	return &Server{
		Server: server,
		Router: router,
	}
}

func (s *Server) Init() error {
	serverCtx, serverStopCtx := context.WithCancel(context.Background())

	sig := make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
	go func() {
		<-sig
		shutdownCtx, _ := context.WithTimeout(serverCtx, 30*time.Second)

		go func() {
			<-shutdownCtx.Done()
			if shutdownCtx.Err() == context.DeadlineExceeded {
				log.Fatal("graceful shutdown timed out.. forcing exit.")
			}
		}()

		if err := s.Server.Shutdown(shutdownCtx); err != nil {
			log.Fatal(err)
		}
		serverStopCtx()
	}()

	fmt.Printf("Starting server on port %s\n", s.Server.Addr)

	chi.Walk(s.Router.Engine, func(method string, route string, handler http.Handler, middlewares ...func(http.Handler) http.Handler) error {
		fmt.Printf("[%s]: '%s' has %d middleware\n", method, route, len(middlewares))
		return nil
	})

	if err := s.Server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatal(err)
	}

	<-serverCtx.Done()

	return nil
}
