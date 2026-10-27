// Package entity berisi entity inti domain dukcapil (data murni, tanpa tag JSON).
// Mengikuti pola event-service: internal/{domain}/domain/entity/.
package entity

// Resident adalah satu record kependudukan TIRUAN yang disimpan Dukcapil mock.
// Verify mencocokkan input ke salah satu record ini (keempat field harus sama).
type Resident struct {
	NIK      string // 16 digit
	Nama     string
	TglLahir string // format YYYY-MM-DD
	NamaIbu  string // nama ibu kandung
}
