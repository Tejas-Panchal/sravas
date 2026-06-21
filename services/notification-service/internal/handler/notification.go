package handler

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"

	"github.com/Tejas-Panchal/sravas/services/notification-service/internal/model"
	"github.com/Tejas-Panchal/sravas/services/notification-service/internal/service"
)

// NotificationHandler holds the service dependency and defines HTTP handler methods
type NotificationHandler struct {
	svc *service.NotificationService
}

// NewNotificationHandler creates a NotificationHandler with the given service
func NewNotificationHandler(svc *service.NotificationService) *NotificationHandler {
	return &NotificationHandler{svc: svc}
}

// GetNotifications handles GET /api/v1/notifications
func (h *NotificationHandler) GetNotifications(w http.ResponseWriter, r *http.Request) {
	userID := r.URL.Query().Get("user_id")
	if userID == "" {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "user_id query parameter is required"})
		return
	}

	notifs, err := h.svc.GetNotifications(userID)
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
		return
	}
	writeJSON(w, http.StatusOK, notifs)
}

// MarkRead handles PUT /api/v1/notifications/{id}/read
func (h *NotificationHandler) MarkRead(w http.ResponseWriter, r *http.Request) {
	notifID := chi.URLParam(r, "id")
	if err := h.svc.MarkRead(notifID); err != nil {
		writeJSON(w, http.StatusNotFound, map[string]string{"error": err.Error()})
		return
	}
	writeJSON(w, http.StatusOK, map[string]string{"message": "marked as read"})
}

// MarkAllRead handles PUT /api/v1/notifications/read-all
func (h *NotificationHandler) MarkAllRead(w http.ResponseWriter, r *http.Request) {
	userID := r.URL.Query().Get("user_id")
	if userID == "" {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "user_id query parameter is required"})
		return
	}
	if err := h.svc.MarkAllRead(userID); err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
		return
	}
	writeJSON(w, http.StatusOK, map[string]string{"message": "all marked as read"})
}

// GetSettings handles GET /api/v1/notifications/settings
func (h *NotificationHandler) GetSettings(w http.ResponseWriter, r *http.Request) {
	userID := r.URL.Query().Get("user_id")
	if userID == "" {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "user_id query parameter is required"})
		return
	}
	s, err := h.svc.GetSettings(userID)
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
		return
	}
	writeJSON(w, http.StatusOK, s)
}

// UpdateSettings handles PUT /api/v1/notifications/settings
func (h *NotificationHandler) UpdateSettings(w http.ResponseWriter, r *http.Request) {
	var s model.NotificationSettings
	if err := json.NewDecoder(r.Body).Decode(&s); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid request body"})
		return
	}
	if s.UserID == "" {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "user_id is required"})
		return
	}
	h.svc.UpdateSettings(s.UserID, &s)
	writeJSON(w, http.StatusOK, map[string]string{"message": "settings updated"})
}

func writeJSON(w http.ResponseWriter, status int, v interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(v)
}
