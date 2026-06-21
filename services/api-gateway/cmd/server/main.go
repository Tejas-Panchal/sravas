package main

import (
	"log"
	"net/http"
	"os"

	"github.com/go-chi/chi/v5"
	chimw "github.com/go-chi/chi/v5/middleware"

	"github.com/Tejas-Panchal/sravas/services/api-gateway/internal/handler"
	"github.com/Tejas-Panchal/sravas/services/api-gateway/internal/middleware"
)

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8000"
	}

	r := chi.NewRouter()

	// Global middleware
	r.Use(chimw.Logger)
	r.Use(chimw.Recoverer)
	r.Use(middleware.CORS)
	r.Use(middleware.RequestLogger)
	r.Use(middleware.ErrorHandler)

	// Mount routes
	handler.Routes(r)

	log.Printf("API Gateway starting on :%s", port)
	log.Fatal(http.ListenAndServe(":"+port, r))
}
