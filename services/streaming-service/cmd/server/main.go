package main

import (
	"log"
	"net/http"
	"os"

	"github.com/go-chi/chi/v5"
	chimw "github.com/go-chi/chi/v5/middleware"

	"github.com/Tejas-Panchal/sravas/services/streaming-service/internal/handler"
	"github.com/Tejas-Panchal/sravas/services/streaming-service/internal/service"
)

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "3004"
	}

	cache := service.NewMemoryCache()
	svc := service.NewStreamingService(cache)
	h := handler.NewStreamingHandler(svc)

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

	r.Route("/api/v1/videos", func(r chi.Router) {
		r.Get("/{id}", h.GetMetadata)
		r.Get("/{id}/manifest.m3u8", h.GetManifest)
		r.Get("/{id}/segment/{seg}", h.GetSegment)
		r.Post("/{id}/watch", h.RecordView)
		r.Put("/{id}/analytics", h.UpdateAnalytics)
	})

	log.Printf("Streaming Service starting on :%s", port)
	log.Fatal(http.ListenAndServe(":"+port, r))
}
