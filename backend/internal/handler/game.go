package handler

import (
	"net/http"

	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
	"github.com/zizouhuweidi/dahaa/internal/domain"
	"github.com/zizouhuweidi/dahaa/internal/service"
)

// GameHandler handles game-related HTTP requests
type GameHandler struct {
	gameService  domain.GameService
	questionRepo domain.QuestionRepository
	validate     *validator.Validate
}

// NewGameHandler creates a new game handler
func NewGameHandler(gameService domain.GameService, questionRepo domain.QuestionRepository) *GameHandler {
	return &GameHandler{
		gameService:  gameService,
		questionRepo: questionRepo,
		validate:     validator.New(),
	}
}

// Register registers the game routes
func (h *GameHandler) Register(e *echo.Echo) {
	g := e.Group("/api/games")
	g.POST("", h.CreateGame)
	g.POST("/join", h.JoinGame)
	g.POST("/:id/start", h.StartGame)
	g.POST("/:id/turns", h.StartTurn)
	g.POST("/:id/turns/category", h.SelectCategory)
	g.POST("/:id/rounds/:round/answers", h.SubmitAnswer)
	g.POST("/:id/rounds/:round/votes", h.SubmitVote)
	g.POST("/:id/rounds/:round/end", h.EndRound)
	g.POST("/:id/end", h.EndGame)
	g.GET("/:code", h.GetGame)
	g.POST("/questions/bulk", h.BulkCreateQuestions)
}

// CreateGameRequest represents the request to create a new game
type CreateGameRequest struct {
	Code     string               `json:"code" validate:"required,min=4,max=6"`
	Player   domain.Player        `json:"player" validate:"required"`
	Settings *domain.GameSettings `json:"settings"`
}

// CreateQuestionRequest represents the request to create a new question
type CreateQuestionRequest struct {
	Category      string   `json:"category" validate:"required"`
	Text          string   `json:"text" validate:"required"`
	Answer        string   `json:"answer" validate:"required"`
	FillerAnswers []string `json:"filler_answers" validate:"required,min=3"`
}

// CreateGame handles the creation of a new game
func (h *GameHandler) CreateGame(c echo.Context) error {
	var req CreateGameRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "Invalid request body",
		})
	}

	if err := c.Validate(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": err.Error(),
		})
	}

	game, err := h.gameService.CreateGame(c.Request().Context(), req.Code, req.Player, req.Settings)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": err.Error(),
		})
	}

	return c.JSON(http.StatusCreated, game)
}

// JoinGameRequest represents the request body for joining a game
type JoinGameRequest struct {
	Player domain.Player `json:"player"`
}

// JoinGame handles a player joining a game
func (h *GameHandler) JoinGame(c echo.Context) error {
	code := c.Param("code")
	var player domain.Player
	if err := c.Bind(&player); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": err.Error(),
		})
	}

	if err := h.gameService.JoinGame(c.Request().Context(), code, player); err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": err.Error(),
		})
	}

	return c.JSON(http.StatusOK, map[string]string{
		"message": "Successfully joined game",
	})
}

// StartGame starts a game
func (h *GameHandler) StartGame(c echo.Context) error {
	gameID := c.Param("id")
	if gameID == "" {
		return echo.NewHTTPError(http.StatusBadRequest, "game ID is required")
	}

	if err := h.gameService.StartGame(c.Request().Context(), gameID); err != nil {
		switch err {
		case service.ErrGameNotFound:
			return echo.NewHTTPError(http.StatusNotFound, "game not found")
		case service.ErrGameInProgress:
			return echo.NewHTTPError(http.StatusConflict, "game is already in progress")
		default:
			return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
		}
	}

	return c.NoContent(http.StatusOK)
}

// StartTurnRequest represents the request body for starting a turn
type StartTurnRequest struct {
	PlayerID string `json:"player_id" validate:"required"`
}

