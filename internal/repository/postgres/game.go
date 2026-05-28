package postgres

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/zizouhuweidi/dahaa/internal/domain"
)

// GameRepository implements the domain.GameRepository interface
type GameRepository struct {
	pool *pgxpool.Pool
}

// NewGameRepository creates a new game repository
func NewGameRepository(pool *pgxpool.Pool) *GameRepository {
	return &GameRepository{
		pool: pool,
	}
}

// Create creates a new game
func (r *GameRepository) Create(ctx context.Context, game *domain.Game) error {
	query := `
		INSERT INTO games (code, status, players, rounds, settings, last_activity, host_id, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
		RETURNING id
	`

	players, err := json.Marshal(game.Players)
	if err != nil {
		return fmt.Errorf("failed to marshal players: %w", err)
	}

	rounds, err := json.Marshal(game.Rounds)
	if err != nil {
		return fmt.Errorf("failed to marshal rounds: %w", err)
	}
	settings, err := json.Marshal(game.Settings)
	if err != nil {
		return fmt.Errorf("failed to marshal settings: %w", err)
	}

	now := time.Now().UTC()
	var id string
	err = r.pool.QueryRow(ctx, query,
		game.Code,
		game.Status,
		players,
		rounds,
		settings,
		game.LastActivity,
		game.HostID,
		now,
		now,
	).Scan(&id)

	if err != nil {
		return fmt.Errorf("failed to create game: %w", err)
	}

	game.ID = id
	return nil
}

// GetByCode retrieves a game by its code
func (r *GameRepository) GetByCode(ctx context.Context, code string) (*domain.Game, error) {
	query := `
		SELECT id, code, status, players, rounds, settings, created_at, updated_at, last_activity, host_id
		FROM games
		WHERE code = $1
	`

	var game domain.Game
	var players, rounds, settings []byte
	err := r.pool.QueryRow(ctx, query, code).Scan(
		&game.ID,
		&game.Code,
		&game.Status,
		&players,
		&rounds,
		&settings,
		&game.CreatedAt,
		&game.UpdatedAt,
		&game.LastActivity,
		&game.HostID,
	)

	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, fmt.Errorf("game not found: %s", code)
		}
		return nil, fmt.Errorf("failed to get game: %w", err)
	}

	if err := json.Unmarshal(players, &game.Players); err != nil {
		return nil, fmt.Errorf("failed to unmarshal players: %w", err)
	}

	if err := json.Unmarshal(rounds, &game.Rounds); err != nil {
		return nil, fmt.Errorf("failed to unmarshal rounds: %w", err)
	}
	if len(settings) > 0 {
		if err := json.Unmarshal(settings, &game.Settings); err != nil {
			return nil, fmt.Errorf("failed to unmarshal settings: %w", err)
		}
	}
	if game.Settings == nil || game.Settings.MaxPlayers == 0 {
		game.Settings = domain.DefaultGameSettings()
	}

	return &game, nil
}

// Update updates a game
func (r *GameRepository) Update(ctx context.Context, game *domain.Game) error {
	query := `
		UPDATE games
		SET status = $1, players = $2, rounds = $3, settings = $4, last_activity = $5, host_id = $6, updated_at = $7
		WHERE code = $8
	`

	players, err := json.Marshal(game.Players)
	if err != nil {
		return fmt.Errorf("failed to marshal players: %w", err)
	}

	rounds, err := json.Marshal(game.Rounds)
	if err != nil {
		return fmt.Errorf("failed to marshal rounds: %w", err)
	}
	settings, err := json.Marshal(game.Settings)
	if err != nil {
		return fmt.Errorf("failed to marshal settings: %w", err)
	}

	now := time.Now().UTC()
	_, err = r.pool.Exec(ctx, query,
		game.Status,
		players,
		rounds,
		settings,
		game.LastActivity,
		game.HostID,
		now,
		game.Code,
	)

	if err != nil {
		return fmt.Errorf("failed to update game: %w", err)
	}

	return nil
}

// Delete deletes a game
func (r *GameRepository) Delete(ctx context.Context, code string) error {
	query := `
		DELETE FROM games
		WHERE code = $1
	`

	_, err := r.pool.Exec(ctx, query, code)
	if err != nil {
		return fmt.Errorf("failed to delete game: %w", err)
	}

	return nil
}

// GetByID retrieves a game by its ID
func (r *GameRepository) GetByID(ctx context.Context, id string) (*domain.Game, error) {
	var game domain.Game
	var players, rounds, settings []byte
	err := r.pool.QueryRow(ctx, `
		SELECT id, code, status, players, rounds, settings, created_at, updated_at, last_activity, host_id
		FROM games
		WHERE id = $1
	`, id).Scan(
		&game.ID,
		&game.Code,
		&game.Status,
		&players,
		&rounds,
		&settings,
		&game.CreatedAt,
		&game.UpdatedAt,
		&game.LastActivity,
		&game.HostID,
	)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, domain.ErrGameNotFound
		}
		return nil, fmt.Errorf("failed to get game: %w", err)
	}

	if err := json.Unmarshal(players, &game.Players); err != nil {
		return nil, fmt.Errorf("failed to unmarshal players: %w", err)
	}

	if err := json.Unmarshal(rounds, &game.Rounds); err != nil {
		return nil, fmt.Errorf("failed to unmarshal rounds: %w", err)
	}
	if len(settings) > 0 {
		if err := json.Unmarshal(settings, &game.Settings); err != nil {
			return nil, fmt.Errorf("failed to unmarshal settings: %w", err)
		}
	}
	if game.Settings == nil || game.Settings.MaxPlayers == 0 {
		game.Settings = domain.DefaultGameSettings()
	}

	return &game, nil
}
