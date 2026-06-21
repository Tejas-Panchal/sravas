package handler

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"

	"github.com/Tejas-Panchal/sravas/services/search-service/internal/middleware"
	"github.com/Tejas-Panchal/sravas/services/search-service/internal/model"
	"github.com/Tejas-Panchal/sravas/services/search-service/internal/service"
)

// SearchHandler holds the service dependency and defines HTTP handler methods
type SearchHandler struct {
	svc *service.SearchService
}

// NewSearchHandler creates a SearchHandler with the given service
func NewSearchHandler(svc *service.SearchService) *SearchHandler {
	return &SearchHandler{svc: svc}
}

// Search handles GET /api/v1/search
func (h *SearchHandler) Search(w http.ResponseWriter, r *http.Request) {
	req := middleware.ParseSearchRequest(r)
	if req.Query == "" {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "query parameter q is required"})
		return
	}
	resp, err := h.svc.Search(req)
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
		return
	}
	writeJSON(w, http.StatusOK, resp)
}

// Suggestions handles GET /api/v1/search/suggestions
func (h *SearchHandler) Suggestions(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query().Get("q")
	if q == "" {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "query parameter q is required"})
		return
	}
	suggestions, err := h.svc.Suggest(q)
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
		return
	}
	writeJSON(w, http.StatusOK, suggestions)
}

// Trending handles GET /api/v1/trending
func (h *SearchHandler) Trending(w http.ResponseWriter, r *http.Request) {
	trending, err := h.svc.Trending()
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
		return
	}
	writeJSON(w, http.StatusOK, trending)
}

// CategoryBrowse handles GET /api/v1/categories/{cat}
func (h *SearchHandler) CategoryBrowse(w http.ResponseWriter, r *http.Request) {
	category := chi.URLParam(r, "cat")
	page, _ := strconv.Atoi(r.URL.Query().Get("page"))
	pageSize, _ := strconv.Atoi(r.URL.Query().Get("page_size"))

	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 50 {
		pageSize = 20
	}

	resp, err := h.svc.CategoryBrowse(category, page, pageSize)
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
		return
	}
	writeJSON(w, http.StatusOK, resp)
}

// IndexVideo handles POST /api/v1/index/{id}
func (h *SearchHandler) IndexVideo(w http.ResponseWriter, r *http.Request) {
	videoID := chi.URLParam(r, "id")
	var video model.SearchVideo
	if err := json.NewDecoder(r.Body).Decode(&video); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid request body"})
		return
	}
	video.ID = videoID
	if err := h.svc.IndexVideo(&video); err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
		return
	}
	writeJSON(w, http.StatusOK, map[string]string{"message": "video indexed"})
}

func writeJSON(w http.ResponseWriter, status int, v interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(v)
}
