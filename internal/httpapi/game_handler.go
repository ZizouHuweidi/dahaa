package httpapi

import (
	"net/http"

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

func (h *GameHandler) CreateGame(w http.ResponseWriter, r *http.Request) {
	var req createGameRequest
	if err := readJSON(r, &req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}
	if req.Player.ID == "" || req.Player.Name == "" {
		writeError(w, http.StatusBadRequest, "player id and name are required")
		return
	}
	game, err := h.gameService.CreateGame(r.Context(), req.Code, req.Player, req.Settings)
	if err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}
	writeJSON(w, http.StatusCreated, game)
}

func (h *GameHandler) GetGame(w http.ResponseWriter, r *http.Request) {
	game, err := h.gameService.GetGame(r.Context(), r.PathValue("code"))
	if err != nil {
		writeError(w, http.StatusNotFound, "game not found")
		return
	}
	writeJSON(w, http.StatusOK, game)
}

func (h *GameHandler) JoinGame(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Player domain.Player `json:"player"`
	}
	if err := readJSON(r, &req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}
	if req.Player.ID == "" || req.Player.Name == "" {
		writeError(w, http.StatusBadRequest, "player id and name are required")
		return
	}
	if err := h.gameService.JoinGame(r.Context(), r.PathValue("code"), req.Player); err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, map[string]string{"message": "successfully joined game"})
}

func (h *GameHandler) StartGame(w http.ResponseWriter, r *http.Request) {
	if err := h.gameService.StartGame(r.Context(), r.PathValue("code")); err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func (h *GameHandler) StartTurn(w http.ResponseWriter, r *http.Request) {
	var req struct {
		PlayerID string `json:"player_id"`
	}
	if err := readJSON(r, &req); err != nil || req.PlayerID == "" {
		writeError(w, http.StatusBadRequest, "player_id is required")
		return
	}
	if err := h.gameService.StartTurn(r.Context(), r.PathValue("code"), req.PlayerID); err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func (h *GameHandler) SelectCategory(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Category string `json:"category"`
	}
	if err := readJSON(r, &req); err != nil || req.Category == "" {
		writeError(w, http.StatusBadRequest, "category is required")
		return
	}
	if err := h.gameService.SelectCategory(r.Context(), r.PathValue("code"), req.Category); err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func (h *GameHandler) SubmitAnswer(w http.ResponseWriter, r *http.Request) {
	var req struct {
		PlayerID string `json:"player_id"`
		Answer   string `json:"answer"`
	}
	if err := readJSON(r, &req); err != nil || req.PlayerID == "" || req.Answer == "" {
		writeError(w, http.StatusBadRequest, "player_id and answer are required")
		return
	}
	if err := h.gameService.SubmitAnswer(r.Context(), r.PathValue("code"), req.PlayerID, req.Answer); err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, map[string]string{"message": "answer submitted successfully"})
}

func (h *GameHandler) SubmitVote(w http.ResponseWriter, r *http.Request) {
	var req struct {
		PlayerID string `json:"player_id"`
		AnswerID string `json:"answer_id"`
	}
	if err := readJSON(r, &req); err != nil || req.PlayerID == "" || req.AnswerID == "" {
		writeError(w, http.StatusBadRequest, "player_id and answer_id are required")
		return
	}
	if err := h.gameService.SubmitVote(r.Context(), r.PathValue("code"), req.PlayerID, req.AnswerID); err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, map[string]string{"message": "vote submitted successfully"})
}

func (h *GameHandler) EndRound(w http.ResponseWriter, r *http.Request) {
	if err := h.gameService.EndRound(r.Context(), r.PathValue("code")); err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, map[string]string{"message": "round ended successfully"})
}

func (h *GameHandler) EndGame(w http.ResponseWriter, r *http.Request) {
	if err := h.gameService.EndGame(r.Context(), r.PathValue("code")); err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

type createQuestionRequest struct {
	Category      string   `json:"category"`
	Text          string   `json:"text"`
	Answer        string   `json:"answer"`
	FillerAnswers []string `json:"filler_answers"`
}

func (h *GameHandler) BulkCreateQuestions(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Questions []createQuestionRequest `json:"questions"`
	}
	if err := readJSON(r, &req); err != nil || len(req.Questions) == 0 {
		writeError(w, http.StatusBadRequest, "questions are required")
		return
	}
	questions := make([]*domain.Question, 0, len(req.Questions))
	for _, q := range req.Questions {
		if q.Category == "" || q.Text == "" || q.Answer == "" || len(q.FillerAnswers) < 3 {
			writeError(w, http.StatusBadRequest, "each question requires category, text, answer, and at least 3 filler answers")
			return
		}
		questions = append(questions, &domain.Question{Text: q.Text, Answer: q.Answer, Category: q.Category, FillerAnswers: q.FillerAnswers})
	}
	if err := h.questionRepo.BulkCreateQuestions(r.Context(), questions); err != nil {
		writeError(w, http.StatusInternalServerError, "failed to create questions: "+err.Error())
		return
	}
	writeJSON(w, http.StatusCreated, map[string]any{"message": "questions created successfully", "count": len(questions)})
}

func (h *GameHandler) GetCategories(w http.ResponseWriter, r *http.Request) {
	categories, err := h.questionRepo.GetCategories(r.Context())
	if err != nil {
		writeError(w, http.StatusInternalServerError, "failed to load categories")
		return
	}
	writeJSON(w, http.StatusOK, map[string]any{"categories": categories})
}
