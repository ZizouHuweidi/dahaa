package session

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
	"github.com/zizouhuweidi/dahaa/internal/domain"
)

const (
	// Session expiration time
	sessionExpiration = 24 * time.Hour

	// Redis key prefixes
	gameKeyPrefix    = "game:"
	playerKeyPrefix  = "player:"
	connectionPrefix = "conn:"
	rateLimitPrefix  = "ratelimit:"
)

// Manager handles game sessions and player connections
type Manager struct {
	redis *redis.Client
}

// NewManager creates a new session manager
func NewManager(redis *redis.Client) *Manager {
	return &Manager{redis: redis}
}

// StoreGame stores a game in Redis
func (m *Manager) StoreGame(ctx context.Context, game *domain.Game) error {
	key := gameKeyPrefix + game.ID
	data, err := json.Marshal(game)
	if err != nil {
		return fmt.Errorf("failed to marshal game: %w", err)
	}

	return m.redis.Set(ctx, key, data, sessionExpiration).Err()
}

// GetGame retrieves a game from Redis
func (m *Manager) GetGame(ctx context.Context, gameID string) (*domain.Game, error) {
	key := gameKeyPrefix + gameID
	data, err := m.redis.Get(ctx, key).Bytes()
	if err != nil {
		if err == redis.Nil {
			return nil, domain.ErrGameNotFound
		}
		return nil, fmt.Errorf("failed to get game: %w", err)
	}

	var game domain.Game
	if err := json.Unmarshal(data, &game); err != nil {
		return nil, fmt.Errorf("failed to unmarshal game: %w", err)
	}

	return &game, nil
}

// DeleteGame deletes a game from Redis
func (m *Manager) DeleteGame(ctx context.Context, gameID string) error {
	key := gameKeyPrefix + gameID
	if err := m.redis.Del(ctx, key).Err(); err != nil {
		return fmt.Errorf("failed to delete game from Redis: %w", err)
	}

	return nil
}

// StorePlayerSession stores a player's session
func (m *Manager) StorePlayerSession(ctx context.Context, gameID string, player *domain.Player) error {
	key := playerKeyPrefix + gameID + ":" + player.ID
	data, err := json.Marshal(player)
	if err != nil {
		return fmt.Errorf("failed to marshal player: %w", err)
	}

	if err := m.redis.Set(ctx, key, data, 1*time.Hour).Err(); err != nil {
		return fmt.Errorf("failed to store player session: %w", err)
	}

	return nil
}

// GetPlayerSession retrieves a player's session
func (m *Manager) GetPlayerSession(ctx context.Context, gameID string, playerID string) (*domain.Player, error) {
	key := playerKeyPrefix + gameID + ":" + playerID
	data, err := m.redis.Get(ctx, key).Bytes()
	if err != nil {
		if err == redis.Nil {
			return nil, domain.ErrPlayerNotFound
		}
		return nil, fmt.Errorf("failed to get player session: %w", err)
	}

	var player domain.Player
	if err := json.Unmarshal(data, &player); err != nil {
		return nil, fmt.Errorf("failed to unmarshal player: %w", err)
	}

	return &player, nil
}

// DeletePlayerSession deletes a player's session
func (m *Manager) DeletePlayerSession(ctx context.Context, gameID string, playerID string) error {
	key := playerKeyPrefix + gameID + ":" + playerID
	if err := m.redis.Del(ctx, key).Err(); err != nil {
		return fmt.Errorf("failed to delete player session: %w", err)
	}
	return nil
}

// StorePlayerConnection stores a player's WebSocket connection
func (m *Manager) StorePlayerConnection(ctx context.Context, gameID, playerID string) error {
	key := connectionPrefix + gameID + ":" + playerID
	return m.redis.Set(ctx, key, "connected", sessionExpiration).Err()
}

// RemovePlayerConnection removes a player's WebSocket connection
func (m *Manager) RemovePlayerConnection(ctx context.Context, gameID, playerID string) error {
	key := connectionPrefix + gameID + ":" + playerID
	return m.redis.Del(ctx, key).Err()
}

