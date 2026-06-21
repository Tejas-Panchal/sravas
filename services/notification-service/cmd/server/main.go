package main

import (
	"log"
	"net/http"
	"os"

	"github.com/go-chi/chi/v5"
	chimw "github.com/go-chi/chi/v5/middleware"

	"github.com/Tejas-Panchal/sravas/services/notification-service/internal/handler"
	"github.com/Tejas-Panchal/sravas/services/notification-service/internal/service"
)

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "3008"
	}

	ws := service.NewLogHub()
	email := service.NewLogSender()
	svc := service.NewNotificationService(ws, email)
	h := handler.NewNotificationHandler(svc)

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

	r.Route("/api/v1/notifications", func(r chi.Router) {
		r.Get("/", h.GetNotifications)
		r.Put("/{id}/read", h.MarkRead)
		r.Put("/read-all", h.MarkAllRead)
		r.Get("/settings", h.GetSettings)
		r.Put("/settings", h.UpdateSettings)
	})

	log.Printf("Notification Service starting on :%s", port)
	log.Fatal(http.ListenAndServe(":"+port, r))
}
