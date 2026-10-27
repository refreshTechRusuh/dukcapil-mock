// Package port berisi interface (kontrak) yang dibutuhkan/ditawarkan domain.
// "outbound" = kontrak yang DIBUTUHKAN usecase dari dunia luar (mis. penyimpanan data).
// Usecase bergantung ke interface ini, BUKAN ke implementasi konkret —
// persis pola usecase->interface di latihan app-db-mu (biar mudah dites/ditukar).
package port

import "github.com/refreshTechRusuh/dukcapil-mock/internal/dukcapil/domain"

// ResidentRepository = sumber data penduduk tiruan.
// Implementasinya nanti: seed data in-memory (produksi mock) + versi mock (test).
type ResidentRepository interface {
	// FindByNIK mengembalikan record untuk sebuah NIK.
	// found=false bila NIK tidak ada (bukan error) — mirip pola comma-ok di map.
	FindByNIK(nik string) (resident domain.Resident, found bool)
}
