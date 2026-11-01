// Package inbound berisi driving port: kontrak yang DITAWARKAN domain ke
// pemanggil (HTTP handler). Handler bergantung ke interface ini, bukan ke
// usecase konkret — mengikuti pola event-service (IEventUsecase).
package inbound

import (
	"context"

	"github.com/refreshTechRusuh/dukcapil-mock/internal/dukcapil/domain/dto"
)

// IDukcapilUsecase adalah kontrak yang dipanggil HTTP handler.
type IDukcapilUsecase interface {
	// Verify mengembalikan valid=true hanya bila keempat field request cocok
	// dengan satu record tersimpan.
	Verify(ctx context.Context, req dto.VerifyRequest) (dto.VerifyResult, error)
}
