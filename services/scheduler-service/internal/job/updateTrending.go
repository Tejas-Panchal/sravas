package job

import (
	"log"
	"math"
	"time"
)

// UpdateTrending recalculates trending scores for all videos
func UpdateTrending() {
	log.Println("[scheduler] update trending: recalculating scores")
	// TODO: query all videos from DB
	// TODO: for each, compute score = viewCount * ln(1 + hoursSinceUpload)
	// TODO: write sorted trending list to Redis for fast lookup
	_ = math.Log(2)
	_ = time.Now()
}
