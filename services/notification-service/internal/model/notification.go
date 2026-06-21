package model

import "time"

// NotificationType represents the category of a notification
type NotificationType string

const (
	NotifNewSubscriber NotificationType = "new_subscriber"
	NotifNewVideo      NotificationType = "new_video"
	NotifLike          NotificationType = "like"
	NotifComment       NotificationType = "comment"
	NotifReply         NotificationType = "reply"
)

// Notification represents a single notification for a user
type Notification struct {
	ID        string                 `json:"id"`
	UserID    string                 `json:"user_id"`
	Type      NotificationType       `json:"type"`
	Title     string                 `json:"title"`
	Message   string                 `json:"message"`
	Data      map[string]interface{} `json:"data,omitempty"`
	Read      bool                   `json:"read"`
	CreatedAt time.Time              `json:"created_at"`
}

// NotificationSettings stores a user's notification preferences
type NotificationSettings struct {
	UserID           string `json:"user_id"`
	EmailEnabled     bool   `json:"email_enabled"`
	PushEnabled      bool   `json:"push_enabled"`
	NewSubscriberOn  bool   `json:"new_subscriber"`
	NewVideoOn       bool   `json:"new_video"`
	LikeOn           bool   `json:"like"`
	CommentOn        bool   `json:"comment"`
	ReplyOn          bool   `json:"reply"`
}
