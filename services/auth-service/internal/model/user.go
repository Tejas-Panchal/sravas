package model

// User represents a user record in the auth service's domain
type User struct {
	ID           string `json:"id"`
	Email        string `json:"email"`
	Username     string `json:"username"`
	PasswordHash string `json:"-"`
	Verified     bool   `json:"verified"`
	CreatedAt    string `json:"created_at"`
}

// RegisterRequest is the payload for POST /api/v1/auth/register
type RegisterRequest struct {
	Email    string `json:"email"`
	Username string `json:"username"`
	Password string `json:"password"`
}

// LoginRequest is the payload for POST /api/v1/auth/login
type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

// PasswordResetRequest is the payload for POST /api/v1/auth/password-reset
type PasswordResetRequest struct {
	Email string `json:"email"`
}

// PasswordResetConfirmRequest is the payload for POST /api/v1/auth/password-reset/confirm
type PasswordResetConfirmRequest struct {
	Token       string `json:"token"`
	NewPassword string `json:"new_password"`
}
