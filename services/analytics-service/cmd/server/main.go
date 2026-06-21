package main

import (
	"log"
	"net/http"
	"os"

	"github.com/go-chi/chi/v5"
	chimw "github.com/go-chi/chi/v5/middleware"

	"github.com/Tejas-Panchal/sravas/services/analytics-service/internal/handler"
	"github.com/Tejas-Panchal/sravas/services/analytics-service/internal/service"
)

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "3007"
	}

	bus := service.NewLogEventBus()
	svc := service.NewAnalyticsService(bus)
	h := handler.NewAnalyticsHandler(svc)

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

	r.Route("/api/v1/analytics", func(r chi.Router) {
		r.Get("/videos/{id}", h.VideoAnalytics)
		r.Get("/channel/{id}", h.ChannelAnalytics)
		r.Post("/events", h.TrackEvent)
		r.Get("/trending", h.Trending)
	})

	log.Printf("Analytics Service starting on :%s", port)
	log.Fatal(http.ListenAndServe(":"+port, r))
}
