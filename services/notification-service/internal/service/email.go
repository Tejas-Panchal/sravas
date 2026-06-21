package service

import (
	"encoding/json"
	"log"
)

// Email represents an outbound email message
type Email struct {
	To      string `json:"to"`
	Subject string `json:"subject"`
	Body    string `json:"body"`
}

// EmailSender defines the interface for sending emails (SendGrid stub)
type EmailSender interface {
	Send(email Email) error
}

// LogSender is a stub implementation that logs emails to stdout
type LogSender struct{}

// NewLogSender creates a new LogSender
func NewLogSender() *LogSender {
	return &LogSender{}
}

// Send logs the email as JSON to stdout (will be replaced with real SendGrid)
func (s *LogSender) Send(email Email) error {
	raw, _ := json.Marshal(email)
	log.Printf("[EMAIL] %s", string(raw))
	return nil
}
