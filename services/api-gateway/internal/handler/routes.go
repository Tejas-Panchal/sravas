package handler

import (
	"github.com/go-chi/chi/v5"
)

// Routes mounts all API Gateway routes on the given chi router
func Routes(r chi.Router) {
	r.Get("/health", HealthHandler)
	r.Get("/ready", ReadyHandler)

	// API v1 — future proxy routes will be mounted here
	r.Route("/api/v1", func(r chi.Router) {
		// r.Mount("/auth", authProxy())
		// r.Mount("/users", userProxy())
		// etc.
	})
}
