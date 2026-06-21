package model

import "time"

// Event represents a single analytics event
type Event struct {
	Type      string    `json:"type"`
	VideoID   string    `json:"video_id"`
	UserID    string    `json:"user_id"`
	Source    string    `json:"source,omitempty"`
	Seconds   float64   `json:"seconds,omitempty"`
	Timestamp time.Time `json:"timestamp"`
}

// VideoAnalytics holds computed analytics for a single video
type VideoAnalytics struct {
	VideoID        string            `json:"video_id"`
	ViewCount      int64             `json:"view_count"`
	TotalWatchTime float64           `json:"total_watch_time"`
	LikeCount      int64             `json:"like_count"`
	CommentCount   int64             `json:"comment_count"`
	Retention      []*RetentionPoint `json:"retention"`
	TrafficSources map[string]int64  `json:"traffic_sources"`
	ActiveViewers  int               `json:"active_viewers"`
}

// ChannelAnalytics holds computed analytics aggregated across all videos for a channel
type ChannelAnalytics struct {
	UserID         string  `json:"user_id"`
	TotalViews     int64   `json:"total_views"`
	TotalWatchTime float64 `json:"total_watch_time"`
	TotalLikes     int64   `json:"total_likes"`
	TotalComments  int64   `json:"total_comments"`
	VideoCount     int     `json:"video_count"`
}

// RetentionPoint represents the audience retention at a given timestamp
type RetentionPoint struct {
	Timestamp float64 `json:"timestamp"`
	ViewCount int64   `json:"view_count"`
}

// TrendingData represents an aggregated trending metric
type TrendingData struct {
	VideoID string `json:"video_id"`
	Metric  string `json:"metric"`
	Value   int64  `json:"value"`
}

// TrackEventRequest is the payload for POST /api/v1/analytics/events
type TrackEventRequest struct {
	Type    string  `json:"type"`
	VideoID string  `json:"video_id"`
	UserID  string  `json:"user_id"`
	Source  string  `json:"source,omitempty"`
	Seconds float64 `json:"seconds,omitempty"`
}
