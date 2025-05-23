package domain

import (
	"context"
	"errors"
	"time"
)

// GameSettings defines the configuration for a game
type GameSettings struct {
	Rounds             int          `json:"rounds"`              // Number of rounds in the game
	TimeLimits         TimeLimits   `json:"time_limits"`         // Time limits for different phases
	SelectedCategories []string     `json:"selected_categories"` // Categories to include in the game
	MaxPlayers         int          `json:"max_players"`         // Maximum number of players
	ScoringRules       ScoringRules `json:"scoring_rules"`       // Custom scoring rules
}

// TimeLimits defines the time limits for different game phases
type TimeLimits struct {
	CategorySelection int `json:"category_selection"` // Time for selecting category (seconds)
	AnswerWriting     int `json:"answer_writing"`     // Time for writing answers (seconds)
	Voting            int `json:"voting"`             // Time for voting (seconds)
}

// DefaultGameSettings returns the default game settings
func DefaultGameSettings() *GameSettings {
	return &GameSettings{
		Rounds: 10,
		TimeLimits: TimeLimits{
			CategorySelection: 30,
			AnswerWriting:     60,
			Voting:            30,
		},
		MaxPlayers: 8,
		ScoringRules: ScoringRules{
			CorrectAnswerPoints: 10,
			VotePoints:          5,
			DuplicatePenalty:    2,
		},
	}
}

// Game represents a game session
type Game struct {
	ID           string        `json:"id"`
	Code         string        `json:"code"` // Join code for the game
	Status       GameStatus    `json:"status"`
	Players      []Player      `json:"players"`
	Rounds       []Round       `json:"rounds"`
	Settings     *GameSettings `json:"settings"` // Game settings
	CreatedAt    time.Time     `json:"created_at"`
	UpdatedAt    time.Time     `json:"updated_at"`
	LastActivity time.Time     `json:"last_activity"` // Added for cleanup
}

// GameStatus represents the current status of a game
type GameStatus string

const (
	GameStatusWaiting GameStatus = "waiting" // Waiting for players to join
	GameStatusPlaying GameStatus = "playing" // Game is in progress
	GameStatusEnded   GameStatus = "ended"   // Game has ended
)

// Player represents a player in the game
type Player struct {
	ID          string    `json:"id"`
	Name        string    `json:"name"`
	Score       int       `json:"score"`
	IsConnected bool      `json:"is_connected"`
	LastSeen    time.Time `json:"last_seen"`
	IsActive    bool      `json:"is_active"`
}

// Round represents a single round in the game
type Round struct {
	Number      int         `json:"number"`
	Category    string      `json:"category"`
	Question    string      `json:"question"`
	QuestionID  string      `json:"question_id"`
	Status      RoundStatus `json:"status"`
	StartTime   time.Time   `json:"start_time"`
	EndTime     time.Time   `json:"end_time"`
	CurrentTurn *Turn       `json:"current_turn"`
	AnswerPool  AnswerPool  `json:"answer_pool"`
	Timer       *Timer      `json:"timer,omitempty"`
}

// RoundStatus represents the current status of a round
type RoundStatus string

const (
	RoundStatusWaiting   RoundStatus = "waiting"
	RoundStatusActive    RoundStatus = "active"
	RoundStatusVoting    RoundStatus = "voting"
	RoundStatusCompleted RoundStatus = "completed"
)

// Answer represents a player's answer
type Answer struct {
	ID        string    `json:"id"`
	PlayerID  string    `json:"player_id"`
	Text      string    `json:"text"`
	Votes     []string  `json:"votes"` // List of player IDs who voted for this answer
	CreatedAt time.Time `json:"created_at"`
}

// GameRepository defines the interface for game-related operations
type GameRepository interface {
	// Create creates a new game
	Create(ctx context.Context, game *Game) error

	// GetByCode retrieves a game by its code
	GetByCode(ctx context.Context, code string) (*Game, error)

	// Update updates a game
	Update(ctx context.Context, game *Game) error

	// Delete deletes a game
	Delete(ctx context.Context, code string) error
}

// GameService defines the interface for game-related operations
type GameService interface {
	// Game management
	CreateGame(ctx context.Context, code string, player Player, settings *GameSettings) (*Game, error)
	GetGame(ctx context.Context, code string) (*Game, error)
	JoinGame(ctx context.Context, code string, player Player) error
	StartGame(ctx context.Context, code string) error
	EndGame(ctx context.Context, code string) error

	// Turn management
	StartTurn(ctx context.Context, gameID string, playerID string) error
	SelectCategory(ctx context.Context, gameID string, category string) error
	SubmitAnswer(ctx context.Context, gameID string, playerID string, answer string) error
	SubmitVote(ctx context.Context, gameID string, playerID string, answerID string) error
	EndRound(ctx context.Context, gameID string) error

	// Session management
	HandlePlayerReconnection(ctx context.Context, gameID string, playerID string) error
	CleanupInactiveGames(ctx context.Context) error
}

// Turn represents a player's turn in the game
type Turn struct {
	PlayerID  string     `json:"player_id"`
	StartTime time.Time  `json:"start_time"`
	EndTime   time.Time  `json:"end_time"`
	Status    TurnStatus `json:"status"`
	Category  string     `json:"category,omitempty"`
	Timer     *Timer     `json:"timer,omitempty"` // Added Timer field
}

type TurnStatus string

const (
	TurnStatusWaiting TurnStatus = "waiting"
	TurnStatusActive  TurnStatus = "active"
	TurnStatusEnded   TurnStatus = "ended"
)

// AnswerPool represents the pool of answers for a round
type AnswerPool struct {
	CorrectAnswer string   `json:"correct_answer"`
	FakeAnswers   []Answer `json:"fake_answers"`   // Player-submitted answers
	FillerAnswers []Answer `json:"filler_answers"` // System-generated filler answers
}

// Timer represents a game timer
type Timer struct {
	Type      TimerType `json:"type"`
	StartTime time.Time `json:"start_time"`
	Duration  int       `json:"duration"` // in seconds
	EndTime   time.Time `json:"end_time"`
}

type TimerType string

const (
	TimerTypeCategorySelection TimerType = "category_selection"
	TimerTypeAnswerWriting     TimerType = "answer_writing"
	TimerTypeVoting            TimerType = "voting"
)

// ScoringRules represents the scoring configuration
type ScoringRules struct {
	CorrectAnswerPoints int `json:"correct_answer_points"`
	VotePoints          int `json:"vote_points"`
	DuplicatePenalty    int `json:"duplicate_penalty"`
}

// Common errors
var (
	ErrGameNotFound    = errors.New("game not found")
	ErrGameNotStarted  = errors.New("game has not started")
	ErrGameInProgress  = errors.New("game is already in progress")
	ErrGameEnded       = errors.New("game has ended")
	ErrInvalidRound    = errors.New("invalid round number")
	ErrAnswerSubmitted = errors.New("answer already submitted")
	ErrVoteSubmitted   = errors.New("vote already submitted")
	ErrInvalidVote     = errors.New("invalid vote")
	ErrPlayerNotFound  = errors.New("player not found")
	ErrPlayerNotInGame = errors.New("player not in game")
	ErrInvalidCategory = errors.New("invalid category")
	ErrInvalidQuestion = errors.New("invalid question")
	ErrInvalidAnswer   = errors.New("invalid answer")
	ErrInvalidSettings = errors.New("invalid game settings")
)
