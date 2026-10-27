// Package delivery adalah composition root: tempat semua lapis dirakit
// (pool -> repo -> usecase -> handler) dan route didaftarkan ke mux.
// Hanya di sini (dan cmd/) kita boleh tahu semua lapis sekaligus.
package delivery

import (
	"net/http"

	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/refreshTechRusuh/dukcapil-mock/internal/dukcapil/application"
	dukcapildb "github.com/refreshTechRusuh/dukcapil-mock/internal/dukcapil/infrastructure/db"
	dukcapilhttp "github.com/refreshTechRusuh/dukcapil-mock/internal/dukcapil/infrastructure/http"
)

// New merakit dependency dan mengembalikan mux siap pakai.
// Alur wiring (luar <- dalam): handler <- usecase <- repo <- pool.
func New(pool *pgxpool.Pool) *http.ServeMux {
	// 1. Repo (implementasi interface pakai DB)
	repo := dukcapildb.NewResidentRepo(pool)

	// 2. Usecase (logika) — disuntik repo lewat interface
	uc := application.New(repo)

	// 3. Handler (HTTP) — disuntik usecase
	handler := dukcapilhttp.NewHandler(uc)

	// 4. Router + daftar route
	mux := http.NewServeMux()
	mux.HandleFunc("POST /internal/dukcapil/verify", handler.Verify)
	mux.HandleFunc("GET /sys/ping", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"status":"00000","message":"dukcapil-mock ok","data":{}}`))
	})

	return mux
}
