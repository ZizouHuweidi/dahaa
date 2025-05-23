package service

import "errors"

// Common service errors
var (
	ErrUserNotFound       = errors.New("user not found")
	ErrUserAlreadyExists  = errors.New("user already exists")
	ErrInvalidCredentials = errors.New("invalid credentials")
	ErrInviteNotFound     = errors.New("invitation not found")
	ErrInviteExpired      = errors.New("invitation has expired")
)
