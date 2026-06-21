package service

import (
	"encoding/json"
	"log"
)

// Event represents a real-time event published to clients
type Event struct {
	Type string      `json:"type"`
	Data interface{} `json:"data"`
}

// EventBus defines the interface for publishing real-time events (Socket.io stub)
type EventBus interface {
	Publish(eventType string, key string, data interface{}) error
}

// LogEventBus is a stub implementation that logs events to stdout
type LogEventBus struct{}

// NewLogEventBus creates a new LogEventBus
func NewLogEventBus() *LogEventBus {
	return &LogEventBus{}
}

// Publish logs the event as JSON to stdout (will be replaced with real Socket.io broadcast)
func (b *LogEventBus) Publish(eventType string, key string, data interface{}) error {
	event := Event{Type: eventType, Data: data}
	raw, _ := json.Marshal(event)
	log.Printf("[WS EVENT] %s", string(raw))
	return nil
}
