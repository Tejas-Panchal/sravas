package service

import (
	"crypto/rand"
	"fmt"
	"sort"
	"sync"
	"time"

	"github.com/Tejas-Panchal/sravas/services/notification-service/internal/middleware"
	"github.com/Tejas-Panchal/sravas/services/notification-service/internal/model"
)

var (
	mu            sync.RWMutex
	notifications = map[string]*model.Notification{}
	stts          = map[string]*model.NotificationSettings{}
)

// NotificationService handles notification CRUD, read tracking, and push/email delivery
type NotificationService struct {
	ws    WebSocketHub
	email EmailSender
}

// NewNotificationService creates a NotificationService with the given hub and email sender
func NewNotificationService(ws WebSocketHub, email EmailSender) *NotificationService {
	return &NotificationService{ws: ws, email: email}
}

func generateID() string {
	b := make([]byte, 16)
	rand.Read(b)
	return fmt.Sprintf("%08x-%04x-%04x-%04x-%012x", b[0:4], b[4:6], b[6:8], b[8:10], b[10:16])
}

// CreateNotification creates a notification for a user and attempts push + email delivery
func (s *NotificationService) CreateNotification(userID string, notifType model.NotificationType, title, message string, data map[string]interface{}) (*model.Notification, error) {
	now := time.Now()
	n := &model.Notification{
		ID:        generateID(),
		UserID:    userID,
		Type:      notifType,
		Title:     title,
		Message:   message,
		Data:      data,
		Read:      false,
		CreatedAt: now,
	}

	mu.Lock()
	notifications[n.ID] = n
	userSettings := stts[userID]
	mu.Unlock()

	// Push via WebSocket if enabled
	if userSettings == nil || userSettings.PushEnabled {
		s.ws.SendToUser(userID, WSMessage{
			UserID: userID,
			Event:  "notification",
			Data:   n,
		})
	}

	// Send email if enabled
	if userSettings == nil || userSettings.EmailEnabled {
		s.email.Send(Email{
			To:      userID,
			Subject: title,
			Body:    message,
		})
	}

	return n, nil
}

// GetNotifications returns all notifications for a user, sorted by newest first
func (s *NotificationService) GetNotifications(userID string) ([]*model.Notification, error) {
	mu.RLock()
	defer mu.RUnlock()

	var result []*model.Notification
	for _, n := range notifications {
		if n.UserID == userID {
			result = append(result, n)
		}
	}

	sort.Slice(result, func(i, j int) bool {
		return result[i].CreatedAt.After(result[j].CreatedAt)
	})

	return result, nil
}

// MarkRead marks a single notification as read
func (s *NotificationService) MarkRead(notifID string) error {
	mu.Lock()
	defer mu.Unlock()

	n, ok := notifications[notifID]
	if !ok {
		return fmt.Errorf("notification not found")
	}
	n.Read = true
	return nil
}

// MarkAllRead marks all notifications for a user as read
func (s *NotificationService) MarkAllRead(userID string) error {
	mu.Lock()
	defer mu.Unlock()

	for _, n := range notifications {
		if n.UserID == userID {
			n.Read = true
		}
	}
	return nil
}

// GetSettings returns notification preferences for a user
func (s *NotificationService) GetSettings(userID string) (*model.NotificationSettings, error) {
	mu.RLock()
	defer mu.RUnlock()

	userSettings, ok := stts[userID]
	if !ok {
		return middleware.DefaultSettings(userID), nil
	}
	return userSettings, nil
}

// UpdateSettings saves notification preferences for a user
func (s *NotificationService) UpdateSettings(userID string, updated *model.NotificationSettings) {
	mu.Lock()
	stts[userID] = updated
	mu.Unlock()
}
