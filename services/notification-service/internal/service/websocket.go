package service

import (
	"encoding/json"
	"log"
)

// WSMessage represents a message pushed to a WebSocket client
type WSMessage struct {
	UserID string      `json:"user_id"`
	Event  string      `json:"event"`
	Data   interface{} `json:"data"`
}

// WebSocketHub defines the interface for real-time push notifications (Socket.io stub)
type WebSocketHub interface {
	SendToUser(userID string, message WSMessage) error
}

// LogHub is a stub implementation that logs WebSocket messages to stdout
type LogHub struct{}

// NewLogHub creates a new LogHub
func NewLogHub() *LogHub {
	return &LogHub{}
}

// SendToUser logs the WebSocket message as JSON to stdout
func (h *LogHub) SendToUser(userID string, message WSMessage) error {
	raw, _ := json.Marshal(message)
	log.Printf("[WS PUSH to %s] %s", userID, string(raw))
	return nil
}
