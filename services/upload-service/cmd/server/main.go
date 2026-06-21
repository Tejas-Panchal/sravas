package main

import (
	"log"
	"net/http"
	"os"

	"github.com/go-chi/chi/v5"
	chimw "github.com/go-chi/chi/v5/middleware"

	"github.com/Tejas-Panchal/sravas/services/upload-service/internal/handler"
	"github.com/Tejas-Panchal/sravas/services/upload-service/internal/service"
)

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "3003"
	}

	// Wire dependencies
	store := service.NewLocalStorage("/tmp/uploads")
	bus := service.NewLogEventBus()
	svc := service.NewUploadService(store, bus)
	h := handler.NewUploadHandler(svc)

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
		r.Post("/upload", h.UploadVideo)
		r.Get("/{id}/status", h.GetStatus)
		r.Put("/{id}", h.UpdateMetadata)
		r.Delete("/{id}", h.DeleteVideo)
	})

	log.Printf("Upload Service starting on :%s", port)
	log.Fatal(http.ListenAndServe(":"+port, r))
}
