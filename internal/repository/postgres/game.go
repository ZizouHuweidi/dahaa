package postgres

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/zizouhuweidi/dahaa/internal/db"
	"github.com/zizouhuweidi/dahaa/internal/db/dbgen"
	"github.com/zizouhuweidi/dahaa/internal/domain"
)

// GameRepository implements the domain.GameRepository interface
type GameRepository struct {
	queries *dbgen.Queries
}

// NewGameRepository creates a new game repository
func NewGameRepository(pool *pgxpool.Pool) *GameRepository {
	return &GameRepository{
		queries: dbgen.New(pool),
	}
}

// Create creates a new game
func (r *GameRepository) Create(ctx context.Context, game *domain.Game) error {
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
	row, err := r.queries.CreateGame(ctx, dbgen.CreateGameParams{
		Code:         game.Code,
		Status:       string(game.Status),
		Players:      players,
		Rounds:       rounds,
		Settings:     settings,
		LastActivity: db.PGTimestamptz(game.LastActivity),
		HostID:       game.HostID,
		CreatedAt:    db.PGTimestamptz(now),
		UpdatedAt:    db.PGTimestamptz(now),
	})

	if err != nil {
		return fmt.Errorf("failed to create game: %w", err)
	}

	game.ID = db.StringFromPGUUID(row.ID)
	game.CreatedAt = db.Time(row.CreatedAt)
	game.UpdatedAt = db.Time(row.UpdatedAt)
	return nil
}

// GetByCode retrieves a game by its code
func (r *GameRepository) GetByCode(ctx context.Context, code string) (*domain.Game, error) {
	row, err := r.queries.GetGameByCode(ctx, code)

	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, fmt.Errorf("game not found: %s", code)
		}
		return nil, fmt.Errorf("failed to get game: %w", err)
	}

	return gameFromParts(
		row.ID,
		row.Code,
		row.Status,
		row.Players,
		row.Rounds,
		row.Settings,
		row.CreatedAt,
		row.UpdatedAt,
		row.LastActivity,
		row.HostID,
	)
}

// Update updates a game
func (r *GameRepository) Update(ctx context.Context, game *domain.Game) error {
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
	row, err := r.queries.UpdateGame(ctx, dbgen.UpdateGameParams{
		Status:       string(game.Status),
		Players:      players,
		Rounds:       rounds,
		Settings:     settings,
		LastActivity: db.PGTimestamptz(game.LastActivity),
		HostID:       game.HostID,
		UpdatedAt:    db.PGTimestamptz(now),
		Code:         game.Code,
	})

	if err != nil {
		return fmt.Errorf("failed to update game: %w", err)
	}

	game.ID = db.StringFromPGUUID(row.ID)
	game.CreatedAt = db.Time(row.CreatedAt)
	game.UpdatedAt = db.Time(row.UpdatedAt)
	return nil
}

// Delete deletes a game
func (r *GameRepository) Delete(ctx context.Context, code string) error {
	err := r.queries.DeleteGameByCode(ctx, code)
	if err != nil {
		return fmt.Errorf("failed to delete game: %w", err)
	}

	return nil
}

// GetByID retrieves a game by its ID
func (r *GameRepository) GetByID(ctx context.Context, id string) (*domain.Game, error) {
	gameID, err := db.PGUUIDFromString(id)
	if err != nil {
		return nil, domain.ErrGameNotFound
	}
	row, err := r.queries.GetGameByID(ctx, gameID)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, domain.ErrGameNotFound
		}
		return nil, fmt.Errorf("failed to get game: %w", err)
	}

	return gameFromParts(
		row.ID,
		row.Code,
		row.Status,
		row.Players,
		row.Rounds,
		row.Settings,
		row.CreatedAt,
		row.UpdatedAt,
		row.LastActivity,
		row.HostID,
	)
}

func gameFromParts(id pgtype.UUID, code string, status string, players []byte, rounds []byte, settings []byte, createdAt pgtype.Timestamptz, updatedAt pgtype.Timestamptz, lastActivity pgtype.Timestamptz, hostID string) (*domain.Game, error) {
	game := domain.Game{
		ID:           db.StringFromPGUUID(id),
		Code:         code,
		Status:       domain.GameStatus(status),
		CreatedAt:    db.Time(createdAt),
		UpdatedAt:    db.Time(updatedAt),
		LastActivity: db.Time(lastActivity),
		HostID:       hostID,
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
