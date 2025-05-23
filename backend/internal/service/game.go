package service

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"math/rand"
	"strings"
	"time"

	"github.com/zizouhuweidi/dahaa/internal/domain"
	"github.com/zizouhuweidi/dahaa/internal/repository/postgres"
	"github.com/zizouhuweidi/dahaa/internal/session"
	"github.com/zizouhuweidi/dahaa/internal/validation"
	"github.com/zizouhuweidi/dahaa/internal/websocket"
)

var (
	ErrGameNotFound    = errors.New("game not found")
	ErrGameFull        = errors.New("game is full")
	ErrGameInProgress  = errors.New("game is already in progress")
	ErrGameNotStarted  = errors.New("game has not started")
	ErrInvalidRound    = errors.New("invalid round number")
	ErrPlayerNotFound  = errors.New("player not found")
	ErrAnswerSubmitted = errors.New("answer already submitted")
	ErrVoteSubmitted   = errors.New("vote already submitted")
	ErrInvalidAnswer   = errors.New("invalid answer")
	ErrInvalidVote     = errors.New("invalid vote")
)

// GameService implements the domain.GameService interface
type GameService struct {
	gameRepo     *postgres.GameRepository
	questionRepo *postgres.QuestionRepository
	hub          *websocket.Hub
	sessionMgr   *session.Manager
}

// NewGameService creates a new game service
func NewGameService(gameRepo *postgres.GameRepository, questionRepo *postgres.QuestionRepository, hub *websocket.Hub, sessionMgr *session.Manager) *GameService {
	return &GameService{
		gameRepo:     gameRepo,
		questionRepo: questionRepo,
		hub:          hub,
		sessionMgr:   sessionMgr,
	}
}

// CreateGame creates a new game session
func (s *GameService) CreateGame(ctx context.Context, code string, player domain.Player, settings *domain.GameSettings) (*domain.Game, error) {
	// Check if game with code already exists
	existingGame, err := s.gameRepo.GetByCode(ctx, code)
	if err == nil && existingGame != nil {
		return nil, errors.New("game code already exists")
	}

	// Use default settings if none provided
	if settings == nil {
		settings = domain.DefaultGameSettings()
	}

	// Validate selected categories
	if len(settings.SelectedCategories) == 0 {
		return nil, errors.New("at least one category must be selected")
	}

	// Verify all selected categories exist
	availableCategories, err := s.questionRepo.GetCategories(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get categories: %w", err)
	}

	categoryMap := make(map[string]bool)
	for _, cat := range availableCategories {
		categoryMap[cat] = true
	}

	for _, cat := range settings.SelectedCategories {
		if !categoryMap[cat] {
			return nil, fmt.Errorf("invalid category: %s", cat)
		}
	}

	// Create new game
	game := &domain.Game{
		ID:           generateID(),
		Code:         code,
		Status:       domain.GameStatusWaiting,
		Players:      []domain.Player{player},
		Rounds:       []domain.Round{},
		Settings:     settings,
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
		LastActivity: time.Now(),
	}

	// Save to database
	if err := s.gameRepo.Create(ctx, game); err != nil {
		return nil, err
	}

	// Store in Redis for active game management
	if err := s.sessionMgr.StoreGame(ctx, game); err != nil {
		return nil, err
	}

	// Notify all clients about new game
	payload, err := json.Marshal(game)
	if err != nil {
		return nil, err
	}

	s.hub.BroadcastToGame(game.ID, "game_created", payload)

	return game, nil
}

// JoinGame allows a player to join an existing game
func (s *GameService) JoinGame(ctx context.Context, code string, player domain.Player) error {
	game, err := s.GetGame(ctx, code)
	if err != nil {
		return err
	}

	if game.Status != domain.GameStatusWaiting {
		return errors.New("game has already started")
	}

	if len(game.Players) >= game.Settings.MaxPlayers {
		return errors.New("game is full")
	}

	// Check if player already exists
	for _, p := range game.Players {
		if p.ID == player.ID {
			return errors.New("player already in game")
		}
	}

	game.Players = append(game.Players, player)
	game.UpdatedAt = time.Now()
	game.LastActivity = time.Now()

	if err := s.UpdateGame(ctx, game); err != nil {
		return err
	}

	// Notify all clients about player joining
	payload, err := json.Marshal(player)
	if err != nil {
		return err
	}

	s.hub.BroadcastToGame(game.ID, "player_joined", payload)

	return nil
}

