package handler

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"

	"github.com/Tejas-Panchal/sravas/services/analytics-service/internal/middleware"
	"github.com/Tejas-Panchal/sravas/services/analytics-service/internal/model"
	"github.com/Tejas-Panchal/sravas/services/analytics-service/internal/service"
)

// AnalyticsHandler holds the service dependency and defines HTTP handler methods
type AnalyticsHandler struct {
	svc *service.AnalyticsService
}

// NewAnalyticsHandler creates an AnalyticsHandler with the given service
func NewAnalyticsHandler(svc *service.AnalyticsService) *AnalyticsHandler {
	return &AnalyticsHandler{svc: svc}
}

// VideoAnalytics handles GET /api/v1/analytics/videos/{id}
func (h *AnalyticsHandler) VideoAnalytics(w http.ResponseWriter, r *http.Request) {
	videoID := chi.URLParam(r, "id")
	analytics, err := h.svc.GetVideoAnalytics(videoID)
	if err != nil {
		writeJSON(w, http.StatusNotFound, map[string]string{"error": err.Error()})
		return
	}
	writeJSON(w, http.StatusOK, analytics)
}

// ChannelAnalytics handles GET /api/v1/analytics/channel/{id}
func (h *AnalyticsHandler) ChannelAnalytics(w http.ResponseWriter, r *http.Request) {
	userID := chi.URLParam(r, "id")
	analytics, err := h.svc.GetChannelAnalytics(userID)
	if err != nil {
		writeJSON(w, http.StatusNotFound, map[string]string{"error": err.Error()})
		return
	}
	writeJSON(w, http.StatusOK, analytics)
}

// TrackEvent handles POST /api/v1/analytics/events
func (h *AnalyticsHandler) TrackEvent(w http.ResponseWriter, r *http.Request) {
	var req model.TrackEventRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid request body"})
		return
	}

	if errMsg := middleware.ValidateEventType(req.Type); errMsg != "" {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": errMsg})
		return
	}
	if req.VideoID == "" {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "video_id is required"})
		return
	}
	if req.UserID == "" {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "user_id is required"})
		return
	}
	if req.Type == "traffic_source" && req.Source == "" {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "source is required for traffic_source events"})
		return
	}

	if err := h.svc.TrackEvent(req); err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
		return
	}
	writeJSON(w, http.StatusAccepted, map[string]string{"message": "event tracked"})
}

// Trending handles GET /api/v1/analytics/trending
func (h *AnalyticsHandler) Trending(w http.ResponseWriter, r *http.Request) {
	trending, err := h.svc.GetTrending()
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
		return
	}
	writeJSON(w, http.StatusOK, trending)
}

func writeJSON(w http.ResponseWriter, status int, v interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(v)
}
