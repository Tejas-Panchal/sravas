package service

import (
	"encoding/json"
	"log"
)

// AnalyticEvent represents a message published to the analytics event bus
type AnalyticEvent struct {
	Type string      `json:"type"`
	Data interface{} `json:"data"`
}

// EventBus defines the interface for publishing analytics events (Kafka stub)
type EventBus interface {
	Publish(eventType string, key string, data interface{}) error
}

// LogEventBus is a stub implementation that logs events to stdout
type LogEventBus struct{}

// NewLogEventBus creates a new LogEventBus
func NewLogEventBus() *LogEventBus {
	return &LogEventBus{}
}

// Publish logs the event as JSON to stdout (will be replaced with real Kafka producer)
func (b *LogEventBus) Publish(eventType string, key string, data interface{}) error {
	event := AnalyticEvent{Type: eventType, Data: data}
	raw, _ := json.Marshal(event)
	log.Printf("[ANALYTICS EVENT] %s", string(raw))
	return nil
}
