// Package outbound berisi driven port: kontrak yang DIBUTUHKAN usecase dari
// dunia luar (penyimpanan data). Mengikuti pola event-service (prefix I,
// context.Context sebagai parameter pertama, return error).
package outbound

import (
	"context"

	"github.com/refreshTechRusuh/dukcapil-mock/internal/dukcapil/domain/entity"
)

// IResidentRepository = sumber data penduduk tiruan.
type IResidentRepository interface {
	// FindByNIK mencari satu penduduk berdasarkan NIK.
	//   - ketemu        -> (resident, nil)
	//   - NIK tidak ada -> (zero, nil)   [found=false]
	//   - error teknis  -> (zero, err)
	// found membedakan "tidak ada" dari "ada"; err untuk kegagalan teknis (DB mati).
	FindByNIK(ctx context.Context, nik string) (resident entity.Resident, found bool, err error)
}
