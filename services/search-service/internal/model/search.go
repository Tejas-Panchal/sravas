package model

import "time"

// SearchVideo represents a video in the search index
type SearchVideo struct {
	ID          string    `json:"id"`
	UserID      string    `json:"user_id"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	Tags        []string  `json:"tags"`
	Category    string    `json:"category"`
	Thumbnail   string    `json:"thumbnail"`
	Duration    float64   `json:"duration"`
	ViewCount   int64     `json:"view_count"`
	LikeCount   int64     `json:"like_count"`
	UploadedAt  time.Time `json:"uploaded_at"`
}

// SearchRequest holds the parameters for a search query
type SearchRequest struct {
	Query    string `json:"q"`
	Category string `json:"category"`
	Sort     string `json:"sort"`
	Page     int    `json:"page"`
	PageSize int    `json:"page_size"`
}

// SearchResponse contains the search results and metadata
type SearchResponse struct {
	Results    []*SearchVideo `json:"results"`
	Total      int            `json:"total"`
	Page       int            `json:"page"`
	PageSize   int            `json:"page_size"`
	TotalPages int            `json:"total_pages"`
}

// Suggestion represents an autocomplete suggestion
type Suggestion struct {
	Text  string  `json:"text"`
	Score float64 `json:"score"`
}

// TrendingVideo represents a video on the trending page
type TrendingVideo struct {
	*SearchVideo
	TrendingScore float64 `json:"trending_score"`
}
