package service

import (
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"

	"github.com/Tejas-Panchal/sravas/services/auth-service/internal/model"
)

var (
	jwtSecret        = getEnv("JWT_SECRET", "change-me")
	jwtRefreshSecret = getEnv("JWT_REFRESH_SECRET", "change-me")
	accessExpiry     = getDuration("JWT_EXPIRY", 15*time.Minute)
	refreshExpiry    = getDuration("JWT_REFRESH_EXPIRY", 7*24*time.Hour)
)

func getEnv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}

func getDuration(key string, fallback time.Duration) time.Duration {
	if v := os.Getenv(key); v != "" {
		d, err := time.ParseDuration(v)
		if err == nil {
			return d
		}
	}
	return fallback
}

// GenerateTokenPair creates both an access token and a refresh token for the given user
func GenerateTokenPair(userID, email string) (*model.TokenPair, error) {
	accessToken, err := generateJWT(userID, email, jwtSecret, accessExpiry)
	if err != nil {
		return nil, err
	}
	refreshToken, err := generateJWT(userID, email, jwtRefreshSecret, refreshExpiry)
	if err != nil {
		return nil, err
	}
	return &model.TokenPair{AccessToken: accessToken, RefreshToken: refreshToken}, nil
}

// ValidateAccessToken parses and validates an access token, returning the claims
func ValidateAccessToken(tokenStr string) (*model.Claims, error) {
	return validateJWT(tokenStr, jwtSecret)
}

// ValidateRefreshToken parses and validates a refresh token, returning the claims
func ValidateRefreshToken(tokenStr string) (*model.Claims, error) {
	return validateJWT(tokenStr, jwtRefreshSecret)
}

func generateJWT(userID, email, secret string, expiry time.Duration) (string, error) {
	claims := model.Claims{
		UserID: userID,
		Email:  email,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(expiry)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(secret))
}

func validateJWT(tokenStr, secret string) (*model.Claims, error) {
	token, err := jwt.ParseWithClaims(tokenStr, &model.Claims{}, func(t *jwt.Token) (interface{}, error) {
		return []byte(secret), nil
	})
	if err != nil {
		return nil, err
	}
	claims, ok := token.Claims.(*model.Claims)
	if !ok || !token.Valid {
		return nil, jwt.ErrSignatureInvalid
	}
	return claims, nil
}
