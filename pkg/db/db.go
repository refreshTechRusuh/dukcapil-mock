// Package db menyediakan koneksi PostgreSQL berupa connection pool (pgxpool).
// Pola mengikuti service lain (event-service) agar konsisten lintas tim.
package db

import (
	"context"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/refreshTechRusuh/dukcapil-mock/internal/config"
)

// NewPool membuat pool koneksi PostgreSQL dari DBConfig, lalu memastikan
// koneksi benar-benar hidup lewat Ping. Mengembalikan pool siap pakai.
func NewPool(ctx context.Context, cfg config.DBConfig) (*pgxpool.Pool, error) {
	// Parse DSN -> konfigurasi pool.
	poolCfg, err := pgxpool.ParseConfig(cfg.DSN())
	if err != nil {
		return nil, fmt.Errorf("db: parse config: %w", err)
	}

	// Tuning pool: berapa koneksi boleh dibuka, berapa lama hidup/idle.
	poolCfg.MaxConns = toInt32(cfg.MaxOpenConns)
	poolCfg.MinConns = toInt32(cfg.MaxIdleConns)
	poolCfg.MaxConnLifetime = time.Hour
	poolCfg.MaxConnIdleTime = 30 * time.Minute
	poolCfg.HealthCheckPeriod = time.Minute

	// Buat pool. Seperti database/sql: belum tentu benar-benar konek di sini.
	pool, err := pgxpool.NewWithConfig(ctx, poolCfg)
	if err != nil {
		return nil, fmt.Errorf("db: connect: %w", err)
	}

	// Ping = paksa satu koneksi nyata untuk membuktikan DB bisa dihubungi.
	// (Sama idenya dengan db.Ping() di latihan database/sql-mu.)
	pingCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()
	if err := pool.Ping(pingCtx); err != nil {
		return nil, fmt.Errorf("db: ping: %w", err)
	}

	return pool, nil
}
