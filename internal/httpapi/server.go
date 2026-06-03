package httpapi

import (
	"log/slog"
	"net/http"

	"github.com/labstack/echo/v5"
	"github.com/labstack/echo/v5/middleware"
	"github.com/zizouhuweidi/dahaa/internal/domain"
	"github.com/zizouhuweidi/dahaa/internal/service"
	"github.com/zizouhuweidi/dahaa/internal/storage"
	ws "github.com/zizouhuweidi/dahaa/internal/websocket"
)

type Server struct {
	e *echo.Echo
}

type Deps struct {
	GameService  domain.GameService
	QuestionRepo domain.QuestionRepository
	UserService  *service.UserService
	ImageStorage *storage.ImageStorage
	Hub          *ws.Hub
}

func NewServer(deps Deps) *Server {
	e := echo.New()
	e.Binder = strictJSONBinder{}
	e.HTTPErrorHandler = httpErrorHandler
	e.Use(middleware.RequestID())
	e.Use(middleware.RequestLoggerWithConfig(middleware.RequestLoggerConfig{
		HandleError:  true,
		LogLatency:   true,
		LogMethod:    true,
		LogURIPath:   true,
		LogRoutePath: true,
		LogRequestID: true,
		LogStatus:    true,
		LogValuesFunc: func(_ *echo.Context, v middleware.RequestLoggerValues) error {
			attrs := []any{
				"id", v.RequestID,
				"method", v.Method,
				"path", v.URIPath,
				"route", v.RoutePath,
				"status", v.Status,
				"duration", v.Latency.String(),
			}
			if v.Error != nil {
				attrs = append(attrs, "error", v.Error)
			}
			slog.Info("request completed", attrs...)
			return nil
		},
	}))
	e.Use(middleware.Recover())
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"*"},
		AllowMethods: []string{http.MethodGet, http.MethodPost, http.MethodOptions},
		AllowHeaders: []string{echo.HeaderContentType, echo.HeaderAuthorization},
	}))

	gameHandler := NewGameHandler(deps.GameService, deps.QuestionRepo)
	userHandler := NewUserHandler(deps.UserService)
	imageHandler := NewImageHandler(deps.ImageStorage)
	wsHandler := ws.NewHandler(deps.Hub)

	api := e.Group("/api")

	games := api.Group("/games")
	games.POST("", gameHandler.CreateGame)
	games.GET("/:code", gameHandler.GetGame)
	games.POST("/:code/join", gameHandler.JoinGame)
	games.POST("/:code/start", gameHandler.StartGame)
	games.POST("/:code/turns", gameHandler.StartTurn)
	games.POST("/:code/turns/category", gameHandler.SelectCategory)
	games.POST("/:code/rounds/:round/answers", gameHandler.SubmitAnswer)
	games.POST("/:code/rounds/:round/votes", gameHandler.SubmitVote)
	games.POST("/:code/rounds/:round/end", gameHandler.EndRound)
	games.POST("/:code/end", gameHandler.EndGame)

	questions := api.Group("/questions")
	questions.POST("/bulk", gameHandler.BulkCreateQuestions)
	questions.GET("/categories", gameHandler.GetCategories)

	users := api.Group("/users")
	users.POST("/register", userHandler.Register)
	users.POST("/login", userHandler.Login)

	images := api.Group("/images")
	images.POST("", imageHandler.UploadImage, middleware.BodyLimit(6<<20))
	images.GET("/:filename", imageHandler.ServeImage)

	e.GET("/ws", echo.WrapHandler(http.HandlerFunc(wsHandler.ServeWS)))
	e.GET("/health", func(c *echo.Context) error {
		return c.JSON(http.StatusOK, map[string]string{"status": "ok"})
	})

	return &Server{e: e}
}

func (s *Server) Handler() http.Handler {
	return s.e
}
