package db

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/refreshTechRusuh/dukcapil-mock/internal/config"
)

// Jalur sukses butuh Postgres asli (di luar unit test). Yang diuji di sini: jalur
// gagal, supaya service gagal start dengan pesan jelas (fail fast), bukan menggantung.

func TestNewPool_ShouldFail_WhenConfigInvalid(t *testing.T) {
	cfg := config.DBConfig{Host: "localhost", Port: 5432, User: "u", Password: "p", Name: "n",
		SSLMode: "mode-ngawur", MaxOpenConns: 2, ConnTimeoutSecs: 1}

	pool, err := NewPool(context.Background(), cfg)

	require.Error(t, err)
	assert.Nil(t, pool)
	assert.Contains(t, err.Error(), "db: parse config")
}

func TestNewPool_ShouldFailFast_WhenDatabaseUnreachable(t *testing.T) {
	cfg := config.DBConfig{Host: "127.0.0.1", Port: 1, User: "u", Password: "p", Name: "n",
		SSLMode: "disable", MaxOpenConns: 2, MaxIdleConns: 0, ConnTimeoutSecs: 1}
	start := time.Now()

	pool, err := NewPool(context.Background(), cfg)

	require.Error(t, err)
	assert.Nil(t, pool, "pool tidak boleh dikembalikan bila ping gagal")
	assert.Contains(t, err.Error(), "db: ping")
	assert.Less(t, time.Since(start), 10*time.Second, "harus gagal cepat, bukan menggantung")
}

func TestNewPool_ShouldFail_WhenContextAlreadyCancelled(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	cfg := config.DBConfig{Host: "127.0.0.1", Port: 1, User: "u", Password: "p", Name: "n",
		SSLMode: "disable", MaxOpenConns: 2, ConnTimeoutSecs: 1}

	pool, err := NewPool(ctx, cfg)

	assert.Error(t, err)
	assert.Nil(t, pool)
}
