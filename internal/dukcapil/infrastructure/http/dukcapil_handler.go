// Package http (infrastructure) = adapter HTTP untuk dukcapil.
// Tugasnya: terjemahkan request HTTP <-> usecase. TIDAK ada logika bisnis di sini.
// Bergantung ke inbound.IDukcapilUsecase (interface), bukan usecase konkret —
// mengikuti pola event-service.
package http

import (
	"encoding/json"
	"net/http"

	"github.com/refreshTechRusuh/dukcapil-mock/internal/dukcapil/domain/dto"
	"github.com/refreshTechRusuh/dukcapil-mock/internal/dukcapil/port/inbound"
	"github.com/refreshTechRusuh/dukcapil-mock/pkg/apinizer"
)

// Handler memegang usecase lewat interface inbound.
type Handler struct {
	uc inbound.IDukcapilUsecase
}

// NewHandler menyuntik usecase saat wiring.
func NewHandler(uc inbound.IDukcapilUsecase) *Handler {
	return &Handler{uc: uc}
}

// Verify menangani POST /internal/dukcapil/verify.
// Alur: decode body -> panggil usecase -> balas envelope {valid: bool}.
func (h *Handler) Verify(w http.ResponseWriter, r *http.Request) {
	var req dto.VerifyRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		apinizer.FailedBadRequest(w, "body tidak valid")
		return
	}

	result, err := h.uc.Verify(r.Context(), req)
	if err != nil {
		apinizer.WriteError(w, err)
		return
	}

	// valid:true maupun valid:false sama-sama 200 (hasil pengecekan sukses).
	apinizer.Success(w, "OK", result)
}
