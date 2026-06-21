package main

import (
	"log"
	"net/http"
	"os"

	"github.com/go-chi/chi/v5"
	chimw "github.com/go-chi/chi/v5/middleware"

	"github.com/Tejas-Panchal/sravas/services/comment-service/internal/handler"
	"github.com/Tejas-Panchal/sravas/services/comment-service/internal/service"
)

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "3006"
	}

	bus := service.NewLogEventBus()
	svc := service.NewCommentService(bus)
	h := handler.NewCommentHandler(svc)

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

	r.Route("/api/v1/videos/{id}/comments", func(r chi.Router) {
		r.Post("/", h.AddComment)
		r.Get("/", h.GetComments)
	})

	r.Route("/api/v1/comments/{id}", func(r chi.Router) {
		r.Put("/", h.EditComment)
		r.Delete("/", h.DeleteComment)
		r.Post("/replies", h.ReplyToComment)
		r.Post("/like", h.LikeComment)
	})

	log.Printf("Comment Service starting on :%s", port)
	log.Fatal(http.ListenAndServe(":"+port, r))
}
