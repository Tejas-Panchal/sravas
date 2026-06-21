package middleware

import "strings"

const maxContentLen = 10000

// ValidateContent checks that comment content is non-empty and within length limits
func ValidateContent(content string) string {
	trimmed := strings.TrimSpace(content)
	if trimmed == "" {
		return "content cannot be empty"
	}
	if len(trimmed) > maxContentLen {
		return "content exceeds maximum length of 10000 characters"
	}
	return ""
}
