// Package application berisi usecase (logika bisnis) dukcapil-mock.
// Lapis ini TAHU domain + port, tapi TIDAK tahu HTTP (net/http tak di-import di sini).
package application

import (
	"strings"

	"github.com/refreshTechRusuh/dukcapil-mock/internal/dukcapil/domain"
	"github.com/refreshTechRusuh/dukcapil-mock/internal/dukcapil/port"
)

// Usecase memegang dependency yang dibutuhkan (lewat interface, bukan konkret).
type Usecase struct {
	repo port.ResidentRepository
}

// New menyuntik repository saat wiring (di delivery/main nanti).
// Mengembalikan *Usecase; pemanggil cukup tahu method Verify.
func New(repo port.ResidentRepository) *Usecase {
	return &Usecase{repo: repo}
}

// Verify = inti Dukcapil mock. valid=true HANYA bila NIK ada DAN keempat field
// (nik, nama, tglLahir, namaIbu) cocok dengan record tersimpan.
func (u *Usecase) Verify(req domain.VerifyRequest) domain.VerifyResult {
	resident, found := u.repo.FindByNIK(req.NIK)
	if !found {
		return domain.VerifyResult{Valid: false}
	}

	valid := matchField(resident.Nama, req.Nama) &&
		matchField(resident.TglLahir, req.TglLahir) &&
		matchField(resident.NamaIbu, req.NamaIbu)

	return domain.VerifyResult{Valid: valid}
}

// matchField membandingkan dua field secara toleran: abaikan spasi di ujung &
// beda huruf besar/kecil (mis. "Siti" == "siti ") supaya mock tidak rewel.
func matchField(stored, input string) bool {
	return strings.EqualFold(strings.TrimSpace(stored), strings.TrimSpace(input))
}
