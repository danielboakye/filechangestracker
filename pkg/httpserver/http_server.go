package httpserver

import (
	"context"
	"log/slog"
	"net/http"

	"github.com/danielboakye/filechangestracker/internal/commandexecutor"
	"github.com/danielboakye/filechangestracker/pkg/filechangestracker"
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
)

// Server represents an HTTP server
type Server struct {
	httpServer *http.Server
	logger     *slog.Logger
	tracker    *filechangestracker.FileChangesTracker
	executor   commandexecutor.CommandExecutor
}

// NewServer creates and returns a new Server instance
func NewServer(
	addr string,
	logger *slog.Logger,
	tracker *filechangestracker.FileChangesTracker,
	executor commandexecutor.CommandExecutor,
) *Server {
	router := chi.NewRouter()
	router.Use(middleware.Logger)

	s := &Server{
		logger:   logger,
		tracker:  tracker,
		executor: executor,
	}

	s.RegisterRoutes(router)

	s.httpServer = &http.Server{
		Addr:    addr,
		Handler: router,
	}

	return s
}

// Start starts the HTTP server
func (s *Server) Start() error {
	s.logger.Info("starting-http-server", slog.String("url", "http://"+s.httpServer.Addr))

	go func() {
		if err := s.httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			s.logger.Error("error-starting-http-server", slog.String("error", err.Error()))
		}
	}()

	return nil
}

func (s *Server) Stop(ctx context.Context) error {
	return s.httpServer.Shutdown(ctx)
}
