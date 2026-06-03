package httpapi

import (
	"net/http"

	"github.com/labstack/echo/v5"
	"github.com/zizouhuweidi/dahaa/internal/domain"
)

type GameHandler struct {
	gameService  domain.GameService
	questionRepo domain.QuestionRepository
}

func NewGameHandler(gameService domain.GameService, questionRepo domain.QuestionRepository) *GameHandler {
	return &GameHandler{gameService: gameService, questionRepo: questionRepo}
}

type createGameRequest struct {
	Code     string               `json:"code"`
	Player   domain.Player        `json:"player"`
	Settings *domain.GameSettings `json:"settings"`
}

func (h *GameHandler) CreateGame(c *echo.Context) error {
	var req createGameRequest
	if err := c.Bind(&req); err != nil {
		return err
	}
	if req.Player.ID == "" || req.Player.Name == "" {
		return echo.NewHTTPError(http.StatusBadRequest, "player id and name are required")
	}
	game, err := h.gameService.CreateGame(c.Request().Context(), req.Code, req.Player, req.Settings)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}
	return c.JSON(http.StatusCreated, game)
}

func (h *GameHandler) GetGame(c *echo.Context) error {
	game, err := h.gameService.GetGame(c.Request().Context(), c.Param("code"))
	if err != nil {
		return echo.NewHTTPError(http.StatusNotFound, "game not found")
	}
	return c.JSON(http.StatusOK, game)
}

func (h *GameHandler) JoinGame(c *echo.Context) error {
	var req struct {
		Player domain.Player `json:"player"`
	}
	if err := c.Bind(&req); err != nil {
		return err
	}
	if req.Player.ID == "" || req.Player.Name == "" {
		return echo.NewHTTPError(http.StatusBadRequest, "player id and name are required")
	}
	if err := h.gameService.JoinGame(c.Request().Context(), c.Param("code"), req.Player); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}
	return c.JSON(http.StatusOK, map[string]string{"message": "successfully joined game"})
}

func (h *GameHandler) StartGame(c *echo.Context) error {
	if err := h.gameService.StartGame(c.Request().Context(), c.Param("code")); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}
	return c.NoContent(http.StatusNoContent)
}

func (h *GameHandler) StartTurn(c *echo.Context) error {
	var req struct {
		PlayerID string `json:"player_id"`
	}
	if err := c.Bind(&req); err != nil {
		return err
	}
	if req.PlayerID == "" {
		return echo.NewHTTPError(http.StatusBadRequest, "player_id is required")
	}
	if err := h.gameService.StartTurn(c.Request().Context(), c.Param("code"), req.PlayerID); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}
	return c.NoContent(http.StatusNoContent)
}

func (h *GameHandler) SelectCategory(c *echo.Context) error {
	var req struct {
		Category string `json:"category"`
	}
	if err := c.Bind(&req); err != nil {
		return err
	}
	if req.Category == "" {
		return echo.NewHTTPError(http.StatusBadRequest, "category is required")
	}
	if err := h.gameService.SelectCategory(c.Request().Context(), c.Param("code"), req.Category); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}
	return c.NoContent(http.StatusNoContent)
}

func (h *GameHandler) SubmitAnswer(c *echo.Context) error {
	var req struct {
		PlayerID string `json:"player_id"`
		Answer   string `json:"answer"`
	}
	if err := c.Bind(&req); err != nil {
		return err
	}
	if req.PlayerID == "" || req.Answer == "" {
		return echo.NewHTTPError(http.StatusBadRequest, "player_id and answer are required")
	}
	if err := h.gameService.SubmitAnswer(c.Request().Context(), c.Param("code"), req.PlayerID, req.Answer); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}
	return c.JSON(http.StatusOK, map[string]string{"message": "answer submitted successfully"})
}

func (h *GameHandler) SubmitVote(c *echo.Context) error {
	var req struct {
		PlayerID string `json:"player_id"`
		AnswerID string `json:"answer_id"`
	}
	if err := c.Bind(&req); err != nil {
		return err
	}
	if req.PlayerID == "" || req.AnswerID == "" {
		return echo.NewHTTPError(http.StatusBadRequest, "player_id and answer_id are required")
	}
	if err := h.gameService.SubmitVote(c.Request().Context(), c.Param("code"), req.PlayerID, req.AnswerID); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}
	return c.JSON(http.StatusOK, map[string]string{"message": "vote submitted successfully"})
}

func (h *GameHandler) EndRound(c *echo.Context) error {
	if err := h.gameService.EndRound(c.Request().Context(), c.Param("code")); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}
	return c.JSON(http.StatusOK, map[string]string{"message": "round ended successfully"})
}

func (h *GameHandler) EndGame(c *echo.Context) error {
	if err := h.gameService.EndGame(c.Request().Context(), c.Param("code")); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}
	return c.NoContent(http.StatusNoContent)
}

type createQuestionRequest struct {
	Category      string   `json:"category"`
	Text          string   `json:"text"`
	Answer        string   `json:"answer"`
	FillerAnswers []string `json:"filler_answers"`
}

func (h *GameHandler) BulkCreateQuestions(c *echo.Context) error {
	var req struct {
		Questions []createQuestionRequest `json:"questions"`
	}
	if err := c.Bind(&req); err != nil {
		return err
	}
	if len(req.Questions) == 0 {
		return echo.NewHTTPError(http.StatusBadRequest, "questions are required")
	}
	questions := make([]*domain.Question, 0, len(req.Questions))
	for _, q := range req.Questions {
		if q.Category == "" || q.Text == "" || q.Answer == "" || len(q.FillerAnswers) < 3 {
			return echo.NewHTTPError(http.StatusBadRequest, "each question requires category, text, answer, and at least 3 filler answers")
		}
		questions = append(questions, &domain.Question{Text: q.Text, Answer: q.Answer, Category: q.Category, FillerAnswers: q.FillerAnswers})
	}
	if err := h.questionRepo.BulkCreateQuestions(c.Request().Context(), questions); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "failed to create questions: "+err.Error())
	}
	return c.JSON(http.StatusCreated, map[string]any{"message": "questions created successfully", "count": len(questions)})
}

func (h *GameHandler) GetCategories(c *echo.Context) error {
	categories, err := h.questionRepo.GetCategories(c.Request().Context())
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "failed to load categories")
	}
	return c.JSON(http.StatusOK, map[string]any{"categories": categories})
}
