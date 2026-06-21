package model

import "time"

// User represents a user profile
type User struct {
	ID           string    `json:"id"`
	Email        string    `json:"email"`
	Username     string    `json:"username"`
	PasswordHash string    `json:"-"`
	Bio          string    `json:"bio"`
	ProfilePic   string    `json:"profile_pic"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

// RegisterRequest is the payload for POST /api/v1/users/register
type RegisterRequest struct {
	Email    string `json:"email"`
	Username string `json:"username"`
	Password string `json:"password"`
}

// LoginRequest is the payload for POST /api/v1/users/login
type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

// UpdateProfileRequest is the payload for PUT /api/v1/users/{id}
type UpdateProfileRequest struct {
	Username   string `json:"username"`
	Bio        string `json:"bio"`
	ProfilePic string `json:"profile_pic"`
}

// ChannelStats holds aggregated channel statistics
type ChannelStats struct {
	UserID          string `json:"user_id"`
	VideoCount      int    `json:"video_count"`
	SubscriberCount int    `json:"subscriber_count"`
	TotalViews      int    `json:"total_views"`
}
