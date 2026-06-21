package main

import (
	"log"
	"net/http"
	"os"

	"github.com/go-chi/chi/v5"
	chimw "github.com/go-chi/chi/v5/middleware"

	"github.com/Tejas-Panchal/sravas/services/user-service/internal/handler"
	"github.com/Tejas-Panchal/sravas/services/user-service/internal/middleware"
)

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "3002"
	}

	r := chi.NewRouter()

	r.Use(chimw.Logger)
	r.Use(chimw.Recoverer)

	r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{"status":"ok"}`))
	})
	r.Get("/ready", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{"status":"ready"}`))
	})

	r.Route("/api/v1/users", func(r chi.Router) {
		r.Post("/register", handler.Register)
		r.Post("/login", handler.Login)

		r.Group(func(r chi.Router) {
			r.Use(middleware.Auth)
			r.Get("/{id}", handler.GetProfile)
			r.Put("/{id}", handler.UpdateProfile)
			r.Delete("/{id}", handler.DeleteAccount)
		})

		r.Group(func(r chi.Router) {
			r.Use(middleware.Auth)
			r.Get("/{id}/subscriptions", handler.GetSubscriptions)
			r.Post("/{id}/subscribe", handler.Subscribe)
			r.Get("/{id}/videos", handler.GetUserVideos)
			r.Get("/{id}/stats", handler.GetChannelStats)
		})
	})

	log.Printf("User Service starting on :%s", port)
	log.Fatal(http.ListenAndServe(":"+port, r))
}