// GetConnectedPlayers returns the number of connected players in a game
func (m *Manager) GetConnectedPlayers(ctx context.Context, gameID string) (int, error) {
	pattern := connectionPrefix + gameID + ":*"
	keys, err := m.redis.Keys(ctx, pattern).Result()
	if err != nil {
		return 0, fmt.Errorf("failed to get connected players: %w", err)
	}
	return len(keys), nil
}

// RateLimit checks if a player has exceeded their rate limit
func (m *Manager) RateLimit(ctx context.Context, playerID string, limit int, window time.Duration) (bool, error) {
	key := rateLimitPrefix + playerID
	count, err := m.redis.Incr(ctx, key).Result()
	if err != nil {
		return false, fmt.Errorf("failed to increment rate limit: %w", err)
	}

	if count == 1 {
		m.redis.Expire(ctx, key, window)
	}

	return count > int64(limit), nil
}

// PublishGameEvent publishes a game event to all connected players
func (m *Manager) PublishGameEvent(ctx context.Context, gameID string, eventType string, payload interface{}) error {
	data, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("failed to marshal event payload: %w", err)
	}

	event := map[string]interface{}{
		"type":    eventType,
		"game_id": gameID,
		"payload": data,
	}

	eventData, err := json.Marshal(event)
	if err != nil {
		return fmt.Errorf("failed to marshal event: %w", err)
	}

	return m.redis.Publish(ctx, "game:"+gameID, eventData).Err()
}

// SubscribeToGame subscribes to game events
func (m *Manager) SubscribeToGame(ctx context.Context, gameID string) *redis.PubSub {
	return m.redis.Subscribe(ctx, "game:"+gameID)
}

// CleanupInactiveGames removes games that haven't been active for a while
func (m *Manager) CleanupInactiveGames(ctx context.Context) error {
	pattern := gameKeyPrefix + "*"
	iter := m.redis.Scan(ctx, 0, pattern, 0).Iterator()
	for iter.Next(ctx) {
		key := iter.Val()
		data, err := m.redis.Get(ctx, key).Bytes()
		if err != nil {
			continue
		}

		var game domain.Game
		if err := json.Unmarshal(data, &game); err != nil {
			continue
		}

		if time.Since(game.LastActivity) > 24*time.Hour {
			if err := m.DeleteGame(ctx, game.ID); err != nil {
				// Log error but continue with other games
				fmt.Printf("Failed to delete inactive game %s: %v\n", game.ID, err)
			}
		}
	}

	if err := iter.Err(); err != nil {
		return fmt.Errorf("failed to scan games: %w", err)
	}

	return nil
}

// StartCleanupJob starts a background job to clean up inactive games
func (m *Manager) StartCleanupJob(ctx context.Context) {
	ticker := time.NewTicker(1 * time.Hour)
	go func() {
		for {
			select {
			case <-ticker.C:
				if err := m.CleanupInactiveGames(ctx); err != nil {
					fmt.Printf("Failed to cleanup inactive games: %v\n", err)
				}
			case <-ctx.Done():
				ticker.Stop()
				return
			}
		}
	}()
}

// GetAllGames retrieves all active game sessions
func (m *Manager) GetAllGames(ctx context.Context) ([]*domain.Game, error) {
	pattern := gameKeyPrefix + "*"
	iter := m.redis.Scan(ctx, 0, pattern, 0).Iterator()
	var games []*domain.Game
	for iter.Next(ctx) {
		key := iter.Val()
		data, err := m.redis.Get(ctx, key).Bytes()
		if err != nil {
			continue
		}

		var game domain.Game
		if err := json.Unmarshal(data, &game); err != nil {
			continue
		}

		games = append(games, &game)
	}

	if err := iter.Err(); err != nil {
		return nil, fmt.Errorf("failed to scan games: %w", err)
	}

	return games, nil
}
