package service

import (
	"errors"
	"sync"
	"time"

	"golang.org/x/crypto/bcrypt"

	"github.com/Tejas-Panchal/sravas/services/user-service/internal/model"
)

var (
	mu    sync.RWMutex
	users = map[string]*model.User{}  // email -> user
	subs  = map[string]map[string]bool{} // userID -> set of channelIDs
)

// Register creates a new user account
func Register(req model.RegisterRequest) (*model.User, error) {
	mu.Lock()
	defer mu.Unlock()

	for _, u := range users {
		if u.Email == req.Email {
			return nil, errors.New("email already registered")
		}
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	now := time.Now()
	user := &model.User{
		ID:           req.Email,
		Email:        req.Email,
		Username:     req.Username,
		PasswordHash: string(hash),
		CreatedAt:    now,
		UpdatedAt:    now,
	}
	users[req.Email] = user
	return user, nil
}

// Login authenticates a user and returns a success message
func Login(req model.LoginRequest) (interface{}, error) {
	mu.RLock()
	user, exists := users[req.Email]
	mu.RUnlock()

	if !exists {
		return nil, errors.New("invalid email or password")
	}
	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(req.Password)); err != nil {
		return nil, errors.New("invalid email or password")
	}
	return map[string]string{"message": "login successful"}, nil
}

// GetProfile returns a user by ID
func GetProfile(userID string) (*model.User, error) {
	mu.RLock()
	defer mu.RUnlock()

	for _, u := range users {
		if u.ID == userID {
			return u, nil
		}
	}
	return nil, errors.New("user not found")
}

// UpdateProfile updates a user's profile fields
func UpdateProfile(userID string, req model.UpdateProfileRequest) (*model.User, error) {
	mu.Lock()
	defer mu.Unlock()

	for _, u := range users {
		if u.ID == userID {
			if req.Username != "" {
				u.Username = req.Username
			}
			if req.Bio != "" {
				u.Bio = req.Bio
			}
			if req.ProfilePic != "" {
				u.ProfilePic = req.ProfilePic
			}
			u.UpdatedAt = time.Now()
			return u, nil
		}
	}
	return nil, errors.New("user not found")
}

// DeleteAccount removes a user
func DeleteAccount(userID string) error {
	mu.Lock()
	defer mu.Unlock()

	for email, u := range users {
		if u.ID == userID {
			delete(users, email)
			return nil
		}
	}
	return errors.New("user not found")
}

// Subscribe adds a channel subscription
func Subscribe(subscriberID, channelID string) error {
	mu.Lock()
	defer mu.Unlock()

	if subs[subscriberID] == nil {
		subs[subscriberID] = make(map[string]bool)
	}
	subs[subscriberID][channelID] = true
	return nil
}

// GetSubscriptions returns the list of channels a user subscribes to
func GetSubscriptions(userID string) []string {
	mu.RLock()
	defer mu.RUnlock()

	var result []string
	for channelID := range subs[userID] {
		result = append(result, channelID)
	}
	return result
}

// GetUserVideos returns a list of video IDs for the user
func GetUserVideos(userID string) []string {
	return []string{} // TODO: query upload-service or DB
}

// GetChannelStats returns placeholder channel statistics
func GetChannelStats(userID string) (*model.ChannelStats, error) {
	mu.RLock()
	defer mu.RUnlock()

	for _, u := range users {
		if u.ID == userID {
			subCount := len(subs[userID])
			return &model.ChannelStats{
				UserID:          userID,
				VideoCount:      0,
				SubscriberCount: subCount,
				TotalViews:      0,
			}, nil
		}
	}
	return nil, errors.New("user not found")
}
