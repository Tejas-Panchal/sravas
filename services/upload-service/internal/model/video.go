package model

import "time"

// VideoStatus represents the processing state of a video
type VideoStatus string

const (
	StatusUploading  VideoStatus = "uploading"
	StatusProcessing VideoStatus = "processing"
	StatusReady      VideoStatus = "ready"
	StatusFailed     VideoStatus = "failed"
)

// Video represents a video's metadata
type Video struct {
	ID          string      `json:"id"`
	UserID      string      `json:"user_id"`
	Title       string      `json:"title"`
	Description string      `json:"description"`
	FileName    string      `json:"file_name"`
	FileSize    int64       `json:"file_size"`
	Status      VideoStatus `json:"status"`
	Duration    float64     `json:"duration,omitempty"`
	ViewCount   int64       `json:"view_count"`
	LikeCount   int64       `json:"like_count"`
	Thumbnail   string      `json:"thumbnail,omitempty"`
	CreatedAt   time.Time   `json:"created_at"`
	UpdatedAt   time.Time   `json:"updated_at"`
}

// UploadProgress tracks the progress of an active upload
type UploadProgress struct {
	VideoID      string `json:"video_id"`
	Progress     int    `json:"progress"`
	FileSize     int64  `json:"file_size"`
	UploadedSize int64  `json:"uploaded_size"`
}
