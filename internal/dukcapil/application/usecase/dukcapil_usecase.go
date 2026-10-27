// Package usecase berisi logika bisnis dukcapil-mock.
// Lapis ini TAHU domain + port, tapi TIDAK tahu HTTP/DB (hexagonal).
// Mengikuti pola event-service: application/usecase/.
package usecase

import (
	"context"
	"strings"

	"github.com/refreshTechRusuh/dukcapil-mock/internal/dukcapil/domain/dto"
	"github.com/refreshTechRusuh/dukcapil-mock/internal/dukcapil/port/outbound"
	"github.com/refreshTechRusuh/dukcapil-mock/pkg/apperror"
)

// DukcapilUsecase memegang dependency (lewat interface outbound, bukan konkret).
type DukcapilUsecase struct {
	repo outbound.IResidentRepository
}

// NewDukcapilUsecase menyuntik repository saat wiring (di delivery).
func NewDukcapilUsecase(repo outbound.IResidentRepository) *DukcapilUsecase {
	return &DukcapilUsecase{repo: repo}
}

// Verify = inti Dukcapil mock. valid=true HANYA bila NIK ada DAN keempat field
// (nik, nama, tglLahir, namaIbu) cocok dengan record tersimpan.
// Error teknis dari repo dibungkus apperror.Internal (dipetakan 500 di HTTP).
func (u *DukcapilUsecase) Verify(ctx context.Context, req dto.VerifyRequest) (dto.VerifyResult, error) {
	resident, found, err := u.repo.FindByNIK(ctx, req.NIK)
	if err != nil {
		return dto.VerifyResult{}, apperror.Internal(err)
	}
	if !found {
		return dto.VerifyResult{Valid: false}, nil
	}

	valid := matchField(resident.Nama, req.Nama) &&
		matchField(resident.TglLahir, req.TglLahir) &&
		matchField(resident.NamaIbu, req.NamaIbu)

	return dto.VerifyResult{Valid: valid}, nil
}

// matchField membandingkan dua field secara toleran: abaikan spasi di ujung &
// beda huruf besar/kecil (mis. "Siti" == "siti ") supaya mock tidak rewel.
func matchField(stored, input string) bool {
	return strings.EqualFold(strings.TrimSpace(stored), strings.TrimSpace(input))
}
