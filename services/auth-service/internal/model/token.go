package model

import "github.com/golang-jwt/jwt/v5"

// TokenPair contains the access and refresh tokens returned after login/refresh
type TokenPair struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

// Claims represents the JWT claims payload
type Claims struct {
	UserID string `json:"user_id"`
	Email  string `json:"email"`
	jwt.RegisteredClaims
}
