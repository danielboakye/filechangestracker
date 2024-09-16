package httpserver

import (
	"log/slog"
	"net/http"

	"github.com/danielboakye/filechangestracker/pkg/filechangestracker"
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
)

// Server represents an HTTP server
type Server struct {
	addr    string
	router  *chi.Mux
	logger  *slog.Logger
	tracker *filechangestracker.FileChangesTracker
}

// NewServer creates and returns a new Server instance
func NewServer(addr string, logger *slog.Logger, tracker *filechangestracker.FileChangesTracker) *Server {
	router := chi.NewRouter()
	router.Use(middleware.Logger)

	s := &Server{
		addr:    addr,
		router:  router,
		logger:  logger,
		tracker: tracker,
	}

	s.RegisterRoutes()

	return s
}

// Start starts the HTTP server
func (s *Server) Start() error {
	s.logger.Info("starting-http-server", slog.String("url", "http://"+s.addr))
	return http.ListenAndServe(s.addr, s.router)
}