// GetGame retrieves a game by its code
func (s *GameService) GetGame(ctx context.Context, gameID string) (*domain.Game, error) {
	// Try to get game from Redis first
	game, err := s.sessionMgr.GetGame(ctx, gameID)
	if err != nil {
		// If not in Redis, try database
		game, err = s.gameRepo.GetByID(ctx, gameID)
		if err != nil {
			return nil, err
		}
		// Store in Redis for future access
		if err := s.sessionMgr.StoreGame(ctx, game); err != nil {
			// Log error but continue
			fmt.Printf("Failed to store game in Redis: %v\n", err)
		}
	}
	return game, nil
}

// UpdateGame updates a game's state
func (s *GameService) UpdateGame(ctx context.Context, game *domain.Game) error {
	// Update in database
	if err := s.gameRepo.Update(ctx, game); err != nil {
		return err
	}

	// Update in Redis
	if err := s.sessionMgr.StoreGame(ctx, game); err != nil {
		// Log error but continue
		fmt.Printf("Failed to update game in Redis: %v\n", err)
	}

	// Notify all clients about game update
	payload, err := json.Marshal(game)
	if err != nil {
		return err
	}

	s.hub.BroadcastToGame(game.ID, "game_updated", payload)

	return nil
}

// DeleteGame deletes a game
func (s *GameService) DeleteGame(ctx context.Context, code string) error {
	// Get game first to get its ID
	game, err := s.GetGame(ctx, code)
	if err != nil {
		return err
	}

	// Delete from database
	if err := s.gameRepo.Delete(ctx, code); err != nil {
		return err
	}

	// Delete from Redis
	if err := s.sessionMgr.DeleteGame(ctx, code); err != nil {
		return err
	}

	// Notify all clients about game deletion
	payload, err := json.Marshal(map[string]string{
		"code": code,
	})
	if err != nil {
		return err
	}

	s.hub.BroadcastToGame(game.ID, "game_deleted", payload)

	return nil
}

// StartGame starts a game session
func (s *GameService) StartGame(ctx context.Context, code string) error {
	game, err := s.GetGame(ctx, code)
	if err != nil {
		return err
	}

	if game.Status != domain.GameStatusWaiting {
		return errors.New("game has already started")
	}

	if len(game.Players) < 2 {
		return errors.New("need at least 2 players to start")
	}

	game.Status = domain.GameStatusPlaying
	game.UpdatedAt = time.Now()
	game.LastActivity = time.Now()

	// Create first round
	round := domain.Round{
		Number:    1,
		Status:    domain.RoundStatusWaiting,
		StartTime: time.Now(),
		AnswerPool: domain.AnswerPool{
			CorrectAnswer: "",
			FakeAnswers:   make([]domain.Answer, 0),
			FillerAnswers: make([]domain.Answer, 0),
		},
	}

	game.Rounds = append(game.Rounds, round)

	if err := s.UpdateGame(ctx, game); err != nil {
		return err
	}

	// Notify all clients about game start
	payload, err := json.Marshal(game)
	if err != nil {
		return err
	}

	s.hub.BroadcastToGame(game.ID, "game_started", payload)

	return nil
}

