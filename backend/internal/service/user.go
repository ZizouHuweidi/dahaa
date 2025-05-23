package service

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"errors"
	"fmt"
	"time"

	"golang.org/x/crypto/bcrypt"

	"github.com/zizouhuweidi/dahaa/internal/domain"
)

// UserService handles user-related operations
type UserService struct {
	userRepo   domain.UserRepository
	inviteRepo domain.GameInviteRepository
}

// NewUserService creates a new user service
func NewUserService(userRepo domain.UserRepository, inviteRepo domain.GameInviteRepository) *UserService {
	return &UserService{
		userRepo:   userRepo,
		inviteRepo: inviteRepo,
	}
}

// RegisterRequest represents a user registration request
type RegisterRequest struct {
	Username    string `json:"username" validate:"required,min=3,max=32"`
	Email       string `json:"email" validate:"required,email"`
	Password    string `json:"password" validate:"required,min=8"`
	DisplayName string `json:"display_name" validate:"required,min=2,max=50"`
}

// Register registers a new user
func (s *UserService) Register(ctx context.Context, req RegisterRequest) (*domain.User, error) {
	// Check if username is taken
	if _, err := s.userRepo.GetByUsername(ctx, req.Username); err == nil {
		return nil, domain.ErrUserAlreadyExists
	}

	// Check if email is taken
	if _, err := s.userRepo.GetByEmail(ctx, req.Email); err == nil {
		return nil, domain.ErrUserAlreadyExists
	}

	// Hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, fmt.Errorf("failed to hash password: %w", err)
	}

	// Create user
	user := &domain.User{
		ID:           generateID(),
		Username:     req.Username,
		Email:        req.Email,
		PasswordHash: string(hashedPassword),
		DisplayName:  req.DisplayName,
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}

	if err := s.userRepo.Create(ctx, user); err != nil {
		return nil, err
	}

	return user, nil
}

// LoginRequest represents a user login request
type LoginRequest struct {
	Username string `json:"username" validate:"required"`
	Password string `json:"password" validate:"required"`
}

// Login authenticates a user and returns a session token
func (s *UserService) Login(ctx context.Context, req LoginRequest) (string, error) {
	user, err := s.userRepo.GetByUsername(ctx, req.Username)
	if err != nil {
		if errors.Is(err, domain.ErrUserNotFound) {
			return "", domain.ErrInvalidCredentials
		}
		return "", err
	}

	// Verify password
	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(req.Password)); err != nil {
		return "", domain.ErrInvalidCredentials
	}

	// Generate session token
	token, err := generateSessionToken()
	if err != nil {
		return "", err
	}

	// Update last login time
	user.LastLoginAt = time.Now()
	if err := s.userRepo.Update(ctx, user); err != nil {
		return "", err
	}

	return token, nil
}

// SendGameInvite sends a game invitation to another user
func (s *UserService) SendGameInvite(ctx context.Context, gameID string, fromUserID string, toUserID string) error {
	// Check if users exist
	fromUser, err := s.userRepo.GetByID(ctx, fromUserID)
	if err != nil {
		return err
	}

	toUser, err := s.userRepo.GetByID(ctx, toUserID)
	if err != nil {
		return err
	}

	// Create invitation
	invite := &domain.GameInvite{
		ID:        generateID(),
		GameID:    gameID,
		FromUser:  fromUser.ID,
		ToUser:    toUser.ID,
		Status:    "pending",
		CreatedAt: time.Now(),
		ExpiresAt: time.Now().Add(24 * time.Hour),
	}

	return s.inviteRepo.Create(ctx, invite)
}

// AcceptGameInvite accepts a game invitation
func (s *UserService) AcceptGameInvite(ctx context.Context, inviteID string) error {
	invite, err := s.inviteRepo.GetByID(ctx, inviteID)
	if err != nil {
		return err
	}

	if invite.Status != "pending" {
		return errors.New("invitation is no longer pending")
	}

	if time.Now().After(invite.ExpiresAt) {
		return errors.New("invitation has expired")
	}

	return s.inviteRepo.UpdateStatus(ctx, inviteID, "accepted")
}

// DeclineGameInvite declines a game invitation
func (s *UserService) DeclineGameInvite(ctx context.Context, inviteID string) error {
	invite, err := s.inviteRepo.GetByID(ctx, inviteID)
	if err != nil {
		return err
	}

	if invite.Status != "pending" {
		return errors.New("invitation is no longer pending")
	}

	return s.inviteRepo.UpdateStatus(ctx, inviteID, "declined")
}

// GetPendingInvites retrieves all pending invitations for a user
func (s *UserService) GetPendingInvites(ctx context.Context, userID string) ([]*domain.GameInvite, error) {
	return s.inviteRepo.GetPendingInvites(ctx, userID)
}

// Helper functions

func generateSessionToken() (string, error) {
	b := make([]byte, 32)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	return base64.URLEncoding.EncodeToString(b), nil
}
