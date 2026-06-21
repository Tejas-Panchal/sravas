package middleware

import "github.com/Tejas-Panchal/sravas/services/notification-service/internal/model"

// DefaultSettings returns a NotificationSettings struct with sensible defaults
func DefaultSettings(userID string) *model.NotificationSettings {
	return &model.NotificationSettings{
		UserID:          userID,
		EmailEnabled:    true,
		PushEnabled:     true,
		NewSubscriberOn: true,
		NewVideoOn:      true,
		LikeOn:          true,
		CommentOn:       true,
		ReplyOn:         true,
	}
}
