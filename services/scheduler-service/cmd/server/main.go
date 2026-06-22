package main

import (
	"encoding/json"
	"log"
	"net/http"
	"os"

	"github.com/robfig/cron/v3"

	"github.com/Tejas-Panchal/sravas/services/scheduler-service/internal/job"
)

// main starts the cron scheduler and an HTTP health endpoint
func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "3009"
	}

	// Start cron jobs
	s := cron.New()
	s.AddFunc("*/5 * * * *", job.ProcessUploads)
	s.AddFunc("*/15 * * * *", job.UpdateTrending)
	s.AddFunc("0 * * * *", job.CleanupTokens)
	s.AddFunc("0 2 * * *", job.AggregateAnalytics)
	s.Start()
	log.Println("Scheduler cron jobs started")

	// Health and ready endpoints
	http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]string{"status": "ok"})
	})
	http.HandleFunc("/ready", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]string{"status": "ready"})
	})

	log.Printf("Scheduler HTTP server on :%s", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}
