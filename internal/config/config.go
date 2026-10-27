// Package config memuat konfigurasi service dari environment variable.
// Versi mini untuk dukcapil-mock: cuma App + Log + DB (tak perlu JWT/Kafka).
// Pola mengikuti service lain (event-service) agar konsisten lintas tim.
package config

import (
	"fmt"
	"log/slog"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

// Config = seluruh konfigurasi service, dikelompokkan per bagian.
type Config struct {
	App AppConfig
	Log LogConfig
	DB  DBConfig
}

type AppConfig struct {
	Env  string
	Port int
}

type LogConfig struct {
	Level  string
	Format string
}

type DBConfig struct {
	Host            string
	Port            int
	User            string
	Password        string
	Name            string
	SSLMode         string
	MaxOpenConns    int
	MaxIdleConns    int
	ConnTimeoutSecs int
}

// Load membaca .env (bila ada) lalu environment, dengan default yang aman
// untuk dev lokal (nyambung ke postgres di docker-compose tools).
func Load() (*Config, error) {
	_ = godotenv.Load() // .env opsional; abaikan error bila file tak ada

	cfg := &Config{
		App: AppConfig{
			Env:  getEnv("APP_ENV", "local"),
			Port: getEnvInt("APP_PORT", 8085), // 8085 = port dukcapil (sesuai tools)
		},
		Log: LogConfig{
			Level:  getEnv("LOG_LEVEL", "info"),
			Format: getEnv("LOG_FORMAT", "json"),
		},
		DB: DBConfig{
			Host:            getEnv("DB_HOST", "localhost"),
			Port:            getEnvInt("DB_PORT", 5432),
			User:            getEnv("DB_USER", "postgres"),
			Password:        getEnv("DB_PASSWORD", "password"),
			Name:            getEnv("DB_NAME", "rusuh_booking"),
			SSLMode:         getEnv("DB_SSL_MODE", "disable"),
			MaxOpenConns:    getEnvInt("DB_MAX_OPEN_CONNS", 10),
			MaxIdleConns:    getEnvInt("DB_MAX_IDLE_CONNS", 5),
			ConnTimeoutSecs: getEnvInt("DB_CONN_TIMEOUT_SECONDS", 5),
		},
	}
	return cfg, nil
}

// DSN membentuk connection string PostgreSQL dari DBConfig.
func (d DBConfig) DSN() string {
	return fmt.Sprintf(
		"host=%s port=%d user=%s password=%s dbname=%s sslmode=%s connect_timeout=%d",
		d.Host, d.Port, d.User, d.Password, d.Name, d.SSLMode, d.ConnTimeoutSecs,
	)
}

// SetupLogger memasang slog global sesuai LogConfig (json/text + level).
func SetupLogger(cfg LogConfig) {
	level := slog.LevelInfo
	switch cfg.Level {
	case "debug":
		level = slog.LevelDebug
	case "warn":
		level = slog.LevelWarn
	case "error":
		level = slog.LevelError
	}

	var h slog.Handler
	if cfg.Format == "text" {
		h = slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: level})
	} else {
		h = slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: level})
	}
	slog.SetDefault(slog.New(h))
}

// --- helper baca env dengan fallback ---

func getEnv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}

func getEnvInt(key string, fallback int) int {
	if v := os.Getenv(key); v != "" {
		if i, err := strconv.Atoi(v); err == nil {
			return i
		}
	}
	return fallback
}
