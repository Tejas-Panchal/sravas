package service

import (
	"crypto/rand"
	"fmt"
	"sort"
	"sync"
	"time"

	"github.com/Tejas-Panchal/sravas/services/comment-service/internal/model"
)

var (
	mu       sync.RWMutex
	comments = map[string]*model.Comment{}
	likes    = map[string]map[string]bool{}
)

// CommentService handles comment CRUD, replies, pagination, and likes
type CommentService struct {
	bus EventBus
}

// NewCommentService creates a CommentService with the given event bus
func NewCommentService(bus EventBus) *CommentService {
	return &CommentService{bus: bus}
}

func generateID() string {
	b := make([]byte, 16)
	rand.Read(b)
	return fmt.Sprintf("%08x-%04x-%04x-%04x-%012x", b[0:4], b[4:6], b[6:8], b[8:10], b[10:16])
}

// AddComment creates a new top-level comment on a video
func (s *CommentService) AddComment(videoID string, req model.AddCommentRequest) (*model.Comment, error) {
	now := time.Now()
	c := &model.Comment{
		ID:        generateID(),
		VideoID:   videoID,
		UserID:    req.UserID,
		Content:   req.Content,
		CreatedAt: now,
		UpdatedAt: now,
	}

	mu.Lock()
	comments[c.ID] = c
	mu.Unlock()

	s.bus.Publish("comment.created", c.ID, c)
	return c, nil
}

// GetComments returns paginated top-level comments for a video, sorted by newest first
func (s *CommentService) GetComments(videoID, cursor string, limit int) (*model.CommentResponse, error) {
	mu.RLock()
	var all []*model.Comment
	for _, c := range comments {
		if c.VideoID == videoID && c.ParentID == "" {
			all = append(all, c)
		}
	}
	mu.RUnlock()

	sort.Slice(all, func(i, j int) bool {
		return all[i].CreatedAt.After(all[j].CreatedAt)
	})

	start := 0
	if cursor != "" {
		for i, c := range all {
			cursorID := fmt.Sprintf("%s:%d", c.ID, c.CreatedAt.UnixNano())
			if cursorID == cursor {
				start = i + 1
				break
			}
		}
	}

	if start > len(all) {
		start = len(all)
	}

	end := start + limit
	if end > len(all) {
		end = len(all)
	}

	var page []*model.Comment
	if start < len(all) {
		page = all[start:end]
	} else {
		page = []*model.Comment{}
	}

	var nextCursor string
	if end < len(all) {
		last := all[end-1]
		nextCursor = fmt.Sprintf("%s:%d", last.ID, last.CreatedAt.UnixNano())
	}

	return &model.CommentResponse{
		Comments:   page,
		NextCursor: nextCursor,
	}, nil
}

// EditComment updates the content of a comment
func (s *CommentService) EditComment(commentID string, content string) (*model.Comment, error) {
	mu.Lock()
	defer mu.Unlock()

	c, ok := comments[commentID]
	if !ok {
		return nil, fmt.Errorf("comment not found")
	}
	c.Content = content
	c.UpdatedAt = time.Now()
	return c, nil
}

// DeleteComment removes a comment and its likes
func (s *CommentService) DeleteComment(commentID string) error {
	mu.Lock()
	defer mu.Unlock()

	if _, ok := comments[commentID]; !ok {
		return fmt.Errorf("comment not found")
	}
	delete(comments, commentID)
	delete(likes, commentID)
	return nil
}

// ReplyToComment creates a reply (child comment) to an existing comment
func (s *CommentService) ReplyToComment(parentID string, req model.ReplyRequest) (*model.Comment, error) {
	mu.RLock()
	parent, ok := comments[parentID]
	mu.RUnlock()
	if !ok {
		return nil, fmt.Errorf("parent comment not found")
	}

	now := time.Now()
	reply := &model.Comment{
		ID:        generateID(),
		VideoID:   parent.VideoID,
		UserID:    req.UserID,
		ParentID:  parentID,
		Content:   req.Content,
		CreatedAt: now,
		UpdatedAt: now,
	}

	mu.Lock()
	comments[reply.ID] = reply
	mu.Unlock()

	s.bus.Publish("comment.reply.created", reply.ID, reply)
	return reply, nil
}

// LikeComment toggles a like on a comment (adds if not liked, removes if already liked)
func (s *CommentService) LikeComment(commentID, userID string) (bool, error) {
	mu.Lock()
	defer mu.Unlock()

	if _, ok := comments[commentID]; !ok {
		return false, fmt.Errorf("comment not found")
	}

	if likes[commentID] == nil {
		likes[commentID] = make(map[string]bool)
	}

	liked := likes[commentID][userID]
	if liked {
		delete(likes[commentID], userID)
		comments[commentID].LikeCount--
	} else {
		likes[commentID][userID] = true
		comments[commentID].LikeCount++
	}

	return !liked, nil
}
