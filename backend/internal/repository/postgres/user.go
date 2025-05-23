package postgres

import (
	"context"
	"errors"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/zizouhuweidi/dahaa/internal/domain"
)

// UserRepository implements domain.UserRepository
type UserRepository struct {
	pool *pgxpool.Pool
}

// NewUserRepository creates a new user repository
func NewUserRepository(pool *pgxpool.Pool) *UserRepository {
	return &UserRepository{pool: pool}
}

// Create creates a new user
func (r *UserRepository) Create(ctx context.Context, user *domain.User) error {
	query := `
		INSERT INTO users (
			id, username, email, password_hash, display_name,
			games_played, games_won, total_points,
			last_login_at, created_at, updated_at
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)
	`

	_, err := r.pool.Exec(ctx, query,
		user.ID,
		user.Username,
		user.Email,
		user.PasswordHash,
		user.DisplayName,
		user.Stats.GamesPlayed,
		user.Stats.GamesWon,
		user.Stats.TotalPoints,
		user.LastLoginAt,
		user.CreatedAt,
		user.UpdatedAt,
	)

	return err
}

// GetByID retrieves a user by ID
func (r *UserRepository) GetByID(ctx context.Context, id string) (*domain.User, error) {
	query := `
		SELECT id, username, email, password_hash, display_name,
			games_played, games_won, total_points,
			last_login_at, created_at, updated_at
		FROM users
		WHERE id = $1
	`

	user := &domain.User{}
	err := r.pool.QueryRow(ctx, query, id).Scan(
		&user.ID,
		&user.Username,
		&user.Email,
		&user.PasswordHash,
		&user.DisplayName,
		&user.Stats.GamesPlayed,
		&user.Stats.GamesWon,
		&user.Stats.TotalPoints,
		&user.LastLoginAt,
		&user.CreatedAt,
		&user.UpdatedAt,
	)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, domain.ErrUserNotFound
		}
		return nil, err
	}

	return user, nil
}

// GetByUsername retrieves a user by username
func (r *UserRepository) GetByUsername(ctx context.Context, username string) (*domain.User, error) {
	query := `
		SELECT id, username, email, password_hash, display_name,
			games_played, games_won, total_points,
			last_login_at, created_at, updated_at
		FROM users
		WHERE username = $1
	`

	user := &domain.User{}
	err := r.pool.QueryRow(ctx, query, username).Scan(
		&user.ID,
		&user.Username,
		&user.Email,
		&user.PasswordHash,
		&user.DisplayName,
		&user.Stats.GamesPlayed,
		&user.Stats.GamesWon,
		&user.Stats.TotalPoints,
		&user.LastLoginAt,
		&user.CreatedAt,
		&user.UpdatedAt,
	)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, domain.ErrUserNotFound
		}
		return nil, err
	}

	return user, nil
}

// GetByEmail retrieves a user by email
func (r *UserRepository) GetByEmail(ctx context.Context, email string) (*domain.User, error) {
	query := `
		SELECT id, username, email, password_hash, display_name,
			games_played, games_won, total_points,
			last_login_at, created_at, updated_at
		FROM users
		WHERE email = $1
	`

	user := &domain.User{}
	err := r.pool.QueryRow(ctx, query, email).Scan(
		&user.ID,
		&user.Username,
		&user.Email,
		&user.PasswordHash,
		&user.DisplayName,
		&user.Stats.GamesPlayed,
		&user.Stats.GamesWon,
		&user.Stats.TotalPoints,
		&user.LastLoginAt,
		&user.CreatedAt,
		&user.UpdatedAt,
	)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, domain.ErrUserNotFound
		}
		return nil, err
	}

	return user, nil
}

// Update updates a user
func (r *UserRepository) Update(ctx context.Context, user *domain.User) error {
	query := `
		UPDATE users
		SET username = $1,
			email = $2,
			password_hash = $3,
			display_name = $4,
			games_played = $5,
			games_won = $6,
			total_points = $7,
			last_login_at = $8,
			updated_at = $9
		WHERE id = $10
	`

	_, err := r.pool.Exec(ctx, query,
		user.Username,
		user.Email,
		user.PasswordHash,
		user.DisplayName,
		user.Stats.GamesPlayed,
		user.Stats.GamesWon,
		user.Stats.TotalPoints,
		user.LastLoginAt,
		time.Now(),
		user.ID,
	)

	return err
}

// Delete deletes a user
func (r *UserRepository) Delete(ctx context.Context, id string) error {
	query := `DELETE FROM users WHERE id = $1`
	_, err := r.pool.Exec(ctx, query, id)
	return err
}

// UpdateStats updates a user's game statistics
func (r *UserRepository) UpdateStats(ctx context.Context, id string, stats domain.UserStats) error {
	query := `
		UPDATE users
		SET games_played = $1,
			games_won = $2,
			total_points = $3,
			updated_at = $4
		WHERE id = $5
	`

	_, err := r.pool.Exec(ctx, query,
		stats.GamesPlayed,
		stats.GamesWon,
		stats.TotalPoints,
		time.Now(),
		id,
	)

	return err
}
