package handler

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"

	"github.com/Tejas-Panchal/sravas/services/user-service/internal/middleware"
	"github.com/Tejas-Panchal/sravas/services/user-service/internal/service"
)

// GetSubscriptions returns the list of channels the user subscribes to
func GetSubscriptions(w http.ResponseWriter, r *http.Request) {
	userID := chi.URLParam(r, "id")
	subs := service.GetSubscriptions(userID)
	middleware.WriteJSON(w, http.StatusOK, subs)
}

// Subscribe subscribes the authenticated user to a channel
func Subscribe(w http.ResponseWriter, r *http.Request) {
	userID := chi.URLParam(r, "id")
	var req struct {
		ChannelID string `json:"channel_id"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		middleware.WriteJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid request body"})
		return
	}
	if req.ChannelID == "" {
		middleware.WriteJSON(w, http.StatusBadRequest, map[string]string{"error": "channel_id is required"})
		return
	}
	if err := service.Subscribe(userID, req.ChannelID); err != nil {
		middleware.WriteJSON(w, http.StatusBadRequest, map[string]string{"error": err.Error()})
		return
	}
	middleware.WriteJSON(w, http.StatusOK, map[string]string{"message": "subscribed"})
}
