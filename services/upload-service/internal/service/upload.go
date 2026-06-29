package service

import (
	"crypto/rand"
	"fmt"
	"mime/multipart"
	"sync"
	"time"

	"github.com/Tejas-Panchal/sravas/services/upload-service/internal/model"
)

var (
	mu      sync.RWMutex
	videos  = map[string]*model.Video{}            // videoID -> video
	uploads = map[string]*model.UploadProgress{}   // videoID -> progress
)

// UploadService handles video upload business logic
type UploadService struct {
	store  Storage
	bus    EventBus
	cdnURL string
}

// NewUploadService creates an UploadService with the given storage, event bus, and optional CDN URL
func NewUploadService(store Storage, bus EventBus, cdnURL string) *UploadService {
	return &UploadService{store: store, bus: bus, cdnURL: cdnURL}
}

// generateID creates a unique video ID using crypto/rand
func generateID() string {
	b := make([]byte, 16)
	rand.Read(b)
	return fmt.Sprintf("%08x-%04x-%04x-%04x-%012x",
		b[0:4], b[4:6], b[6:8], b[8:10], b[10:16])
}

// ProcessUpload handles a complete upload: validates, stores, creates record, publishes event
func (s *UploadService) ProcessUpload(fileHeader *multipart.FileHeader, userID, title, description string) (*model.Video, error) {
	videoID := generateID()

	file, err := fileHeader.Open()
	if err != nil {
		return nil, fmt.Errorf("failed to open uploaded file: %w", err)
	}
	defer file.Close()

	// Save file to storage
	filename := videoID + "_" + fileHeader.Filename
	if err := s.store.Save(filename, file); err != nil {
		return nil, fmt.Errorf("failed to save file: %w", err)
	}

	now := time.Now()
	video := &model.Video{
		ID:          videoID,
		UserID:      userID,
		Title:       title,
		Description: description,
		FileName:    filename,
		FileSize:    fileHeader.Size,
		Status:      model.StatusUploading,
		CreatedAt:   now,
		UpdatedAt:   now,
	}

	if s.cdnURL != "" {
		video.URL = s.cdnURL + "/" + filename
	}

	mu.Lock()
	videos[videoID] = video
	uploads[videoID] = &model.UploadProgress{
		VideoID:      videoID,
		Progress:     0,
		FileSize:     fileHeader.Size,
		UploadedSize: fileHeader.Size,
	}
	mu.Unlock()

	// Publish event to Kafka
	s.bus.Publish("video.upload.queued", videoID, map[string]interface{}{
		"video_id":  videoID,
		"user_id":   userID,
		"file_name": fileHeader.Filename,
		"file_size": fileHeader.Size,
	})

	return video, nil
}

// GetStatus returns the current upload status and progress for a video
func (s *UploadService) GetStatus(videoID string) (*model.Video, *model.UploadProgress, error) {
	mu.RLock()
	defer mu.RUnlock()

	video, ok := videos[videoID]
	if !ok {
		return nil, nil, fmt.Errorf("video not found")
	}
	progress := uploads[videoID]
	return video, progress, nil
}

// UpdateMetadata updates the title and description of a video
func (s *UploadService) UpdateMetadata(videoID, title, description string) (*model.Video, error) {
	mu.Lock()
	defer mu.Unlock()

	video, ok := videos[videoID]
	if !ok {
		return nil, fmt.Errorf("video not found")
	}
	if title != "" {
		video.Title = title
	}
	if description != "" {
		video.Description = description
	}
	video.UpdatedAt = time.Now()
	return video, nil
}

// DeleteVideo removes a video from storage and the registry
func (s *UploadService) DeleteVideo(videoID string) error {
	mu.Lock()
	video, ok := videos[videoID]
	if ok {
		delete(videos, videoID)
		delete(uploads, videoID)
	}
	mu.Unlock()

	if !ok {
		return fmt.Errorf("video not found")
	}

	// Remove the file from storage
	if err := s.store.Delete(video.FileName); err != nil {
		return fmt.Errorf("failed to delete file: %w", err)
	}

	return nil
}
