package httpserver

import (
	"github.com/danielboakye/filechangestracker/internal/commandexecutor"
	"github.com/danielboakye/filechangestracker/internal/filechangestracker"
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
)

type Handler struct {
	tracker  filechangestracker.FileChangesTracker
	executor commandexecutor.CommandExecutor
}

func NewHandler(
	tracker filechangestracker.FileChangesTracker,
	executor commandexecutor.CommandExecutor,
) *Handler {
	return &Handler{
		tracker:  tracker,
		executor: executor,
	}
}

// RegisterRoutes setups routes for http server
func (h *Handler) RegisterRoutes() *chi.Mux {
	router := chi.NewRouter()
	router.Use(middleware.Logger)

	router.Route("/v1", func(r chi.Router) {
		r.Post("/commands", h.HandleSubmitCommands)
		r.Get("/health", h.HandleHealthCheck)
		r.Get("/logs", h.HandleGetLogs)
	})

	router.NotFound(h.NotFoundHandler)

	return router
}
