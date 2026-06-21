package handler

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"

	"github.com/Tejas-Panchal/sravas/services/comment-service/internal/middleware"
	"github.com/Tejas-Panchal/sravas/services/comment-service/internal/model"
	"github.com/Tejas-Panchal/sravas/services/comment-service/internal/service"
)

// CommentHandler holds the service dependency and defines HTTP handler methods
type CommentHandler struct {
	svc *service.CommentService
}

// NewCommentHandler creates a CommentHandler with the given service
func NewCommentHandler(svc *service.CommentService) *CommentHandler {
	return &CommentHandler{svc: svc}
}

// AddComment handles POST /api/v1/videos/{id}/comments
func (h *CommentHandler) AddComment(w http.ResponseWriter, r *http.Request) {
	videoID := chi.URLParam(r, "id")
	var req model.AddCommentRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid request body"})
		return
	}
	if req.UserID == "" {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "user_id is required"})
		return
	}
	if errMsg := middleware.ValidateContent(req.Content); errMsg != "" {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": errMsg})
		return
	}
	comment, err := h.svc.AddComment(videoID, req)
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
		return
	}
	writeJSON(w, http.StatusCreated, comment)
}

// GetComments handles GET /api/v1/videos/{id}/comments
func (h *CommentHandler) GetComments(w http.ResponseWriter, r *http.Request) {
	videoID := chi.URLParam(r, "id")
	cursor := r.URL.Query().Get("cursor")
	limit, _ := strconv.Atoi(r.URL.Query().Get("limit"))
	if limit < 1 || limit > 100 {
		limit = 20
	}

	resp, err := h.svc.GetComments(videoID, cursor, limit)
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
		return
	}
	writeJSON(w, http.StatusOK, resp)
}

// EditComment handles PUT /api/v1/comments/{id}
func (h *CommentHandler) EditComment(w http.ResponseWriter, r *http.Request) {
	commentID := chi.URLParam(r, "id")
	var req model.EditCommentRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid request body"})
		return
	}
	if errMsg := middleware.ValidateContent(req.Content); errMsg != "" {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": errMsg})
		return
	}
	comment, err := h.svc.EditComment(commentID, req.Content)
	if err != nil {
		writeJSON(w, http.StatusNotFound, map[string]string{"error": err.Error()})
		return
	}
	writeJSON(w, http.StatusOK, comment)
}

// DeleteComment handles DELETE /api/v1/comments/{id}
func (h *CommentHandler) DeleteComment(w http.ResponseWriter, r *http.Request) {
	commentID := chi.URLParam(r, "id")
	if err := h.svc.DeleteComment(commentID); err != nil {
		writeJSON(w, http.StatusNotFound, map[string]string{"error": err.Error()})
		return
	}
	writeJSON(w, http.StatusOK, map[string]string{"message": "comment deleted"})
}

// ReplyToComment handles POST /api/v1/comments/{id}/replies
func (h *CommentHandler) ReplyToComment(w http.ResponseWriter, r *http.Request) {
	parentID := chi.URLParam(r, "id")
	var req model.ReplyRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid request body"})
		return
	}
	if req.UserID == "" {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "user_id is required"})
		return
	}
	if errMsg := middleware.ValidateContent(req.Content); errMsg != "" {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": errMsg})
		return
	}
	reply, err := h.svc.ReplyToComment(parentID, req)
	if err != nil {
		writeJSON(w, http.StatusNotFound, map[string]string{"error": err.Error()})
		return
	}
	writeJSON(w, http.StatusCreated, reply)
}

// LikeComment handles POST /api/v1/comments/{id}/like
func (h *CommentHandler) LikeComment(w http.ResponseWriter, r *http.Request) {
	commentID := chi.URLParam(r, "id")
	var req struct {
		UserID string `json:"user_id"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid request body"})
		return
	}
	if req.UserID == "" {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "user_id is required"})
		return
	}
	liked, err := h.svc.LikeComment(commentID, req.UserID)
	if err != nil {
		writeJSON(w, http.StatusNotFound, map[string]string{"error": err.Error()})
		return
	}
	writeJSON(w, http.StatusOK, map[string]bool{"liked": liked})
}

func writeJSON(w http.ResponseWriter, status int, v interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(v)
}
