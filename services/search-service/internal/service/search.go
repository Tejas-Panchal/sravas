package service

import (
	"math"
	"sort"
	"time"

	"github.com/Tejas-Panchal/sravas/services/search-service/internal/model"
)

// SearchService handles search, suggestion, trending, and category logic
type SearchService struct {
	index Indexer
}

// NewSearchService creates a SearchService with the given index backend
func NewSearchService(index Indexer) *SearchService {
	return &SearchService{index: index}
}

// Search performs a full-text search with filters and pagination
func (s *SearchService) Search(req *model.SearchRequest) (*model.SearchResponse, error) {
	results := s.index.Search(req.Query)

	if req.Category != "" {
		var filtered []*model.SearchVideo
		for _, v := range results {
			if v.Category == req.Category {
				filtered = append(filtered, v)
			}
		}
		results = filtered
	}

	switch req.Sort {
	case "date":
		sort.Slice(results, func(i, j int) bool {
			return results[i].UploadedAt.After(results[j].UploadedAt)
		})
	case "views":
		sort.Slice(results, func(i, j int) bool {
			return results[i].ViewCount > results[j].ViewCount
		})
	}

	total := len(results)
	totalPages := int(math.Ceil(float64(total) / float64(req.PageSize)))
	start := (req.Page - 1) * req.PageSize
	if start > total {
		start = total
	}
	end := start + req.PageSize
	if end > total {
		end = total
	}

	var pageResults []*model.SearchVideo
	if start < total {
		pageResults = results[start:end]
	} else {
		pageResults = []*model.SearchVideo{}
	}

	return &model.SearchResponse{
		Results:    pageResults,
		Total:      total,
		Page:       req.Page,
		PageSize:   req.PageSize,
		TotalPages: totalPages,
	}, nil
}

// Suggest returns autocomplete suggestions
func (s *SearchService) Suggest(prefix string) ([]*model.Suggestion, error) {
	return s.index.Suggest(prefix), nil
}

// Trending returns videos sorted by a trending score (views x recency)
func (s *SearchService) Trending() ([]*model.TrendingVideo, error) {
	all := s.index.GetAll()
	now := time.Now()

	trending := make([]*model.TrendingVideo, 0, len(all))
	for _, v := range all {
		hoursSinceUpload := now.Sub(v.UploadedAt).Hours()
		score := float64(v.ViewCount) * math.Log1p(hoursSinceUpload)
		trending = append(trending, &model.TrendingVideo{
			SearchVideo:   v,
			TrendingScore: score,
		})
	}

	sort.Slice(trending, func(i, j int) bool {
		return trending[i].TrendingScore > trending[j].TrendingScore
	})

	if len(trending) > 50 {
		trending = trending[:50]
	}

	return trending, nil
}

// CategoryBrowse returns videos in a specific category, sorted by date
func (s *SearchService) CategoryBrowse(category string, page, pageSize int) (*model.SearchResponse, error) {
	all := s.index.GetAll()

	var filtered []*model.SearchVideo
	for _, v := range all {
		if v.Category == category {
			filtered = append(filtered, v)
		}
	}

	sort.Slice(filtered, func(i, j int) bool {
		return filtered[i].UploadedAt.After(filtered[j].UploadedAt)
	})

	total := len(filtered)
	totalPages := int(math.Ceil(float64(total) / float64(pageSize)))
	start := (page - 1) * pageSize
	if start > total {
		start = total
	}
	end := start + pageSize
	if end > total {
		end = total
	}

	var pageResults []*model.SearchVideo
	if start < total {
		pageResults = filtered[start:end]
	} else {
		pageResults = []*model.SearchVideo{}
	}

	return &model.SearchResponse{
		Results:    pageResults,
		Total:      total,
		Page:       page,
		PageSize:   pageSize,
		TotalPages: totalPages,
	}, nil
}

// IndexVideo adds or updates a video in the search index
func (s *SearchService) IndexVideo(video *model.SearchVideo) error {
	return s.index.Index(video)
}
