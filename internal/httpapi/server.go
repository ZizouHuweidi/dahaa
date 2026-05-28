package httpapi

import (
	"net/http"

	"github.com/zizouhuweidi/dahaa/internal/domain"
	"github.com/zizouhuweidi/dahaa/internal/service"
	"github.com/zizouhuweidi/dahaa/internal/storage"
	ws "github.com/zizouhuweidi/dahaa/internal/websocket"
)

type Server struct {
	mux *http.ServeMux
}

type Deps struct {
	GameService  domain.GameService
	QuestionRepo domain.QuestionRepository
	UserService  *service.UserService
	ImageStorage *storage.ImageStorage
	Hub          *ws.Hub
}

func NewServer(deps Deps) *Server {
	mux := http.NewServeMux()

	gameHandler := NewGameHandler(deps.GameService, deps.QuestionRepo)
	userHandler := NewUserHandler(deps.UserService)
	imageHandler := NewImageHandler(deps.ImageStorage)
	wsHandler := ws.NewHandler(deps.Hub)

	mux.HandleFunc("POST /api/games", gameHandler.CreateGame)
	mux.HandleFunc("GET /api/games/{code}", gameHandler.GetGame)
	mux.HandleFunc("POST /api/games/{code}/join", gameHandler.JoinGame)
	mux.HandleFunc("POST /api/games/{code}/start", gameHandler.StartGame)
	mux.HandleFunc("POST /api/games/{code}/turns", gameHandler.StartTurn)
	mux.HandleFunc("POST /api/games/{code}/turns/category", gameHandler.SelectCategory)
	mux.HandleFunc("POST /api/games/{code}/rounds/{round}/answers", gameHandler.SubmitAnswer)
	mux.HandleFunc("POST /api/games/{code}/rounds/{round}/votes", gameHandler.SubmitVote)
	mux.HandleFunc("POST /api/games/{code}/rounds/{round}/end", gameHandler.EndRound)
	mux.HandleFunc("POST /api/games/{code}/end", gameHandler.EndGame)
	mux.HandleFunc("POST /api/questions/bulk", gameHandler.BulkCreateQuestions)
	mux.HandleFunc("GET /api/questions/categories", gameHandler.GetCategories)

	mux.HandleFunc("POST /api/users/register", userHandler.Register)
	mux.HandleFunc("POST /api/users/login", userHandler.Login)

	mux.HandleFunc("POST /api/images", imageHandler.UploadImage)
	mux.HandleFunc("GET /api/images/{filename}", imageHandler.ServeImage)

	mux.HandleFunc("GET /ws", wsHandler.ServeWS)
	mux.HandleFunc("GET /health", func(w http.ResponseWriter, r *http.Request) {
		writeJSON(w, http.StatusOK, map[string]string{"status": "ok"})
	})

	return &Server{mux: mux}
}

func (s *Server) Handler() http.Handler {
	return Chain(s.mux, Recover, Logger, CORS)
}
