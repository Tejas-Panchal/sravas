package model

import "time"

// Comment represents a user comment on a video
type Comment struct {
	ID        string    `json:"id"`
	VideoID   string    `json:"video_id"`
	UserID    string    `json:"user_id"`
	ParentID  string    `json:"parent_id,omitempty"`
	Content   string    `json:"content"`
	LikeCount int       `json:"like_count"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// AddCommentRequest is the payload for POST /api/v1/videos/{id}/comments
type AddCommentRequest struct {
	UserID  string `json:"user_id"`
	Content string `json:"content"`
}

// EditCommentRequest is the payload for PUT /api/v1/comments/{id}
type EditCommentRequest struct {
	Content string `json:"content"`
}

// ReplyRequest is the payload for POST /api/v1/comments/{id}/replies
type ReplyRequest struct {
	UserID  string `json:"user_id"`
	Content string `json:"content"`
}

// CommentResponse wraps a list of comments with pagination cursor
type CommentResponse struct {
	Comments   []*Comment `json:"comments"`
	NextCursor string     `json:"next_cursor,omitempty"`
}
