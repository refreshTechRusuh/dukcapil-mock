package application

import (
	"testing"

	"github.com/refreshTechRusuh/dukcapil-mock/internal/dukcapil/domain"
)

// mockRepo = implementasi ketiga dari interface ResidentRepository, khusus test.
// Isinya satu record tetap; kalau NIK cocok -> kembalikan record, else found=false.
// (Sama seperti mock repository di latihan app-db/usecase mu.)
type mockRepo struct {
	stored domain.Resident
}

func (m mockRepo) FindByNIK(nik string) (domain.Resident, bool) {
	if nik == m.stored.NIK {
		return m.stored, true
	}
	return domain.Resident{}, false
}

func TestVerify(t *testing.T) {
	// Satu record acuan yang "tersimpan" di Dukcapil.
	stored := domain.Resident{
		NIK:      "3273000000000001",
		Nama:     "Budi",
		TglLahir: "1990-01-01",
		NamaIbu:  "Siti",
	}
	uc := New(mockRepo{stored: stored})

	tests := []struct {
		nama  string
		input domain.VerifyRequest
		valid bool
	}{
		{
			nama:  "semua 4 field cocok -> valid",
			input: domain.VerifyRequest{NIK: "3273000000000001", Nama: "Budi", TglLahir: "1990-01-01", NamaIbu: "Siti"},
			valid: true,
		},
		{
			nama:  "cocok walau beda huruf besar/kecil & spasi -> valid",
			input: domain.VerifyRequest{NIK: "3273000000000001", Nama: "budi", TglLahir: "1990-01-01", NamaIbu: " siti "},
			valid: true,
		},
		{
			nama:  "nik tidak ada -> invalid",
			input: domain.VerifyRequest{NIK: "9999999999999999", Nama: "Budi", TglLahir: "1990-01-01", NamaIbu: "Siti"},
			valid: false,
		},
		{
			nama:  "nama salah -> invalid",
			input: domain.VerifyRequest{NIK: "3273000000000001", Nama: "Andi", TglLahir: "1990-01-01", NamaIbu: "Siti"},
			valid: false,
		},
		{
			nama:  "tgl lahir salah -> invalid",
			input: domain.VerifyRequest{NIK: "3273000000000001", Nama: "Budi", TglLahir: "1991-12-31", NamaIbu: "Siti"},
			valid: false,
		},
		{
			nama:  "nama ibu salah -> invalid",
			input: domain.VerifyRequest{NIK: "3273000000000001", Nama: "Budi", TglLahir: "1990-01-01", NamaIbu: "Rina"},
			valid: false,
		},
	}

	for _, tc := range tests {
		t.Run(tc.nama, func(t *testing.T) {
			got := uc.Verify(tc.input)
			if got.Valid != tc.valid {
				t.Errorf("Verify(%+v).Valid = %v; mau %v", tc.input, got.Valid, tc.valid)
			}
		})
	}
}
