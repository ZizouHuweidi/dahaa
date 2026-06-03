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

// GameInviteRepository implements domain.GameInviteRepository
type GameInviteRepository struct {
	queries *dbgen.Queries
}

// NewGameInviteRepository creates a new game invite repository
func NewGameInviteRepository(pool *pgxpool.Pool) *GameInviteRepository {
	return &GameInviteRepository{queries: dbgen.New(pool)}
}

// Create creates a new game invitation
func (r *GameInviteRepository) Create(ctx context.Context, invite *domain.GameInvite) error {
	row, err := r.queries.CreateGameInvite(ctx, dbgen.CreateGameInviteParams{
		ID:        invite.ID,
		GameID:    invite.GameID,
		FromUser:  invite.FromUser,
		ToUser:    invite.ToUser,
		Status:    invite.Status,
		CreatedAt: db.PGTimestamptz(invite.CreatedAt),
		ExpiresAt: db.PGTimestamptz(invite.ExpiresAt),
	})
	if err != nil {
		return err
	}
	applyDBGameInvite(invite, row)
	return err
}

// GetByID retrieves a game invitation by ID
func (r *GameInviteRepository) GetByID(ctx context.Context, id string) (*domain.GameInvite, error) {
	row, err := r.queries.GetGameInviteByID(ctx, id)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, errors.New("invitation not found")
		}
		return nil, err
	}

	return dbGameInviteToDomain(row), nil
}

// GetPendingInvites retrieves all pending invitations for a user
func (r *GameInviteRepository) GetPendingInvites(ctx context.Context, userID string) ([]*domain.GameInvite, error) {
	rows, err := r.queries.GetPendingGameInvites(ctx, dbgen.GetPendingGameInvitesParams{
		ToUser:    userID,
		ExpiresAt: db.PGTimestamptz(time.Now()),
	})
	if err != nil {
		return nil, err
	}

	var invites []*domain.GameInvite
	for _, row := range rows {
		invites = append(invites, dbGameInviteToDomain(row))
	}

	return invites, nil
}

// UpdateStatus updates the status of a game invitation
func (r *GameInviteRepository) UpdateStatus(ctx context.Context, id string, status string) error {
	_, err := r.queries.UpdateGameInviteStatus(ctx, dbgen.UpdateGameInviteStatusParams{Status: status, ID: id})
	return err
}

// Delete deletes a game invitation
func (r *GameInviteRepository) Delete(ctx context.Context, id string) error {
	return r.queries.DeleteGameInvite(ctx, id)
}

func dbGameInviteToDomain(row dbgen.GameInvite) *domain.GameInvite {
	invite := &domain.GameInvite{}
	applyDBGameInvite(invite, row)
	return invite
}

func applyDBGameInvite(invite *domain.GameInvite, row dbgen.GameInvite) {
	invite.ID = row.ID
	invite.GameID = row.GameID
	invite.FromUser = row.FromUser
	invite.ToUser = row.ToUser
	invite.Status = row.Status
	invite.CreatedAt = db.Time(row.CreatedAt)
	invite.ExpiresAt = db.Time(row.ExpiresAt)
}