// StartTurn starts a new turn for a player
func (s *GameService) StartTurn(ctx context.Context, gameID string, playerID string) error {
	game, err := s.GetGame(ctx, gameID)
	if err != nil {
		return err
	}

	if game.Status != domain.GameStatusPlaying {
		return domain.ErrGameNotStarted
	}

	currentRound := &game.Rounds[len(game.Rounds)-1]
	if currentRound.Status != domain.RoundStatusWaiting {
		return errors.New("round is not in waiting state")
	}

	// Create new turn
	turn := &domain.Turn{
		PlayerID:  playerID,
		StartTime: time.Now(),
		Status:    domain.TurnStatusActive,
	}

	// Set timer for category selection using game settings
	turn.Timer = &domain.Timer{
		Type:      domain.TimerTypeCategorySelection,
		StartTime: time.Now(),
		Duration:  game.Settings.TimeLimits.CategorySelection,
		EndTime:   time.Now().Add(time.Duration(game.Settings.TimeLimits.CategorySelection) * time.Second),
	}

	currentRound.CurrentTurn = turn
	return s.UpdateGame(ctx, game)
}

// SelectCategory handles category selection during a turn
func (s *GameService) SelectCategory(ctx context.Context, gameID string, category string) error {
	game, err := s.GetGame(ctx, gameID)
	if err != nil {
		return err
	}

	currentRound := &game.Rounds[len(game.Rounds)-1]
	if currentRound.CurrentTurn == nil || currentRound.CurrentTurn.Status != domain.TurnStatusActive {
		return errors.New("no active turn")
	}

	// Validate category is in selected categories
	validCategory := false
	for _, c := range game.Settings.SelectedCategories {
		if c == category {
			validCategory = true
			break
		}
	}

	if !validCategory {
		return errors.New("invalid category")
	}

	// Get random question from category
	question, err := s.questionRepo.GetRandomQuestion(ctx, category)
	if err != nil {
		return err
	}

	// Update round with question and category
	currentRound.Category = category
	currentRound.Question = question.Text
	currentRound.QuestionID = question.ID
	currentRound.AnswerPool.CorrectAnswer = question.Answer

	// Start answer writing timer using game settings
	currentRound.Timer = &domain.Timer{
		Type:      domain.TimerTypeAnswerWriting,
		StartTime: time.Now(),
		Duration:  game.Settings.TimeLimits.AnswerWriting,
		EndTime:   time.Now().Add(time.Duration(game.Settings.TimeLimits.AnswerWriting) * time.Second),
	}

	return s.UpdateGame(ctx, game)
}

// SubmitAnswer submits a player's answer for the current round
func (s *GameService) SubmitAnswer(ctx context.Context, gameID string, playerID string, answer string) error {
	game, err := s.GetGame(ctx, gameID)
	if err != nil {
		return err
	}

	currentRound := &game.Rounds[len(game.Rounds)-1]
	if currentRound.Status != domain.RoundStatusWaiting {
		return errors.New("round is not accepting answers")
	}

	// Check if answer is similar to any existing answer
	for _, ans := range currentRound.AnswerPool.FakeAnswers {
		if validation.IsSimilarAnswer(ans.Text, answer) {
			return errors.New("answer is too similar to an existing answer")
		}
	}

	// Check if answer is similar to correct answer
	if validation.IsSimilarAnswer(currentRound.AnswerPool.CorrectAnswer, answer) {
		return errors.New("answer is too similar to the correct answer")
	}

	// Add answer to pool
	newAnswer := domain.Answer{
		ID:        generateID(),
		PlayerID:  playerID,
		Text:      answer,
		Votes:     make([]string, 0),
		CreatedAt: time.Now(),
	}

	currentRound.AnswerPool.FakeAnswers = append(currentRound.AnswerPool.FakeAnswers, newAnswer)

	// Check if we need to generate filler answers
	if err := s.ensureAnswerPool(ctx, game); err != nil {
		return err
	}

	// If all players have submitted answers, start voting
	if len(currentRound.AnswerPool.FakeAnswers) == len(game.Players)-1 {
		currentRound.Status = domain.RoundStatusVoting
		currentRound.Timer = &domain.Timer{
			Type:      domain.TimerTypeVoting,
			StartTime: time.Now(),
			Duration:  30, // 30 seconds for voting
			EndTime:   time.Now().Add(30 * time.Second),
		}
	}

	return s.UpdateGame(ctx, game)
}

