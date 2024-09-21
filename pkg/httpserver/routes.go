package httpserver

import "github.com/go-chi/chi"

// RegisterRoutes setups routes for http server
func (s *Server) RegisterRoutes(router *chi.Mux) {
	router.Route("/v1", func(r chi.Router) {
		r.Post("/commands", s.HandleSubmitCommands)
		r.Get("/health", s.HandleHealthCheck)
		r.Get("/logs", s.HandleGetLogs)
	})

	router.NotFound(s.NotFoundHandler)
}
