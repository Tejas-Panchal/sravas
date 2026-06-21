package middleware

import (
	"encoding/json"
	"net/http"
	"regexp"
	"strings"
)

var emailRegex = regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)

// Validate is a middleware that validates common request constraints
func Validate(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ct := r.Header.Get("Content-Type")
		if r.Method == "POST" || r.Method == "PUT" {
			if !strings.HasPrefix(ct, "application/json") {
				http.Error(w, `{"error":"content-type must be application/json"}`, http.StatusUnsupportedMediaType)
				return
			}
		}
		next.ServeHTTP(w, r)
	})
}

// ValidateEmail checks whether the provided string is a valid email format
func ValidateEmail(email string) bool {
	return emailRegex.MatchString(email)
}

// ValidatePassword enforces minimum password requirements (8+ chars, at least 1 letter and 1 digit)
func ValidatePassword(password string) bool {
	if len(password) < 8 {
		return false
	}
	hasLetter := false
	hasDigit := false
	for _, c := range password {
		if (c >= 'a' && c <= 'z') || (c >= 'A' && c <= 'Z') {
			hasLetter = true
		}
		if c >= '0' && c <= '9' {
			hasDigit = true
		}
	}
	return hasLetter && hasDigit
}

// WriteJSON sends a JSON response with the given status code
func WriteJSON(w http.ResponseWriter, status int, v interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(v)
}
