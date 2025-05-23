package domain

import (
	"context"
	"errors"
	"time"
)

// Common errors
var (
	ErrQuestionNotFound = errors.New("question not found")
)

// QuestionRepository defines the interface for question-related operations
type QuestionRepository interface {
	// GetRandomQuestion retrieves a random question from a category
	GetRandomQuestion(ctx context.Context, category string) (*Question, error)

	// GetCategories retrieves all available categories
	GetCategories(ctx context.Context) ([]string, error)

	// GetByID retrieves a question by its ID
	GetByID(ctx context.Context, id string) (*Question, error)

	// CreateQuestion creates a new question
	CreateQuestion(ctx context.Context, question *Question) error

	// UpdateQuestion updates an existing question
	UpdateQuestion(ctx context.Context, question *Question) error

	// DeleteQuestion deletes a question
	DeleteQuestion(ctx context.Context, id string) error

	// BulkCreateQuestions creates multiple questions in a single transaction
	BulkCreateQuestions(ctx context.Context, questions []*Question) error

	// ValidateQuestion validates a question's data
	ValidateQuestion(ctx context.Context, question *Question) error
}

// Question represents a game question
type Question struct {
	ID            string    `json:"id"`
	Text          string    `json:"text"`
	Answer        string    `json:"answer"`
	Category      string    `json:"category"`
	FillerAnswers []string  `json:"filler_answers"` // Pre-defined plausible but incorrect answers
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
}
