package config

import (
	"context"
	"log/slog"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestLoad_ShouldUseSafeDefaults_WhenEnvEmpty(t *testing.T) {
	cfg, err := Load()

	require.NoError(t, err)
	assert.Equal(t, 8085, cfg.App.Port)
	assert.Equal(t, "local", cfg.App.Env)
	assert.Equal(t, "info", cfg.Log.Level)
	assert.Equal(t, "json", cfg.Log.Format)
	assert.Equal(t, "localhost", cfg.DB.Host)
	assert.Equal(t, 5432, cfg.DB.Port)
	assert.Equal(t, 10, cfg.DB.MaxOpenConns)
}

func TestLoad_ShouldReadEnvOverrides(t *testing.T) {
	t.Setenv("APP_PORT", "9090")
	t.Setenv("DB_HOST", "10.0.0.5")
	t.Setenv("DB_SSL_MODE", "require")
	t.Setenv("LOG_LEVEL", "debug")

	cfg, err := Load()

	require.NoError(t, err)
	assert.Equal(t, 9090, cfg.App.Port)
	assert.Equal(t, "10.0.0.5", cfg.DB.Host)
	assert.Equal(t, "require", cfg.DB.SSLMode)
	assert.Equal(t, "debug", cfg.Log.Level)
}

func TestGetEnvInt_ShouldFallBack_WhenNotANumber(t *testing.T) {
	t.Setenv("RUSAK", "abc")
	t.Setenv("BAGUS", "42")

	assert.Equal(t, 7, getEnvInt("RUSAK", 7))
	assert.Equal(t, 42, getEnvInt("BAGUS", 7))
	assert.Equal(t, 7, getEnvInt("TIDAK_ADA", 7))
}

func TestDSN_ShouldContainAllConnectionParts(t *testing.T) {
	d := DBConfig{Host: "h", Port: 5432, User: "u", Password: "p", Name: "n", SSLMode: "require", ConnTimeoutSecs: 3}

	assert.Equal(t, "host=h port=5432 user=u password=p dbname=n sslmode=require connect_timeout=3", d.DSN())
}

func TestSetupLogger_ShouldApplyLevelAndFormat(t *testing.T) {
	prev := slog.Default()
	t.Cleanup(func() { slog.SetDefault(prev) })
	ctx := context.Background()

	for level, minLevel := range map[string]slog.Level{
		"debug": slog.LevelDebug, "info": slog.LevelInfo, "warn": slog.LevelWarn,
		"error": slog.LevelError, "xxx": slog.LevelInfo,
	} {
		for _, format := range []string{"json", "text"} {
			SetupLogger(LogConfig{Level: level, Format: format})
			for _, l := range []slog.Level{slog.LevelDebug, slog.LevelInfo, slog.LevelWarn, slog.LevelError} {
				assert.Equal(t, l >= minLevel, slog.Default().Enabled(ctx, l), "level=%s format=%s cek=%v", level, format, l)
			}
		}
	}
}
