// Package domain berisi data & aturan murni untuk dukcapil-mock.
// Lapis paling dalam: TIDAK tahu HTTP, DB, atau framework apa pun.
// (Sama seperti folder model/ di latihan app-db-mu, tapi istilah hexagonal.)
package domain

// Resident adalah satu record kependudukan TIRUAN yang disimpan Dukcapil mock.
// Verify nanti mencocokkan input ke salah satu record ini (keempat field harus sama).
type Resident struct {
	NIK      string // 16 digit
	Nama     string
	TglLahir string // format YYYY-MM-DD
	NamaIbu  string // nama ibu kandung
}

// VerifyRequest = data yang dikirim pemanggil (User/Booking) untuk dicek.
// Tag json dipakai saat nanti di-decode dari body HTTP (lapis infrastructure).
type VerifyRequest struct {
	NIK      string `json:"nik"`
	Nama     string `json:"nama"`
	TglLahir string `json:"tglLahir"`
	NamaIbu  string `json:"namaIbu"`
}

// VerifyResult = hasil pengecekan. valid=true hanya bila keempat field cocok.
type VerifyResult struct {
	Valid bool `json:"valid"`
}
