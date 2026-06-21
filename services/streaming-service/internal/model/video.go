package model

import "time"

// StreamVideo represents video metadata for the streaming layer
type StreamVideo struct {
	ID          string    `json:"id"`
	UserID      string    `json:"user_id"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	Duration    float64   `json:"duration"`
	ViewCount   int64     `json:"view_count"`
	LikeCount   int64     `json:"like_count"`
	Thumbnail   string    `json:"thumbnail"`
	Status      string    `json:"status"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// HLSManifest represents the master playlist structure
type HLSManifest struct {
	VideoID    string  `json:"video_id"`
	Playlist   string  `json:"playlist"`
	SegmentDur float64 `json:"segment_duration"`
	SegCount   int     `json:"segment_count"`
}

// SegmentInfo identifies a single HLS segment
type SegmentInfo struct {
	VideoID   string `json:"video_id"`
	SegmentID string `json:"segment_id"`
	Data      []byte `json:"-"`
	Size      int64  `json:"size"`
}
