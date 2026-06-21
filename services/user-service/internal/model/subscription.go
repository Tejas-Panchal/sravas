package model

import "time"

// Subscription represents a user subscribing to a channel
type Subscription struct {
	SubscriberID string    `json:"subscriber_id"`
	ChannelID    string    `json:"channel_id"`
	CreatedAt    time.Time `json:"created_at"`
}
