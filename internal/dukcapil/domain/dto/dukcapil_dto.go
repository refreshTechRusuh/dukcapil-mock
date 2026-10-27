// Package dto berisi bentuk data request/response dukcapil (dengan tag JSON).
// Dipisah dari entity mengikuti pola event-service: domain/entity vs domain/dto.
package dto

// VerifyRequest = data yang dikirim pemanggil (User/Booking) untuk dicek.
type VerifyRequest struct {
	NIK      string `json:"nik"      validate:"required,len=16,numeric"`
	Nama     string `json:"nama"     validate:"required"`
	TglLahir string `json:"tglLahir" validate:"required"`
	NamaIbu  string `json:"namaIbu"  validate:"required"`
}

// VerifyResult = hasil pengecekan. valid=true hanya bila keempat field cocok.
type VerifyResult struct {
	Valid bool `json:"valid"`
}
