package main

import (
	"log"
	"net/http"
	"os"

	"github.com/go-chi/chi/v5"
	chimw "github.com/go-chi/chi/v5/middleware"

	"github.com/Tejas-Panchal/sravas/services/search-service/internal/handler"
	"github.com/Tejas-Panchal/sravas/services/search-service/internal/service"
)

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "3005"
	}

	index := service.NewMemoryIndex()
	svc := service.NewSearchService(index)
	h := handler.NewSearchHandler(svc)

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

	r.Route("/api/v1", func(r chi.Router) {
		r.Get("/search", h.Search)
		r.Get("/search/suggestions", h.Suggestions)
		r.Get("/trending", h.Trending)
		r.Get("/categories/{cat}", h.CategoryBrowse)
		r.Post("/index/{id}", h.IndexVideo)
	})

	log.Printf("Search Service starting on :%s", port)
	log.Fatal(http.ListenAndServe(":"+port, r))
}