// ensureAnswerPool ensures we have n+1 answers in the pool
func (s *GameService) ensureAnswerPool(ctx context.Context, game *domain.Game) error {
	currentRound := &game.Rounds[len(game.Rounds)-1]
	requiredAnswers := len(game.Players)

	// If we have enough answers, no need for fillers
	if len(currentRound.AnswerPool.FakeAnswers) >= requiredAnswers {
		return nil
	}

	// Get the current question to access its filler answers
	question, err := s.questionRepo.GetByID(ctx, currentRound.QuestionID)
	if err != nil {
		return fmt.Errorf("failed to get question: %w", err)
	}

	// Shuffle filler answers to randomize selection
	fillerAnswers := make([]string, len(question.FillerAnswers))
	copy(fillerAnswers, question.FillerAnswers)
	rand.Shuffle(len(fillerAnswers), func(i, j int) {
		fillerAnswers[i], fillerAnswers[j] = fillerAnswers[j], fillerAnswers[i]
	})

	// Add filler answers until we have enough
	neededFillers := requiredAnswers - len(currentRound.AnswerPool.FakeAnswers)
	for i := 0; i < neededFillers && i < len(fillerAnswers); i++ {
		fillerAnswer := fillerAnswers[i]

		// Check if filler answer is similar to any existing answer
		isSimilar := false
		for _, ans := range currentRound.AnswerPool.FakeAnswers {
			if validation.IsSimilarAnswer(ans.Text, fillerAnswer) {
				isSimilar = true
				break
			}
		}
		if validation.IsSimilarAnswer(currentRound.AnswerPool.CorrectAnswer, fillerAnswer) {
			isSimilar = true
		}

		// If similar, skip this filler
		if isSimilar {
			continue
		}

		newAnswer := domain.Answer{
			ID:        generateID(),
			PlayerID:  "system",
			Text:      fillerAnswer,
			Votes:     make([]string, 0),
			CreatedAt: time.Now(),
		}

		currentRound.AnswerPool.FillerAnswers = append(currentRound.AnswerPool.FillerAnswers, newAnswer)
	}

	// If we still need more fillers, generate them using templates
	if len(currentRound.AnswerPool.FillerAnswers) < neededFillers {
		remaining := neededFillers - len(currentRound.AnswerPool.FillerAnswers)
		for i := 0; i < remaining; i++ {
			fillerAnswer, err := s.generateFillerAnswer(ctx, currentRound.Question, currentRound.Category)
			if err != nil {
				return err
			}

			// Check if filler answer is similar to any existing answer
			isSimilar := false
			for _, ans := range currentRound.AnswerPool.FakeAnswers {
				if validation.IsSimilarAnswer(ans.Text, fillerAnswer) {
					isSimilar = true
					break
				}
			}
			if validation.IsSimilarAnswer(currentRound.AnswerPool.CorrectAnswer, fillerAnswer) {
				isSimilar = true
			}

			// If similar, try again
			if isSimilar {
				i--
				continue
			}

			newAnswer := domain.Answer{
				ID:        generateID(),
				PlayerID:  "system",
				Text:      fillerAnswer,
				Votes:     make([]string, 0),
				CreatedAt: time.Now(),
			}

			currentRound.AnswerPool.FillerAnswers = append(currentRound.AnswerPool.FillerAnswers, newAnswer)
		}
	}

	return nil
}

// generateFillerAnswer generates a contextually appropriate filler answer
func (s *GameService) generateFillerAnswer(ctx context.Context, question string, category string) (string, error) {
	// For now, use a simple template-based approach
	templates := map[string][]string{
		"movies": {
			"A classic film about %s",
			"The story of %s",
			"A movie featuring %s",
			"A film starring %s",
		},
		"music": {
			"A song by %s",
			"A hit from %s",
			"A track featuring %s",
			"A collaboration with %s",
		},
		"books": {
			"A novel about %s",
			"A story featuring %s",
			"A book by %s",
			"A tale of %s",
		},
		"general": {
			"Something related to %s",
			"A thing about %s",
			"An item connected to %s",
			"A concept involving %s",
		},
	}

	// Get templates for category or use general ones
	tmpls, ok := templates[category]
	if !ok {
		tmpls = templates["general"]
	}

	// Extract key terms from the question
	terms := extractKeyTerms(question)
	if len(terms) == 0 {
		terms = []string{category}
	}

	// Select a random template and term
	tmpl := tmpls[rand.Intn(len(tmpls))]
	term := terms[rand.Intn(len(terms))]

	return fmt.Sprintf(tmpl, term), nil
}

