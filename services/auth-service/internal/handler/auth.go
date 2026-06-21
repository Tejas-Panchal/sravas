package handler

import (
	"encoding/json"
	"net/http"

	"github.com/Tejas-Panchal/sravas/services/auth-service/internal/middleware"
	"github.com/Tejas-Panchal/sravas/services/auth-service/internal/model"
	"github.com/Tejas-Panchal/sravas/services/auth-service/internal/service"
)

// Register handles user registration
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

// Refresh issues a new access token using a valid refresh token
func Refresh(w http.ResponseWriter, r *http.Request) {
	var req struct {
		RefreshToken string `json:"refresh_token"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		middleware.WriteJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid request body"})
		return
	}
	if req.RefreshToken == "" {
		middleware.WriteJSON(w, http.StatusBadRequest, map[string]string{"error": "refresh_token is required"})
		return
	}
	tokens, err := service.Refresh(req.RefreshToken)
	if err != nil {
		middleware.WriteJSON(w, http.StatusUnauthorized, map[string]string{"error": err.Error()})
		return
	}
	middleware.WriteJSON(w, http.StatusOK, tokens)
}

// PasswordReset initiates a password reset flow by sending a reset token to the user's email
func PasswordReset(w http.ResponseWriter, r *http.Request) {
	var req model.PasswordResetRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		middleware.WriteJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid request body"})
		return
	}
	if req.Email == "" || !middleware.ValidateEmail(req.Email) {
		middleware.WriteJSON(w, http.StatusBadRequest, map[string]string{"error": "valid email is required"})
		return
	}
	if err := service.PasswordReset(req.Email); err != nil {
		middleware.WriteJSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
		return
	}
	middleware.WriteJSON(w, http.StatusOK, map[string]string{"message": "reset token sent if email exists"})
}

// PasswordResetConfirm completes the password reset with a token and new password
func PasswordResetConfirm(w http.ResponseWriter, r *http.Request) {
	var req model.PasswordResetConfirmRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		middleware.WriteJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid request body"})
		return
	}
	if req.Token == "" || req.NewPassword == "" {
		middleware.WriteJSON(w, http.StatusBadRequest, map[string]string{"error": "token and new_password are required"})
		return
	}
	if !middleware.ValidatePassword(req.NewPassword) {
		middleware.WriteJSON(w, http.StatusBadRequest, map[string]string{"error": "password must be 8+ chars with at least 1 letter and 1 digit"})
		return
	}
	if err := service.PasswordResetConfirm(req.Token, req.NewPassword); err != nil {
		middleware.WriteJSON(w, http.StatusBadRequest, map[string]string{"error": err.Error()})
		return
	}
	middleware.WriteJSON(w, http.StatusOK, map[string]string{"message": "password updated successfully"})
}
