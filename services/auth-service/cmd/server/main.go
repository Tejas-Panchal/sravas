package main

import (
	"log"
	"net/http"
	"os"

	"github.com/go-chi/chi/v5"
	chimw "github.com/go-chi/chi/v5/middleware"

	"github.com/Tejas-Panchal/sravas/services/auth-service/internal/handler"
	"github.com/Tejas-Panchal/sravas/services/auth-service/internal/middleware"
)

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "3001"
	}

	r := chi.NewRouter()

	r.Use(chimw.Logger)
	r.Use(chimw.Recoverer)
	r.Use(middleware.Validate)

	r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{"status":"ok"}`))
	})
	r.Get("/ready", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{"status":"ready"}`))
	})

	r.Route("/api/v1/auth", func(r chi.Router) {
		r.Post("/register", handler.Register)
		r.Post("/login", handler.Login)
		r.Post("/refresh", handler.Refresh)
		r.Post("/password-reset", handler.PasswordReset)
		r.Post("/password-reset/confirm", handler.PasswordResetConfirm)
	})

	log.Printf("Auth Service starting on :%s", port)
	log.Fatal(http.ListenAndServe(":"+port, r))
}
