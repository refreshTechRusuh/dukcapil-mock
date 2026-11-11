// Package main = entry point dukcapil-mock.
// Tugas main: startup (load config, konek DB, rakit server) + graceful shutdown.
// Semua wiring ada di internal/delivery; semua logika di internal/dukcapil.
package main

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/refreshTechRusuh/dukcapil-mock/internal/config"
	"github.com/refreshTechRusuh/dukcapil-mock/internal/delivery"
	pkgdb "github.com/refreshTechRusuh/dukcapil-mock/pkg/db"
	pkgotel "github.com/refreshTechRusuh/dukcapil-mock/pkg/otel"
)

func main() {
	// 1. Konfigurasi dari env (+ .env bila ada).
	cfg, err := config.Load()
	if err != nil {
		slog.Error("gagal load config", "err", err)
		os.Exit(1)
	}
	config.SetupLogger(cfg.Log)

	// rootCtx: dibatalkan saat service mau mati, supaya operasi panjang ikut berhenti.
	rootCtx, cancelRoot := context.WithCancel(context.Background())
	defer cancelRoot()

	// OpenTelemetry -> GCP Cloud Trace. OTEL_ENABLED=false mematikan (lokal/test).
	// Gagal setup TIDAK menggagalkan service: tracing hanya observabilitas.
	if otelShutdown, err := pkgotel.SetupTracerProvider(rootCtx, "dukcapil-mock", "1.0.0"); err != nil {
		slog.Warn("otel: gagal setup tracer provider, tracing nonaktif", "err", err)
	} else {
		defer func() {
			shutCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
			defer cancel()
			_ = otelShutdown(shutCtx)
		}()
	}

	// 2. Koneksi DB (pool). Gagal konek = gagal start (fail fast).
	pool, err := pkgdb.NewPool(rootCtx, cfg.DB)
	if err != nil {
		slog.Error("gagal konek postgres", "err", err)
		os.Exit(1)
	}
	defer pool.Close()

	// 3. Rakit server (wiring semua lapis) + pasang ke http.Server.
	mux := delivery.New(pool)
	srv := &http.Server{
		Addr:              fmt.Sprintf(":%d", cfg.App.Port),
		Handler:           mux,
		ReadHeaderTimeout: 5 * time.Second,
	}

	// 4. Jalankan di goroutine supaya main bisa menunggu sinyal berhenti.
	go func() {
		slog.Info("dukcapil-mock listening", "port", cfg.App.Port)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			slog.Error("server error", "err", err)
			os.Exit(1)
		}
	}()

	// 5. Tunggu Ctrl+C / SIGTERM, lalu matikan rapi (graceful shutdown).
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	slog.Info("shutting down...")
	cancelRoot()
	shutCtx, cancelShut := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancelShut()
	if err := srv.Shutdown(shutCtx); err != nil {
		slog.Error("shutdown error", "err", err)
	}
	slog.Info("dukcapil-mock stopped")
}
