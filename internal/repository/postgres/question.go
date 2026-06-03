package postgres

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/zizouhuweidi/dahaa/internal/db"
	"github.com/zizouhuweidi/dahaa/internal/db/dbgen"
	"github.com/zizouhuweidi/dahaa/internal/domain"
)

// QuestionRepository implements the domain.QuestionRepository interface
type QuestionRepository struct {
	pool    *pgxpool.Pool
	queries *dbgen.Queries
}

// NewQuestionRepository creates a new question repository
func NewQuestionRepository(pool *pgxpool.Pool) *QuestionRepository {
	return &QuestionRepository{
		pool:    pool,
		queries: dbgen.New(pool),
	}
}

// GetRandomQuestion retrieves a random question from a category
func (r *QuestionRepository) GetRandomQuestion(ctx context.Context, category string) (*domain.Question, error) {
	row, err := r.queries.GetRandomQuestion(ctx, category)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, fmt.Errorf("no questions found for category: %s", category)
		}
		return nil, fmt.Errorf("failed to get random question: %w", err)
	}
	question := domain.Question{
		ID:            db.StringFromPGUUID(row.ID),
		Text:          row.Text,
		Answer:        row.Answer,
		Category:      row.Category,
		FillerAnswers: row.FillerAnswers,
		CreatedAt:     db.Time(row.CreatedAt),
		UpdatedAt:     db.Time(row.UpdatedAt),
	}
	return &question, nil
}

// GetCategories retrieves all available categories
func (r *QuestionRepository) GetCategories(ctx context.Context) ([]string, error) {
	categories, err := r.queries.GetQuestionCategories(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get categories: %w", err)
	}
	return categories, nil
}

// GetByID retrieves a question by its ID
func (r *QuestionRepository) GetByID(ctx context.Context, id string) (*domain.Question, error) {
	questionID, err := db.PGUUIDFromString(id)
	if err != nil {
		return nil, domain.ErrQuestionNotFound
	}
	row, err := r.queries.GetQuestionByID(ctx, questionID)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, domain.ErrQuestionNotFound
		}
		return nil, fmt.Errorf("failed to get question: %w", err)
	}
	question := domain.Question{
		ID:            db.StringFromPGUUID(row.ID),
		Text:          row.Text,
		Answer:        row.Answer,
		Category:      row.Category,
		FillerAnswers: row.FillerAnswers,
		CreatedAt:     db.Time(row.CreatedAt),
		UpdatedAt:     db.Time(row.UpdatedAt),
	}
	return &question, nil
}

// CreateQuestion creates a new question
func (r *QuestionRepository) CreateQuestion(ctx context.Context, question *domain.Question) error {
	row, err := r.queries.CreateQuestion(ctx, dbgen.CreateQuestionParams{
		Text:          question.Text,
		Answer:        question.Answer,
		Category:      question.Category,
		FillerAnswers: question.FillerAnswers,
	})
	if err != nil {
		return err
	}
	question.ID = db.StringFromPGUUID(row.ID)
	question.CreatedAt = db.Time(row.CreatedAt)
	question.UpdatedAt = db.Time(row.UpdatedAt)
	return nil
}

// UpdateQuestion updates an existing question
func (r *QuestionRepository) UpdateQuestion(ctx context.Context, question *domain.Question) error {
	questionID, err := db.PGUUIDFromString(question.ID)
	if err != nil {
		return domain.ErrQuestionNotFound
	}
	row, err := r.queries.UpdateQuestion(ctx, dbgen.UpdateQuestionParams{
		Text:          question.Text,
		Answer:        question.Answer,
		Category:      question.Category,
		FillerAnswers: question.FillerAnswers,
		ID:            questionID,
	})
	if err != nil {
		if err == pgx.ErrNoRows {
			return domain.ErrQuestionNotFound
		}
		return err
	}
	question.UpdatedAt = db.Time(row.UpdatedAt)
	return nil
}

// DeleteQuestion deletes a question
func (r *QuestionRepository) DeleteQuestion(ctx context.Context, id string) error {
	questionID, err := db.PGUUIDFromString(id)
	if err != nil {
		return domain.ErrQuestionNotFound
	}
	rowsAffected, err := r.queries.DeleteQuestion(ctx, questionID)
	if err != nil {
		return fmt.Errorf("failed to delete question: %w", err)
	}
	if rowsAffected == 0 {
		return domain.ErrQuestionNotFound
	}
	return nil
}

// BulkCreateQuestions creates multiple questions in a single transaction
func (r *QuestionRepository) BulkCreateQuestions(ctx context.Context, questions []*domain.Question) error {
	tx, err := r.pool.Begin(ctx)
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback(ctx)

	queries := r.queries.WithTx(tx)

	for _, question := range questions {
		row, err := queries.CreateQuestion(ctx, dbgen.CreateQuestionParams{
			Text:          question.Text,
			Answer:        question.Answer,
			Category:      question.Category,
			FillerAnswers: question.FillerAnswers,
		})
		if err != nil {
			return fmt.Errorf("failed to create question: %w", err)
		}
		question.ID = db.StringFromPGUUID(row.ID)
		question.CreatedAt = db.Time(row.CreatedAt)
		question.UpdatedAt = db.Time(row.UpdatedAt)
	}

	if err := tx.Commit(ctx); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	return nil
}

// ValidateQuestion validates a question's data
func (r *QuestionRepository) ValidateQuestion(ctx context.Context, question *domain.Question) error {
	if question.Text == "" {
		return fmt.Errorf("question text cannot be empty")
	}
	if question.Answer == "" {
		return fmt.Errorf("question answer cannot be empty")
	}
	if question.Category == "" {
		return fmt.Errorf("question category cannot be empty")
	}
	if len(question.Category) < 3 {
		return fmt.Errorf("category must be at least 3 characters long")
	}
	if len(question.Category) > 50 {
		return fmt.Errorf("category cannot be longer than 50 characters")
	}
	return nil
}
