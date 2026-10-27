// Package db (infrastructure) = adapter yang mengimplementasikan port.ResidentRepository
// memakai PostgreSQL (pgxpool). Ini satu-satunya tempat yang tahu soal SQL;
// usecase tetap bersih (cuma kenal interface).
package db

import (
	"context"
	"errors"
	"log/slog"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/refreshTechRusuh/dukcapil-mock/internal/dukcapil/domain"
)

// ResidentRepo adalah implementasi ResidentRepository berbasis PostgreSQL.
type ResidentRepo struct {
	pool *pgxpool.Pool
}

// NewResidentRepo menyuntik pool koneksi (dibuat di main/delivery).
func NewResidentRepo(pool *pgxpool.Pool) *ResidentRepo {
	return &ResidentRepo{pool: pool}
}

// FindByNIK mencari satu penduduk berdasarkan NIK.
//   - ketemu            -> (resident, true)
//   - NIK tidak ada     -> (zero, false)   [bukan error]
//   - error teknis DB   -> (zero, false) + log  [DB mati, dsb]
//
// Catatan: interface port tidak mengembalikan error (pola comma-ok), jadi
// kegagalan teknis diperlakukan sebagai "tidak ketemu" + dicatat ke log.
func (r *ResidentRepo) FindByNIK(nik string) (domain.Resident, bool) {
	const q = `
		SELECT nik, nama, to_char(tgl_lahir, 'YYYY-MM-DD'), nama_ibu
		FROM residents
		WHERE nik = $1`

	var res domain.Resident
	err := r.pool.QueryRow(context.Background(), q, nik).
		Scan(&res.NIK, &res.Nama, &res.TglLahir, &res.NamaIbu)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			// NIK tidak terdaftar -> bukan error, cukup "tidak ketemu".
			return domain.Resident{}, false
		}
		// Error teknis (DB mati, query salah) -> log, anggap tidak ketemu.
		slog.Error("dukcapil: query residents gagal", "err", err)
		return domain.Resident{}, false
	}

	return res, true
}
