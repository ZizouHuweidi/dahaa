package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"path/filepath"
	"syscall"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/redis/go-redis/v9"
	"github.com/zizouhuweidi/dahaa/internal/handler"
	"github.com/zizouhuweidi/dahaa/internal/repository/postgres"
	"github.com/zizouhuweidi/dahaa/internal/service"
	"github.com/zizouhuweidi/dahaa/internal/session"
	"github.com/zizouhuweidi/dahaa/internal/storage"
	"github.com/zizouhuweidi/dahaa/internal/websocket"
)

func main() {
	// Initialize database connection
	pool, err := postgres.NewDB()
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer pool.Close()

	// Initialize Redis client
	redisClient := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "", // no password set
		DB:       0,  // use default DB
	})

	// Initialize image storage
	imageStorage, err := storage.NewImageStorage(filepath.Join("uploads", "images"))
	if err != nil {
		log.Fatalf("Failed to initialize image storage: %v", err)
	}

	// Initialize repositories
	userRepo := postgres.NewUserRepository(pool)
	gameInviteRepo := postgres.NewGameInviteRepository(pool)
	gameRepo := postgres.NewGameRepository(pool)
	questionRepo := postgres.NewQuestionRepository(pool)

	// Initialize session manager
	sessionManager := session.NewManager(redisClient)

	// Initialize websocket hub
	hub := websocket.NewHub()
	go hub.Run()

	// Initialize services
	userService := service.NewUserService(userRepo, gameInviteRepo)
	gameService := service.NewGameService(gameRepo, questionRepo, hub, sessionManager)

	// Initialize handlers
	userHandler := handler.NewUserHandler(userService)
	gameHandler := handler.NewGameHandler(gameService, questionRepo)
	wsHandler := handler.NewWebSocketHandler(hub)
	imageHandler := handler.NewImageHandler(imageStorage)

	// Initialize Echo
	e := echo.New()

	// Middleware
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Use(middleware.CORS())

	// Routes
	api := e.Group("/api")

	// User routes
	users := api.Group("/users")
	users.POST("/register", userHandler.Register)
	users.POST("/login", userHandler.Login)
	users.POST("/invites/:game_id/:to_user_id", userHandler.SendGameInvite)
	users.POST("/invites/:invite_id/accept", userHandler.AcceptGameInvite)
	users.POST("/invites/:invite_id/decline", userHandler.DeclineGameInvite)
	users.GET("/invites", userHandler.GetPendingInvites)

	// Game routes
	games := api.Group("/games")
	games.POST("", gameHandler.CreateGame)
	games.GET("/:code", gameHandler.GetGame)
	games.POST("/:code/join", gameHandler.JoinGame)
	games.POST("/:code/start", gameHandler.StartGame)
	games.POST("/:code/rounds/:round/answers", gameHandler.SubmitAnswer)
	games.POST("/:code/rounds/:round/votes", gameHandler.SubmitVote)
	games.POST("/:code/rounds/:round/end", gameHandler.EndRound)
	games.POST("/:code/end", gameHandler.EndGame)

	// WebSocket route
	e.GET("/ws", wsHandler.HandleWebSocket)

	// Image routes
	e.POST("/api/images", imageHandler.UploadImage)
	e.GET("/api/images/:filename", imageHandler.ServeImage)

	// Health check endpoint
	e.GET("/health", func(c echo.Context) error {
		return c.JSON(http.StatusOK, map[string]string{
			"status": "ok",
		})
	})

	// Start server
	go func() {
		if err := e.Start(":8080"); err != nil && err != http.ErrServerClosed {
			e.Logger.Fatal("shutting down the server")
		}
	}()

	// Wait for interrupt signal to gracefully shutdown the server
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := e.Shutdown(ctx); err != nil {
		e.Logger.Fatal(err)
	}
}
