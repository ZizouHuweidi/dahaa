package config

import (
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"
)

type Config struct {
	Environment string
	Server      ServerConfig
	Database    DatabaseConfig
	Redis       RedisConfig
	Storage     StorageConfig
}

type ServerConfig struct {
	Port               string
	ReadHeaderTimeout  time.Duration
	ReadTimeout        time.Duration
	WriteTimeout       time.Duration
	CORSAllowedOrigins []string
}

type DatabaseConfig struct {
	URL               string
	MaxConns          int32
	MinConns          int32
	MaxConnLifetime   time.Duration
	MaxConnIdleTime   time.Duration
	HealthCheckPeriod time.Duration
}

type RedisConfig struct {
	Host     string
	Port     string
	Password string
	DB       int
	Timeout  time.Duration
}

type StorageConfig struct {
	ImageDir string
}

func Load() (*Config, error) {
	readHeaderTimeout, err := getDurationEnv("SERVER_READ_HEADER_TIMEOUT", 5*time.Second)
	if err != nil {
		return nil, err
	}
	readTimeout, err := getDurationEnv("SERVER_READ_TIMEOUT", 30*time.Second)
	if err != nil {
		return nil, err
	}
	writeTimeout, err := getDurationEnv("SERVER_WRITE_TIMEOUT", 30*time.Second)
	if err != nil {
		return nil, err
	}
	dbMaxConns, err := getInt32Env("DB_MAX_CONNS", 25)
	if err != nil {
		return nil, err
	}
	dbMinConns, err := getInt32Env("DB_MIN_CONNS", 5)
	if err != nil {
		return nil, err
	}
	dbMaxConnLifetime, err := getDurationEnv("DB_MAX_CONN_LIFETIME", time.Hour)
	if err != nil {
		return nil, err
	}
	dbMaxConnIdleTime, err := getDurationEnv("DB_MAX_CONN_IDLE_TIME", 30*time.Minute)
	if err != nil {
		return nil, err
	}
	dbHealthCheckPeriod, err := getDurationEnv("DB_HEALTH_CHECK_PERIOD", time.Minute)
	if err != nil {
		return nil, err
	}
	redisDB, err := getIntEnv("REDIS_DB", 0)
	if err != nil {
		return nil, err
	}
	redisTimeout, err := getDurationEnv("REDIS_TIMEOUT", 10*time.Second)
	if err != nil {
		return nil, err
	}

	return &Config{
		Environment: getEnv("ENVIRONMENT", "development"),
		Server: ServerConfig{
			Port:               getEnv("PORT", "8080"),
			ReadHeaderTimeout:  readHeaderTimeout,
			ReadTimeout:        readTimeout,
			WriteTimeout:       writeTimeout,
			CORSAllowedOrigins: getEnvSlice("CORS_ALLOWED_ORIGINS", []string{"http://localhost:3000", "http://127.0.0.1:3000", "http://localhost:5173", "http://127.0.0.1:5173"}),
		},
		Database: DatabaseConfig{
			URL:               databaseURL(),
			MaxConns:          dbMaxConns,
			MinConns:          dbMinConns,
			MaxConnLifetime:   dbMaxConnLifetime,
			MaxConnIdleTime:   dbMaxConnIdleTime,
			HealthCheckPeriod: dbHealthCheckPeriod,
		},
		Redis: RedisConfig{
			Host:     getEnv("REDIS_HOST", "localhost"),
			Port:     getEnv("REDIS_PORT", "6379"),
			Password: getEnv("REDIS_PASSWORD", ""),
			DB:       redisDB,
			Timeout:  redisTimeout,
		},
		Storage: StorageConfig{
			ImageDir: getEnv("IMAGE_UPLOAD_DIR", "uploads/images"),
		},
	}, nil
}

func databaseURL() string {
	if value := os.Getenv("DATABASE_URL"); value != "" {
		return value
	}
	return fmt.Sprintf(
		"postgres://%s:%s@%s:%s/%s?sslmode=disable",
		getEnv("POSTGRES_USER", "postgres"),
		getEnv("POSTGRES_PASSWORD", "postgres"),
		getEnv("POSTGRES_HOST", "localhost"),
		getEnv("POSTGRES_PORT", "5432"),
		getEnv("POSTGRES_DB", "dahaa"),
	)
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func getIntEnv(key string, defaultValue int) (int, error) {
	if value := os.Getenv(key); value != "" {
		parsed, err := strconv.Atoi(value)
		if err != nil {
			return 0, fmt.Errorf("invalid %s: %w", key, err)
		}
		return parsed, nil
	}
	return defaultValue, nil
}

func getInt32Env(key string, defaultValue int32) (int32, error) {
	if value := os.Getenv(key); value != "" {
		parsed, err := strconv.ParseInt(value, 10, 32)
		if err != nil {
			return 0, fmt.Errorf("invalid %s: %w", key, err)
		}
		return int32(parsed), nil
	}
	return defaultValue, nil
}

func getDurationEnv(key string, defaultValue time.Duration) (time.Duration, error) {
	if value := os.Getenv(key); value != "" {
		parsed, err := time.ParseDuration(value)
		if err != nil {
			return 0, fmt.Errorf("invalid %s: %w", key, err)
		}
		return parsed, nil
	}
	return defaultValue, nil
}

func getEnvSlice(key string, defaultValue []string) []string {
	if value := os.Getenv(key); value != "" {
		parts := strings.Split(value, ",")
		result := make([]string, 0, len(parts))
		for _, part := range parts {
			part = strings.TrimSpace(part)
			if part != "" {
				result = append(result, part)
			}
		}
		return result
	}
	return defaultValue
}
