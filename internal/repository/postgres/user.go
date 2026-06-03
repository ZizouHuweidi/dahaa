package postgres

import (
	"context"
	"errors"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/zizouhuweidi/dahaa/internal/db"
	"github.com/zizouhuweidi/dahaa/internal/db/dbgen"
	"github.com/zizouhuweidi/dahaa/internal/domain"
)

// UserRepository implements domain.UserRepository
type UserRepository struct {
	queries *dbgen.Queries
}

// NewUserRepository creates a new user repository
func NewUserRepository(pool *pgxpool.Pool) *UserRepository {
	return &UserRepository{queries: dbgen.New(pool)}
}

// Create creates a new user
func (r *UserRepository) Create(ctx context.Context, user *domain.User) error {
	row, err := r.queries.CreateUser(ctx, dbgen.CreateUserParams{
		ID:           user.ID,
		Username:     user.Username,
		Email:        user.Email,
		PasswordHash: user.PasswordHash,
		DisplayName:  user.DisplayName,
		GamesPlayed:  int32(user.Stats.GamesPlayed),
		GamesWon:     int32(user.Stats.GamesWon),
		TotalPoints:  int32(user.Stats.TotalPoints),
		LastLoginAt:  db.PGTimestamptz(user.LastLoginAt),
		CreatedAt:    db.PGTimestamptz(user.CreatedAt),
		UpdatedAt:    db.PGTimestamptz(user.UpdatedAt),
	})
	if err != nil {
		return err
	}
	applyDBUser(user, row)
	return err
}

// GetByID retrieves a user by ID
func (r *UserRepository) GetByID(ctx context.Context, id string) (*domain.User, error) {
	row, err := r.queries.GetUserByID(ctx, id)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, domain.ErrUserNotFound
		}
		return nil, err
	}

	return dbUserToDomain(row), nil
}

// GetByUsername retrieves a user by username
func (r *UserRepository) GetByUsername(ctx context.Context, username string) (*domain.User, error) {
	row, err := r.queries.GetUserByUsername(ctx, username)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, domain.ErrUserNotFound
		}
		return nil, err
	}

	return dbUserToDomain(row), nil
}

// GetByEmail retrieves a user by email
func (r *UserRepository) GetByEmail(ctx context.Context, email string) (*domain.User, error) {
	row, err := r.queries.GetUserByEmail(ctx, email)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, domain.ErrUserNotFound
		}
		return nil, err
	}

	return dbUserToDomain(row), nil
}

// Update updates a user
func (r *UserRepository) Update(ctx context.Context, user *domain.User) error {
	row, err := r.queries.UpdateUser(ctx, dbgen.UpdateUserParams{
		Username:     user.Username,
		Email:        user.Email,
		PasswordHash: user.PasswordHash,
		DisplayName:  user.DisplayName,
		GamesPlayed:  int32(user.Stats.GamesPlayed),
		GamesWon:     int32(user.Stats.GamesWon),
		TotalPoints:  int32(user.Stats.TotalPoints),
		LastLoginAt:  db.PGTimestamptz(user.LastLoginAt),
		UpdatedAt:    db.PGTimestamptz(time.Now()),
		ID:           user.ID,
	})
	if err != nil {
		return err
	}
	applyDBUser(user, row)
	return err
}

// Delete deletes a user
func (r *UserRepository) Delete(ctx context.Context, id string) error {
	return r.queries.DeleteUser(ctx, id)
}

// UpdateStats updates a user's game statistics
func (r *UserRepository) UpdateStats(ctx context.Context, id string, stats domain.UserStats) error {
	_, err := r.queries.UpdateUserStats(ctx, dbgen.UpdateUserStatsParams{
		GamesPlayed: int32(stats.GamesPlayed),
		GamesWon:    int32(stats.GamesWon),
		TotalPoints: int32(stats.TotalPoints),
		UpdatedAt:   db.PGTimestamptz(time.Now()),
		ID:          id,
	})
	return err
}

func dbUserToDomain(row dbgen.User) *domain.User {
	user := &domain.User{}
	applyDBUser(user, row)
	return user
}

func applyDBUser(user *domain.User, row dbgen.User) {
	user.ID = row.ID
	user.Username = row.Username
	user.Email = row.Email
	user.PasswordHash = row.PasswordHash
	user.DisplayName = row.DisplayName
	user.Stats.GamesPlayed = int(row.GamesPlayed)
	user.Stats.GamesWon = int(row.GamesWon)
	user.Stats.TotalPoints = int(row.TotalPoints)
	user.LastLoginAt = db.Time(row.LastLoginAt)
	user.CreatedAt = db.Time(row.CreatedAt)
	user.UpdatedAt = db.Time(row.UpdatedAt)
}