// extractKeyTerms extracts key terms from a question
func extractKeyTerms(question string) []string {
	// Simple implementation: split by spaces and filter out common words
	words := strings.Fields(strings.ToLower(question))
	commonWords := map[string]bool{
		"the": true, "a": true, "an": true, "in": true, "on": true,
		"at": true, "to": true, "for": true, "with": true, "by": true,
		"what": true, "who": true, "where": true, "when": true, "why": true,
		"how": true, "is": true, "are": true, "was": true, "were": true,
	}

	var terms []string
	for _, word := range words {
		if !commonWords[word] && len(word) > 2 {
			terms = append(terms, word)
		}
	}

	return terms
}

// SubmitVote submits a player's vote for an answer
func (s *GameService) SubmitVote(ctx context.Context, gameID string, playerID string, answerID string) error {
	game, err := s.GetGame(ctx, gameID)
	if err != nil {
		return err
	}

	currentRound := &game.Rounds[len(game.Rounds)-1]
	if currentRound.Status != domain.RoundStatusVoting {
		return errors.New("round is not in voting phase")
	}

	// Check if player has already voted
	for _, answer := range currentRound.AnswerPool.FakeAnswers {
		for _, vote := range answer.Votes {
			if vote == playerID {
				return domain.ErrVoteSubmitted
			}
		}
	}

	// Find the answer and add the vote
	found := false
	for i := range currentRound.AnswerPool.FakeAnswers {
		if currentRound.AnswerPool.FakeAnswers[i].ID == answerID {
			currentRound.AnswerPool.FakeAnswers[i].Votes = append(currentRound.AnswerPool.FakeAnswers[i].Votes, playerID)
			found = true
			break
		}
	}

	if !found {
		return domain.ErrInvalidVote
	}

	// Check if all players have voted
	totalVotes := 0
	for _, answer := range currentRound.AnswerPool.FakeAnswers {
		totalVotes += len(answer.Votes)
	}

	if totalVotes == len(game.Players) {
		// Group answers by content to find duplicates
		answerGroups := make(map[string][]domain.Answer)
		for _, answer := range currentRound.AnswerPool.FakeAnswers {
			answerGroups[answer.Text] = append(answerGroups[answer.Text], answer)
		}

		// Calculate scores
		for _, group := range answerGroups {
			if len(group) > 1 {
				// This is a group of duplicate answers
				// Each player in the group gets points for each vote their shared answer received
				totalVotesForGroup := 0
				for _, answer := range group {
					totalVotesForGroup += len(answer.Votes)
				}

				// Each player in the group gets points equal to the total votes
				pointsPerPlayer := totalVotesForGroup
				for _, answer := range group {
					for i := range game.Players {
						if game.Players[i].ID == answer.PlayerID {
							game.Players[i].Score += pointsPerPlayer
							break
						}
					}
				}
			} else {
				// Single answer - normal scoring
				answer := group[0]
				points := len(answer.Votes)
				for i := range game.Players {
					if game.Players[i].ID == answer.PlayerID {
						game.Players[i].Score += points
						break
					}
				}
			}
		}

		currentRound.Status = domain.RoundStatusCompleted
		currentRound.EndTime = time.Now()

		// Notify all players of round end and scores
		payload, err := json.Marshal(game)
		if err != nil {
			return err
		}
		s.hub.BroadcastToGame(game.ID, "round_ended", payload)
	}

	return s.UpdateGame(ctx, game)
}

