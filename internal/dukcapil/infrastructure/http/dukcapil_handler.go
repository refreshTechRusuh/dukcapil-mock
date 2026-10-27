// Package http (infrastructure) = adapter HTTP untuk dukcapil.
// Tugasnya: terjemahkan request HTTP <-> usecase. TIDAK ada logika bisnis di sini.
package http

import (
	"encoding/json"
	"net/http"

	"github.com/refreshTechRusuh/dukcapil-mock/internal/dukcapil/application"
	"github.com/refreshTechRusuh/dukcapil-mock/internal/dukcapil/domain"
	"github.com/refreshTechRusuh/dukcapil-mock/pkg/apinizer"
)

// Handler memegang usecase (lewat tipe konkret *application.Usecase; untuk
// service sekecil ini belum perlu interface inbound).
type Handler struct {
	uc *application.Usecase
}

// NewHandler menyuntik usecase saat wiring.
func NewHandler(uc *application.Usecase) *Handler {
	return &Handler{uc: uc}
}

// Verify menangani POST /internal/dukcapil/verify.
// Alur: decode body -> panggil usecase -> balas envelope {valid: bool}.
func (h *Handler) Verify(w http.ResponseWriter, r *http.Request) {
	var req domain.VerifyRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		apinizer.FailedBadRequest(w, "body tidak valid")
		return
	}

	result := h.uc.Verify(req)

	// valid:true maupun valid:false sama-sama 200 (hasil pengecekan sukses).
	apinizer.Success(w, "OK", result)
}
