package handler

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"

	"github.com/Tejas-Panchal/sravas/services/user-service/internal/middleware"
	"github.com/Tejas-Panchal/sravas/services/user-service/internal/model"
	"github.com/Tejas-Panchal/sravas/services/user-service/internal/service"
)

// Register creates a new user account
func Register(w http.ResponseWriter, r *http.Request) {
	var req model.RegisterRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		middleware.WriteJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid request body"})
		return
	}
	if req.Email == "" || req.Username == "" || req.Password == "" {
		middleware.WriteJSON(w, http.StatusBadRequest, map[string]string{"error": "email, username, and password are required"})
		return
	}
	if !middleware.ValidateEmail(req.Email) {
		middleware.WriteJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid email format"})
		return
	}
	if !middleware.ValidatePassword(req.Password) {
		middleware.WriteJSON(w, http.StatusBadRequest, map[string]string{"error": "password must be 8+ chars with at least 1 letter and 1 digit"})
		return
	}
	user, err := service.Register(req)
	if err != nil {
		middleware.WriteJSON(w, http.StatusConflict, map[string]string{"error": err.Error()})
		return
	}
	middleware.WriteJSON(w, http.StatusCreated, user)
}

// Login authenticates a user and returns JWT tokens
func Login(w http.ResponseWriter, r *http.Request) {
	var req model.LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		middleware.WriteJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid request body"})
		return
	}
	if req.Email == "" || req.Password == "" {
		middleware.WriteJSON(w, http.StatusBadRequest, map[string]string{"error": "email and password are required"})
		return
	}
	tokens, err := service.Login(req)
	if err != nil {
		middleware.WriteJSON(w, http.StatusUnauthorized, map[string]string{"error": err.Error()})
		return
	}
	middleware.WriteJSON(w, http.StatusOK, tokens)
}

// GetProfile returns the authenticated user's profile
func GetProfile(w http.ResponseWriter, r *http.Request) {
	userID := chi.URLParam(r, "id")
	user, err := service.GetProfile(userID)
	if err != nil {
		middleware.WriteJSON(w, http.StatusNotFound, map[string]string{"error": err.Error()})
		return
	}
	middleware.WriteJSON(w, http.StatusOK, user)
}

// UpdateProfile updates the authenticated user's profile fields
func UpdateProfile(w http.ResponseWriter, r *http.Request) {
	userID := chi.URLParam(r, "id")
	var req model.UpdateProfileRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		middleware.WriteJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid request body"})
		return
	}
	user, err := service.UpdateProfile(userID, req)
	if err != nil {
		middleware.WriteJSON(w, http.StatusBadRequest, map[string]string{"error": err.Error()})
		return
	}
	middleware.WriteJSON(w, http.StatusOK, user)
}

// DeleteAccount permanently deletes the user's account
func DeleteAccount(w http.ResponseWriter, r *http.Request) {
	userID := chi.URLParam(r, "id")
	if err := service.DeleteAccount(userID); err != nil {
		middleware.WriteJSON(w, http.StatusBadRequest, map[string]string{"error": err.Error()})
		return
	}
	middleware.WriteJSON(w, http.StatusOK, map[string]string{"message": "account deleted"})
}

// GetUserVideos returns a list of videos uploaded by the user
func GetUserVideos(w http.ResponseWriter, r *http.Request) {
	userID := chi.URLParam(r, "id")
	videos := service.GetUserVideos(userID)
	middleware.WriteJSON(w, http.StatusOK, videos)
}

// GetChannelStats returns aggregated statistics for the user's channel
func GetChannelStats(w http.ResponseWriter, r *http.Request) {
	userID := chi.URLParam(r, "id")
	stats, err := service.GetChannelStats(userID)
	if err != nil {
		middleware.WriteJSON(w, http.StatusNotFound, map[string]string{"error": err.Error()})
		return
	}
	middleware.WriteJSON(w, http.StatusOK, stats)
}