// StartTurn starts a new turn for a player
func (h *GameHandler) StartTurn(c echo.Context) error {
	gameID := c.Param("id")
	if gameID == "" {
		return echo.NewHTTPError(http.StatusBadRequest, "game ID is required")
	}

	var req StartTurnRequest
	if err := c.Bind(&req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid request body")
	}

	if err := c.Validate(&req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	if err := h.gameService.StartTurn(c.Request().Context(), gameID, req.PlayerID); err != nil {
		switch err {
		case service.ErrGameNotFound:
			return echo.NewHTTPError(http.StatusNotFound, "game not found")
		case service.ErrGameNotStarted:
			return echo.NewHTTPError(http.StatusConflict, "game has not started")
		default:
			return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
		}
	}

	return c.NoContent(http.StatusOK)
}

// SelectCategoryRequest represents the request body for selecting a category
type SelectCategoryRequest struct {
	Category string `json:"category" validate:"required"`
}

// SelectCategory handles category selection during a turn
func (h *GameHandler) SelectCategory(c echo.Context) error {
	gameID := c.Param("id")
	if gameID == "" {
		return echo.NewHTTPError(http.StatusBadRequest, "game ID is required")
	}

	var req SelectCategoryRequest
	if err := c.Bind(&req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid request body")
	}

	if err := c.Validate(&req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	if err := h.gameService.SelectCategory(c.Request().Context(), gameID, req.Category); err != nil {
		switch err {
		case service.ErrGameNotFound:
			return echo.NewHTTPError(http.StatusNotFound, "game not found")
		case service.ErrGameNotStarted:
			return echo.NewHTTPError(http.StatusConflict, "game has not started")
		default:
			return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
		}
	}

	return c.NoContent(http.StatusOK)
}

// SubmitAnswer handles a player submitting their answer
func (h *GameHandler) SubmitAnswer(c echo.Context) error {
	code := c.Param("code")
	var req struct {
		PlayerID string `json:"player_id" validate:"required"`
		Answer   string `json:"answer" validate:"required"`
	}

	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": err.Error(),
		})
	}

	if err := h.validate.Struct(req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": err.Error(),
		})
	}

	if err := h.gameService.SubmitAnswer(c.Request().Context(), code, req.PlayerID, req.Answer); err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": err.Error(),
		})
	}

	return c.JSON(http.StatusOK, map[string]string{
		"message": "Answer submitted successfully",
	})
}

// SubmitVote handles a player submitting their vote
func (h *GameHandler) SubmitVote(c echo.Context) error {
	code := c.Param("code")
	var req struct {
		PlayerID string `json:"player_id" validate:"required"`
		AnswerID string `json:"answer_id" validate:"required"`
	}

	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": err.Error(),
		})
	}

	if err := h.validate.Struct(req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": err.Error(),
		})
	}

	if err := h.gameService.SubmitVote(c.Request().Context(), code, req.PlayerID, req.AnswerID); err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": err.Error(),
		})
	}

	return c.JSON(http.StatusOK, map[string]string{
		"message": "Vote submitted successfully",
	})
}

// EndRound handles ending the current round
func (h *GameHandler) EndRound(c echo.Context) error {
	code := c.Param("code")
	if err := h.gameService.EndRound(c.Request().Context(), code); err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": err.Error(),
		})
	}

	return c.JSON(http.StatusOK, map[string]string{
		"message": "Round ended successfully",
	})
}

// EndGame ends the game session
func (h *GameHandler) EndGame(c echo.Context) error {
	gameID := c.Param("id")
	if gameID == "" {
		return echo.NewHTTPError(http.StatusBadRequest, "game ID is required")
	}

	if err := h.gameService.EndGame(c.Request().Context(), gameID); err != nil {
		switch err {
		case service.ErrGameNotFound:
			return echo.NewHTTPError(http.StatusNotFound, "game not found")
		default:
			return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
		}
	}

	return c.NoContent(http.StatusOK)
}

// BulkCreateQuestionsRequest represents the request body for bulk creating questions
type BulkCreateQuestionsRequest struct {
	Questions []CreateQuestionRequest `json:"questions" validate:"required,min=1,dive"`
}

// BulkCreateQuestions handles bulk creation of questions
func (h *GameHandler) BulkCreateQuestions(c echo.Context) error {
	var req BulkCreateQuestionsRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "Invalid request body",
		})
	}

	// Validate the request
	if err := h.validate.Struct(req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": err.Error(),
		})
	}

	// Convert to domain.Question slice
	questions := make([]*domain.Question, 0, len(req.Questions))
	for _, q := range req.Questions {
		questions = append(questions, &domain.Question{
			Text:          q.Text,
			Answer:        q.Answer,
			Category:      q.Category,
			FillerAnswers: q.FillerAnswers,
		})
	}

	// Create questions in bulk
	if err := h.questionRepo.BulkCreateQuestions(c.Request().Context(), questions); err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": "Failed to create questions: " + err.Error(),
		})
	}

	return c.JSON(http.StatusCreated, map[string]interface{}{
		"message":  "Questions created successfully",
		"count":   len(questions),
	})
}

// GetGame handles retrieving a game by code
func (h *GameHandler) GetGame(c echo.Context) error {
	code := c.Param("code")
	if code == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "Game code is required",
		})
	}

	game, err := h.gameService.GetGame(c.Request().Context(), code)
	if err != nil {
		return c.JSON(http.StatusNotFound, map[string]string{
			"error": "Game not found",
		})
	}

	return c.JSON(http.StatusOK, game)
}
