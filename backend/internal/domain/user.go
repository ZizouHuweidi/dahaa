package domain

import (
	"context"
	"errors"
	"time"
)

// Common errors
var (
	ErrUserNotFound       = errors.New("user not found")
	ErrUserAlreadyExists  = errors.New("user already exists")
	ErrInvalidCredentials = errors.New("invalid credentials")
)

// User represents a registered user
type User struct {
	ID           string    `json:"id"`
	Username     string    `json:"username"`
	Email        string    `json:"email"`
	PasswordHash string    `json:"-"` // Never expose password hash
	DisplayName  string    `json:"display_name"`
	AvatarURL    string    `json:"avatar_url,omitempty"`
	Stats        UserStats `json:"stats"`
	LastLoginAt  time.Time `json:"last_login_at"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

// UserStats represents user's game statistics
type UserStats struct {
	GamesPlayed   int `json:"games_played"`
	GamesWon      int `json:"games_won"`
	TotalPoints   int `json:"total_points"`
	TotalScore    int `json:"total_score"`
	HighestScore  int `json:"highest_score"`
	PerfectRounds int `json:"perfect_rounds"`
	FooledPlayers int `json:"fooled_players"`
	CorrectVotes  int `json:"correct_votes"`
}

// UserRepository defines the interface for user-related operations
type UserRepository interface {
	// Create creates a new user
	Create(ctx context.Context, user *User) error

	// GetByID retrieves a user by their ID
	GetByID(ctx context.Context, id string) (*User, error)

	// GetByUsername retrieves a user by their username
	GetByUsername(ctx context.Context, username string) (*User, error)

	// GetByEmail retrieves a user by their email
	GetByEmail(ctx context.Context, email string) (*User, error)

	// Update updates a user's information
	Update(ctx context.Context, user *User) error

	// Delete deletes a user
	Delete(ctx context.Context, id string) error

	// UpdateStats updates a user's game statistics
	UpdateStats(ctx context.Context, id string, stats UserStats) error
}

// GameInvite represents a game invitation
type GameInvite struct {
	ID        string    `json:"id"`
	GameID    string    `json:"game_id"`
	FromUser  string    `json:"from_user"` // User ID of the sender
	ToUser    string    `json:"to_user"`   // User ID of the recipient
	Status    string    `json:"status"`    // pending, accepted, declined
	CreatedAt time.Time `json:"created_at"`
	ExpiresAt time.Time `json:"expires_at"`
}

// GameInviteRepository defines the interface for game invitation operations
type GameInviteRepository interface {
	// Create creates a new game invitation
	Create(ctx context.Context, invite *GameInvite) error

	// GetByID retrieves an invitation by its ID
	GetByID(ctx context.Context, id string) (*GameInvite, error)

	// GetPendingInvites retrieves all pending invitations for a user
	GetPendingInvites(ctx context.Context, userID string) ([]*GameInvite, error)

	// UpdateStatus updates an invitation's status
	UpdateStatus(ctx context.Context, id string, status string) error

	// Delete deletes an invitation
	Delete(ctx context.Context, id string) error
}
