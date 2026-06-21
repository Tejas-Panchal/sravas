package service

import (
	"fmt"
	"sync"
	"time"

	"github.com/Tejas-Panchal/sravas/services/analytics-service/internal/model"
)

var (
	mu     sync.RWMutex
	events []*model.Event
)

const activeViewerWindow = 5 * time.Minute

// AnalyticsService handles event ingestion and analytics computation
type AnalyticsService struct {
	bus EventBus
}

// NewAnalyticsService creates an AnalyticsService with the given event bus
func NewAnalyticsService(bus EventBus) *AnalyticsService {
	return &AnalyticsService{bus: bus}
}

// TrackEvent stores an analytics event and publishes it to the event bus
func (s *AnalyticsService) TrackEvent(req model.TrackEventRequest) error {
	evt := &model.Event{
		Type:      req.Type,
		VideoID:   req.VideoID,
		UserID:    req.UserID,
		Source:    req.Source,
		Seconds:   req.Seconds,
		Timestamp: time.Now(),
	}

	mu.Lock()
	events = append(events, evt)
	mu.Unlock()

	s.bus.Publish("analytics.event", evt.VideoID, evt)
	return nil
}

// GetVideoAnalytics computes analytics for a specific video
func (s *AnalyticsService) GetVideoAnalytics(videoID string) (*model.VideoAnalytics, error) {
	mu.RLock()
	defer mu.RUnlock()

	var (
		viewCount        int64
		totalWatchTime   float64
		likeCount        int64
		commentCount     int64
		activeCount      int
		trafficSources   = make(map[string]int64)
		retentionBuckets = make(map[float64]int64)
		now              = time.Now()
	)

	for _, e := range events {
		if e.VideoID != videoID {
			continue
		}

		switch e.Type {
		case "view":
			viewCount++
			if now.Sub(e.Timestamp) <= activeViewerWindow {
				activeCount++
			}
		case "watch_time":
			totalWatchTime += e.Seconds
			bucket := float64(int(e.Seconds/10)) * 10
			retentionBuckets[bucket]++
		case "like":
			likeCount++
		case "comment":
			commentCount++
		case "traffic_source":
			trafficSources[e.Source]++
		}
	}

	if viewCount == 0 {
		return nil, fmt.Errorf("no analytics data for video %s", videoID)
	}

	var retention []*model.RetentionPoint
	for ts, count := range retentionBuckets {
		retention = append(retention, &model.RetentionPoint{Timestamp: ts, ViewCount: count})
	}

	return &model.VideoAnalytics{
		VideoID:        videoID,
		ViewCount:      viewCount,
		TotalWatchTime: totalWatchTime,
		LikeCount:      likeCount,
		CommentCount:   commentCount,
		Retention:      retention,
		TrafficSources: trafficSources,
		ActiveViewers:  activeCount,
	}, nil
}

// GetChannelAnalytics computes analytics aggregated across all videos belonging to a user
func (s *AnalyticsService) GetChannelAnalytics(userID string) (*model.ChannelAnalytics, error) {
	mu.RLock()
	defer mu.RUnlock()

	var (
		totalViews int64
		totalWatch float64
		totalLikes int64
		totalComms int64
		videoIDs   = make(map[string]bool)
	)

	for _, e := range events {
		if e.UserID != userID {
			continue
		}
		videoIDs[e.VideoID] = true

		switch e.Type {
		case "view":
			totalViews++
		case "watch_time":
			totalWatch += e.Seconds
		case "like":
			totalLikes++
		case "comment":
			totalComms++
		}
	}

	if len(videoIDs) == 0 {
		return nil, fmt.Errorf("no analytics data for channel %s", userID)
	}

	return &model.ChannelAnalytics{
		UserID:         userID,
		TotalViews:     totalViews,
		TotalWatchTime: totalWatch,
		TotalLikes:     totalLikes,
		TotalComments:  totalComms,
		VideoCount:     len(videoIDs),
	}, nil
}

// GetTrending returns top videos by view count (used by the trending endpoint)
func (s *AnalyticsService) GetTrending() ([]*model.TrendingData, error) {
	mu.RLock()
	defer mu.RUnlock()

	videoViews := make(map[string]int64)
	for _, e := range events {
		if e.Type == "view" {
			videoViews[e.VideoID]++
		}
	}

	var trending []*model.TrendingData
	for videoID, count := range videoViews {
		trending = append(trending, &model.TrendingData{
			VideoID: videoID,
			Metric:  "views",
			Value:   count,
		})
	}

	for i := 0; i < len(trending); i++ {
		for j := i + 1; j < len(trending); j++ {
			if trending[j].Value > trending[i].Value {
				trending[i], trending[j] = trending[j], trending[i]
			}
		}
	}

	if len(trending) > 50 {
		trending = trending[:50]
	}

	return trending, nil
}
