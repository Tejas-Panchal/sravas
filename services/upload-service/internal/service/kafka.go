package service

import (
	"encoding/json"
	"log"
)

// Event represents a message published to the event bus
type Event struct {
	Topic string      `json:"topic"`
	Key   string      `json:"key"`
	Data  interface{} `json:"data"`
}

// EventBus defines the interface for publishing events (Kafka stub)
type EventBus interface {
	Publish(topic string, key string, data interface{}) error
}

// LogEventBus is a stub implementation that logs events to stdout
type LogEventBus struct{}

// NewLogEventBus creates a new LogEventBus
func NewLogEventBus() *LogEventBus {
	return &LogEventBus{}
}

// Publish logs the event as JSON to stdout (will be replaced with real Kafka producer)
func (b *LogEventBus) Publish(topic string, key string, data interface{}) error {
	event := Event{Topic: topic, Key: key, Data: data}
	raw, _ := json.Marshal(event)
	log.Printf("[EVENT] %s", string(raw))
	return nil
}
