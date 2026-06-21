package handler

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"

	"github.com/Tejas-Panchal/sravas/services/streaming-service/internal/middleware"
	"github.com/Tejas-Panchal/sravas/services/streaming-service/internal/service"
)

// StreamingHandler holds the service dependency and defines HTTP handler methods
type StreamingHandler struct {
	svc *service.StreamingService
}

// NewStreamingHandler creates a StreamingHandler with the given service
func NewStreamingHandler(svc *service.StreamingService) *StreamingHandler {
	return &StreamingHandler{svc: svc}
}

// GetMetadata handles GET /api/v1/videos/{id}
func (h *StreamingHandler) GetMetadata(w http.ResponseWriter, r *http.Request) {
	videoID := chi.URLParam(r, "id")
	video, err := h.svc.GetMetadata(videoID)
	if err != nil {
		writeJSON(w, http.StatusNotFound, map[string]string{"error": err.Error()})
		return
	}
	writeJSON(w, http.StatusOK, video)
}

// GetManifest handles GET /api/v1/videos/{id}/manifest.m3u8
func (h *StreamingHandler) GetManifest(w http.ResponseWriter, r *http.Request) {
	videoID := chi.URLParam(r, "id")
	manifest, err := h.svc.GetManifest(videoID)
	if err != nil {
		writeJSON(w, http.StatusNotFound, map[string]string{"error": err.Error()})
		return
	}
	w.Header().Set("Content-Type", "application/vnd.apple.mpegurl")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(manifest.Playlist))
}

// GetSegment handles GET /api/v1/videos/{id}/segment/{seg}
// Supports HTTP Range header for video seeking
func (h *StreamingHandler) GetSegment(w http.ResponseWriter, r *http.Request) {
	videoID := chi.URLParam(r, "id")
	segID := chi.URLParam(r, "seg")

	segment, err := h.svc.GetSegment(videoID, segID)
	if err != nil {
		writeJSON(w, http.StatusNotFound, map[string]string{"error": err.Error()})
		return
	}

	data := segment.Data
	contentLen := segment.Size

	// Handle Range header for seeking
	if rng := middleware.ParseRange(r, contentLen); rng != nil {
		if rng.Start >= contentLen {
			w.Header().Set("Content-Range", fmt.Sprintf("bytes */%d", contentLen))
			http.Error(w, "", http.StatusRequestedRangeNotSatisfiable)
			return
		}
		if rng.End >= contentLen {
			rng.End = contentLen - 1
		}
		data = data[rng.Start : rng.End+1]
		w.Header().Set("Content-Range", fmt.Sprintf("bytes %d-%d/%d", rng.Start, rng.End, contentLen))
		w.Header().Set("Content-Length", strconv.FormatInt(int64(len(data)), 10))
		w.WriteHeader(http.StatusPartialContent)
	} else {
		w.Header().Set("Content-Length", strconv.FormatInt(contentLen, 10))
		w.WriteHeader(http.StatusOK)
	}

	w.Header().Set("Content-Type", "video/MP2T")
	w.Header().Set("Accept-Ranges", "bytes")
	w.Header().Set("Cache-Control", "public, max-age=3600")
	w.Write(data)
}

// RecordView handles POST /api/v1/videos/{id}/watch
func (h *StreamingHandler) RecordView(w http.ResponseWriter, r *http.Request) {
	videoID := chi.URLParam(r, "id")
	if err := h.svc.RecordView(videoID); err != nil {
		writeJSON(w, http.StatusNotFound, map[string]string{"error": err.Error()})
		return
	}
	writeJSON(w, http.StatusOK, map[string]string{"message": "view recorded"})
}

// UpdateAnalytics handles PUT /api/v1/videos/{id}/analytics
func (h *StreamingHandler) UpdateAnalytics(w http.ResponseWriter, r *http.Request) {
	videoID := chi.URLParam(r, "id")
	var req struct {
		WatchSeconds float64 `json:"watch_seconds"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid request body"})
		return
	}
	if err := h.svc.UpdateAnalytics(videoID, req.WatchSeconds); err != nil {
		writeJSON(w, http.StatusNotFound, map[string]string{"error": err.Error()})
		return
	}
	writeJSON(w, http.StatusOK, map[string]string{"message": "analytics updated"})
}

func writeJSON(w http.ResponseWriter, status int, v interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(v)
}
