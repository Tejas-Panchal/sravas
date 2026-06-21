package model

// Preference stores user-specific settings
type Preference struct {
	UserID               string `json:"user_id"`
	Theme                string `json:"theme"`
	NotificationsEnabled bool   `json:"notifications_enabled"`
}
