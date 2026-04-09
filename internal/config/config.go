package config

import (
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/joho/godotenv"
)

// Config holds all application configuration loaded from environment variables.
type Config struct {
	App      AppConfig
	DB       DBConfig
	Redis    RedisConfig
	JWT      JWTConfig
}

// AppConfig holds application-level settings.
type AppConfig struct {
	Env  string
	Port string
	Name string
}

// DBConfig holds PostgreSQL connection settings.
type DBConfig struct {
	Host     string
	Port     string
	User     string
	Password string
	Name     string
	SSLMode  string
	TimeZone string
}

// RedisConfig holds Redis connection settings.
type RedisConfig struct {
	Host     string
	Port     string
	Password string
	DB       int
}

// JWTConfig holds JWT settings.
type JWTConfig struct {
	Secret          string
	AccessTokenTTL  time.Duration
	RefreshTokenTTL time.Duration
}

// Load reads environment variables (from .env if present) and returns Config.
func Load() (*Config, error) {
	// Load .env if present; ignore error so it works without the file in production.
	_ = godotenv.Load()

	redisDB, _ := strconv.Atoi(getEnv("REDIS_DB", "0"))
	accessTTL, err := time.ParseDuration(getEnv("JWT_ACCESS_TTL", "15m"))
	if err != nil {
		return nil, fmt.Errorf("invalid JWT_ACCESS_TTL: %w", err)
	}
	refreshTTL, err := time.ParseDuration(getEnv("JWT_REFRESH_TTL", "168h"))
	if err != nil {
		return nil, fmt.Errorf("invalid JWT_REFRESH_TTL: %w", err)
	}

	cfg := &Config{
		App: AppConfig{
			Env:  getEnv("APP_ENV", "development"),
			Port: getEnv("APP_PORT", "8080"),
			Name: getEnv("APP_NAME", "task-api"),
		},
		DB: DBConfig{
			Host:     getEnv("DB_HOST", "localhost"),
			Port:     getEnv("DB_PORT", "5432"),
			User:     getEnv("DB_USER", "postgres"),
			Password: getEnv("DB_PASSWORD", "postgres"),
			Name:     getEnv("DB_NAME", "task_db"),
			SSLMode:  getEnv("DB_SSLMODE", "disable"),
			TimeZone: getEnv("DB_TIMEZONE", "UTC"),
		},
		Redis: RedisConfig{
			Host:     getEnv("REDIS_HOST", "localhost"),
			Port:     getEnv("REDIS_PORT", "6379"),
			Password: getEnv("REDIS_PASSWORD", ""),
			DB:       redisDB,
		},
		JWT: JWTConfig{
			Secret:          getEnv("JWT_SECRET", "change-me-in-production"),
			AccessTokenTTL:  accessTTL,
			RefreshTokenTTL: refreshTTL,
		},
	}

	return cfg, nil
}

// DSN returns a PostgreSQL connection string.
func (d *DBConfig) DSN() string {
	return fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=%s TimeZone=%s",
		d.Host, d.Port, d.User, d.Password, d.Name, d.SSLMode, d.TimeZone,
	)
}

// RedisAddr returns the Redis address as "host:port".
func (r *RedisConfig) RedisAddr() string {
	return fmt.Sprintf("%s:%s", r.Host, r.Port)
}

func getEnv(key, defaultValue string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return defaultValue
}
