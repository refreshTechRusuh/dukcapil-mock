// Package apinizer menyediakan envelope response standar lintas service:
//   { "status", "message", "data" }  (+ "partial" opsional)
// status "00000" = sukses; selain itu = string kode HTTP ("400","500",...).
// Bentuk ini identik dengan service lain (event/booking/payment) agar klien
// memperlakukan semua response seragam.
//
// Dukcapil-mock hanya butuh sedikit: Success (200) + FailedBadRequest (400,
// untuk body JSON rusak). verify yang menghasilkan valid:false TETAP memakai
// Success — karena "tidak valid" itu DATA, bukan error HTTP.
package apinizer

import (
	"encoding/json"
	"log/slog"
	"net/http"
)

// Envelope adalah bentuk response tunggal untuk setiap endpoint.
type Envelope struct {
	Status  string `json:"status"`
	Message string `json:"message"`
	Data    any    `json:"data"`
	Partial bool   `json:"partial,omitempty"`
}

// SuccessStatus = kode status aplikasi untuk response sukses.
const SuccessStatus = "00000"

// Success menulis response 200 OK dengan payload data.
func Success(w http.ResponseWriter, message string, data any) {
	writeJSON(w, http.StatusOK, Envelope{
		Status:  SuccessStatus,
		Message: message,
		Data:    data,
	})
}

// FailedBadRequest menulis response 400 (mis. body JSON tidak bisa di-parse).
// data diisi objek kosong agar bentuk tetap konsisten.
func FailedBadRequest(w http.ResponseWriter, message string) {
	writeJSON(w, http.StatusBadRequest, Envelope{
		Status:  "400",
		Message: message,
		Data:    struct{}{},
	})
}

func writeJSON(w http.ResponseWriter, status int, body any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if err := json.NewEncoder(w).Encode(body); err != nil {
		slog.Error("apinizer: gagal encode response", "err", err)
	}
}
