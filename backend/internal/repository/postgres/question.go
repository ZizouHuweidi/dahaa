package postgres

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/zizouhuweidi/dahaa/internal/domain"
)

// QuestionRepository implements the domain.QuestionRepository interface
type QuestionRepository struct {
	pool *pgxpool.Pool
}

// NewQuestionRepository creates a new question repository
func NewQuestionRepository(pool *pgxpool.Pool) *QuestionRepository {
	return &QuestionRepository{
		pool: pool,
	}
}

// GetRandomQuestion retrieves a random question from a category
func (r *QuestionRepository) GetRandomQuestion(ctx context.Context, category string) (*domain.Question, error) {
	var question domain.Question
	err := r.pool.QueryRow(ctx, `
		SELECT id, text, answer, category, created_at, updated_at
		FROM questions
		WHERE category = $1
		ORDER BY RANDOM()
		LIMIT 1
	`, category).Scan(
		&question.ID,
		&question.Text,
		&question.Answer,
		&question.Category,
		&question.CreatedAt,
		&question.UpdatedAt,
	)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, fmt.Errorf("no questions found for category: %s", category)
		}
		return nil, fmt.Errorf("failed to get random question: %w", err)
	}
	return &question, nil
}

// GetCategories retrieves all available categories
func (r *QuestionRepository) GetCategories(ctx context.Context) ([]string, error) {
	query := `
		SELECT DISTINCT category
		FROM questions
		ORDER BY category
	`

	rows, err := r.pool.Query(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("failed to get categories: %w", err)
	}
	defer rows.Close()

	var categories []string
	for rows.Next() {
		var category string
		if err := rows.Scan(&category); err != nil {
			return nil, fmt.Errorf("failed to scan category: %w", err)
		}
		categories = append(categories, category)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating categories: %w", err)
	}

	return categories, nil
}

// GetDifficulties retrieves all available difficulty levels
func (r *QuestionRepository) GetDifficulties(ctx context.Context) ([]string, error) {
	query := `
		SELECT DISTINCT difficulty
		FROM questions
		ORDER BY difficulty
	`

	rows, err := r.pool.Query(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("failed to get difficulties: %w", err)
	}
	defer rows.Close()

	var difficulties []string
	for rows.Next() {
		var difficulty string
		if err := rows.Scan(&difficulty); err != nil {
			return nil, fmt.Errorf("failed to scan difficulty: %w", err)
		}
		difficulties = append(difficulties, difficulty)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating difficulties: %w", err)
	}

	return difficulties, nil
}

// GetByID retrieves a question by its ID
func (r *QuestionRepository) GetByID(ctx context.Context, id string) (*domain.Question, error) {
	var question domain.Question
	var fillerAnswers []string
	err := r.pool.QueryRow(ctx, `
		SELECT id, text, answer, category, filler_answers, created_at, updated_at
		FROM questions
		WHERE id = $1
	`, id).Scan(
		&question.ID,
		&question.Text,
		&question.Answer,
		&question.Category,
		&fillerAnswers,
		&question.CreatedAt,
		&question.UpdatedAt,
	)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, domain.ErrQuestionNotFound
		}
		return nil, fmt.Errorf("failed to get question: %w", err)
	}
	question.FillerAnswers = fillerAnswers
	return &question, nil
}

// CreateQuestion creates a new question
func (r *QuestionRepository) CreateQuestion(ctx context.Context, question *domain.Question) error {
	query := `
		INSERT INTO questions (text, answer, category, filler_answers)
		VALUES ($1, $2, $3, $4)
		RETURNING id, created_at, updated_at
	`
	return r.pool.QueryRow(ctx, query,
		question.Text,
		question.Answer,
		question.Category,
		question.FillerAnswers,
	).Scan(&question.ID, &question.CreatedAt, &question.UpdatedAt)
}

// UpdateQuestion updates an existing question
func (r *QuestionRepository) UpdateQuestion(ctx context.Context, question *domain.Question) error {
	query := `
		UPDATE questions
		SET text = $1, answer = $2, category = $3, filler_answers = $4, updated_at = CURRENT_TIMESTAMP
		WHERE id = $5
		RETURNING updated_at
	`
	return r.pool.QueryRow(ctx, query,
		question.Text,
		question.Answer,
		question.Category,
		question.FillerAnswers,
		question.ID,
	).Scan(&question.UpdatedAt)
}

// DeleteQuestion deletes a question
func (r *QuestionRepository) DeleteQuestion(ctx context.Context, id string) error {
	query := `DELETE FROM questions WHERE id = $1`
	result, err := r.pool.Exec(ctx, query, id)
	if err != nil {
		return fmt.Errorf("failed to delete question: %w", err)
	}
	if result.RowsAffected() == 0 {
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

	query := `
		INSERT INTO questions (text, answer, category, filler_answers)
		VALUES ($1, $2, $3, $4)
		RETURNING id, created_at, updated_at
	`

	for _, question := range questions {
		err := tx.QueryRow(ctx, query,
			question.Text,
			question.Answer,
			question.Category,
			question.FillerAnswers,
		).Scan(&question.ID, &question.CreatedAt, &question.UpdatedAt)
		if err != nil {
			return fmt.Errorf("failed to create question: %w", err)
		}
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
