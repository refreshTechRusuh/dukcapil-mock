// Package delivery adalah composition root: tempat semua lapis dirakit
// (pool -> repo -> usecase -> handler) dan route didaftarkan ke mux.
// Hanya di sini (dan cmd/) kita boleh tahu semua lapis sekaligus.
package delivery

import (
	"net/http"

	"github.com/refreshTechRusuh/dukcapil-mock/internal/dukcapil/application/usecase"
	dukcapildb "github.com/refreshTechRusuh/dukcapil-mock/internal/dukcapil/infrastructure/db"
	dukcapilhttp "github.com/refreshTechRusuh/dukcapil-mock/internal/dukcapil/infrastructure/http"
	"github.com/refreshTechRusuh/dukcapil-mock/pkg/middleware"
)

// New merakit dependency dan mengembalikan handler siap pakai (mux dibungkus tracing).
// Alur wiring (luar <- dalam): handler <- usecase <- repo <- pool.
func New(pool dukcapildb.Querier) http.Handler {
	// 1. Repo (implementasi outbound port pakai DB)
	repo := dukcapildb.NewResidentRepository(pool)

	// 2. Usecase (logika) — disuntik repo lewat interface
	uc := usecase.NewDukcapilUsecase(repo)

	// 3. Handler (HTTP) — disuntik usecase lewat interface inbound
	handler := dukcapilhttp.NewHandler(uc)

	// 4. Router + daftar route
	mux := http.NewServeMux()
	mux.HandleFunc("POST /api/dukcapil/v1/verify", handler.Verify)
	mux.HandleFunc("GET /sys/ping", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"status":"00000","message":"dukcapil-mock ok","data":{}}`))
	})

	// Tracing: span per request, melanjutkan trace dari pemanggil (user-service) lewat
	// header traceparent. Health check dikecualikan.
	return middleware.OtelMiddleware("dukcapil-mock", "/sys/ping")(mux)
}