// EndRound ends the current round and starts a new one
func (s *GameService) EndRound(ctx context.Context, code string) error {
	game, err := s.GetGame(ctx, code)
	if err != nil {
		return err
	}

	if len(game.Rounds) == 0 {
		return errors.New("no rounds found in game")
	}

	currentRound := &game.Rounds[len(game.Rounds)-1]
	if currentRound.Status == domain.RoundStatusCompleted {
		return errors.New("round already completed")
	}

	// Mark round as completed
	currentRound.Status = domain.RoundStatusCompleted
	currentRound.EndTime = time.Now()

	// Notify all players of round end and scores
	payload, err := json.Marshal(game)
	if err != nil {
		return err
	}
	s.hub.BroadcastToGame(game.ID, "round_ended", payload)

	return s.UpdateGame(ctx, game)
}

// EndGame ends a game session
func (s *GameService) EndGame(ctx context.Context, code string) error {
	game, err := s.GetGame(ctx, code)
	if err != nil {
		return err
	}

	if game.Status == domain.GameStatusEnded {
		return errors.New("game has already ended")
	}

	game.Status = domain.GameStatusEnded
	game.UpdatedAt = time.Now()
	game.LastActivity = time.Now()

	if err := s.UpdateGame(ctx, game); err != nil {
		return err
	}

	// Notify all clients about game end
	payload, err := json.Marshal(game)
	if err != nil {
		return err
	}

	s.hub.BroadcastToGame(game.ID, "game_ended", payload)

	// Clean up game session
	if err := s.sessionMgr.DeleteGame(ctx, game.ID); err != nil {
		return err
	}

	return nil
}

// HandlePlayerReconnection handles a player reconnecting to the game
func (s *GameService) HandlePlayerReconnection(ctx context.Context, gameID string, playerID string) error {
	game, err := s.GetGame(ctx, gameID)
	if err != nil {
		return err
	}

	// Find player
	var player *domain.Player
	for i, p := range game.Players {
		if p.ID == playerID {
			player = &game.Players[i]
			break
		}
	}

	if player == nil {
		return errors.New("player not found in game")
	}

	// Update player status
	player.IsConnected = true
	player.LastSeen = time.Now()

	if err := s.UpdateGame(ctx, game); err != nil {
		return err
	}

	// Notify all clients about player reconnection
	payload, err := json.Marshal(player)
	if err != nil {
		return err
	}

	s.hub.BroadcastToGame(game.ID, "player_reconnected", payload)

	return nil
}

// HandlePlayerDisconnection handles a player disconnecting from the game
func (s *GameService) HandlePlayerDisconnection(ctx context.Context, gameID string, playerID string) error {
	game, err := s.GetGame(ctx, gameID)
	if err != nil {
		return err
	}

	// Update player's active status
	for i := range game.Players {
		if game.Players[i].ID == playerID {
			game.Players[i].IsConnected = false
			game.Players[i].LastSeen = time.Now()
			break
		}
	}

	// Update game state
	if err := s.UpdateGame(ctx, game); err != nil {
		return err
	}

	// Notify other players
	payload, err := json.Marshal(map[string]interface{}{
		"player_id": playerID,
		"status":    "disconnected",
	})
	if err != nil {
		return err
	}
	s.hub.BroadcastToGame(gameID, "player_disconnected", payload)

	return nil
}

// CleanupInactiveGames cleans up inactive game sessions
func (s *GameService) CleanupInactiveGames(ctx context.Context) error {
	games, err := s.sessionMgr.GetAllGames(ctx)
	if err != nil {
		return err
	}

	now := time.Now()
	for _, game := range games {
		// Check if game has been inactive for more than 24 hours
		if now.Sub(game.LastActivity) > 24*time.Hour {
			if err := s.EndGame(ctx, game.Code); err != nil {
				return err
			}
		}
	}

	return nil
}

// Helper functions

func generateID() string {
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	b := make([]byte, 8)
	for i := range b {
		b[i] = charset[rand.Intn(len(charset))]
	}
	return string(b)
}

func generateGameCode() string {
	const charset = "ABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	b := make([]byte, 6)
	for i := range b {
		b[i] = charset[rand.Intn(len(charset))]
	}
	return string(b)
}
