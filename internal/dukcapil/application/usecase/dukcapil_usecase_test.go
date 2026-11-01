package usecase

import (
	"context"
	"errors"
	"testing"

	"github.com/refreshTechRusuh/dukcapil-mock/internal/dukcapil/domain/dto"
	"github.com/refreshTechRusuh/dukcapil-mock/internal/dukcapil/domain/entity"
)

// mockRepo = implementasi ketiga IResidentRepository (selain DB), untuk test.
// data = "DB" palsu di memori; errOnCall memaksa error teknis bila di-set.
type mockRepo struct {
	data      map[string]entity.Resident
	errOnCall error
}

func (m mockRepo) FindByNIK(_ context.Context, nik string) (entity.Resident, bool, error) {
	if m.errOnCall != nil {
		return entity.Resident{}, false, m.errOnCall
	}
	r, ok := m.data[nik]
	return r, ok, nil
}

func newSeededUsecase() *DukcapilUsecase {
	repo := mockRepo{data: map[string]entity.Resident{
		"3273000000000001": {
			NIK:      "3273000000000001",
			Nama:     "Budi Santoso",
			TglLahir: "1990-01-01",
			NamaIbu:  "Siti Aminah",
		},
	}}
	return NewDukcapilUsecase(repo)
}

func TestVerify(t *testing.T) {
	uc := newSeededUsecase()

	cases := []struct {
		name string
		req  dto.VerifyRequest
		want bool
	}{
		{"semua 4 field cocok -> valid",
			dto.VerifyRequest{NIK: "3273000000000001", Nama: "Budi Santoso", TglLahir: "1990-01-01", NamaIbu: "Siti Aminah"}, true},
		{"cocok walau beda huruf besar/kecil & spasi -> valid",
			dto.VerifyRequest{NIK: "3273000000000001", Nama: "budi santoso ", TglLahir: "1990-01-01", NamaIbu: " siti aminah"}, true},
		{"nik tidak ada -> invalid",
			dto.VerifyRequest{NIK: "9999999999999999", Nama: "Budi Santoso", TglLahir: "1990-01-01", NamaIbu: "Siti Aminah"}, false},
		{"nama salah -> invalid",
			dto.VerifyRequest{NIK: "3273000000000001", Nama: "SALAH", TglLahir: "1990-01-01", NamaIbu: "Siti Aminah"}, false},
		{"tgl lahir salah -> invalid",
			dto.VerifyRequest{NIK: "3273000000000001", Nama: "Budi Santoso", TglLahir: "2000-12-31", NamaIbu: "Siti Aminah"}, false},
		{"nama ibu salah -> invalid",
			dto.VerifyRequest{NIK: "3273000000000001", Nama: "Budi Santoso", TglLahir: "1990-01-01", NamaIbu: "SALAH"}, false},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			res, err := uc.Verify(context.Background(), c.req)
			if err != nil {
				t.Fatalf("Verify(%+v) error tak terduga: %v", c.req, err)
			}
			if res.Valid != c.want {
				t.Errorf("Verify(%+v).Valid = %v; mau %v", c.req, res.Valid, c.want)
			}
		})
	}
}

func TestVerify_DBError_ReturnsError(t *testing.T) {
	repo := mockRepo{errOnCall: errors.New("db mati")}
	uc := NewDukcapilUsecase(repo)

	_, err := uc.Verify(context.Background(), dto.VerifyRequest{NIK: "3273000000000001"})
	if err == nil {
		t.Fatal("mau error saat repo gagal teknis, dapat nil")
	}
}
