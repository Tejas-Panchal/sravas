package service

import (
	"strings"
	"sync"
	"unicode"

	"github.com/Tejas-Panchal/sravas/services/search-service/internal/model"
)

// Indexer defines the interface for the search index (in-memory or Elasticsearch)
type Indexer interface {
	Index(video *model.SearchVideo) error
	Delete(videoID string) error
	Search(query string) []*model.SearchVideo
	Suggest(prefix string) []*model.Suggestion
	GetAll() []*model.SearchVideo
}

// MemoryIndex is an in-memory implementation of Indexer with term-based inverted index
type MemoryIndex struct {
	mu     sync.RWMutex
	videos map[string]*model.SearchVideo
	index  map[string][]string
}

// NewMemoryIndex creates a new MemoryIndex
func NewMemoryIndex() *MemoryIndex {
	return &MemoryIndex{
		videos: make(map[string]*model.SearchVideo),
		index:  make(map[string][]string),
	}
}

// tokenize splits text into lowercase terms (removes punctuation, keeps alphanumeric words)
func tokenize(text string) []string {
	var terms []string
	var current strings.Builder
	for _, r := range text {
		if unicode.IsLetter(r) || unicode.IsDigit(r) {
			current.WriteRune(unicode.ToLower(r))
		} else if current.Len() > 0 {
			terms = append(terms, current.String())
			current.Reset()
		}
	}
	if current.Len() > 0 {
		terms = append(terms, current.String())
	}
	return terms
}

// Index adds or updates a video in the in-memory index
func (m *MemoryIndex) Index(video *model.SearchVideo) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	if existing, ok := m.videos[video.ID]; ok {
		m.removeFromIndex(existing)
	}

	m.videos[video.ID] = video

	// Index title terms (weighted higher via duplicate entries for scoring)
	titleTerms := tokenize(video.Title)
	for _, t := range titleTerms {
		m.index[t] = append(m.index[t], video.ID)
		m.index[t] = append(m.index[t], video.ID)
	}

	// Index description terms
	for _, t := range tokenize(video.Description) {
		m.index[t] = append(m.index[t], video.ID)
	}

	// Index tags
	for _, tag := range video.Tags {
		terms := tokenize(tag)
		for _, t := range terms {
			m.index[t] = append(m.index[t], video.ID)
		}
	}

	return nil
}

// removeFromIndex deletes all index entries for a video
func (m *MemoryIndex) removeFromIndex(video *model.SearchVideo) {
	for term, ids := range m.index {
		filtered := ids[:0]
		for _, id := range ids {
			if id != video.ID {
				filtered = append(filtered, id)
			}
		}
		if len(filtered) == 0 {
			delete(m.index, term)
		} else {
			m.index[term] = filtered
		}
	}
}

// Search performs a multi-term search across the index and returns ranked results
func (m *MemoryIndex) Search(query string) []*model.SearchVideo {
	m.mu.RLock()
	defer m.mu.RUnlock()

	terms := tokenize(query)
	if len(terms) == 0 {
		return nil
	}

	scores := make(map[string]float64)
	for _, term := range terms {
		for _, id := range m.index[term] {
			scores[id]++
		}
	}

	type scoredVideo struct {
		video *model.SearchVideo
		score float64
	}
	var scored []scoredVideo
	for id, score := range scores {
		if v, ok := m.videos[id]; ok {
			scored = append(scored, scoredVideo{video: v, score: score})
		}
	}

	for i := 0; i < len(scored); i++ {
		for j := i + 1; j < len(scored); j++ {
			if scored[j].score > scored[i].score {
				scored[i], scored[j] = scored[j], scored[i]
			}
		}
	}

	results := make([]*model.SearchVideo, 0, len(scored))
	for _, s := range scored {
		results = append(results, s.video)
	}
	return results
}

// Suggest returns autocomplete suggestions for a prefix by finding matching indexed terms
func (m *MemoryIndex) Suggest(prefix string) []*model.Suggestion {
	m.mu.RLock()
	defer m.mu.RUnlock()

	lowerPrefix := strings.ToLower(prefix)
	if lowerPrefix == "" {
		return nil
	}

	var suggestions []*model.Suggestion
	seen := make(map[string]bool)

	for term := range m.index {
		if strings.HasPrefix(term, lowerPrefix) && !seen[term] {
			seen[term] = true
			suggestions = append(suggestions, &model.Suggestion{
				Text:  term,
				Score: float64(len(m.index[term])),
			})
		}
	}

	for i := 0; i < len(suggestions); i++ {
		for j := i + 1; j < len(suggestions); j++ {
			if suggestions[j].Score > suggestions[i].Score {
				suggestions[i], suggestions[j] = suggestions[j], suggestions[i]
			}
		}
	}

	if len(suggestions) > 10 {
		suggestions = suggestions[:10]
	}

	return suggestions
}

// Delete removes a video from the index
func (m *MemoryIndex) Delete(videoID string) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	if video, ok := m.videos[videoID]; ok {
		m.removeFromIndex(video)
		delete(m.videos, videoID)
	}
	return nil
}

// GetAll returns all indexed videos (used for trending and category browsing)
func (m *MemoryIndex) GetAll() []*model.SearchVideo {
	m.mu.RLock()
	defer m.mu.RUnlock()

	results := make([]*model.SearchVideo, 0, len(m.videos))
	for _, v := range m.videos {
		results = append(results, v)
	}
	return results
}
