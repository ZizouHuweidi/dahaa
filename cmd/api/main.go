package main

import (
	"context"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/labstack/echo/v5"
	"github.com/zizouhuweidi/dahaa/internal/config"
	"github.com/zizouhuweidi/dahaa/internal/database"
	"github.com/zizouhuweidi/dahaa/internal/handler"
	"github.com/zizouhuweidi/dahaa/internal/repository/postgres"
	"github.com/zizouhuweidi/dahaa/internal/service"
	"github.com/zizouhuweidi/dahaa/internal/session"
	"github.com/zizouhuweidi/dahaa/internal/storage"
	"github.com/zizouhuweidi/dahaa/internal/websocket"
)

func main() {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}))
	slog.SetDefault(logger)
	cfg, err := config.Load()
	if err != nil {
		logger.Error("failed to load config", "error", err)
		os.Exit(1)
	}

	// Initialize database connection
	pool, err := database.ConnectPostgres(cfg.Database)
	if err != nil {
		logger.Error("failed to connect to database", "error", err)
		os.Exit(1)
	}
	defer pool.Close()

	// Initialize Redis client
	redisClient, err := database.ConnectRedis(cfg.Redis)
	if err != nil {
		logger.Error("failed to connect to redis", "error", err)
		os.Exit(1)
	}
	defer redisClient.Close()

	// Initialize image storage
	imageStorage, err := storage.NewImageStorage(cfg.Storage.ImageDir)
	if err != nil {
		logger.Error("failed to initialize image storage", "error", err)
		os.Exit(1)
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

	apiServer := handler.NewServer(handler.Deps{
		GameService:        gameService,
		QuestionRepo:       questionRepo,
		UserService:        userService,
		ImageStorage:       imageStorage,
		Hub:                hub,
		CORSAllowedOrigins: cfg.Server.CORSAllowedOrigins,
		Ready: func(c *echo.Context) error {
			ctx := c.Request().Context()
			if err := pool.Ping(ctx); err != nil {
				return c.JSON(http.StatusServiceUnavailable, map[string]string{"status": "not_ready", "database": "unavailable"})
			}
			if err := redisClient.Ping(ctx).Err(); err != nil {
				return c.JSON(http.StatusServiceUnavailable, map[string]string{"status": "not_ready", "redis": "unavailable"})
			}
			return c.JSON(http.StatusOK, map[string]string{"status": "ready"})
		},
	})

	e := apiServer.Echo()
	server := &http.Server{
		Addr:              ":" + cfg.Server.Port,
		Handler:           e,
		ReadHeaderTimeout: cfg.Server.ReadHeaderTimeout,
		ReadTimeout:       cfg.Server.ReadTimeout,
		WriteTimeout:      cfg.Server.WriteTimeout,
	}

	// Start server
	go func() {
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Error("server failed", "error", err)
			os.Exit(1)
		}
	}()
	logger.Info("server started", "addr", server.Addr)

	// Wait for interrupt signal to gracefully shutdown the server
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		logger.Error("server shutdown failed", "error", err)
		os.Exit(1)
	}
	logger.Info("server stopped")
}
