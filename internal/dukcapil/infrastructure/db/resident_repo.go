// Package db (infrastructure) mengimplementasikan outbound.IResidentRepository
// memakai PostgreSQL (pgxpool). Satu-satunya tempat yang tahu SQL.
// Pola mengikuti event-service: pakai context, kembalikan error, bungkus dengan %w.
package db

import (
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/refreshTechRusuh/dukcapil-mock/internal/dukcapil/domain/entity"
)

// ResidentRepository adalah implementasi IResidentRepository berbasis PostgreSQL.
type ResidentRepository struct {
	pool *pgxpool.Pool
}

// NewResidentRepository menyuntik pool koneksi (dibuat di delivery).
func NewResidentRepository(pool *pgxpool.Pool) *ResidentRepository {
	return &ResidentRepository{pool: pool}
}

// FindByNIK mencari satu penduduk berdasarkan NIK.
//   - ketemu        -> (resident, true, nil)
//   - NIK tidak ada -> (zero, false, nil)
//   - error teknis  -> (zero, false, err)   [DB mati, query gagal]
func (r *ResidentRepository) FindByNIK(ctx context.Context, nik string) (entity.Resident, bool, error) {
	const q = `
		SELECT nik, nama, to_char(tgl_lahir, 'YYYY-MM-DD'), nama_ibu
		FROM residents
		WHERE nik = $1`

	var res entity.Resident
	err := r.pool.QueryRow(ctx, q, nik).
		Scan(&res.NIK, &res.Nama, &res.TglLahir, &res.NamaIbu)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			// NIK tidak terdaftar -> bukan error, cukup "tidak ketemu".
			return entity.Resident{}, false, nil
		}
		// Error teknis -> teruskan ke usecase (dibungkus apperror.Internal di sana).
		return entity.Resident{}, false, fmt.Errorf("query residents: %w", err)
	}

	return res, true, nil
}
