package service

import (
	"crypto/rand"
	"fmt"
	"sync"
	"time"

	"github.com/Tejas-Panchal/sravas/services/streaming-service/internal/model"
)

const (
	cacheTTL = time.Hour
	segDur   = 10.0
	segCount = 30
)

var (
	mu      sync.RWMutex
	videos  = map[string]*model.StreamVideo{} // videoID -> metadata
	views   = map[string]int64{}              // videoID -> view count
	watches = map[string]float64{}            // videoID -> total watch seconds
)

// StreamingService handles video streaming business logic
type StreamingService struct {
	cache Cacher
}

// NewStreamingService creates a StreamingService with the given cache backend
func NewStreamingService(cache Cacher) *StreamingService {
	return &StreamingService{cache: cache}
}

// GetMetadata returns cached or fresh video metadata
func (s *StreamingService) GetMetadata(videoID string) (*model.StreamVideo, error) {
	if cached, ok := s.cache.Get("meta:" + videoID); ok {
		return cached.(*model.StreamVideo), nil
	}

	mu.RLock()
	v, ok := videos[videoID]
	mu.RUnlock()
	if !ok {
		return nil, fmt.Errorf("video not found")
	}

	s.cache.Set("meta:"+videoID, v, cacheTTL)
	return v, nil
}

// GetManifest generates a stub HLS master playlist for the given video
func (s *StreamingService) GetManifest(videoID string) (*model.HLSManifest, error) {
	if cached, ok := s.cache.Get("manifest:" + videoID); ok {
		return cached.(*model.HLSManifest), nil
	}

	mu.RLock()
	_, ok := videos[videoID]
	mu.RUnlock()
	if !ok {
		return nil, fmt.Errorf("video not found")
	}

	m3u8 := "#EXTM3U\n"
	m3u8 += "#EXT-X-STREAM-INF:BANDWIDTH=2800000,RESOLUTION=1920x1080\n"
	m3u8 += fmt.Sprintf("/api/v1/videos/%s/segment/", videoID)
	m3u8 += "\n"

	manifest := &model.HLSManifest{
		VideoID:    videoID,
		Playlist:   m3u8,
		SegmentDur: segDur,
		SegCount:   segCount,
	}

	s.cache.Set("manifest:"+videoID, manifest, cacheTTL)
	return manifest, nil
}

// GetSegment returns a stub video segment (simulated .ts data)
func (s *StreamingService) GetSegment(videoID, segmentID string) (*model.SegmentInfo, error) {
	mu.RLock()
	_, ok := videos[videoID]
	mu.RUnlock()
	if !ok {
		return nil, fmt.Errorf("video not found")
	}

	// Generate a unique 1KB segment so each request feels like a real segment
	data := make([]byte, 1024)
	rand.Read(data)

	return &model.SegmentInfo{
		VideoID:   videoID,
		SegmentID: segmentID,
		Data:      data,
		Size:      int64(len(data)),
	}, nil
}

// RecordView increments the view counter for a video
func (s *StreamingService) RecordView(videoID string) error {
	mu.Lock()
	defer mu.Unlock()

	if _, ok := videos[videoID]; !ok {
		return fmt.Errorf("video not found")
	}
	views[videoID]++
	s.cache.Delete("meta:" + videoID) // bust cache so metadata reflects new count
	return nil
}

// UpdateAnalytics records watch time data for a video
func (s *StreamingService) UpdateAnalytics(videoID string, watchSeconds float64) error {
	mu.Lock()
	defer mu.Unlock()

	if _, ok := videos[videoID]; !ok {
		return fmt.Errorf("video not found")
	}
	watches[videoID] += watchSeconds
	return nil
}

// IngestVideo is called by upload-service (or a sync job) to make a video available for streaming
func (s *StreamingService) IngestVideo(v *model.StreamVideo) {
	mu.Lock()
	videos[v.ID] = v
	mu.Unlock()
}
