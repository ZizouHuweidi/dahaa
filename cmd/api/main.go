package main

import (
	"context"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"path/filepath"
	"syscall"
	"time"

	"github.com/redis/go-redis/v9"
	"github.com/zizouhuweidi/dahaa/internal/httpapi"
	"github.com/zizouhuweidi/dahaa/internal/repository/postgres"
	"github.com/zizouhuweidi/dahaa/internal/service"
	"github.com/zizouhuweidi/dahaa/internal/session"
	"github.com/zizouhuweidi/dahaa/internal/storage"
	"github.com/zizouhuweidi/dahaa/internal/websocket"
)

func main() {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}))
	slog.SetDefault(logger)

	// Initialize database connection
	pool, err := postgres.NewDB()
	if err != nil {
		logger.Error("failed to connect to database", "error", err)
		os.Exit(1)
	}
	defer pool.Close()

	// Initialize Redis client
	redisHost := getEnv("REDIS_HOST", "localhost")
	redisPort := getEnv("REDIS_PORT", "6379")
	redisPassword := getEnv("REDIS_PASSWORD", "")

	redisClient := redis.NewClient(&redis.Options{
		Addr:     redisHost + ":" + redisPort,
		Password: redisPassword,
		DB:       0, // use default DB
	})

	// Test Redis connection
	ctx := context.Background()
	if err := redisClient.Ping(ctx).Err(); err != nil {
		logger.Error("failed to connect to redis", "error", err)
		os.Exit(1)
	}

	// Initialize image storage
	imageStorage, err := storage.NewImageStorage(filepath.Join("uploads", "images"))
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

	apiServer := httpapi.NewServer(httpapi.Deps{
		GameService:  gameService,
		QuestionRepo: questionRepo,
		UserService:  userService,
		ImageStorage: imageStorage,
		Hub:          hub,
	})

	server := &http.Server{
		Addr:              ":8080",
		Handler:           apiServer.Handler(),
		ReadHeaderTimeout: 5 * time.Second,
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

// getEnv gets an environment variable or returns a default value
func getEnv(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}
