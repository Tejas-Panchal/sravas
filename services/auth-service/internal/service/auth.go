package service

import (
	"errors"
	"sync"

	"github.com/Tejas-Panchal/sravas/services/auth-service/internal/model"
)

// In-memory stores (will be replaced by PostgreSQL in Week 2)
var (
	usersMu       sync.RWMutex
	users         = map[string]*model.User{} // email -> user
	resetTokensMu sync.RWMutex
	resetTokens   = map[string]string{} // token -> email
)

// Register creates a new user account
func Register(req model.RegisterRequest) (*model.User, error) {
	usersMu.Lock()
	defer usersMu.Unlock()

	if _, exists := users[req.Email]; exists {
		return nil, errors.New("email already registered")
	}

	hash, err := HashPassword(req.Password)
	if err != nil {
		return nil, err
	}

	user := &model.User{
		ID:           req.Email, // placeholder until PostgreSQL
		Email:        req.Email,
		Username:     req.Username,
		PasswordHash: hash,
		Verified:     false,
		CreatedAt:    "now", // placeholder
	}
	users[req.Email] = user
	return user, nil
}

// Login authenticates a user and returns JWT tokens
func Login(req model.LoginRequest) (*model.TokenPair, error) {
	usersMu.RLock()
	user, exists := users[req.Email]
	usersMu.RUnlock()

	if !exists {
		return nil, errors.New("invalid email or password")
	}
	if !CheckPassword(req.Password, user.PasswordHash) {
		return nil, errors.New("invalid email or password")
	}

	return GenerateTokenPair(user.ID, user.Email)
}

// Refresh validates a refresh token and returns a new token pair
func Refresh(refreshToken string) (*model.TokenPair, error) {
	claims, err := ValidateRefreshToken(refreshToken)
	if err != nil {
		return nil, errors.New("invalid or expired refresh token")
	}
	return GenerateTokenPair(claims.UserID, claims.Email)
}

// PasswordReset generates a reset token and "sends" it (stored in-memory for now)
func PasswordReset(email string) error {
	usersMu.RLock()
	_, exists := users[email]
	usersMu.RUnlock()

	if !exists {
		return nil // silently succeed to avoid email enumeration
	}

	token := "reset-token-placeholder" // TODO: generate real token
	resetTokensMu.Lock()
	resetTokens[token] = email
	resetTokensMu.Unlock()

	// TODO: send email via service.Email.SendPasswordReset
	return nil
}

// PasswordResetConfirm validates the reset token and updates the password
func PasswordResetConfirm(token, newPassword string) error {
	resetTokensMu.Lock()
	email, exists := resetTokens[token]
	delete(resetTokens, token)
	resetTokensMu.Unlock()

	if !exists {
		return errors.New("invalid or expired reset token")
	}

	hash, err := HashPassword(newPassword)
	if err != nil {
		return err
	}

	usersMu.Lock()
	if user, ok := users[email]; ok {
		user.PasswordHash = hash
	}
	usersMu.Unlock()

	return nil
}
