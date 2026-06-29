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

	// Wire dependencies — use S3 when bucket is configured, local filesystem otherwise
	var store service.Storage
	if bucket := os.Getenv("S3_BUCKET"); bucket != "" {
		var err error
		store, err = service.NewS3Storage(bucket, os.Getenv("AWS_REGION"))
		if err != nil {
			log.Fatalf("failed to create S3 storage: %v", err)
		}
		log.Printf("Using S3 storage: bucket=%s", bucket)
	} else {
		store = service.NewLocalStorage("/tmp/uploads")
		log.Printf("Using local storage: /tmp/uploads")
	}

	bus := service.NewLogEventBus()
	cdnURL := os.Getenv("CDN_URL")
	svc := service.NewUploadService(store, bus, cdnURL)
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
