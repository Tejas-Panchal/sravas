package handler

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"

	"github.com/Tejas-Panchal/sravas/services/upload-service/internal/middleware"
	"github.com/Tejas-Panchal/sravas/services/upload-service/internal/service"
)

// UploadHandler holds the service dependency and defines HTTP handler methods
type UploadHandler struct {
	svc *service.UploadService
}

// NewUploadHandler creates an UploadHandler with the given service
func NewUploadHandler(svc *service.UploadService) *UploadHandler {
	return &UploadHandler{svc: svc}
}

// UploadVideo handles POST /api/v1/videos/upload — multipart file upload
func (h *UploadHandler) UploadVideo(w http.ResponseWriter, r *http.Request) {
	// Limit request body to 2GB + 1MB metadata buffer
	r.Body = http.MaxBytesReader(w, r.Body, middleware.MaxUploadSize+(1<<20))

	if err := r.ParseMultipartForm(32 << 20); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "file too large or invalid form data"})
		return
	}

	file, fileHeader, err := r.FormFile("video")
	if err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "video file is required"})
		return
	}
	defer file.Close()

	if errMsg := middleware.ValidateFile(fileHeader); errMsg != "" {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": errMsg})
		return
	}

	userID := r.FormValue("user_id")
	title := r.FormValue("title")
	description := r.FormValue("description")

	if userID == "" {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "user_id is required"})
		return
	}

	video, err := h.svc.ProcessUpload(fileHeader, userID, title, description)
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
		return
	}

	writeJSON(w, http.StatusAccepted, video)
}

// GetStatus handles GET /api/v1/videos/{id}/status
func (h *UploadHandler) GetStatus(w http.ResponseWriter, r *http.Request) {
	videoID := chi.URLParam(r, "id")
	video, progress, err := h.svc.GetStatus(videoID)
	if err != nil {
		writeJSON(w, http.StatusNotFound, map[string]string{"error": err.Error()})
		return
	}
	writeJSON(w, http.StatusOK, map[string]interface{}{
		"video":    video,
		"progress": progress,
	})
}

// UpdateMetadata handles PUT /api/v1/videos/{id}
func (h *UploadHandler) UpdateMetadata(w http.ResponseWriter, r *http.Request) {
	videoID := chi.URLParam(r, "id")
	var req struct {
		Title       string `json:"title"`
		Description string `json:"description"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid request body"})
		return
	}
	video, err := h.svc.UpdateMetadata(videoID, req.Title, req.Description)
	if err != nil {
		writeJSON(w, http.StatusNotFound, map[string]string{"error": err.Error()})
		return
	}
	writeJSON(w, http.StatusOK, video)
}

// DeleteVideo handles DELETE /api/v1/videos/{id}
func (h *UploadHandler) DeleteVideo(w http.ResponseWriter, r *http.Request) {
	videoID := chi.URLParam(r, "id")
	if err := h.svc.DeleteVideo(videoID); err != nil {
		writeJSON(w, http.StatusNotFound, map[string]string{"error": err.Error()})
		return
	}
	writeJSON(w, http.StatusOK, map[string]string{"message": "video deleted"})
}

func writeJSON(w http.ResponseWriter, status int, v interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(v)
}
