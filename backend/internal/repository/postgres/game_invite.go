package postgres

import (
	"context"
	"errors"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/zizouhuweidi/dahaa/internal/domain"
)

// GameInviteRepository implements domain.GameInviteRepository
type GameInviteRepository struct {
	pool *pgxpool.Pool
}

// NewGameInviteRepository creates a new game invite repository
func NewGameInviteRepository(pool *pgxpool.Pool) *GameInviteRepository {
	return &GameInviteRepository{pool: pool}
}

// Create creates a new game invitation
func (r *GameInviteRepository) Create(ctx context.Context, invite *domain.GameInvite) error {
	query := `
		INSERT INTO game_invites (
			id, game_id, from_user, to_user,
			status, created_at, expires_at
		) VALUES ($1, $2, $3, $4, $5, $6, $7)
	`

	_, err := r.pool.Exec(ctx, query,
		invite.ID,
		invite.GameID,
		invite.FromUser,
		invite.ToUser,
		invite.Status,
		invite.CreatedAt,
		invite.ExpiresAt,
	)

	return err
}

// GetByID retrieves a game invitation by ID
func (r *GameInviteRepository) GetByID(ctx context.Context, id string) (*domain.GameInvite, error) {
	query := `
		SELECT id, game_id, from_user, to_user,
			status, created_at, expires_at
		FROM game_invites
		WHERE id = $1
	`

	invite := &domain.GameInvite{}
	err := r.pool.QueryRow(ctx, query, id).Scan(
		&invite.ID,
		&invite.GameID,
		&invite.FromUser,
		&invite.ToUser,
		&invite.Status,
		&invite.CreatedAt,
		&invite.ExpiresAt,
	)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, errors.New("invitation not found")
		}
		return nil, err
	}

	return invite, nil
}

// GetPendingInvites retrieves all pending invitations for a user
func (r *GameInviteRepository) GetPendingInvites(ctx context.Context, userID string) ([]*domain.GameInvite, error) {
	query := `
		SELECT id, game_id, from_user, to_user,
			status, created_at, expires_at
		FROM game_invites
		WHERE to_user = $1
			AND status = 'pending'
			AND expires_at > $2
		ORDER BY created_at DESC
	`

	rows, err := r.pool.Query(ctx, query, userID, time.Now())
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var invites []*domain.GameInvite
	for rows.Next() {
		invite := &domain.GameInvite{}
		err := rows.Scan(
			&invite.ID,
			&invite.GameID,
			&invite.FromUser,
			&invite.ToUser,
			&invite.Status,
			&invite.CreatedAt,
			&invite.ExpiresAt,
		)
		if err != nil {
			return nil, err
		}
		invites = append(invites, invite)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return invites, nil
}

// UpdateStatus updates the status of a game invitation
func (r *GameInviteRepository) UpdateStatus(ctx context.Context, id string, status string) error {
	query := `
		UPDATE game_invites
		SET status = $1
		WHERE id = $2
	`

	_, err := r.pool.Exec(ctx, query, status, id)
	return err
}

// Delete deletes a game invitation
func (r *GameInviteRepository) Delete(ctx context.Context, id string) error {
	query := `DELETE FROM game_invites WHERE id = $1`
	_, err := r.pool.Exec(ctx, query, id)
	return err
}
